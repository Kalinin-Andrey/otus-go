package hw05

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"sync"
)

// GoroutineData is data for goroutine starting
type GoroutineData struct{
	TasksChannel       chan func()error
	IfNeed2StartChannel chan struct{}
	ErrorsChannel       chan error
	WaitGroup          *sync.WaitGroup
}

func stoperByErrorsQuote(ctx context.Context, errorsCh chan error, errorsQuota uint) (errorsLimit chan struct{}) {
	errorsLimit = make(chan struct{}, 1)

	go func(){
		var errorsCount uint
	L1:
		for {
			select {
			case <-ctx.Done():
				break L1
			case <- errorsCh:
				errorsCount++

				if errorsCount == errorsQuota {
					errorsLimit <- struct{}{}
					log.Printf("error limit %v has exceed", errorsCount)
					break L1
				}
			}
		}
	}()
	return errorsLimit
}

// Run function starts proceed of tasks in N gorutines with M maximum possible errors
func Run(tasks []func()error, N int, M int) (err error) {

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
	errorsLimit := stoperByErrorsQuote(ctx, errorsCh, uint(M))

L1:
	for _, task := range tasks {
		select {
		case <- errorsLimit:
			err = errors.New("error limit exceeded")
			break L1
		case tasksCh <- task:
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
	dataSlice	:= getGoroutineDataSlice(goroutinesQuota, tasksNum)
	tasksCh     = make(chan func() error)
	errorsCh	= make(chan error, errorsQuota)
	wg 			= &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer close(tasksCh)
		log.Print("funcs start: ")
L1:
		for i := uint(0); ; {
			select {
			case <-ctx.Done():
				break L1
			case <- (*dataSlice)[i].IfNeed2StartChannel:
				wg.Add(1)
				(*dataSlice)[i].TasksChannel = tasksCh
				(*dataSlice)[i].ErrorsChannel = errorsCh
				(*dataSlice)[i].WaitGroup = wg
				go funcProceed(&(*dataSlice)[i])
			default:
			}

			i++
			if i == goroutinesQuota {
				i = 0
			}
		}
		log.Print("goroutinesStarter has stop")
		wg.Done()
	}()
	return tasksCh, wg, errorsCh
}

// getGoroutineDataSlice is get slice of GoroutineData
func getGoroutineDataSlice(N uint, tasksNum int) *[]GoroutineData {
	dataSlice := make([]GoroutineData, 0, N)

	for i := uint(0); i < N; i++ {
		data := GoroutineData{
			IfNeed2StartChannel: make(chan struct{}, 1),
		}
		data.IfNeed2StartChannel <- struct{}{}
		dataSlice = append(dataSlice, data)
	}
	return &dataSlice
}

// funcProceed proceeds a given func
func funcProceed(data *GoroutineData) {
	defer func(){
		if err := recover(); err != nil {
			data.ErrorsChannel <- errors.Wrap(err.(error), "panic happened in given func")
			data.IfNeed2StartChannel <- struct{}{}
		}
		data.WaitGroup.Done()
	}()

	for task := range data.TasksChannel {
		err := task()
		log.Print(".")

		if err != nil {
			data.ErrorsChannel <- err
		}
	}
}