package webstachk

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseConfigFile(configPath string) ([]target, error) {
	errorWrap := func(err error) error {
		return errors.Wrap(err, "cause in parseConfigFile")
	}
	var targets = []target{}
	file, err := os.Open(configPath)
	if err != nil {
		return targets, errorWrap(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	for index := 0; true; index++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return targets, errorWrap(err)
		}
		url := record[0]
		intStatuses, err := parseStatuses(record[1])
		if err != nil {
			return targets, errors.Wrap(err, fmt.Sprintf("cause in parseConfigFile %dth record", index+1))
		}
		interval, err := parseTime(record[2])
		if err != nil {
			return targets, errors.Wrap(err, fmt.Sprintf("cause in parseConfigFile %dth record", index+1))
		}
		targets = append(targets, target{url, intStatuses, interval})
	}
	return targets, nil
}

func parseStatuses(statuses string) ([]int, error) {
	strStatuses := strings.Split(statuses, "|")
	intStatuses := make([]int, 0)
	for i, str := range strStatuses {
		var err error
		if str == "" {
			continue
		}
		statusCode, err := strconv.Atoi(str)
		if err != nil {
			return intStatuses, errors.Wrap(err, fmt.Sprintf("%dth item can't convert int", i+1))
		}
		intStatuses = append(intStatuses, statusCode)
	}
	return intStatuses, nil
}

func parseTime(timeString string) (int, error) {
	timeNum, err := strconv.Atoi(timeString[:len(timeString)-1])
	if err != nil {
		return timeNum, err
	}
	switch timeString[len(timeString)-1] {
	case 'd':
		timeNum *= 60 * 60 * 24
	case 'h':
		timeNum *= 60 * 60
	case 'm':
		timeNum *= 60
	case 's':
		timeNum *= 1
	default:
		return 0, fmt.Errorf("the format of the access time interval is incorrect")
	}
	return timeNum, nil
}
