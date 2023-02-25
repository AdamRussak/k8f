package provider

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"k8f/core"
	"log"
	"os"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/imdario/mergo"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

func (mc CommandOptions) runMerge(newConf Config) error {
	var kconfigs []*clientcmdapi.Config
	confToupdate, err := toClientConfig(&newConf)
	kconfigs = append(kconfigs, confToupdate)
	core.OnErrorFail(err, "failed to convert Config struct to clientcmdapi.Config")
	outConfigs := clientcmdapi.NewConfig()
	printString(os.Stdout, "Loading KubeConfig file: "+mc.Path+" \n")
	loadConfig, err := loadKubeConfig(mc.Path)
	core.OnErrorFail(err, "File "+mc.Path+" is not kubeconfig\n")
	kconfigs = append(kconfigs, loadConfig)
	for _, conf := range kconfigs {
		kco := &KubeConfigOption{
			config:   conf,
			fileName: getFileName(mc.Path),
		}
		outConfigs, err = kco.handleContexts(outConfigs, mc)
		if err != nil {
			return err
		}
	}
	var y []byte
	y, _ = yaml.Marshal(outConfigs)
	err = os.WriteFile(mc.Path, y, 0666)
	core.OnErrorFail(err, "failed to save config")
	return nil
}

func loadKubeConfig(yaml string) (*clientcmdapi.Config, error) {
	loadConfig, err := clientcmd.LoadFromFile(yaml)
	if err != nil {
		return nil, err
	}
	if len(loadConfig.Contexts) == 0 {
		return nil, fmt.Errorf("no kubeconfig in %s ", yaml)
	}
	return loadConfig, err
}

func CheckValidContext(clear bool, config *clientcmdapi.Config) *clientcmdapi.Config {
	for key, obj := range config.Contexts {
		if _, ok := config.AuthInfos[obj.AuthInfo]; !ok {
			if clear {
				printString(os.Stdout, fmt.Sprintf("clear lapsed AuthInfo [%s]\n", obj.AuthInfo))
			} else {
				printYellow(os.Stdout, fmt.Sprintf("WARNING: AuthInfo 「%s」 has no matching context 「%s」, please run `kubecm clear` to clean up this Context.\n", obj.AuthInfo, key))
			}
			delete(config.Contexts, key)
			delete(config.Clusters, obj.Cluster)
		}
		if _, ok := config.Clusters[obj.Cluster]; !ok {
			if clear {
				printString(os.Stdout, fmt.Sprintf("clear lapsed Cluster [%s]\n", obj.Cluster))
			} else {
				printYellow(os.Stdout, fmt.Sprintf("WARNING: Cluster 「%s」 has no matching context 「%s」, please run `kubecm clear` to clean up this Context.\n", obj.Cluster, key))
			}
			delete(config.Contexts, key)
			delete(config.AuthInfos, obj.AuthInfo)
		}
	}
	return config
}
func printYellow(out io.Writer, content string) {
	ct.ChangeColor(ct.Yellow, false, ct.None, false)
	fmt.Fprint(out, content)
	ct.ResetColor()
}
func printString(out io.Writer, name string) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Fprint(out, name)
	ct.ResetColor()
}

// the actule writing of the config to the OS
func (mc CommandOptions) WriteConfig(outConfig *clientcmdapi.Config) {
	err := clientcmd.WriteToFile(*outConfig, mc.Path)
	core.OnErrorFail(err, "failed to save new config")
	fmt.Printf("「%s」 write successful!\n", mc.Path)
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
		if len(kc.config.Contexts) > 1 {
			newName = name
		} else {
			newName = kc.fileName
		}
		if checkContextName(newName, oldConfig) {
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
		fmt.Printf("Add Context: %s \n", newName)
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
		cluster := clientcmdapi.Cluster{
			Server:                   c.Cluster.Server,
			CertificateAuthorityData: []byte(c.Cluster.CertificateAuthorityData),
		}
		clientConfig.Clusters[c.Name] = &cluster
	}

	// Set users
	for _, u := range cfg.Users {
		user := clientcmdapi.AuthInfo{
			ClientCertificateData: []byte(u.User.ClientCertificateData),
			ClientKeyData:         []byte(u.User.ClientKeyData),
			Token:                 u.User.Token,
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
