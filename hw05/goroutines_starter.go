package hw05

import (
	"github.com/pkg/errors"
	"log"
	"strconv"
	"sync"
)

// stopperByErrorsQuote is ther stopper by errors quota
func stopperByErrorsQuota(errorsCh chan error, errorsQuota uint) (errorsLimit chan struct{}) {
	errorsLimit = make(chan struct{}, 1)

	go func() {
		defer close(errorsLimit)
		var errorsCount uint

		for range errorsCh {
			errorsCount++

			if errorsCount == errorsQuota {
				errorsLimit <- struct{}{}
				log.Printf("error limit %v has exceed", errorsCount)
			}
		}
	}()
	return errorsLimit
}

// Run function starts proceed of tasks in N gorutines with M maximum possible errors
func Run(tasks []func() error, N int, M int) (err error) {

	if N < 0 || M < 0 {
		return errors.New("N and M params mast be more then 0")
	}

	if N == 0 && len(tasks) > 0 {
		return errors.New("there is no worker for given " + strconv.Itoa(len(tasks)) + " tasks")
	}

	if N == 0 || M == 0 || len(tasks) == 0 {
		return nil
	}

	tasksCh, wg, errorsCh := goroutinesStarter(uint(N), uint(M), len(tasks))
	defer close(errorsCh)
	errorsLimit := stopperByErrorsQuota(errorsCh, uint(M))

L1:
	for _, task := range tasks {

		select {
		case <-errorsLimit:
			err = errors.New("error limit exceeded")
			break L1
		default:
		}

		select {
		case <-errorsLimit:
			err = errors.New("error limit exceeded")
			break L1
		case tasksCh <- task:
		}
	}

	if err == nil {
		log.Print("tasks are over")
	}
	close(tasksCh)
	log.Print("close task channel")
	wg.Wait()

	if err == nil {
		//	check for error limit exceeded after finish
		select {
		case <-errorsLimit:
			err = errors.New("error limit exceeded")
			log.Print("error limit exceeded in the end")
		default:
		}
	}
	log.Print("Done!")
	return err
}

// GorutinesStarter starts proceed of tasks in N gorutines with M maximum possible errors
func goroutinesStarter(goroutinesQuota uint, errorsQuota uint, tasksNum int) (tasksCh chan func() error, wg *sync.WaitGroup, errorsCh chan error) {
	tasksCh = make(chan func() error)
	errorsCh = make(chan error, goroutinesQuota)
	wg = &sync.WaitGroup{}
	log.Print("goroutines are starting:")
	wg.Add(1)

	go func() {
		for i := 0; i < int(goroutinesQuota); i++ {
			wg.Add(1)
			go taskPropceed(wg, tasksCh, errorsCh)
		}
		wg.Done()
		log.Print("all goroutines has start")
	}()
	return tasksCh, wg, errorsCh
}

// taskPropceed proceeds a given task
func taskPropceed(wg *sync.WaitGroup, tasksCh chan func() error, errorsCh chan error) {
	defer func() {
		wg.Done()
	}()

	for task := range tasksCh {
		funcProceed(task, errorsCh)
	}
}

// funcProceed proceeds a given func
func funcProceed(f func() error, ErrorsChannel chan error) {
	defer func() {
		if err := recover(); err != nil {
			ErrorsChannel <- errors.New("panic happened in given func: " + err.(string))
		}
	}()

	err := f()
	log.Print(".")

	if err != nil {
		ErrorsChannel <- err
	}
}
