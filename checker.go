package webstachk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	maxConnectionNum = 200
)

// Code is type of Message codes
type Code string

const (
	// WarningCode indicates a warning
	WarningCode = Code("Warning")

	// ErrorCode indicates a error
	ErrorCode = Code("Error")

	// AliveCode indicates a alive
	AliveCode = Code("Alive")
)

// Target is a item of checking targets
type Target struct {
	Url      string
	Statuses []int
	Interval int64
	Start    int64
}

// Checker is a webstachk object
type Checker struct {
	Targets    []Target
	OutputPath string
	IsJSON     bool
	MaxCon     chan bool
	Interval   int64
}

// Warning stores information of mismatched status
type Warning struct {
	WCode      Code      `json:"code"`
	Url        string    `json:"url,omitempty"`
	Expected   []int     `json:"expected,omitempty"`
	Found      int       `json:"found,omitempty"`
	TimeRecord time.Time `json:"time"`
}

// NewChecker create a new webstachk object
func NewChecker(configPath, outputPath string, interval int, maxCon int, isJSON bool) (*Checker, error) {
	errorWrap := func(err error) error {
		return errors.Wrap(err, "cause in NewChecker")
	}
	chk := new(Checker)
	targets, err := parseConfigFile(configPath)
	if err != nil {
		return chk, errorWrap(err)
	}
	chk.Targets = targets
	chk.OutputPath = outputPath
	chk.IsJSON = isJSON
	chk.MaxCon = make(chan bool, maxCon)
	chk.Interval = int64(interval)
	return chk, nil
}

func (c *Checker) Check() error {
	ch := make(chan Warning)
	go func() {
		for warning := range ch {
			c.write(warning)
		}
	}()
	go c.CallAlive()
	for {
		now := time.Now().Unix()
		for i, t := range c.Targets {
			if now-t.Start >= t.Interval {
				c.MaxCon <- true
				c.Targets[i].Start = now
				// If a variable is not substituted, t is overwritten and the object that calls the function is overwritten
				call := t
				go call.access(ch, c.MaxCon)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (t *Target) access(ch chan Warning, maxCon chan bool) {
	<-maxCon
	// fmt.Printf("%s %v is started at %s\n", t.Url, t.Statuses, time.Unix(t.Start, 0).Format("2006/01/02 15:04:05"))
	resp, err := http.Get(t.Url)
	if err != nil {
		logPrint(errors.Wrapf(err, "couse for %s", t.Url))
		return
	}
	defer resp.Body.Close()
	if ok := t.isMatch(resp.StatusCode); !ok {

		w := new(Warning)
		w.Url = t.Url
		w.Expected = t.Statuses
		w.Found = resp.StatusCode
		w.TimeRecord = time.Now()
		w.WCode = WarningCode
		ch <- *w
	}
	return
}

func (t *Target) isMatch(status int) bool {
	for _, code := range t.Statuses {
		if code == status {
			return true
		}
	}
	return false
}

// CallAlive send Alive message
func (c *Checker) CallAlive() {
	for {
		w := Warning{WCode: AliveCode, TimeRecord: time.Now()}
		c.write(w)
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}
func (c *Checker) write(w Warning) {
	var str string
	if c.IsJSON {
		buf, err := json.Marshal(&w)
		if err != nil {
			logPrint(err)
		}
		str = string(buf)
	} else {
		codeString := joinCode(w.Expected, "_or_")
		switch w.WCode {
		case WarningCode, ErrorCode:
			str = fmt.Sprintf("%s %s: %s status %d found, but expected %s\n", w.TimeRecord.Format("2006/01/02 15:04:05"), w.WCode, w.Url, w.Found, codeString)
		case AliveCode:
			str = fmt.Sprintf("%s %s\n", w.TimeRecord.Format("2006/01/02 15:04:05"), w.WCode)
		}
	}
	appendFile(c.OutputPath, str)
}
