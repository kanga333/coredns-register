package register

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

type MockRegister struct {
	SRVRegisterCount   int
	SRVUnregisterCount int
}

func (r *MockRegister) SRVRegister(record SRVRecord) error {
	r.SRVRegisterCount++
	return nil
}

func (r *MockRegister) SRVUnregister(record SRVRecord) error {
	fmt.Println("call")
	r.SRVUnregisterCount++
	return nil
}

func TestScheduler(t *testing.T) {
	mock := MockRegister{}

	recordA := SRVRecord{Domain: "a"}
	recordB := SRVRecord{Domain: "b"}
	records := &Records{
		SRV: []SRVRecord{recordA, recordB},
	}

	s := Scheduler{
		interval: 2,
		register: &mock,
		records:  records,
	}

	sig := make(chan os.Signal, 1)
	go s.Start(sig)

	time.Sleep(3 * time.Second)
	sig <- syscall.SIGTERM
	time.Sleep(1 * time.Second)

	if mock.SRVRegisterCount != 4 {
		t.Errorf("Load got: %v, want: %v", mock.SRVRegisterCount, 4)
	}

	fmt.Println("hoge")

	if mock.SRVUnregisterCount != 2 {
		t.Errorf("Load got: %v, want: %v", mock.SRVUnregisterCount, 2)
	}

}
