package hw05

import (
	"errors"
	"log"
	"testing"
	"time"
)

var NumberOfFuncs int = 30
var TimeOfFunc int = 2 // sec.
var NumberOfGoroutines int = 10
var QuoteOfErrors int = 3

func funcWithError() error {
	time.Sleep(time.Duration(TimeOfFunc) * time.Second)
	return errors.New("time out")
}

func funcWithoutError() error {
	time.Sleep(10 * time.Second)
	return nil
}

func funcWithPanic() error {
	log.Panic("Panic!")
	panic("A am panic!")
	//panic(errors.New("A am panic!"))
}


func TestRunPositive(t *testing.T) {
	funcsSlice := make([]func() error, 0, NumberOfFuncs)
	countOfFuncsChannel := make(chan struct{}, NumberOfFuncs)

	function := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithoutError()
	}

	for i := 0; i < NumberOfFuncs; i++ {
		funcsSlice = append(funcsSlice, function)
	}
	err := Run(funcsSlice, NumberOfGoroutines, QuoteOfErrors)

	if err != nil {
		t.Error("Result must be positive")
	}

	if len(countOfFuncsChannel) != NumberOfFuncs {
		t.Errorf("Count of executed funcs is %v, expected %v", len(countOfFuncsChannel), NumberOfFuncs)
	}
}


func TestRunNegative(t *testing.T) {
	funcsSlice := make([]func() error, 0, NumberOfFuncs)
	countOfFuncsChannel := make(chan struct{}, NumberOfFuncs)

	function := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithoutError()
	}

	functionFail := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithError()
	}

	for i := 0; i < QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, functionFail)
	}

	for i := 0; i < NumberOfFuncs-QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, function)
	}
	err := Run(funcsSlice, NumberOfGoroutines, QuoteOfErrors)

	if err == nil {
		t.Error("Result must be negative")
	}

	if len(countOfFuncsChannel) > QuoteOfErrors+NumberOfGoroutines {
		t.Errorf("Count of executed funcs is %v, expected <= %v", len(countOfFuncsChannel), QuoteOfErrors+NumberOfGoroutines)
	}
}


func TestRunWithPanics(t *testing.T) {
	funcsSlice := make([]func() error, 0, NumberOfFuncs)
	countOfFuncsChannel := make(chan struct{}, NumberOfFuncs)

	function := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithoutError()
	}

	functionFail := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithPanic()
	}

	for i := 0; i < QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, functionFail)
	}

	for i := 0; i < NumberOfFuncs-QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, function)
	}
	err := Run(funcsSlice, NumberOfGoroutines, QuoteOfErrors)

	if err == nil {
		t.Error("Result must be negative")
	}

	if len(countOfFuncsChannel) > QuoteOfErrors+NumberOfGoroutines {
		t.Errorf("Count of executed funcs is %v, expected <= %v", len(countOfFuncsChannel), QuoteOfErrors+NumberOfGoroutines)
	}
}



func TestRunNegativeAfterFinish(t *testing.T) {
	funcsSlice := make([]func() error, 0, NumberOfFuncs)
	countOfFuncsChannel := make(chan struct{}, NumberOfFuncs)

	function := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithoutError()
	}

	functionFail := func() error {
		countOfFuncsChannel <- struct{}{}
		return funcWithError()
	}

	for i := 0; i < NumberOfFuncs-QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, function)
	}

	for i := 0; i < QuoteOfErrors; i++ {
		funcsSlice = append(funcsSlice, functionFail)
	}
	err := Run(funcsSlice, NumberOfGoroutines, QuoteOfErrors)

	if err == nil {
		t.Error("Result must be negative")
	}

	if len(countOfFuncsChannel) != NumberOfFuncs {
		t.Errorf("Count of executed funcs is %v, expected %v", len(countOfFuncsChannel), NumberOfFuncs)
	}
}


