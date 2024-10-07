package code

import (
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	parse := NewParser()
	WalkDir("/Users/gaowei7/code/go/src/gitlabee.com/licloud-workflow-service", func(path string) {
		fmt.Print(path)
		file, err := parse.ParseByFile(path)
		if err != nil {
			return
		}
		file.PrintResults()

	})
}
