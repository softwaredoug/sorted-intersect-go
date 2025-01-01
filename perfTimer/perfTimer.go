package perfTimer

import (
	"fmt"
	"time"
	"strings"
	"strconv"
)


type Timing struct {
	funcName string
	execTime time.Duration
}


type TimerLog struct {
	title string
	description string
	timings []Timing
}


func NewTimerLog(title string, description string) *TimerLog {
	return &TimerLog{title: title, description: description}
}


func (t *TimerLog) Timer(name string) func() {
    start := time.Now()
    return func() {
		execTime := time.Since(start)
		t.timings = append(t.timings, Timing{funcName: name, execTime: execTime})
        // fmt.Printf("%s took %v\n", name, execTime)
    }
}


// Scale a graph to maxVal
func (t *TimerLog) GraphTimings(width int) {
	// Find the max time
	maxTime := time.Duration(0)
	minTime := time.Duration(0x7FFFFFFFFFFFFFFF)
	maxFuncNameLen := 0
	for _, timing := range t.timings {
		if len(timing.funcName) > maxFuncNameLen {
			maxFuncNameLen = len(timing.funcName)
		}
		if timing.execTime > maxTime {
			maxTime = timing.execTime
		}
		if timing.execTime < minTime {
			minTime = timing.execTime
		}
	}

	fmt.Printf("\t\t\t\t%s\n", t.title)
	fmt.Printf("\t\t\t\t%s\n", strings.Repeat("-", len(t.title)))
	if t.description != "" {
		// split into lines before indenting
		descLines := strings.Split(t.description, "\n")
		for _, line := range descLines {
			fmt.Printf("\t\t\t\t%s\n", line)
		}
		fmt.Printf("\t\t\t\t%s\n", strings.Repeat("-", len(t.description)))
	}

	// Dump an ASCII graph, with maxVal = width
	for _, timing := range t.timings {
		drawWidth := int(float64(width) * float64(timing.execTime) / float64(maxTime))
		fmt.Printf("%" + strconv.Itoa(maxFuncNameLen) + "s: %s (%v", timing.funcName, strings.Repeat("=", drawWidth), timing.execTime)
		if timing.execTime == minTime {
			fmt.Printf("*)\n")
		} else {
			fmt.Printf(")\n")
		}
	}
}


