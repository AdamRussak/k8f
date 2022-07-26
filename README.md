<img src="https://raw.githubusercontent.com/AdamRussak/k8f/main/examples/background.png" data-canonical-src="https://raw.githubusercontent.com/AdamRussak/k8f/main/examples/background.png"  width="500" height="200" />

> image created Using Dall-E<br>
> 
[![CodeQL](https://github.com/AdamRussak/k8f/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/AdamRussak/k8f/actions/workflows/codeql-analysis.yml)  [![release-artifacts](https://github.com/AdamRussak/k8f/actions/workflows/release-new-version.yaml/badge.svg)](https://github.com/AdamRussak/k8f/actions/workflows/release-new-version.yaml) ![GitHub](https://img.shields.io/github/license/AdamRussak/k8f) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AdamRussak/k8f) ![GitHub all releases](https://img.shields.io/github/downloads/AdamRussak/k8f/total) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=AdamRussak_k8f&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=AdamRussak_k8f) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=AdamRussak_k8f&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=AdamRussak_k8f)
# k8f
A CLI tool to *find*, *list*, *connect* and check version for K8S Clusters in all your resources at once,
in a single command
this tool supports **Azure AKS**, **AWS EKS** and Partily supports **GCP GKE**  
currently this tool supports the following commands:

* **list** - to list all the Managed K8S in your accounts and info about there Version  
* **find** - to find a cluster in an unknown region/account. 
>currently only supports Azure and AWS
* **connect** - to generate Kubeconfig to all the managed K8S in your accounts.  
>currently only supports Azure and AWS

## prerequisite:
- for Azure: installed and logged in azure cli  
- for AWS: install AWS cli and Profiles for each Account at `~/.aws/credentials`  
- for GCP: Installed gcloud cli and logged in

## Supported Platform:
- [ ] AWS  
- [ ] Azure
- [ ] GCP
>GCP currently only supports List command
## Commands
###  list
```sh
List all K8S in Azure/AWS or Both

Usage:
  k8f list [flags]

Examples:
k8f list {aws/azure/all}

Flags:
  -h, --help            help for list
  -o, --output string   Set output type(json or yaml) (default "json")

Global Flags:
      --aws-region string   Set Default AWS Region (default "eu-west-1")
  -v, --verbose             verbose logging
```

List Command Sample Output:

[![Sample of List command output](https://raw.githubusercontent.com/AdamRussak/k8f/main/examples/k8f-list.jpg "Sample of List command output")](https://raw.githubusercontent.com/AdamRussak/k8f/main/examples/k8f-list.jpg "Sample of List command output")

###  connect
```sh
Connect to all the clusters of a provider or all Supported Providers

Usage:
  k8f connect [flags]

Examples:
k8f connect aws -p ./testfiles/config --backup -v
k8f connect aws --isEnv -p ./testfiles/config --overwrite --backup --role-name "test role" -v

Flags:
      --auth               change from CLI Auth to AMI Auth, Default set to CLI
      --backup             If true, backup config file to $HOME/.kube/config.bk
      --dry-run            If true, only run a dry-run with cli output
  -h, --help               help for connect
      --isEnv              Add AWS Env Profile to the AWS Config
  -o, --output string      Merged kubeconfig output type(json or yaml) (default "yml")
      --overwrite          If true, force merge kubeconfig
  -p, --path string        Merged kubeconfig output path (default "/home/vscode/.kube/config")
      --role-name string   Set Role Name (Example: 'myRoleName')

Global Flags:
      --aws-region string   Set Default AWS Region (default "eu-west-1")
  -v, --verbose             verbose logging
```

###  find
```sh
Find if a specific K8S exist in Azure or AWS

Usage:
  k8f find [flags]

Examples:
k8f find {aws/azure/all} my-k8s-cluster

Flags:
  -h, --help   help for find

Global Flags:
      --aws-region string   Set Default AWS Region (default "eu-west-1")
  -v, --verbose             verbose logging
```