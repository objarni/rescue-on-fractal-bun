package tests

import (
	"github.com/approvals/go-approval-tests/reporters"
	"os/exec"
)

type beyondCompare struct{}

func NewBCompare() reporters.Reporter {
	return &beyondCompare{}
}

func (s *beyondCompare) Report(approved, received string) bool {
	cmd := exec.Command("/usr/bin/bcompare", approved, received)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return true
}
