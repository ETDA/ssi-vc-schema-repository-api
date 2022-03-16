package helpers

import (
	"log"
	"os"
	"time"
)

func RunAtTime(callback func(), firstRun *time.Time, repeatAt time.Duration, runNow bool) {
	if runNow {
		callback()
	}

	now := time.Now()

	if now.Sub(*firstRun) > (24 * time.Hour) {
		log.Println("firstRun time must not be earlier than 24 hours from time.Now(). Please set your firstRun date to today and set time as you want to schedule, if the time already pass the firstRun will be set to that timw in tomorrow")
		os.Exit(1)
	}

	if firstRun.Before(now) {
		tmr := firstRun.Add(24 * time.Hour)
		firstRun = &tmr
	}

	log.Printf("scheduled first run at %v \n", firstRun.String())

	// wait until time to execute first run
	time.Sleep(firstRun.Sub(now))

	for {
		callback()
		log.Printf("scheduled to next run at %v \n", time.Now().Add(repeatAt).String())
		time.Sleep(repeatAt)
	}
}
