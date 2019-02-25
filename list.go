package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

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
		if pod.Namespace == ns || pod.Name == os.Getenv("MY_POD_NAME") {
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
