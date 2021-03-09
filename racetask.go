package racetask

import (
	"context"
	"errors"
	"sync"
	"time"
)

type NilResq struct {
}

const TaskError = string("task: err")
const TimeOutError = string("task: timeout")

type RaceTask interface {
	Add(...func() (interface{}, error))
	AddWithCtx(...func(context.Context) (interface{}, error))
	Run(...int) (interface{}, error)
	SetTimeOut(time.Duration)
	SetErrIgnore(bool)
	TimeOut(func() (interface{}, error))
	TimeOutWithCtx(func(context.Context) (interface{}, error))
}

type export struct {
	itf interface{}
	err error
}

type task struct {
	ctx         context.Context
	once        sync.Once
	errIgnore   bool
	cancelFunc  func()
	jobs        []func(context.Context) (interface{}, error)
	timeout     time.Duration //default 10 minute
	timeoutFunc func() (interface{}, error)
	export      export
}

// New create task pool with context cancel if job has error
func New(ctx context.Context) RaceTask {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancelFunc := context.WithCancel(ctx)

	return &task{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		jobs:       make([]func(context.Context) (interface{}, error), 0),
		timeout:    time.Minute * 15,
		timeoutFunc: func() (interface{}, error) {
			return NilResq{}, errors.New(TimeOutError)
		},
		errIgnore: true,
	}
}

// Add job to task pool
func (t *task) SetTimeOut(td time.Duration) {
	t.timeout = td
}

/*
true:    ignore err and return fastest results
false:   return fastest no err results
         if all err return task err
*/
func (t *task) SetErrIgnore(b bool) {
	t.errIgnore = b
}

// Add job to task pool
func (t *task) Add(jobs ...func() (interface{}, error)) {
	jw := make([]func(context.Context) (interface{}, error), len(jobs))

	for i := range jobs {
		i := i // fix shallow
		jw[i] = func(_ context.Context) (interface{}, error) {
			return jobs[i]()
		}
	}

	t.AddWithCtx(jw...)
}

// AddWithCtx add job with context
func (t *task) AddWithCtx(job ...func(context.Context) (interface{}, error)) {
	t.jobs = append(t.jobs, job...)
}

// Add TimeOut job to task pool
func (t *task) TimeOut(jobs func() (interface{}, error)) {
	jw := func(_ context.Context) (interface{}, error) {
		return jobs()
	}

	t.AddWithCtx(jw)
}

// TimeOutCtx add job with context
func (t *task) TimeOutWithCtx(job func(context.Context) (interface{}, error)) {
	t.jobs = append(t.jobs, job)
}

func (t *task) Run(n ...int) (interface{}, error) {
	jl := len(t.jobs)
	if jl == 0 {
		return NilResq{}, nil
	}
	ret := make(chan export)
	dones := 0
	onceBody := func() {
		ret <- t.export
	}
	for i := 0; i < jl; i++ {
		go func(job func(context.Context) (interface{}, error)) {
			itf, err := job(t.ctx)
			dones++
			if err == nil || t.errIgnore {
				t.export = export{
					err: err,
					itf: itf,
				}
				t.once.Do(onceBody)
			} else if dones == jl {
				t.export = export{
					err: errors.New(TaskError),
				}
				t.once.Do(onceBody)
			}
		}(t.jobs[i])
	}
	select {
	case r := <-ret:
		t.cancelFunc()
		return r.itf, r.err
	case <-time.After(t.timeout):
		return t.timeoutFunc()
	}
}
