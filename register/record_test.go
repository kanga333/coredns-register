package register

import (
	"reflect"
	"testing"
)

func TestRecordAdd(t *testing.T) {
	recordA := SRVRecord{Domain: "a"}
	recordB := SRVRecord{Domain: "b"}

	want := Records{
		SRV: []SRVRecord{recordA, recordB},
	}

	a := Records{SRV: []SRVRecord{recordA}}
	b := Records{SRV: []SRVRecord{recordB}}

	a.Add(&b)
	if !reflect.DeepEqual(a, want) {
		t.Errorf("Load got: %v, want: %v", a, want)
	}
}

func TestRecordAddSRV(t *testing.T) {
	recordA := SRVRecord{Domain: "a", Address: "127.0.0.1", Port: 80}
	recordB := SRVRecord{Domain: "b", Address: "127.0.0.1", Port: 80}

	want := Records{
		SRV: []SRVRecord{recordA, recordB},
	}

	a := Records{SRV: []SRVRecord{recordA}}
	a.AddSRV("b", "127.0.0.1", 80)
	if !reflect.DeepEqual(a, want) {
		t.Errorf("Load got: %v, want: %v", a, want)
	}
}
