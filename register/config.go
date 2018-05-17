package register

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	yaml "gopkg.in/yaml.v2"
)

// Config represents setting of coredns-register.
type Config struct {
	Hostname    string     `yaml:"hostname"`
	Address     string     `yaml:"address"`
	Interval    int        `yaml:"interval"`
	Etcd        EtcdConfig `yaml:"etcd"`
	RecordFiles []string   `yaml:"record_files,omitempty"`
	Records     Records    `yaml:"records,omitempty"`
}

// LoadFile reads yaml in filename and unmarshal it to v.
func LoadFile(filename string, v interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	expandContent := []byte(os.ExpandEnv(string(content)))
	err = yaml.UnmarshalStrict(expandContent, v)
	if err != nil {
		return fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}

	return nil
}

// CreateScheduler creates a scheduler based on the contents of config.
func (c *Config) CreateScheduler(lg *zap.Logger) (*Scheduler, error) {
	records, err := c.loadRecords()
	if err != nil {
		return nil, err
	}

	s := &Scheduler{
		interval: c.Interval,
		register: &EtcdRegister{
			hostname: c.Hostname,
			etcd:     c.Etcd,
		},
		records: records,
		logger:  lg,
	}

	return s, nil
}

func (c *Config) loadRecords() (*Records, error) {
	recordFiles := []string{}
	for _, path := range c.RecordFiles {
		files, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		recordFiles = append(recordFiles, files...)
	}

	records := &Records{}
	for _, file := range recordFiles {
		tmp := &Records{}
		err := LoadFile(file, tmp)
		if err != nil {
			return nil, err
		}
		records.Add(tmp)
	}

	records.Add(&c.Records)
	records.InitAddress(c.Address)
	return records, nil
}
