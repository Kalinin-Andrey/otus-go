package main

import (
	"os"
	"testing"
)

type TestCase struct {
	from string
	to string
	limit uint
	offset uint
}

var testCases []TestCase = []TestCase{
	TestCase{
		from    : "go.mod",
		to      : "go.mod.copy",
		limit   : 0,
		offset  : 0,
	},
	TestCase{
		from    : "go.mod",
		to      : "go.mod.copy",
		limit   : 0,
		offset  : 47,
	},
	TestCase{
		from    : "go.mod",
		to      : "go.mod.copy",
		limit   : 207,
		offset  : 47,
	},
	TestCase{
		from    : "go.mod",
		to      : "go.mod.copy",
		limit   : 18,
		offset  : 47,
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

		var sizeToCopy uint = uint(fileInfoFrom.Size()) - testCase.offset

		if testCase.limit > 0 {
			if sizeToCopy > testCase.limit {
				sizeToCopy = testCase.limit
			}
		}

		if int64(sizeToCopy) != fileInfoTo.Size() {
			t.Errorf("Size of the resulting file is %v, expected %v", fileInfoTo.Size(), sizeToCopy)
		}

		err = os.Remove(testCase.to)
		if err != nil {
			t.Errorf("os.Remove error: %s", err)
		}
	}
}



