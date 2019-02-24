package main

import "log"

func main() {
	status := make(chan string, 1)
	go kubeMonkey(status)
	go healthCheck(status)
	msg := <-status
	if msg == "stop" {
		log.Println("Received the stop signal.")
	}
}

// Remove the self pod from deletablePods
// Add Cron scheduler
// Send Events
// Add Metric Endpoint
