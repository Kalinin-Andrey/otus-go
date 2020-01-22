package main

import (
	"os"
	"testing"
)

type TestCase struct {
	from string
	to string
	limit int
	offset int
}

var testCases []TestCase = []TestCase{
	TestCase{
		from    : "go.mod",
		to      : "go.mod.copy",
		limit   : -1,
		offset  : 0,
	},
}


func TestCopy(t *testing.T) {

	for _, testCase := range testCases {

		err := Copy(testCase.from, testCase.to, testCase.limit, testCase.offset)
		if err != nil {
			t.Errorf("Copy error: %s", err)
		}

		fileInfoFrom, err := os.Stat(testCase.from)
		if err != nil {
			t.Errorf("os.Stat error: %s", err)
		}

		fileInfoTo, err := os.Stat(testCase.to)
		if err != nil {
			t.Errorf("os.Stat error: %s", err)
		}

		if fileInfoFrom.Size() != fileInfoTo.Size() {
			t.Errorf("Size of the resulting file is %v, expected %v", fileInfoTo.Size(), fileInfoFrom.Size())
		}
	}
}



