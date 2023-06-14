<img src="https://raw.githubusercontent.com/AdamRussak/public-images/main/k8f/k8f_logo.png" data-canonical-src="https://raw.githubusercontent.com/AdamRussak/public-images/main/k8f/k8f_logo.png"  width="500" height="200" />

> image created Using MidJorny<br>
> 
[![CodeQL](https://github.com/AdamRussak/k8f/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/AdamRussak/k8f/actions/workflows/codeql-analysis.yml)  [![release-artifacts](https://github.com/AdamRussak/k8f/actions/workflows/release-new-version.yaml/badge.svg)](https://github.com/AdamRussak/k8f/actions/workflows/release-new-version.yaml) ![GitHub](https://img.shields.io/github/license/AdamRussak/k8f) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AdamRussak/k8f) ![GitHub all releases](https://img.shields.io/github/downloads/AdamRussak/k8f/total) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=AdamRussak_k8f&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=AdamRussak_k8f) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=AdamRussak_k8f&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=AdamRussak_k8f)

[badge-info]: https://raw.githubusercontent.com/AdamRussak/public-images/main/badges/info.svg 'Info'

> ![badge-info][badge-info]<br>
> Tested with:<br>
> AWS CLI: 2.9.17 <br>
> AZ CLI: 2.44.1 <br>
> Kubectl: v1.26.1 <br>

<br>

# k8f

A CLI tool to *find*, *list*, *connect* and check version for K8S Clusters in all your resources at once,
in a single command

## Table of Contents

* [Prerequisite](#prerequisite)
* [Supported Platform](#supported-platform)
* [Commands](#commands)
  * [List](#list)
  * [Connect](#connect)
  * [Find](#find)
* [How to install](#how-to-install)
  * [Windows](#windows)
  * [Linux](#linux)
  * [MacOS](#macos)



## prerequisite
- for Azure: installed and logged in azure cli  
- for AWS: install AWS cli and Profiles for each Account at `~/.aws/credentials`  and `~/.aws/config`
- for GCP: Installed gcloud cli and logged in

Supported Platform
-----------
- [ ] AWS  
- [ ] Azure
- [ ] GCP
>GCP currently only supports List command
Commands
------
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

[![Sample of List command output](https://raw.githubusercontent.com/AdamRussak/public-images/main/k8f/k8f-list.jpg "Sample of List command output")](https://raw.githubusercontent.com/AdamRussak/public-images/main/k8f/k8f-list.jpg "Sample of List command output")

###  connect
```sh
Connect to all the clusters of a provider or all Supported Providers

Usage:
  k8f connect [flags]

Examples:
k8f connect aws -p ./testfiles/config --backup -v
k8f connect aws --isEnv -p ./testfiles/config --overwrite --backup --role-name "test role" -v

Flags:
      --auth               change from AWS CLI Auth to AWS IAM Authenticator, Default set to AWS CLI
      --backup             If true, backup config file to $HOME/.kube/config.bk
      --dry-run            If true, only run a dry-run with cli output
      --force-merge        If set, all duplication will be merged without prompt, default is interactive
  -h, --help               help for connect
      --isEnv              Add AWS Profile as Env to the Kubeconfig
      --merge              If true, add new K8s to the existing kubeconfig path
  -o, --output string      kubeconfig output type format(json or yaml) (default "yaml")
      --overwrite          If true, force overwrite kubeconfig
  -p, --path string        kubeconfig output path (default "/home/vscode/.kube/config")
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
How to install
-----------
### Windows
Latest:
```ps
$downloads = "$env:USERPROFILE\Downloads"
$source = "$downloads\k8f.exe"
$destination = "C:\tool"
Invoke-WebRequest -Uri "https://github.com/AdamRussak/k8f/releases/latest/download/k8f.exe" -OutFile $source
New-Item -ItemType Directory -Path $destination
Copy-Item -Path $source -Destination $destination
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$destination\k8f.exe", "Machine")
```
Version:
```ps
$downloads = "$env:USERPROFILE\Downloads"
$source = "$downloads\k8f.exe"
$destination = "C:\tool"
$version = "0.3.1"
Invoke-WebRequest -Uri "https://github.com/AdamRussak/k8f/releases/download/$version/k8f.exe" -OutFile "$downloads\k8f.exe"
New-Item -ItemType Directory -Path $destination
Copy-Item -Path $source -Destination $destination
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$destination\k8f.exe", "Machine")
```
### Linux
Latest:
```sh
cd ~ && wget https://github.com/AdamRussak/k8f/releases/latest/download/k8f
cp ~/k8f /usr/local/bin/k8f
sudo chmod 777 /usr/local/bin/k8f
```
Version:
```sh
cd ~ && wget https://github.com/AdamRussak/k8f/releases/download/<version>/k8f
cp ~/k8f /usr/local/bin/k8f
sudo chmod 777 /usr/local/bin/k8f
```
### MacOS
Latest:
```sh
cd ~ && wget https://github.com/AdamRussak/k8f/releases/latest/download/k8f_darwin-arm64
mv k8f_darwin-arm64 ./k8f
cp ~/k8f /usr/local/bin/k8f
sudo chmod 777 /usr/local/bin/k8f
```
Version:
```sh
cd ~ && wget https://github.com/AdamRussak/k8f/releases/download/<version>/k8f_darwin-arm64
mv k8f_darwin-arm64 ./k8f
cp ~/k8f /usr/local/bin/k8f
sudo chmod 777 /usr/local/bin/k8f
```