package webstachk

import (
	"fmt"
	"testing"
)

// type Checker struct {
// 	Targets    []Target
// 	OutputPath string
// 	IsJSON     bool
// 	MaxCon     chan bool
// 	Interval   int64
// }

func TestNewChecker(t *testing.T) {
	type data struct {
		configPath string
		outputPath string
		interval   int
		maxc       int
		isJSON     bool
	}
	tests := []data{
		data{
			"./test/testTarget.csv", "", 0, 0, false,
		}, data{
			"./test/testTarget.csv", "testOut", 1, 0, false,
		},
		data{
			"./test/testTarget.csv", "out", 100, 10, false,
		},
	}
	for i, test := range tests {
		c, err := NewChecker(test.configPath, test.outputPath, test.interval, test.maxc, test.isJSON)
		if err != nil {
			t.Errorf("NewChecker: %dth item %s", i, err)
		}

		chk := new(Checker)
		targets, err := parseConfigFile(test.configPath)
		if err != nil {
			t.Errorf("NewChecker: %dth item %s", i, err)
		}
		chk.Targets = targets
		chk.OutputPath = test.outputPath
		chk.IsJSON = test.isJSON
		chk.MaxCon = make(chan bool, test.maxc)
		chk.Interval = int64(test.interval)
		if fmt.Sprintf("%v %v %v %v", chk.Targets, chk.OutputPath, chk.IsJSON, chk.Interval) != fmt.Sprintf("%v %v %v %v", c.Targets, c.OutputPath, c.IsJSON, c.Interval) {
			t.Errorf("NewChecker: %dth item %v found but expect %v", i, c, chk)
		}
	}
}

func TestIsMatch(t *testing.T) {

	type data struct {
		target Target
		input  int
		expect bool
	}
	tests := []data{
		{Target{Statuses: []int{200}}, 200, true},
		{Target{Statuses: []int{200}}, 301, false},
		{Target{Statuses: []int{200, 301, 302}}, 301, true},
		{Target{Statuses: []int{200, 301, 302}}, 404, false},
	}

	for i, test := range tests {
		if ok := test.target.isMatch(test.input); ok != test.expect {
			t.Errorf("[isMatch] %dth item found %v but expect %v", i, ok, test.expect)
		}
	}
}
