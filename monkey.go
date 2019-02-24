package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

// Kill pods at random
func kubeMonkey(status chan string) {
	defer func() { status <- "stop" }()
	waitingPeriod := getWaitingPeriod()
	for {
		kconfig, err := rest.InClusterConfig()
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

		numDeletePods := getDeleteNum(len(deletablePods))
		log.Printf("Deleting %d pods\n", numDeletePods)

		err = deletePods(clientset, deletablePods, numDeletePods)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(waitingPeriod) * time.Minute)
	}
}
