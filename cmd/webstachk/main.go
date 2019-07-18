package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/xapima/webstachk"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func logFatal(err error) {
	log.Fatalf("Error: webstachk %s %s", time.Now().Format("2006/01/02 15:04:05"), err)
}

func main() {
	log.SetPrefix("webstachk: ")
	log.SetFlags(0)
	configPath := flag.String("t", "", "path to config.csv")
	outputPath := flag.String("o", "", "output file path. If not set, it will be output to standard output")
	interval := flag.Int("i", 60, "interval to self health check(second). In case of 0 it does not check")
	maxCon := flag.Int("n", 200, "Parallel number")
	isJSON := flag.Bool("j", false, "change output format to json")
	flag.Parse()

	if *configPath == "" {
		logFatal(fmt.Errorf("-t option is required"))
	}
	if !exists(*configPath) {
		logFatal(fmt.Errorf("%s is not exist", *configPath))
	}
	checker, err := webstachk.NewChecker(*configPath, *outputPath, *interval, *maxCon, *isJSON)
	if err != nil {
		logFatal(errors.Wrap(err, "cant create new Checker"))
	}
	checker.Check()
}
