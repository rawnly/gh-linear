package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	Spinner *spinner.Spinner
}

func NewSpinner(text string) Spinner {
	s := spinner.New(spinner.CharSets[39], 200*time.Millisecond)
	s.Suffix = " " + text

	return Spinner{
		Spinner: s,
	}
}

func (s *Spinner) Fail(message string) {
	s.Spinner.FinalMSG = fmt.Sprintf("❌ %s\n", message)
	s.Spinner.Stop()
}

func (s *Spinner) Succeed(message string) {
	s.Spinner.FinalMSG = fmt.Sprintf("%s\n", message)
	s.Spinner.Stop()
}
