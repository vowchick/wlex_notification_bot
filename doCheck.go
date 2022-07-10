package main

import (
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

func runBot(stop chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			myLog("PANIC IN BOT!!!!\n%s\n%#v\n", string(debug.Stack()), r)
			log.Printf("closed runBot")
			time.Sleep(2 * time.Second)
			os.Exit(1)
		}
	}()

	period := 10 * time.Second
	log.Printf("period : %d minute", period)

	ticker := time.NewTicker(time.Until(time.Now().Truncate(period).Add(period)))
	first := true
	for {
		select {
		case curTime := <-ticker.C:
			if first == true {
				first = false
				ticker.Stop()
				period := 10 * time.Second // тут изменить period
				ticker = time.NewTicker(period)
				log.Printf("period : %d minute", period)
			}

			log.Printf("started %s", curTime.Format("2 15:04:05"))

			doCheck()

			log.Printf("ended %s", curTime.Format("2 15:04:05"))
		case <-stop:
			log.Printf("closed runBot")
			myLog("stopped bot")
			return
		}

	}

}

func doCheck() {
	myLog("Hello!")
}
