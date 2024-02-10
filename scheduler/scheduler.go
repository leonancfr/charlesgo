package scheduler

import (
	"fmt"
	"reflect"
	"time"

	cron "github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var scheduler cron.Scheduler

func InitScheduler() error {
	if scheduler != nil {
		return nil
	}
	newScheduler, err := cron.NewScheduler()
	if err != nil {
		return fmt.Errorf("fail to create scheduler. %v", err)
	}
	scheduler = newScheduler
	scheduler.Start()
	return nil
}

func RegisterFunctionToSchedule(duration time.Duration, function interface{}, parameters ...interface{}) (interface{}, error) {
	if scheduler == nil {
		return nil, fmt.Errorf("the scheduler needs to be initialized")
	}

	funcType := reflect.TypeOf(function)

	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("provided 'function' is not a function")
	}

	expectedParams := funcType.NumIn()
	actualParams := len(parameters)

	if expectedParams != actualParams {
		return nil, fmt.Errorf("incorrect number of parameters. Expected %d, got %d", expectedParams, actualParams)
	}

	job, err := scheduler.NewJob(
		cron.DurationJob(duration),
		cron.NewTask(function, parameters...),
	)

	if err != nil {
		return nil, fmt.Errorf("cannot create the job. %v", err)
	}
	return job.ID(), nil
}

func RemoveFunctionFromSchedule(jobId interface{}) error {
	if scheduler == nil {
		return fmt.Errorf("the scheduler needs to be initialized")
	}

	decodedJobId, ok := jobId.(uuid.UUID)
	if !ok {
		return fmt.Errorf("invalid jobId type")
	}

	if err := scheduler.RemoveJob(decodedJobId); err != nil {
		return fmt.Errorf("failed to remove job. %v", err)
	}

	return nil
}
