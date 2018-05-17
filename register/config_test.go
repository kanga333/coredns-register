package register

import (
	"reflect"
	"testing"
)

func TestLoadFile(t *testing.T) {
	want := &Config{
		Hostname: "host",
		Address:  "127.0.0.1",
		Interval: 60,
		Etcd: EtcdConfig{
			DiscoverySRV: "dns.domain.test",
			Basepath:     "/base",
			Endpoints:    []string{"http://127.0.0.1:2379"},
		},
		RecordFiles: []string{"fixtures/record.d/*yml"},
		Records: Records{
			SRV: []SRVRecord{
				SRVRecord{
					Domain: "a.domain.test",
					Port:   80,
				},
			},
		},
	}
	got := &Config{}
	err := LoadFile("fixtures/coredns-register.yml", got)
	if err != nil {
		t.Fatalf("Load return error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load got: %v, want: %v", got, want)
	}
}

func TestInvalidLoadFile(t *testing.T) {
	v := &Config{}
	err := LoadFile("fixtures/invalid.yml", v)
	if err == nil {
		t.Error("Load should return error")
	}
}

func TestBrokenPathLoadFile(t *testing.T) {
	v := &Config{}
	err := LoadFile("fixtures/broken-path.yml", v)
	if err == nil {
		t.Error("Load should return error")
	}
}

func TestLoadRecords(t *testing.T) {
	recordsA := Records{
		SRV: []SRVRecord{
			SRVRecord{
				Domain:  "a.domain.test",
				Address: "127.0.0.2",
			},
		},
	}

	want := &Records{
		SRV: []SRVRecord{
			SRVRecord{
				Domain:  "b.domain.test",
				Address: "127.0.0.1",
			},
			SRVRecord{
				Domain:  "a.domain.test",
				Address: "127.0.0.2",
			},
		},
	}
	conf := &Config{
		Hostname: "host",
		Address:  "127.0.0.1",
		Etcd: EtcdConfig{
			DiscoverySRV: "dns.domain.test",
			Basepath:     "/base",
			Endpoints:    []string{"http://127.0.0.1:2379"},
		},
		RecordFiles: []string{"fixtures/record.d/*yml"},
		Records:     recordsA,
	}
	got, err := conf.LoadRecords()
	if err != nil {
		t.Fatalf("LoadRecords return error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRecords got: %v, want: %v", got, want)
	}
	if got, want := conf.Records, recordsA; !reflect.DeepEqual(got, want) {
		t.Errorf("Records after LoadRecords want: %v, got: %v", got, want)
	}

}
