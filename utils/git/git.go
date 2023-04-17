package git

import (
	"strings"

	"github.com/rawnly/gitgud/run"
)

func GetBranches() ([]string, error) {
	output, err := run.Git("branch").Output()
	branches := strings.Split(string(output), "\n")

	for i, branch := range branches {
		b := strings.Replace(branch, "*", "", -1)
		b = strings.TrimSpace(b)

		branches[i] = b
	}

	return branches, err
}
