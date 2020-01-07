package hw05

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"sync"
)

// GoroutineData is data for goroutine starting
type GoroutineData struct{
	TasksChannel       chan func()error
	IfIsFreeChannel     chan struct{}
	IfNeed2StartChannel chan struct{}
	ErrorsChannel       chan error
	WaitGroup          *sync.WaitGroup
}

// Run function starts proceed of tasks in N gorutines with M maximum possible errors
func Run(tasks []func()error, N int, M int) (err error) {

	if N < 0 || M < 0 {
		return errors.New("N and M params mast be more then 0")
	}

	if N == 0 || M == 0 || len(tasks) == 0 {
		return nil
	}

	ctx, finish		:= context.WithCancel(context.Background())
	dataSlice, wg, errorsCh	:= goroutinesStarter(ctx, uint(N), uint(M), len(tasks))
	taskNum			:= 0

	for i := 0; ; {

		if len(errorsCh) == M {
			err = errors.New("error limit exceeded")
			log.Print("error limit exceeded")
			break
		}

		if taskNum == len(tasks) {
			log.Print("tasks are over")
			break
		}
		_, ok := <- (*dataSlice)[i].IfIsFreeChannel

		if ok {
			(*dataSlice)[i].TasksChannel <- tasks[taskNum]
			taskNum++
		}

		i++
		if i == N {
			i = 0
		}
	}
	log.Print("send finish signal")
	finish()
	wg.Wait()
	//	check again
	if err == nil && len(errorsCh) == M {
		log.Print("error limit exceeded in the end")
	}
	log.Print("Done!")
	return err
}

// GorutinesStarter starts proceed of tasks in N gorutines with M maximum possible errors
func goroutinesStarter(ctx context.Context, goroutinesQuota uint, errorsQuota uint, tasksNum int) (dataSlice *[]GoroutineData, wg *sync.WaitGroup, errorsCh chan error) {
	dataSlice	= getGoroutineDataSlice(goroutinesQuota, tasksNum)
	errorsCh	= make(chan error, errorsQuota)
	wg 			= &sync.WaitGroup{}

	go func() {
		defer closeTasksChannels(dataSlice)
		log.Print("funcs start: ")
L1:
		for i := uint(0); ; {
			select {
			case <-ctx.Done():
				break L1
			case <- (*dataSlice)[i].IfNeed2StartChannel:
				wg.Add(1)
				(*dataSlice)[i].ErrorsChannel = errorsCh
				(*dataSlice)[i].WaitGroup = wg
				(*dataSlice)[i].IfIsFreeChannel <- struct{}{}
				go funcProceed(&(*dataSlice)[i])
			default:
			}

			i++
			if i == goroutinesQuota {
				i = 0
			}
		}
		log.Print("goroutinesStarter has stop")
	}()
	return dataSlice, wg, errorsCh
}

// getGoroutineDataSlice is get slice of GoroutineData
func getGoroutineDataSlice(N uint, tasksNum int) *[]GoroutineData {
	dataSlice := make([]GoroutineData, 0, N)

	for i := uint(0); i < N; i++ {
		data := GoroutineData{
			TasksChannel:       make(chan func()error, tasksNum),
			IfIsFreeChannel:     make(chan struct{}, 1),
			IfNeed2StartChannel: make(chan struct{}, 1),
		}
		data.IfNeed2StartChannel <- struct{}{}
		dataSlice = append(dataSlice, data)
	}
	return &dataSlice
}

// closeTasksChannels is close tasks channels for a given []GoroutineData
func closeTasksChannels(dataSlice *[]GoroutineData) {
	for i, data := range *dataSlice{
		close(data.TasksChannel)
		log.Print("Close tasks channel ", i)
	}
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
		data.IfIsFreeChannel <- struct{}{}
	}
}