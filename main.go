package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/rest" -- To be Added when moving to in cluster setup
	"log"
)

func getRunningPods(clientset *kubernetes.Clientset) (podsRunning []v1.Pod, err error) {
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
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

	deletablePods, err := getRunningPods(clientset)
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range deletablePods {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}
}

// Initialise the clientset
// Get the whitelisted namespaces
// Get the running pods
// Get the killable pods
// Get the percentage of pods to kill
// Kill random pods from the killable pods list
// Wait for an 10 minutes and repeat the process
// Get a health endpoint
// Get metric endpoint
// Send events to pods
