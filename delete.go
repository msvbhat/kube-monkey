package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"time"
)

// Delete the pod and return error value
func deletePod(clientset *kubernetes.Clientset, pod v1.Pod) error {
	return clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
}

// Delete the random pods and return error value
func deletePods(clientset *kubernetes.Clientset, deletablePods []v1.Pod, numDeletePods int) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numDeletePods; i++ {
		podToBeDeleted := deletablePods[r.Intn(len(deletablePods))]
		err := deletePod(clientset, podToBeDeleted)
		if err != nil {
			return err
		}
	}
	return nil
}
