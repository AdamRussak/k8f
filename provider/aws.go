package provider

import (
	"context"
	"k8f/core"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/stscreds/#:~:text=or%20service%20clients.-,Assume%20Role,-To%20assume%20an
func (c CommandOptions) FullAwsList() Provider {
	var f []Account
	core.CheckEnvVarOrSitIt("AWS_REGION", c.AwsRegion)
	profiles := GetLocalAwsProfiles()
	addOnVersion := profiles[0].getVersion()
	l := getLatestEKS(getEKSversionsList(addOnVersion))
	log.Trace(profiles)
	c0 := make(chan Account)
	for _, profile := range profiles {
		go func(c0 chan Account, profile AwsProfiles, l string) {
			var re []Cluster
			log.Info(string("Using AWS profile: " + profile.Name))
			// opt := session.Options{Profile: profile.Name}
			opt := config.WithSharedConfigProfile(profile.Name)
			conf, err := config.LoadDefaultConfig(context.TODO(), opt)
			core.OnErrorFail(err, awsErrorMessage)
			// s := session.Must(conf, err)
			regions := profile.listRegions(conf)
			c2 := make(chan []Cluster)
			for _, reg := range regions {
				go printOutResult(reg, l, profile, addOnVersion, c2)
			}
			for i := 0; i < len(regions); i++ {
				aRegion := <-c2
				if len(aRegion) > 0 {
					re = append(re, aRegion...)
				}
			}
			c0 <- Account{profile.Name, re, len(re), ""}
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
func (p AwsProfiles) getVersion() *eks.DescribeAddonVersionsOutput {
	var svc *eks.Client
	opt := config.WithSharedConfigProfile(p.Name)
	conf, err := config.LoadDefaultConfig(context.TODO(), opt)
	core.OnErrorFail(err, "Failed to get Version")
	if p.IsRole {
		svc = eks.NewFromConfig(conf, func(o *eks.Options) { modifyOptions(p, conf) })
	} else {
		svc = eks.NewFromConfig(conf)
	}
	input2 := &eks.DescribeAddonVersionsInput{}
	r, err := svc.DescribeAddonVersions(context.TODO(), input2)
	core.OnErrorFail(err, "Failed to get Describe Version")
	return r
}

// gets the latest form suppported Addons
func getLatestEKS(addons []string) string {
	return evaluateVersion(addons)
}

// create Version list
func getEKSversionsList(addons *eks.DescribeAddonVersionsOutput) []string {
	var supportList []string
	for _, a := range addons.Addons {
		for _, c := range a.AddonVersions {
			for _, v := range c.Compatibilities {
				supportList = append(supportList, *v.ClusterVersion)
			}
		}
	}
	return supportList
}

// get installed Version on existing Clusters
func (p AwsProfiles) getEksCurrentVersion(cluster string, s aws.Config, reg string, c3 chan []string) {
	var svc *eks.Client
	if p.IsRole {
		svc = eks.NewFromConfig(s, func(o *eks.Options) { modifyOptions(p, s) })
	} else {
		svc = eks.NewFromConfig(s)
	}
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	}
	result, err := svc.DescribeCluster(context.TODO(), input)
	core.OnErrorFail(err, "Failed to Get Cluster Info")
	c3 <- []string{cluster, *result.Cluster.Version}
}

// get all Regions avilable
func (p AwsProfiles) listRegions(s aws.Config) []string {
	core.CheckEnvVarOrSitIt("AWS_REGION", Kregion)
	var reg []string
	var svc *ec2.Client
	if p.IsRole {
		svc = ec2.NewFromConfig(aws.Config{}, func(o *ec2.Options) { modifyOptions(p, s) })
	} else {
		svc = ec2.NewFromConfig(s)
	}
	input := &ec2.DescribeRegionsInput{}
	result, err := svc.DescribeRegions(context.TODO(), input, func(o *ec2.Options) { o.Region = Kregion })
	core.OnErrorFail(err, "Failed Get Region info")
	for _, r := range result.Regions {
		reg = append(reg, *r.RegionName)
	}
	return reg
}

func printOutResult(reg string, latest string, profile AwsProfiles, addons *eks.DescribeAddonVersionsOutput, c chan []Cluster) {
	var loc []Cluster
	var svc *eks.Client
	conf, err := config.LoadDefaultConfig(context.TODO())
	core.OnErrorFail(err, awsErrorMessage)
	if profile.IsRole {
		svc = eks.NewFromConfig(conf, func(o *eks.Options) {
			stsAssumeRole(profile, conf)
			o.Region = reg
		})
	} else {
		svc = eks.NewFromConfig(conf)
	}
	input := &eks.ListClustersInput{}
	result, err := svc.ListClusters(context.TODO(), input)
	core.OnErrorFail(err, "Failed to list Clusters")
	log.Debug(string("We are In Region: " + reg + " Profile " + profile.Name))
	// URGENT: need to fix region settings in search
	if len(result.Clusters) > 0 {
		c3 := make(chan []string)
		for _, element := range result.Clusters {
			go profile.getEksCurrentVersion(element, conf, reg, c3)
		}
		for i := 0; i < len(result.Clusters); i++ {
			res := <-c3
			loc = append(loc, Cluster{res[0], res[1], latest, reg, "", "", HowManyVersionsBack(getEKSversionsList(addons), res[1])})
		}
	}
	c <- loc
}

func GetLocalAwsProfiles() []AwsProfiles {
	var arr []AwsProfiles
	fname := config.DefaultSharedCredentialsFilename() // Get aws.config default shared credentials file name
	f, err := ini.Load(fname)
	core.OnErrorFail(err, "Failed to load profile")
	for _, v := range f.Sections() {
		if len(v.Keys()) != 0 {
			kbool, karn := checkIfItsAssumeRole(v.Keys())
			if kbool {
				arr = append(arr, AwsProfiles{Name: v.Name(), IsRole: true, Arn: karn})
			} else {
				arr = append(arr, AwsProfiles{Name: v.Name(), IsRole: false})
			}

		}
	}
	return (arr) // Create JSON string response
}

// Connect Logic
func (c CommandOptions) ConnectAllEks() AllConfig {
	var auth []Users
	var contexts []Contexts
	var clusters []Clusters
	var arnContext string
	core.CheckEnvVarOrSitIt("AWS_REGION", c.AwsRegion)
	p := c.FullAwsList()
	for _, a := range p.Accounts {
		r := make(chan LocalConfig)
		for _, clus := range a.Clusters {
			go func(r chan LocalConfig, clus Cluster, a Account, commandOptions CommandOptions) {
				cfg, err := config.LoadDefaultConfig(context.TODO())
				core.OnErrorFail(err, awsErrorMessage)
				eksSvc := eks.NewFromConfig(cfg, func(o *eks.Options) {
					o.Region = clus.Region
				})
				input := &eks.DescribeClusterInput{
					Name: aws.String(clus.Name),
				}
				result, err := eksSvc.DescribeCluster(context.TODO(), input)
				core.OnErrorFail(err, "Error calling DescribeCluster")
				r <- GenerateKubeConfiguration(result.Cluster, clus.Region, a, commandOptions)
			}(r, clus, a, c)
		}
		for i := 0; i < len(a.Clusters); i++ {
			result := <-r
			arnContext = result.Context.Cluster
			auth = append(auth, Users{Name: arnContext, User: result.Authinfo})
			contexts = append(contexts, Contexts{Name: arnContext, Context: result.Context})
			clusters = append(clusters, Clusters{Name: arnContext, Cluster: result.Cluster})
		}
	}
	if !c.Combined {
		log.Println("Started aws only config creation")
		c.Merge(AllConfig{auth, contexts, clusters}, arnContext)
		return AllConfig{}
	}
	log.Println("Started aws combined config creation")
	return AllConfig{auth, contexts, clusters}

}

// Create AWS Config
func GenerateKubeConfiguration(cluster *types.Cluster, r string, a Account, c CommandOptions) LocalConfig {
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
			Args:       c.AwsArgs(r, *cluster.Name, *cluster.Arn),
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

func (c CommandOptions) AwsArgs(region string, clusterName string, arn string) []string {
	var args []string
	if c.AwsRoleString != "" && !c.AwsAuth {
		args = []string{"--region", region, "eks", "get-token", "--cluster-name", clusterName, "--role-arn", "arn:aws:iam::" + SplitAzIDAndGiveItem(arn, ":", 4) + ":role/" + c.AwsRoleString}
	} else if c.AwsRoleString != "" && c.AwsAuth {
		args = []string{"token", "-i", clusterName, "--role-arn", "arn:aws:iam::" + SplitAzIDAndGiveItem(arn, ":", 4) + ":role/" + c.AwsRoleString}
	} else {
		args = []string{"--region", region, "eks", "get-token", "--cluster-name", clusterName}
	}
	return args
}

func (c CommandOptions) AwsEnvs(profile string) interface{} {
	if c.AwsEnvProfile {
		var envArray []Env
		envArray = append(envArray, Env{Name: "AWS_PROFILE", Value: profile})
		return envArray
	}
	return nil
}

func (c CommandOptions) GetSingleAWSCluster(clusterToFind string) Cluster {
	log.Info("Starting AWS find cluster named: " + clusterToFind)
	core.CheckEnvVarOrSitIt("AWS_REGION", c.AwsRegion)
	//get Profiles//search this name in account
	var f Cluster
	profiles := GetLocalAwsProfiles()
	c0 := make(chan Cluster)
	// search each profile
	for _, profile := range profiles {
		go getAwsClusters(c0, profile, clusterToFind)
	}
	for i := 0; i < len(profiles); i++ {
		res := <-c0
		if res.Name == clusterToFind {
			f = res
		}
	}
	return f
	//search this name in region
	//once it is found erturn info to the user
}

func getAwsClusters(c0 chan Cluster, profile AwsProfiles, clusterToFind string) {
	var re Cluster
	log.Info(string("Using AWS profile: " + profile.Name))
	// opt := session.Options{Profile: profile.Name}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	core.OnErrorFail(err, awsErrorMessage)
	regions := profile.listRegions(cfg)
	profiles := GetLocalAwsProfiles()
	addOnVersion := profiles[0].getVersion()
	c2 := make(chan []Cluster)
	for _, reg := range regions {
		go printOutResult(reg, clusterToFind, profile, addOnVersion, c2)
	}
	for i := 0; i < len(regions); i++ {
		aRegion := <-c2
		if len(aRegion) > 0 {
			for _, cluster := range aRegion {
				if cluster.Name == clusterToFind {
					re = cluster
				}
			}
		}
	}
	c0 <- re
}

func checkIfItsAssumeRole(keys []*ini.Key) (bool, string) {
	log.Debug(keys)
	var ARNRegexp = regexp.MustCompile(`^arn:(\w|-)*:iam::\d+:role\/?(\w+|-|\/|\.)*$`)
	for _, a := range keys {
		if ARNRegexp.MatchString(a.String()) {
			log.Debug("Is ARN: " + a.String())
			return true, a.String()
		}
	}
	return false, ""
}

func stsAssumeRole(awsProfile AwsProfiles, session aws.Config) *stscreds.AssumeRoleProvider {
	log.Debug("checking if this profile is a role")
	if awsProfile.IsRole {
		log.Debug("it is a role")
		client := sts.NewFromConfig(session)
		appCreds := stscreds.NewAssumeRoleProvider(client, awsProfile.Arn)
		return appCreds
	}
	log.Debug("it is NOT a role")
	return nil

}

func modifyOptions(p AwsProfiles, s aws.Config) *ec2.Options {
	// Modify the options as needed.
	options := ec2.Options{}
	options.Credentials = stsAssumeRole(p, s)

	return &options
}
