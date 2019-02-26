package main

import (
	"github.com/robfig/cron"
	"log"
)

func main() {
	status := make(chan string, 1)
	schedule := getSchedule()
	c := cron.New()
	err := c.AddFunc(schedule, func() { kubeMonkey(status) })
	if err != nil {
		log.Fatal("Error Adding the fucntion to Cron")
	}
	c.Start()
	go healthCheck(status)
	msg := <-status
	if msg == "stop" {
		log.Println("Received the stop signal. Exiting...")
	}
}
