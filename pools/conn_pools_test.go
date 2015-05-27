package pools

import (
	"github.com/openinx/muker/mysql"
	"testing"
	//"time"
)

func TestNewConnPool(t *testing.T) {
	cp, err := NewConnPool(100)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if cp == nil {
		t.Errorf("got nil conn pool")
	}

	/*
		cli, err2 := mysql.NewClient()
		if err2 != nil {
			t.Errorf("new client error :%v", err2)
		}

		err3 := cp.Put(cli)
		if err3 != nil {
			t.Errorf("put error: %v", err3)
		}
	*/

	cli, err := cp.Get()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = cp.Put(cli)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = cp.Close()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	//time.Sleep(10000 * time.Second)
}

func TestConnPoolLargeSize(t *testing.T) {
	cp, err := NewConnPool(10000)
	if err == nil {
		t.Errorf("max connection error is not found.")
	}
	if cp.Cap() != 10000 {
		t.Errorf("cap is not 10000")
	}
	err = cp.Close()
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestConnPoolGetError(t *testing.T) {
	cp, err := NewConnPool(50)

	if err != nil {
		t.Errorf("%v", err)
	}

	connList := make([]*mysql.Client, 50)
	for i := 0; i < 50; i++ {
		mycli, err := cp.Get()
		if err != nil {
			t.Errorf("error %d: %v", i, err)
		}
		connList[i] = mycli
	}
	if len(connList) != 50 {
		t.Errorf("error: connList len is not 50")
	}
	for i := 0; i < 50; i++ {
		err = cp.Put(connList[i])
		if err != nil {
			t.Errorf("put conn error: %v", err)
		}
	}

	if err = cp.Close(); err != nil {
		t.Errorf("close conn pool error: %v", err)
	}
}
