package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"gopkg.in/ini.v1"
)

func MainAWS() {
	var f []account
	l := getLatestEKS(getVersion())
	profiles := GetLocalAwsProfiles()
	for _, profile := range profiles {
		var re []region
		var count int
		fmt.Println("Using this profile: ", profile)
		opt := session.Options{Profile: profile}
		conf, err := session.NewSessionWithOptions(opt)
		if err != nil {
			fmt.Println(err)
		}
		s := session.Must(conf, err)
		regions := listRegions(s)
		for _, reg := range regions {
			aRegion := printOutResult(reg, l, profile)
			if aRegion.TotalCount > 0 {
				log.Println("Region:", reg, " Cluster Count: ", aRegion.TotalCount)
				re = append(re, aRegion)
				count = count + aRegion.TotalCount
			}
		}
		f = append(f, account{profile, report{re, count}})
	}
	kJson, _ := json.Marshal(f)
	fmt.Println(string(kJson))
}

//get Addons Supported EKS versions
func getVersion() *eks.DescribeAddonVersionsOutput {
	s, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	svc := eks.New(s)
	input2 := &eks.DescribeAddonVersionsInput{}
	r, err := svc.DescribeAddonVersions(input2)
	if err != nil {
		fmt.Println(err)
	}
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
func getEksCurrentVersion(cluster string, s *session.Session) string {
	svc := eks.New(s)
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "err"
	}

	return *result.Cluster.Version
}

//get all Regions avilable
func listRegions(s *session.Session) []string {
	var reg []string
	svc := ec2.New(s)
	input := &ec2.DescribeRegionsInput{}

	result, err := svc.DescribeRegions(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	for _, r := range result.Regions {
		reg = append(reg, *r.RegionName)
	}
	return reg
}

func printOutResult(reg string, latest string, profile string) region {
	var loc []cluster
	opt := session.Options{Profile: profile, Config: aws.Config{Region: aws.String(reg)}}
	conf, err := session.NewSessionWithOptions(opt)
	if err != nil {
		fmt.Println(err)
	}
	sess := session.Must(conf, err)
	svc := eks.New(sess)
	input := &eks.ListClustersInput{}
	result, err := svc.ListClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println("We are In Region: ", reg)
	if len(result.Clusters) > 0 {
		for _, element := range result.Clusters {
			c := getEksCurrentVersion(*element, sess)
			loc = append(loc, cluster{*element, c, latest})
		}
	}
	return region{reg, loc, len(loc)}
}

func GetLocalAwsProfiles() []string {
	arr := []string{}
	fname := config.DefaultSharedCredentialsFilename() // Get aws.config default shared credentials file name
	f, err := ini.Load(fname)                          // Load ini file
	if err != nil {
		fmt.Println(err)
	} else {
		for _, v := range f.Sections() {
			if len(v.Keys()) != 0 {
				arr = append(arr, v.Name()) // Get only the sections having Keys. Not sure why this is returning DEFAULT here
			}
		}
	}
	return (arr) // Create JSON string response
}
