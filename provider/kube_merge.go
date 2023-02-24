package provider

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"k8f/core"
	"log"
	"os"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// The Flow:
// - get 2 paths (will need 1 path and 1 from memoroy (vriable))

// MergeCommand merge cmd struct
type MergeCommand struct {
	cobra.Command
}

// KubeConfigOption kubeConfig option
type KubeConfigOption struct {
	config   *clientcmdapi.Config
	fileName string
}

func (mc CommandOptions) runMerge(newConf *clientcmdapi.Config) error {
	outConfigs := clientcmdapi.NewConfig()
	printString(os.Stdout, "Loading KubeConfig file: "+mc.Path+" \n")
	loadConfig, err := loadKubeConfig(mc.Path)
	if err != nil {
		core.OnErrorFail(err, "File "+mc.Path+" is not kubeconfig\n")
	}
	kco := &KubeConfigOption{
		config:   loadConfig,
		fileName: getFileName(mc.Path),
	}
	kco = &KubeConfigOption{
		config:   newConf,
		fileName: "Var",
	}
	outConfigs, err = kco.handleContexts(outConfigs)
	if err != nil {
		return err
	}
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

func listFile(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var fileList []string
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		if file.IsDir() {
			listFile(folder + "/" + file.Name())
		} else {
			fileList = append(fileList, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return fileList
}

func mergeExample() string {
	return `
# Merge multiple kubeconfig
kubecm merge 1st.yaml 2nd.yaml 3rd.yaml
# Merge KubeConfig in the dir directory
kubecm merge -f dir
# Merge KubeConfig in the dir directory to the specified file.
kubecm merge -f dir --config kubecm.config
`
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
func WriteConfig(cover bool, file string, outConfig *clientcmdapi.Config) error {
	if cover {
		err := clientcmd.WriteToFile(*outConfig, cfgFile)
		if err != nil {
			return err
		}
		fmt.Printf("「%s」 write successful!\n", file)
		err = PrintTable(outConfig)
		if err != nil {
			return err
		}
	} else {
		err := clientcmd.WriteToFile(*outConfig, "kubecm.config")
		if err != nil {
			return err
		}
		printString(os.Stdout, "generate ./kubecm.config\n")
	}
	return nil
}
func getFileName(path string) string {
	n := strings.Split(path, "/")
	result := strings.Split(n[len(n)-1], ".")
	return result[0]
}

func (kc *KubeConfigOption) handleContexts(oldConfig *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	for name, ctx := range kc.config.Contexts {
		var newName string
		if len(kc.config.Contexts) > 1 {
			newName = fmt.Sprintf("%s-%s", kc.fileName, HashSufString(name))
		} else {
			newName = kc.fileName
		}
		if checkContextName(newName, oldConfig) {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it. (If you select `False`, this context will not be merged)", newName))
			if nameConfirm == "True" {
				newName = PromptUI("Rename", newName)
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
func BoolUI(label string) string {
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
		Size:      uiSize,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}
