package register

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Hostname    string     `yaml:"hostname"`
	Address     string     `yaml:"address"`
	Interval    int        `yaml:"interval"`
	Etcd        EtcdConfig `yaml:"etcd"`
	RecordFiles []string   `yaml:"record_files,omitempty"`
	Records     Records    `yaml:"records,omitempty"`
}

func (c *Config) LoadRecords() (*Records, error) {
	recordFiles := []string{}
	for _, path := range c.RecordFiles {
		files, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		recordFiles = append(recordFiles, files...)
	}
	records, err := loadRecords(recordFiles)
	if err != nil {
		return nil, err
	}
	records.Add(&c.Records)
	records.InitAddress(c.Address)
	return records, nil
}

func (c *Config) CreateScheduler(lg *zap.Logger) (*Scheduler, error) {
	records, err := c.LoadRecords()
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

func LoadFile(filename string, v interface{}) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.UnmarshalStrict(content, v)
	if err != nil {
		return fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}

	return nil
}

func loadRecords(files []string) (*Records, error) {
	records := &Records{}
	for _, file := range files {
		tmp := &Records{}
		err := LoadFile(file, tmp)
		if err != nil {
			return nil, err
		}
		records.Add(tmp)
	}
	return records, nil
}
