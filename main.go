package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
	// "github.com/pkg/profile"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

var logFile = "./logs"

const maxLogFileSize = 200000000

type logWriter struct {
	logFile *os.File
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	stat, err := writer.logFile.Stat()
	if err != nil {
		return 0, err
	}
	if stat.Size() > maxLogFileSize {
		myLog("shorter log file, its too large")
		writer.logFile.Seek(-maxLogFileSize/2, io.SeekEnd)
		var buf [maxLogFileSize / 2]byte
		writer.logFile.Read(buf[:])
		writer.logFile.Seek(0, io.SeekStart)
		writer.logFile.Write(buf[:])
		writer.logFile.Truncate(maxLogFileSize / 2)
	}
	writer.logFile.Seek(0, io.SeekEnd)

	return fmt.Fprint(writer.logFile, "\033[33m"+time.Now().Format("2006/01/02 15:04:05")+"\033[0m "+string(bytes))
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			myLog("PANIC IN MAIN!!!!\n%s\n%#v\n", string(debug.Stack()), r)
			log.Printf("panic in main %s\n%#v\n", string(debug.Stack()), r)
		}
	}()
	// defer profile.Start(profile.MemProfile).Stop()
	// defer profile.Start().Stop()

	err := errors.New("")
	log.SetFlags(0)
	logger := new(logWriter)
	logger.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("cant open log file", err)
		return
	}
	log.SetOutput(logger)
	os.Stdout = logger.logFile
	log.Printf("\nrun\n")

	var wg sync.WaitGroup
	var stop []chan bool
	telLogChan := make(chan logMessage)
	wg.Add(1)
	stop = append(stop, make(chan bool))
	go cpuProfileRun(stop[len(stop)-1], &wg)

	var settings Settings
	err = readFile(&settings, "settings.json")
	if err != nil {
		log.Println("cant read settings", err)
		return
	}

	wg.Add(1)
	stop = append(stop, make(chan bool))
	go startTelegramLog(telLogChan, stop[len(stop)-1], &wg, &settings)

	time.Sleep(time.Second * 2)

	wg.Add(1)
	stop = append(stop, make(chan bool))
	go runBot(stop[len(stop)-1], &wg)

	time.Sleep(time.Second * 2)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)

		done <- true

	}()
	<-done
	time.Sleep(time.Millisecond * 200)
	log.Printf("closes\n")
	for i := len(stop) - 1; i >= 0; i-- {
		stop[i] <- true
	}
	wg.Wait()

	err = syncFile(&settings)
	if err != nil {
		log.Println("cant sync settings", err)
	}
	log.Printf("closed")
}

func tickerUpdateFiles(saves ...JSONSaving) {

	for _, save := range saves {
		err := syncFile(save)
		if err != nil {
			myLog("cant write file %s %s", save.GetFile(), err)
		}
	}
}

func cpuProfileRun(stop chan bool, wg *sync.WaitGroup) {
	defer func() {
		log.Printf("closed cpuProfileRun")
		wg.Done()
	}()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		if pprof.StartCPUProfile(f) != nil {
			return
		}
		defer pprof.StopCPUProfile()
	}
	<-stop
}
