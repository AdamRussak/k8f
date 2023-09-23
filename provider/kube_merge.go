package provider

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"k8f/core"
	"strings"

	"github.com/imdario/mergo"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// MergeCommand merge cmd struct
type MergeCommand struct {
	cobra.Command
}

// KubeConfigOption kubeConfig option
type KubeConfigOption struct {
	config   *clientcmdapi.Config
	fileName string
}

func (mc CommandOptions) runMerge(newConf Config) ([]byte, error) {
	var kconfigs []*clientcmdapi.Config
	confToupdate, err := toClientConfig(&newConf)
	kconfigs = append(kconfigs, confToupdate)
	core.FailOnError(err, "failed to convert Config struct to clientcmdapi.Config")
	outConfigs := clientcmdapi.NewConfig()
	log.Infof("Loading KubeConfig file: %s\n", mc.Path)
	loadConfig, err := loadKubeConfig(mc.Path)
	core.FailOnError(err, "File "+mc.Path+" is not kubeconfig\n")
	kconfigs = append(kconfigs, loadConfig)
	for _, conf := range kconfigs {
		kco := &KubeConfigOption{
			config:   conf,
			fileName: getFileName(mc.Path),
		}
		outConfigs, err = kco.handleContexts(outConfigs, mc)
		if err != nil {
			return nil, err
		}
	}
	outConfigs.APIVersion = "v1"
	outConfigs.Kind = "Config"
	var contxetcItem *clientcmdapi.Context
	for _, value := range outConfigs.Contexts {
		contxetcItem = value
		break
	}
	outConfigs.CurrentContext = contxetcItem.Cluster
	confByte, err := clientcmd.Write(*outConfigs)
	if err != nil {
		return nil, err
	}
	return confByte, nil
}

func loadKubeConfig(yaml string) (*clientcmdapi.Config, error) {
	loadConfig, err := clientcmd.LoadFromFile(yaml)
	if err != nil {
		return nil, err
	}
	return loadConfig, err
}

func getFileName(path string) string {
	n := strings.Split(path, "/")
	result := strings.Split(n[len(n)-1], ".")
	return result[0]
}

func (kc *KubeConfigOption) handleContexts(oldConfig *clientcmdapi.Config, mc CommandOptions) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	for name, ctx := range kc.config.Contexts {
		var newName string
		if len(kc.config.Contexts) >= 1 {
			newName = name
		} else {
			newName = kc.fileName
		}
		if checkContextName(newName, oldConfig) && !mc.ForceMerge {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it. (If you select `False`, this context will not be merged)", newName), mc)
			if nameConfirm == "True" {
				newName = core.PromptUI("Rename", newName)
				if newName == kc.fileName {
					return nil, errors.New("need to rename")
				}
			} else {
				continue
			}
		}
		itemConfig := kc.handleContext(oldConfig, newName, ctx)
		newConfig = appendConfig(newConfig, itemConfig)
		log.Infof("Add Context: %s \n", newName)
	}
	outConfig := appendConfig(oldConfig, newConfig)
	return outConfig, nil
}
func checkContextName(name string, oldConfig *clientcmdapi.Config) bool {
	if _, ok := oldConfig.Contexts[name]; ok {
		return true
	}
	return false
}
func BoolUI(label string, mc CommandOptions) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F37A {{ . | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F47B {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     []string{"False", "True"},
		Templates: templates,
		Size:      mc.UiSize,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}

// HashSufString return the string of HashSuf.
func HashSufString(data string) string {
	sum, _ := hEncode(Hash(data))
	return sum
}
func (kc *KubeConfigOption) handleContext(oldConfig *clientcmdapi.Config,
	name string, ctx *clientcmdapi.Context) *clientcmdapi.Config {

	var (
		clusterNameSuffix string
		userNameSuffix    string
	)

	isClusterNameExist, isUserNameExist := checkClusterAndUserName(oldConfig, ctx.Cluster, ctx.AuthInfo)
	newConfig := clientcmdapi.NewConfig()
	suffix := HashSufString(name)

	if isClusterNameExist {
		clusterNameSuffix = "-" + suffix
	}
	if isUserNameExist {
		userNameSuffix = "-" + suffix
	}

	userName := fmt.Sprintf("%v%v", ctx.AuthInfo, userNameSuffix)
	clusterName := fmt.Sprintf("%v%v", ctx.Cluster, clusterNameSuffix)
	newCtx := ctx.DeepCopy()
	newConfig.AuthInfos[userName] = kc.config.AuthInfos[newCtx.AuthInfo]
	newConfig.Clusters[clusterName] = kc.config.Clusters[newCtx.Cluster]
	newConfig.Contexts[name] = newCtx
	newConfig.Contexts[name].AuthInfo = userName
	newConfig.Contexts[name].Cluster = clusterName

	return newConfig
}
func checkClusterAndUserName(oldConfig *clientcmdapi.Config, newClusterName, newUserName string) (bool, bool) {
	var (
		isClusterNameExist bool
		isUserNameExist    bool
	)

	for _, ctx := range oldConfig.Contexts {
		if ctx.Cluster == newClusterName {
			isClusterNameExist = true
		}
		if ctx.AuthInfo == newUserName {
			isUserNameExist = true
		}
	}

	return isClusterNameExist, isUserNameExist
}

// Copied from https://github.com/kubernetes/kubernetes
// /blob/master/pkg/kubectl/util/hash/hash.go
func hEncode(hex string) (string, error) {
	if len(hex) < 10 {
		return "", fmt.Errorf(
			"input length must be at least 10")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}

func appendConfig(c1, c2 *clientcmdapi.Config) *clientcmdapi.Config {
	config := clientcmdapi.NewConfig()
	_ = mergo.Merge(config, c1)
	_ = mergo.Merge(config, c2)
	return config
}

// Hash returns the hex form of the sha256 of the argument.
func Hash(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

// Convert merged config to clientcmdapi.Config
func toClientConfig(cfg *Config) (*clientcmdapi.Config, error) {
	clientConfig := clientcmdapi.NewConfig()

	// Set API version
	clientConfig.APIVersion = cfg.APIVersion

	// Set current context
	clientConfig.CurrentContext = cfg.CurrentContext

	// Set clusters
	for _, c := range cfg.Clusters {
		decodedBytes, err := base64.StdEncoding.DecodeString(c.Cluster.CertificateAuthorityData)
		core.FailOnError(err, decodeError)
		cluster := clientcmdapi.Cluster{
			Server:                   c.Cluster.Server,
			CertificateAuthorityData: decodedBytes,
		}
		clientConfig.Clusters[c.Name] = &cluster
	}

	// Set users
	for _, u := range cfg.Users {
		user := getUserForCluster(u)
		clientConfig.AuthInfos[u.Name] = &user
	}

	// Set contexts
	for _, c := range cfg.Contexts {
		context := clientcmdapi.Context{
			Cluster:  c.Context.Cluster,
			AuthInfo: c.Context.User,
		}
		clientConfig.Contexts[c.Name] = &context
	}

	return clientConfig, nil
}

func getUserForCluster(u Users) clientcmdapi.AuthInfo {
	var user clientcmdapi.AuthInfo
	if !checkIfStructInit(u.User, "exec") {
		clientCertificateDataBytes, err := base64.StdEncoding.DecodeString(u.User.ClientCertificateData)
		core.FailOnError(err, decodeError)
		ClientKeyDataBytes, err := base64.StdEncoding.DecodeString(u.User.ClientKeyData)
		core.FailOnError(err, decodeError)
		user = clientcmdapi.AuthInfo{
			ClientCertificateData: []byte(clientCertificateDataBytes),
			ClientKeyData:         []byte(ClientKeyDataBytes),
			Token:                 u.User.Token,
		}
	} else {
		user = clientcmdapi.AuthInfo{
			Exec: &clientcmdapi.ExecConfig{
				APIVersion: u.User.Exec.APIVersion,
				Command:    u.User.Exec.Command,
				Args:       u.User.Exec.Args,
				Env:        []clientcmdapi.ExecEnvVar{},
			},
		}
		envs, _ := u.User.Exec.Env.([]Env)
		for _, env := range envs {
			user.Exec.Env = append(user.Exec.Env, clientcmdapi.ExecEnvVar{Name: env.Name, Value: env.Value})
		}
	}
	return user
}
