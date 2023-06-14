package cmd

const (
	// global
	providerError     = "requires cloud provider"
	providerListError = "invalid cloud provider specified: %s"
	// connect
	connectCMD     = "connect"
	connectShort   = "Connect to all the clusters of a provider or all Supported Providers"
	connectExample = `k8f connect aws -p ./testfiles/config --backup -v
	k8f connect aws --isEnv -p ./testfiles/config --overwrite --backup --role-name "test role" -v`

	//find
	findCMD     = "find"
	findShort   = "Find if a specific K8S exist in Azure or AWS"
	findExample = `k8f find {aws/azure/all} my-k8s-cluster`
	//list
	listCMD     = "list"
	listShort   = "List all K8S in Azure/AWS or Both"
	listExample = `k8f list {aws/azure/all}`
)
