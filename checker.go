package webstachk

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	maxConnectionNum = 200
)

// StatusCheck start checking web status
func StatusCheck(configPath, outputPath string, maxConnectionNum int, interval int) error {
	errorWrap := func(err error) error {
		return errors.Wrap(err, "cause in StatusCheck")
	}
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

	targets, err := parseConfigFile(configPath)
	if err != nil {
		return errorWrap(err)
	}
	go func() {
		if interval != 0 {
			for {
				str := fmt.Sprintf("Alive: %s\n", time.Now().Format("2006/01/02 15:04:05"))
				appendFile(outputPath, str)
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}

	}()
	if err := check(targets, outputPath); err != nil {
		return errorWrap(err)
	}
	return nil

}

type target struct {
	url      string
	statuses []int
	interval int
}

type timeRecord struct {
	interval int64
	start    int64
}

func check(targets []target, outputPath string) error {
	errorWrap := func(err error) error {
		return errors.Wrap(err, "cause in check")
	}

	maxConnection := make(chan bool, maxConnectionNum)
	timeRecords := make([]timeRecord, len(targets))
	for i, item := range targets {
		timeRecords[i].interval = int64(item.interval)
		timeRecords[i].start = int64(0)
	}
	for {
		var nowTime int64
		nowTime = time.Now().Unix()
		for i := range timeRecords {
			if nowTime-timeRecords[i].start >= timeRecords[i].interval {
				maxConnection <- true
				timeRecords[i].start = nowTime

				go func(i int) {
					url := targets[i].url
					resp, err := http.Get(url)
					if err != nil {
						logPrint(errorWrap(err))
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
						outputString := fmt.Sprintf("Warning: %s %s status %d found, but expected %s\n", time.Unix(nowTime, 0).Format("2006/01/02 15:04:05"), targets[i].url, resp.StatusCode, codeString)
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
