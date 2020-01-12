package hw02

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type TestCase struct {
	In      string
	Out     string
	IsError bool
}

var TestCases []TestCase = []TestCase{
	TestCase{
		In:  "a4bc2d5e",
		Out: "aaaabccddddde",
	},
	TestCase{
		In:  "abcd",
		Out: "abcd",
	},
	TestCase{
		In:      "45",
		Out:     "",
		IsError: true,
	},
	TestCase{
		In:  `qwe\4\5`,
		Out: "qwe45",
	},
	TestCase{
		In:  `qwe\45`,
		Out: "qwe44444",
	},
	TestCase{
		In:  `qwe\\5`,
		Out: `qwe\\\\\`,
	},
}

func TestSimpleUnpackString(t *testing.T) {

	for _, c := range TestCases {
		res, err := SimpleUnpackString(c.In)
		require.Equal(t, c.IsError, err != nil, "error not match")
		require.Equal(t, c.Out, res, "result not match")
	}
}
