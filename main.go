package main

import (
	"fmt"
	"sync"
	"time"

	// To install this package run
	// "go get github.com/go-co-op/gocron"
	"github.com/go-co-op/gocron"
)


func task(wg *sync.WaitGroup, c chan int, count *int, s *gocron.Scheduler) {

	*count += 1

	// Took a random-number(5) to limit number of time scheduler1 should work
	if *count <=5{
		// Sending data into channel the limit is reached
		c <- *count
	}
	// Once the limit is reached:
	if *count == 5 {
		// Closing channel
		close(c)
		// Sending done response to wait group
		wg.Done()
		// Removing scheduler1 from main scheduler1, to stop the scheduling
		s.RemoveByTag("scheduler1")
	}
}

func task1(wg *sync.WaitGroup, c chan int, count1 *int, s *gocron.Scheduler) {
	
	*count1 += 1

	fmt.Printf("|------ Simple scheduler count %v -------------|\n", *count1)
	
	// Took a random-number(10) to limit number of time scheduler1 should work
	if *count1 == 10 {
		// Once the limit is reached:
		// Sending done response to wait group
		wg.Done()
		// Removing scheduler1 from main scheduler2, to stop the scheduling
		s.RemoveByTag("scheduler2")
	}
}

func main() {
	var wg sync.WaitGroup
	var count = 0
	var count1 = 0
	var channel = make(chan int)

	// Initialising new Scheduler using gocron.NewScheduler(time.UTC) 
	var scheduler = gocron.NewScheduler(time.UTC)

	// Scheduling two tasks:
	// 1. Runs normally
	// 2. Passes data into channel every second
	scheduler.Every(1).Seconds().Tag("scheduler1").Do(task, &wg, channel, &count, scheduler)
	scheduler.Every(1).Seconds().Tag("scheduler2").Do(task1, &wg, channel, &count1, scheduler)

	// As I'm using 2 schedulers adding 2 to Wait Group
	wg.Add(2)

	// Starting schedulers as go routine
	go scheduler.StartBlocking()
	
	// this loop iterates untill channel is supplied with data
	for {
		msg, open := <- channel

		if !open {
			break
		}
		fmt.Printf("|------------- Data From channel %v -----------|\n", msg)
	}

	fmt.Println("**************************************************************")
	fmt.Println("Scheduler2 Task Finished")
	fmt.Println("Channel Closed")
	fmt.Println("Scheduler1 Still Running")
	fmt.Println("**************************************************************")
	
	wg.Wait()

	fmt.Println("**************************************************************")
	fmt.Println("Scheduler1 Task Finished")
	fmt.Println("**************************************************************")
	fmt.Println("\n🥳 Task Done 🥳")
}