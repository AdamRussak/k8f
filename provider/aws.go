package provider

//TODO: get regions from the AWS CLI CONFIG
import (
	"k8-upgrade/core"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"gopkg.in/ini.v1"
)

func MainAWS() string {
	var f []Account
	l := getLatestEKS(getVersion())
	profiles := GetLocalAwsProfiles()
	c := make(chan Account)
	for _, profile := range profiles {
		go func(c chan Account, profile string, l string) {
			var re []Cluster
			log.Println("Using this profile: ", profile)
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
			c <- Account{profile, re, len(re)}
		}(c, profile, l)

	}
	for i := 0; i < len(profiles); i++ {
		res := <-c
		if len(res.Clusters) != 0 {
			f = append(f, res)
		}
	}
	var count int
	for _, a := range f {
		count = count + a.TotalCount
	}
	return runResult(Provider{"aws", f, count})
}

//get Addons Supported EKS versions
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
func getEksCurrentVersion(cluster string, s *session.Session, c3 chan []string) {
	svc := eks.New(s)
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	}
	result, err := svc.DescribeCluster(input)
	core.OnErrorFail(err, "Failed to Get Cluster Info")
	c3 <- []string{cluster, *result.Cluster.Version}
}

//get all Regions avilable
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
	log.Println("We are In Region: ", reg, "Profile", profile)
	if len(result.Clusters) > 0 {
		c3 := make(chan []string)
		for _, element := range result.Clusters {
			go getEksCurrentVersion(*element, sess, c3)
		}
		for i := 0; i < len(result.Clusters); i++ {
			res := <-c3
			loc = append(loc, Cluster{res[0], res[1], latest, reg})
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
