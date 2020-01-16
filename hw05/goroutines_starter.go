package hw05

import (
	"context"
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

	ctx, finish := context.WithCancel(context.Background())
	tasksCh, wg, errorsCh := goroutinesStarter(ctx, uint(N), uint(M), len(tasks))
	defer close(errorsCh)
	errorsLimit := stopperByErrorsQuota(errorsCh, uint(M))

L1:
	for _, task := range tasks {

		select {
		case <-errorsLimit:
			err = errors.New("error limit exceeded")
			break L1
		default:

			select {
			case <-errorsLimit:
				err = errors.New("error limit exceeded")
				break L1
			case tasksCh <- task:
			}
		}
	}

	if err == nil {
		log.Print("tasks are over")
	}
	finish()
	log.Print("send finish signal")
	wg.Wait()
	//	check again
	if err == nil && len(errorsCh) == M {
		log.Print("error limit exceeded in the end")
	}
	log.Print("Done!")
	return err
}

// GorutinesStarter starts proceed of tasks in N gorutines with M maximum possible errors
func goroutinesStarter(ctx context.Context, goroutinesQuota uint, errorsQuota uint, tasksNum int) (tasksCh chan func() error, wg *sync.WaitGroup, errorsCh chan error) {
	tasksCh = make(chan func() error)
	errorsCh = make(chan error, goroutinesQuota)
	wg = &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer close(tasksCh)
		log.Print("funcs start: ")

		for i := 0; i < int(goroutinesQuota); i++ {
			wg.Add(1)
			go taskPropceed(wg, tasksCh, errorsCh)
		}
		log.Print("all goroutines has start")
		<-ctx.Done()
		log.Print("goroutinesStarter has stop")
		wg.Done()
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
