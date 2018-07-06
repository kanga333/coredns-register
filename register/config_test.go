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
	got, err := conf.loadRecords()
	if err != nil {
		t.Fatalf("LoadRecords return error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRecords got: %v, want: %v", got, want)
	}
	if got, want := conf.Records, recordsA; !reflect.DeepEqual(got, want) {
		t.Errorf("Records after LoadRecords got: %v, want: %v", got, want)
	}

}

func TestCreateScheduker(t *testing.T) {
	records := Records{
		SRV: []SRVRecord{
			SRVRecord{
				Domain:  "a.domain.test",
				Address: "127.0.0.1",
				Port:    80,
			},
		},
	}
	etcd := EtcdConfig{
		DiscoverySRV: "dns.domain.test",
		Basepath:     "/base",
		Endpoints:    []string{"http://127.0.0.1:2379"},
	}
	c := &Config{
		Hostname: "host",
		Address:  "127.0.0.1",
		Interval: 60,
		Etcd:     etcd,
		Records:  records,
	}

	want := &Scheduler{
		interval: 60,
		register: &EtcdRegister{
			hostname: "host",
			etcd:     etcd,
		},
		records: &records,
		logger:  nil,
	}

	got, err := c.CreateScheduler(nil)
	if err != nil {
		t.Fatalf("CreateScheduler return error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateScheduler got: %v, want: %v", got, want)
	}
}

func TestParseSRV(t *testing.T) {
	records := Records{
		SRV: []SRVRecord{
			SRVRecord{
				Domain:  "a.domain.test",
				Address: "127.0.0.1",
				Port:    80,
			},
		},
	}
	c := &Config{
		Address:    "127.0.0.1",
		SRVRecords: "b.domain.test:81,c.domain.test:82",
		Records:    records,
	}

	want := Records{
		SRV: []SRVRecord{
			SRVRecord{
				Domain:  "a.domain.test",
				Address: "127.0.0.1",
				Port:    80,
			},
			SRVRecord{
				Domain:  "b.domain.test",
				Address: "127.0.0.1",
				Port:    81,
			},
			SRVRecord{
				Domain:  "c.domain.test",
				Address: "127.0.0.1",
				Port:    82,
			},
		},
	}

	err := c.parseSRVRecords()
	if err != nil {
		t.Fatalf("parseSRVRecords return error: %v", err)
	}

	if !reflect.DeepEqual(c.Records, want) {
		t.Errorf("CreateScheduler got: %v, want: %v", c.Records, want)
	}
}
