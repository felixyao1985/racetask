package racetask

import (
	"context"
	"errors"
	"testing"
	"time"
)

type taskTest struct {
	name string
}

func TestTask(t *testing.T) {
	ti := New(context.Background())

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 300)
		t.Log(1)
		return taskTest{"felix1"}, nil
	})

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 200)
		t.Log(2)
		return taskTest{"felix2"}, errors.New("fake test")
	})

	ti.AddWithCtx(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		t.Log(3)
		return taskTest{"felix3"}, errors.New("fake test")
	})

	ret, err := ti.Run()
	t.Log("Run  : ", ret, err)
	time.Sleep(time.Second)
	r := ret.(taskTest)
	if r.name != "felix3" {
		t.Error("want felix3 return ", r.name)
	}
}

func TestTaskSetError(t *testing.T) {
	ti := New(context.Background())
	ti.SetErrIgnore(false)

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 300)
		t.Log(1)
		return taskTest{"felix1"}, nil
	})

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 200)
		t.Log(2)
		return taskTest{"felix2"}, errors.New("fake test")
	})

	ti.AddWithCtx(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		t.Log(3)
		return taskTest{"felix3"}, errors.New("fake test")
	})

	ret, err := ti.Run()
	t.Log("Run  : ", ret, err)
	time.Sleep(time.Second)
	r := ret.(taskTest)
	if r.name != "felix1" {
		t.Error("want felix1 return ", r.name)
	}
}

func TestTaskAllError(t *testing.T) {
	ti := New(context.Background())
	ti.SetErrIgnore(false)

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 300)
		t.Log(1)
		return taskTest{"felix1"}, errors.New("fake test")
	})

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 200)
		t.Log(2)
		return taskTest{"felix2"}, errors.New("fake test")
	})

	ti.AddWithCtx(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		t.Log(3)
		return taskTest{"felix3"}, errors.New("fake test")
	})

	ret, err := ti.Run()
	t.Log("Run  : ", ret, err)
	time.Sleep(time.Second)
	if err.Error() != TaskError {
		t.Error("want  TaskError  return ", err.Error())
	}
}

func TestTaskTimeOut(t *testing.T) {
	ti := New(context.Background())

	ti.Add(func() (interface{}, error) {
		time.Sleep(time.Second * 10)
		t.Log(1)
		return taskTest{"felix"}, nil
	})

	ti.SetTimeOut(time.Millisecond * 300)

	ret, err := ti.Run()
	t.Log("Run  : ", ret, err)
	if err.Error() != TimeOutError {
		t.Error("want  TimeOutError  return ", err.Error())
	}
}
