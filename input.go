package main

import (
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

// Get the number of pods to be deleted from the environment variable
func getDeleteNum(numRunningPods int) int {
	deletePercentage, _ := strconv.ParseInt(os.Getenv("DELETE_PERCENTAGE"), 0, 64)
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

// Get the Waiting period between each invocation of Pod deletion
func getWaitingPeriod() int {
	waitMinutes, _ := strconv.ParseInt(os.Getenv("WAIT_MINUTES"), 0, 64)
	if waitMinutes <= 0 {
		log.Println("Can not wait for zero 0r less. Setting the value as 1")
		return 1
	}
	return int(waitMinutes)
}
