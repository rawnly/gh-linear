package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	spinner *spinner.Spinner
}

var spinnerIdx = 39
var spinnerSpeed time.Duration = 200

func NewSpinner(text string) Spinner {
	s := spinner.New(spinner.CharSets[spinnerIdx], spinnerSpeed*time.Millisecond)
	s.Suffix = " " + text

	return Spinner{
		spinner: s,
	}
}

func (s *Spinner) Start() {
	s.spinner.Start()
}

func (s *Spinner) Fail(message string) {
	s.spinner.FinalMSG = fmt.Sprintf("❌ %s\n", message)
	s.spinner.Stop()
}

func (s *Spinner) Succeed(message string) {
	s.spinner.FinalMSG = fmt.Sprintf("%s\n", message)
	s.spinner.Stop()
}
