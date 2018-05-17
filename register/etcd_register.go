package register

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
)

// Register registers DNS record.
type Register interface {
	SRVRegister(record SRVRecord) error
	SRVUnregister(record SRVRecord) error
}

// EtcdConfig represents connection setting with etcd v2.
type EtcdConfig struct {
	DiscoverySRV string   `yaml:"discovery-srv,omitempty"`
	Endpoints    []string `yaml:"endpoints,omitempty"`
	Basepath     string   `yaml:"basepath"`
}

func (c *EtcdConfig) newClient() (etcd.KeysAPI, error) {
	var endpoints []string
	endpoints = c.Endpoints

	if c.DiscoverySRV != "" {
		d := etcd.NewSRVDiscover()
		result, err := d.Discover(c.DiscoverySRV)
		if err != nil {
			return nil, err
		}
		endpoints = result
	}

	if endpoints == nil {
		return nil, errors.New("Endpoint is not set")
	}

	cfg := etcd.Config{
		Endpoints:               endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}

	return etcd.NewKeysAPI(client), nil
}

// EtcdRegister manages DNS registration to the etcd backend of CoreDNS.
type EtcdRegister struct {
	hostname string
	etcd     EtcdConfig
}

// SRVRegister registers SRV record.
func (r *EtcdRegister) SRVRegister(record SRVRecord) error {
	client, err := r.etcd.newClient()
	if err != nil {
		return err
	}
	key := r.generateKey(record.Domain)
	value := generateSRVValue(record)
	_, err = client.Set(context.Background(), key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

// SRVUnregister unregisters SRV record.
func (r *EtcdRegister) SRVUnregister(record SRVRecord) error {
	client, err := r.etcd.newClient()
	if err != nil {
		return err
	}
	key := r.generateKey(record.Domain)
	_, err = client.Delete(context.Background(), key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *EtcdRegister) generateKey(domain string) string {
	domains := strings.Split(domain, ".")
	reverse(domains)

	var paths []string
	paths = append(paths, r.etcd.Basepath)
	paths = append(paths, domains...)
	paths = append(paths, r.hostname)
	return strings.Join(paths, "/")
}

func generateSRVValue(r SRVRecord) string {
	return fmt.Sprintf("{\"host\":\"%s\",\"port\":%d}", r.Address, r.Port)
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
