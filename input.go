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
	deletePercentage, _ := strconv.ParseInt(os.Getenv("DELETE_PERCENTAGE"),
		0, 64)
	if deletePercentage <= 0 {
		log.Println("Delete percentage is <= 0. Nothing to delete")
		return 0
	}
	if deletePercentage >= 100 {
		log.Println("Delete percentage is >= 100. Deleting all pods")
		return numRunningPods
	}
	return int(math.Floor(float64(numRunningPods) *
		float64(deletePercentage) / 100))
}

// Get the Waiting period between each invocation of Pod deletion
func getSchedule() string {
	schedule := os.Getenv("KM_SCHEDULE")
	if schedule == "" {
		log.Println("Schedule wasn't specified. Using default of every 1m")
		return "@every 1m"
	}
	return schedule
}
