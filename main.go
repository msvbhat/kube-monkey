package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/rest" -- To be Added when moving to in cluster setup
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Get a list of whitelisted Namespaces. The pods in the whitelisted
// will not be killed
func getWhitelistedNS() (whitelistedNS []string) {
	whitelistedNS = strings.Fields(os.Getenv("NAMESPACE_WHITELIST"))

	// Append kube-system to list of whitelisted namespaces if not present
	var flag = false
	for _, ns := range whitelistedNS {
		if "kube-system" == ns {
			flag = true
		}
	}
	if !flag {
		whitelistedNS = append(whitelistedNS, "kube-system")
	}
	return whitelistedNS
}

// Get a list of all running pods in the cluster
func getRunningPods(clientset *kubernetes.Clientset) (podsRunning []v1.Pod, err error) {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			podsRunning = append(podsRunning, pod)
		}
	}
	return podsRunning, nil

}

// Check if the pod belongs to any whitelisted Namespaces
func isWhitelisted(whitelistedNS []string, pod v1.Pod) bool {
	for _, ns := range whitelistedNS {
		if pod.Namespace == ns {
			return true
		}
	}
	return false
}

// Get a list of all pods which are eligible for delete
func getDeletablePods(clientset *kubernetes.Clientset, whitelistedNS []string) (deletablePods []v1.Pod, err error) {
	runningPods, err := getRunningPods(clientset)
	if err != nil {
		return nil, err
	}
	for _, pod := range runningPods {
		if !isWhitelisted(whitelistedNS, pod) {
			deletablePods = append(deletablePods, pod)
		}
	}
	return deletablePods, nil
}

func getDeleteNum(numRunningPods int) int {
	deletePercentage, _ := strconv.ParseInt(os.Getenv("DELETE_NUM_PERCENTAGE"), 0, 64)
	if deletePercentage <= 0 {
		log.Println("Delete percentage is set to 0 or less. Nothing to delete")
		return 0
	}
	if deletePercentage >= 100 {
		log.Println("Delete percentage is set to 100 or more. Deleting all pods")
		return numRunningPods
	}
	return int(math.Floor(float64(numRunningPods) * float64(deletePercentage) / 100))
}

func main() {
	// TODO: Use in cluster config after initial testing phase
	//kconfig, err := rest.InClusterConfig()
	kconfig, err := clientcmd.BuildConfigFromFlags("", "/Users/msvbhat/.kube/config")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		log.Fatal(err)
	}

	whitelistedNS := getWhitelistedNS()
	deletablePods, err := getDeletablePods(clientset, whitelistedNS)
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range deletablePods {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}

	fmt.Println(getDeleteNum(len(deletablePods)))
}

// Kill random pods from the killable pods list
// Wait for an 10 minutes and repeat the process
// Refactor
// Get a health endpoint
// Get metric endpoint
// Send events to pods
