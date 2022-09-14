package provider

import (
	"k8f/core"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func (c CommandOptions) FullAwsList() Provider {
	var f []Account
	core.CheckEnvVarOrSitIt("AWS_REGION", c.AwsRegion)
	l := getLatestEKS(getVersion())
	profiles := GetLocalAwsProfiles()
	c0 := make(chan Account)
	for _, profile := range profiles {
		go func(c0 chan Account, profile string, l string) {
			var re []Cluster
			log.Info(string("Using AWS profile: " + profile))
			opt := session.Options{Profile: profile}
			conf, err := session.NewSessionWithOptions(opt)
			core.OnErrorFail(err, "Failed to create new session")
			s := session.Must(conf, err)
			regions := listRegions(s)
			c2 := make(chan []Cluster)
			for _, reg := range regions {
				go printOutResult(reg, l, profile, c2)
			}
			for i := 0; i < len(regions); i++ {
				aRegion := <-c2
				if len(aRegion) > 0 {
					re = append(re, aRegion...)
				}
			}
			c0 <- Account{profile, re, len(re), ""}
		}(c0, profile, l)

	}
	for i := 0; i < len(profiles); i++ {
		res := <-c0
		if len(res.Clusters) != 0 {
			f = append(f, res)
		}
	}
	return Provider{"aws", f, countTotal(f)}
}

// get Addons Supported EKS versions
func getVersion() *eks.DescribeAddonVersionsOutput {
	s, err := session.NewSession()
	core.OnErrorFail(err, "Failed to get Version")
	svc := eks.New(s)
	input2 := &eks.DescribeAddonVersionsInput{}
	r, err := svc.DescribeAddonVersions(input2)
	core.OnErrorFail(err, "Failed to get Describe Version")
	return r
}

// gets the latest form suppported Addons
func getLatestEKS(addons *eks.DescribeAddonVersionsOutput) string {
	var supportList []string
	for _, a := range addons.Addons {
		for _, c := range a.AddonVersions {
			for _, v := range c.Compatibilities {
				supportList = append(supportList, *v.ClusterVersion)
			}
		}
	}
	return evaluateVersion(supportList)
}

// get installed Version on existing Clusters
func getEksCurrentVersion(cluster string, s *session.Session, reg string, c3 chan []string) {
	svc := eks.New(s)
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	}
	result, err := svc.DescribeCluster(input)
	core.OnErrorFail(err, "Failed to Get Cluster Info")
	c3 <- []string{cluster, *result.Cluster.Version}
}

// get all Regions avilable
func listRegions(s *session.Session) []string {
	var reg []string
	svc := ec2.New(s)
	input := &ec2.DescribeRegionsInput{}

	result, err := svc.DescribeRegions(input)
	core.OnErrorFail(err, "Failed Get Region info")
	for _, r := range result.Regions {
		reg = append(reg, *r.RegionName)
	}
	return reg
}

func printOutResult(reg string, latest string, profile string, c chan []Cluster) {
	var loc []Cluster
	opt := session.Options{Profile: profile, Config: aws.Config{Region: aws.String(reg)}}
	conf, err := session.NewSessionWithOptions(opt)
	core.OnErrorFail(err, "Failed to create new session")
	sess := session.Must(conf, err)
	svc := eks.New(sess)
	input := &eks.ListClustersInput{}
	result, err := svc.ListClusters(input)
	core.OnErrorFail(err, "Failed to list Clusters")
	log.Debug(string("We are In Region: " + reg + " Profile " + profile))
	if len(result.Clusters) > 0 {
		c3 := make(chan []string)
		for _, element := range result.Clusters {
			go getEksCurrentVersion(*element, sess, reg, c3)
		}
		for i := 0; i < len(result.Clusters); i++ {
			res := <-c3
			loc = append(loc, Cluster{res[0], res[1], latest, reg, "", ""})
		}
	}
	c <- loc
}

func GetLocalAwsProfiles() []string {
	arr := []string{}
	fname := config.DefaultSharedCredentialsFilename() // Get aws.config default shared credentials file name
	f, err := ini.Load(fname)                          // Load ini file
	core.OnErrorFail(err, "Failed to load profile")
	for _, v := range f.Sections() {
		if len(v.Keys()) != 0 {
			arr = append(arr, v.Name()) // Get only the sections having Keys. Not sure why this is returning DEFAULT here
		}
	}

	return (arr) // Create JSON string response
}

// Connect Logic
func (c CommandOptions) ConnectAllEks() AllConfig {
	var auth []Users
	var context []Contexts
	var clusters []Clusters
	var arnContext string
	core.CheckEnvVarOrSitIt("AWS_REGION", c.AwsRegion)
	p := c.FullAwsList()
	for _, a := range p.Accounts {
		r := make(chan LocalConfig)
		for _, clus := range a.Clusters {
			go func(r chan LocalConfig, clus Cluster, a Account, commandOptions CommandOptions) {
				opt := session.Options{Profile: a.Name, Config: aws.Config{
					Region: aws.String(clus.Region),
				}}
				sess := session.Must(session.NewSessionWithOptions(opt))
				eksSvc := eks.New(sess)
				input := &eks.DescribeClusterInput{
					Name: aws.String(clus.Name),
				}
				result, err := eksSvc.DescribeCluster(input)
				core.OnErrorFail(err, "Error calling DescribeCluster")
				r <- GenerateKubeConfiguration(result.Cluster, clus.Region, a, commandOptions)
			}(r, clus, a, c)
		}
		for i := 0; i < len(a.Clusters); i++ {
			result := <-r
			arnContext = result.Context.Cluster
			auth = append(auth, Users{Name: arnContext, User: result.Authinfo})
			context = append(context, Contexts{Name: arnContext, Context: result.Context})
			clusters = append(clusters, Clusters{Name: arnContext, Cluster: result.Cluster})
		}
	}
	if !c.Combined {
		log.Println("Started aws only config creation")
		c.Merge(AllConfig{auth, context, clusters}, arnContext)
		return AllConfig{}
	}
	log.Println("Started aws combined config creation")
	return AllConfig{auth, context, clusters}

}

// Create AWS Config
func GenerateKubeConfiguration(cluster *eks.Cluster, r string, a Account, c CommandOptions) LocalConfig {
	clusters := CCluster{
		Server:                   *cluster.Endpoint,
		CertificateAuthorityData: *cluster.CertificateAuthority.Data,
	}
	contexts := Context{
		Cluster: *cluster.Arn,
		User:    *cluster.Arn,
	}

	authinfos := User{
		Exec: Exec{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Args:       c.AwsArgs(r, *cluster.Name),
			Env:        c.AwsEnvs(a.Name),
			Command:    c.setCommand(),
		},
	}
	return LocalConfig{authinfos, contexts, clusters}
}

func (c CommandOptions) setCommand() string {
	if c.AwsAuth {
		return "aws-iam-authenticator"
	}
	return "aws"
}
func (c CommandOptions) AwsArgs(region string, clusterName string) []string {
	var args []string
	if c.AwsRoleString != "" && !c.AwsAuth {
		args = []string{"--region", region, "eks", "get-token", "--cluster-name", clusterName, "- --role-arn", c.AwsRoleString}
	} else if c.AwsRoleString != "" && c.AwsAuth {
		args = []string{"token", "-i", clusterName, "- --role-arn", c.AwsRoleString}
	} else {
		args = []string{"--region", region, "eks", "get-token", "--cluster-name", clusterName}
	}
	return args
}

func (c CommandOptions) AwsEnvs(profile string) interface{} {
	if c.AwsEnvProfile {
		env := Env{Name: "AWS_PROFILE", Value: profile}
		return env
	}
	return nil
}

// The format for the config
// https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
// args:
// - --region
// - $region_code
// - eks
// - get-token
// - --cluster-name
// - $cluster_name
// # - "- --role-arn"
// # - "arn:aws:iam::$account_id:role/my-role"
// # env:
// # - name: "AWS_PROFILE"
// #   value: "aws-profile"
