package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/xapima/webstatuschecker"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func logFatal(err error) {
	log.Fatalf("Error: webStatusChecker %s %s", time.Now().Format("2006/01/02 15:04:05"), err)
}

func main() {
	errorWrap := func(err error) error {
		return errors.Wrap(err, "cause in main")
	}

	log.SetPrefix("webStatusChecker: ")
	log.SetFlags(0)
	configPath := flag.String("t", "", "path to config.csv")
	outputPath := flag.String("o", "", "output file path. If not set, it will be output to standard output")
	intervalTime := flag.Int("i", 60, "interval to self health check(second). In case of 0 it does not check")
	maxConnectionNum := flag.Int("n", 200, "Parallel number")
	// verbose := flag.Bool("v", false, "show verbose")
	flag.Parse()

	// if !*verbose {
	// 	log.SetOutput(ioutil.Discard)
	// }

	if *configPath == "" {
		logFatal(fmt.Errorf("-t option is required"))
	}
	if !exists(*configPath) {
		logFatal(fmt.Errorf("%s is not exist", *configPath))
	}
	if err := webstatuschecker.StatusCheck(*configPath, *outputPath, *maxConnectionNum, *intervalTime); err != nil {
		logFatal(errorWrap(err))
	}
}
