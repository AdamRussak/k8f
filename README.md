# k8f
A CLI tool to find, list, connect, search and check version for K8S Clusters in all your resources at once,
this tool supports **Azure AKS** and **AWS EKS**.  
currently this tool supports the following commands:  
**list** - to list all the Managed K8S in your accounts  
**connect** - to generate Kubeconfig to all the managed K8S in your accounts.  

## prerequisite:
- for Azure: installed and logged in azure cli  
- for AWS: install AWS cli and Profiles for each Account at `~/.aws/credentials`  

## Supported Platform:
AWS  
Azure<br>

## Commands

### list
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
      --role-name string   Set Role Name (Example: '')

Global Flags:
      --aws-region string   Set Default AWS Region (default "eu-west-1")
  -v, --verbose             verbose logging
```
