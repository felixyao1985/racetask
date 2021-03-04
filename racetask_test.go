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
		return taskTest{"felix2"}, errors.New("fake test")
	})

	ti.AddWithCtx(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		t.Log(2)
		return taskTest{"felix3"}, nil
	})

	ret, err := ti.Run()
	t.Log("Run  : ", ret, err)
	time.Sleep(time.Second)
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
}
