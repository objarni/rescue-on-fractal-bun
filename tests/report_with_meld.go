package tests

import (
	"github.com/approvals/go-approval-tests/reporters"
	"os/exec"
)

type reportWithMeld struct{}

func ReportWithMeld() reporters.Reporter {
	return &reportWithMeld{}
}

func (s *reportWithMeld) Report(approved, received string) bool {
	cmd := exec.Command("/usr/bin/meld", approved, received)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return true
}
