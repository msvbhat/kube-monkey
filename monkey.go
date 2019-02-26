package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

// Kill pods at random
func kubeMonkey(status chan string) {
	log.Println("kubeMonkey started...")
	kconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
		status <- "stop"
	}

	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		log.Fatal(err)
		status <- "stop"
	}

	whitelistedNS := getWhitelistedNS()
	deletablePods, err := getDeletablePods(clientset, whitelistedNS)
	if err != nil {
		log.Fatal(err)
		status <- "stop"
	}

	numDeletePods := getDeleteNum(len(deletablePods))
	log.Printf("Deleting %d pods\n", numDeletePods)

	err = deletePods(clientset, deletablePods, numDeletePods)
	if err != nil {
		log.Fatal(err)
		status <- "stop"
	}
	log.Println("Deleted the specific pods")
}
