package register

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientEndpoints(t *testing.T) {
	ec := EtcdConfig{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	_, err := ec.newClient()
	if err != nil {
		t.Fatalf("NewClient return error: %v", err)
	}
}

func TestNewClient_Nil(t *testing.T) {
	ec := EtcdConfig{}
	_, err := ec.newClient()
	if err == nil || err.Error() != "Endpoint is not set" {
		t.Errorf("Empty EtcdConfig does not return error")
	}
}

func TestGenerateKey(t *testing.T) {
	bc := Config{
		Hostname: "host",
		Etcd: EtcdConfig{
			Basepath: "/base",
		},
	}
	er := EtcdRegister{
		hostname: bc.Hostname,
		etcd:     bc.Etcd,
	}
	if got, want := er.generateKey("domain.test"), "/base/test/domain/host"; got != want {
		t.Errorf("generateKey returns: %v, want: %v", got, want)

	}
}

func TestSRVRegister(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/keys/base/test/domain/host" {
			t.Errorf("Invalid URL: %v", r.URL)
		}
		w.Write([]byte(`{"host":"127.0.0.1","port":80}`))
	})
	server := httptest.NewServer(handler)
	endpoints := []string{server.URL}

	bc := Config{
		Hostname: "host",
		Address:  "127.0.0.1",
		Etcd: EtcdConfig{
			Basepath:  "/base",
			Endpoints: endpoints,
		},
	}
	er := EtcdRegister{
		hostname: bc.Hostname,
		etcd:     bc.Etcd,
	}
	r := SRVRecord{
		Domain:  "domain.test",
		Address: "127.0.0.1",
		Port:    80,
	}
	err := er.SRVRegister(r)
	if err != nil {
		t.Fatalf("SRVRegister return err: %v", err)
	}
}

func TestSRVUnregister(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/keys/base/test/domain/host" {
			t.Errorf("Invalid URL: %v", r.URL)
		}
		w.Write([]byte(`{"host":"127.0.0.1","port":80}`))
	})
	server := httptest.NewServer(handler)
	endpoints := []string{server.URL}

	bc := Config{
		Hostname: "host",
		Address:  "127.0.0.1",
		Etcd: EtcdConfig{
			Basepath:  "/base",
			Endpoints: endpoints,
		},
	}
	er := EtcdRegister{
		hostname: bc.Hostname,
		etcd:     bc.Etcd,
	}
	r := SRVRecord{
		Domain:  "domain.test",
		Address: "127.0.0.1",
		Port:    80,
	}
	err := er.SRVUnregister(r)
	if err != nil {
		t.Fatalf("SRVUnregister return err: %v", err)
	}
}

func Test_generateSRVValue(t *testing.T) {
	r := SRVRecord{
		Domain:  "domain.test",
		Address: "127.0.0.1",
		Port:    80,
	}

	if got, want := generateSRVValue(r), `{"host":"127.0.0.1","port":80}`; got != want {
		t.Errorf("generateKey returns: %v, want: %v", got, want)

	}
}
