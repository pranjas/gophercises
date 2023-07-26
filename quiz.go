package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	_ "fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var timeout int
	var csvFileName string

	flag.IntVar(&timeout, "timeout", 5, "an integer value for timeout")
	flag.StringVar(&csvFileName, "csv", "problems.csv", "csv file containing 2 values")
	flag.Parse()

	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Printf("error opening file %s", err)
		os.Exit(127)
	}
	csvReader := csv.NewReader(csvFile)
	scanner := bufio.NewScanner(os.Stdin)
	var timer *time.Timer
	score := 0
	totalCount := 0

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				log.Println("Done!")
			} else {
				log.Println("Error occured in csv reader")
			}
			break
		}
		if len(record) != 2 {
			log.Fatalf("invalid csv format, expected 2 values per row, got %d", len(record))
		}
		timerNeedsDrain := false
		totalCount += 1

		if timer == nil {
			timer = time.NewTimer(time.Second * time.Duration(timeout))
		} else {
			timer.Reset(time.Second * time.Duration(timeout))
		}

		log.Printf("What is %s = ", record[0])
		if scanner.Scan() {
			input := scanner.Text()
			sum, _ := strconv.ParseInt(input, 10, 64)
			recordSum, _ := strconv.ParseInt(record[1], 10, 64)
			if sum == recordSum {
				//If the timer was NOT cancelled
				//then we need to drain it.
				if !timer.Stop() {
					<-timer.C
					log.Println("Oops better luck next time!")
				} else { //Timer was stopped won't be anything on channel.
					log.Println("You're right!")
					score += 1
				}
				timerNeedsDrain = false
			} else {
				//When the answer is incorrect we need
				//to drain whatever the case may be
				//for time.Stop since timerNeedsDrain
				//remain true.
				log.Println("No that's not correct!")
			}
			//Drain the timer channel
			if !timer.Stop() && timerNeedsDrain {
				<-timer.C
			}
		}
	}
	log.Printf("You scored %d out of %d", score, totalCount)
}
