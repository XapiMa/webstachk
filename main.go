package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	maxConnectionNum = 200
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {

	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	log.SetPrefix("webStatusChecker: ")
	log.SetFlags(0)
	targetPath := flag.String("t", filepath.Join(filepath.Dir(execPath), "target.csv"), "path to target.csv")
	outputPath := flag.String("o", "", "output file path. If not set, it will be output to standard output")
	timeLimit := flag.Int("l", 0, "Monitoring time (second). In the case of 0, it is infinite")
	maxConnectionNum := flag.Int("n", 200, "Parallel number")
	verbose := flag.Bool("v", false, "show verbose")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	if !exists(*targetPath) {
		fmt.Fprint(os.Stderr, ("Error: target.csv is not exist.\n"))
		flag.Usage()
		os.Exit(2)
	}
	if err := StatusCheck(*targetPath, *outputPath, int64(*timeLimit), *maxConnectionNum); err != nil {
		log.Fatal(err)
	}
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// StatusCheck start checking web status
func StatusCheck(targetPath, outputPath string, timeLimit int64, maxConnectionNum int) error {
	if outputPath != "" {
		if !exists(outputPath) {
			dir, file := filepath.Split(outputPath)
			if file != "" {
				if !exists(dir) {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return err
					}
				}
			}
		}
	}

	targets, err := parseTargetFile(targetPath)
	if err != nil {
		return err
	}

	if err := check(targets, outputPath, timeLimit); err != nil {
		return err
	}
	return nil

}

type target struct {
	url      string
	statuses []int
	interval int
}

func parseTargetFile(targetPath string) ([]target, error) {
	var targets = []target{}
	file, err := os.Open(targetPath)
	failOnError(err)
	defer file.Close()
	reader := csv.NewReader(file) // utf8
	// reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	// reader := csv.NewReader(transform.NewReader(file, japanese.EUCJP.NewDecoder()))
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}
		url := record[0]
		strStatuses := strings.Split(record[1], "|")
		intStatuses := make([]int, len(strStatuses))
		for i, str := range strStatuses {
			var err error
			intStatuses[i], err = strconv.Atoi(str)
			if err != nil {
				return targets, err
			}
		}
		var coefficient int
		switch record[2][len(record[2])-1] {
		case 'd':
			coefficient = 60 * 60 * 24
		case 'h':
			coefficient = 60 * 60
		case 'm':
			coefficient = 60
		case 's':
			coefficient = 1
		default:
			return targets, fmt.Errorf("the format of the access time interval is incorrect")
		}
		timeNum, err := strconv.Atoi(record[2][:len(record[2])-1])
		if err != nil {
			return targets, err
		}
		timeNum *= coefficient

		targets = append(targets, target{url, intStatuses, timeNum})
	}
	return targets, nil
}

type timeRecord struct {
	interval int64
	start    int64
}

func check(targets []target, outputPath string, limit int64) error {

	maxConnection := make(chan bool, maxConnectionNum)
	timeRecords := make([]timeRecord, len(targets))
	for i, item := range targets {
		timeRecords[i].interval = int64(item.interval)
		timeRecords[i].start = int64(0)
	}
	allStart := time.Now().Unix()
	for {
		var nowTime int64
		nowTime = time.Now().Unix()
		if limit > 0 {
			if nowTime-allStart >= limit {
				break
			}
		}
		for i := range timeRecords {
			if nowTime-timeRecords[i].start >= timeRecords[i].interval {
				maxConnection <- true
				timeRecords[i].start = nowTime

				go func(i int) {
					url := targets[i].url
					resp, err := http.Get(url)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error:%s %s %s\n", time.Unix(nowTime, 0), url, err)
						return
					}
					defer resp.Body.Close()
					flag := false
					for _, code := range targets[i].statuses {
						if code == resp.StatusCode {
							flag = true
							break
						}
					}
					if !flag {
						codeString := ""
						for j, code := range targets[i].statuses {
							codeString += strconv.Itoa(code)
							if j != len(targets[i].statuses)-1 {
								codeString += "_or_"
							}
						}
						outputString := fmt.Sprintf("Warning: %s %s return status %d but expected %s\n", time.Unix(nowTime, 0), targets[i].url, resp.StatusCode, codeString)
						appendFile(outputPath, outputString)
					}
					<-maxConnection
				}(i)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
func appendFile(outputPath, outputString string) error {
	if outputPath == "" {
		fmt.Printf("%s", outputString)
	} else {
		file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Write(([]byte)(outputString))
		file.Close()
	}
	return nil
}
