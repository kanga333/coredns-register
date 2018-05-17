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
