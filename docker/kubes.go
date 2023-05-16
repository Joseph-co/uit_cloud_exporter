package docker

import (
	"flag"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"strings"
)


func GetDeployStatus(clientset *kubernetes.Clientset)bool{
	list, err := clientset.AppsV1().Deployments("monitoring").List(metav1.ListOptions{})
	if err != nil{
		log.Fatal(err)
	}
	for _, l := range list.Items {
		if l.Name == "kube-state-metrics" && l.Status.AvailableReplicas >= 1 {
			return true
		}
	}
	return false
}

func GetK8sConf()(strc string, kubeconfig *string){
	//var kubeconfig *string
	if home:= homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig",filepath.Join(home, ".kube","config"),"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig","","absolute path to the kubeconfig file")
	}
	flag.Parse()

	filePath := *kubeconfig // Replace with the actual file path

	// Check if the file exists
	if fileExists(filePath) {
		// Read the file contents
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}

		// Convert file content to string
		fileContent := string(content)

		// Find strings following "server:"
		stringsBehindServer := findStringsBehindServer(fileContent)

		// Print the extracted strings
		for _, str := range stringsBehindServer {
			strc = str
		}
	} else {
		log.Fatal("File does not exist.")
	}
	return strc,kubeconfig
}


// Check if a file exists
func fileExists(filePath string) bool {
	_, err := ioutil.ReadFile(filePath)
	return err == nil
}

// Find strings behind "server:" keyword
func findStringsBehindServer(content string) []string {
	var result []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "server:") {
			parts := strings.Split(line, "server:")
			if len(parts) > 1 {
				str := strings.TrimSpace(parts[1])
				result = append(result, str)
			}
		}
	}
	return result
}