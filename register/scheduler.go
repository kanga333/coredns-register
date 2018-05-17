package register

import (
	"os"
	"time"

	"go.uber.org/zap"
)

type Scheduler struct {
	interval int
	register Register
	records  *Records
	logger   *zap.Logger
}

func (s *Scheduler) Start(signal chan os.Signal) {
	t := time.NewTicker(time.Duration(s.interval) * time.Second)
	defer t.Stop()
	s.srvRegister(s.records.SRV)
	defer s.srvUnregister(s.records.SRV)
LOOP:
	for {
		select {
		case <-t.C:
			s.srvRegister(s.records.SRV)
		case receive := <-signal:
			if s.logger != nil {
				s.logger.Sugar().Infof("receive signal %v", receive)
			}
			break LOOP
		}
	}
}

func (s *Scheduler) srvRegister(records []SRVRecord) {
	for _, record := range records {
		if s.logger != nil {
			s.logger.Info("register srv record", zap.String("domain", record.Domain))
		}
		err := s.register.SRVRegister(record)
		if err != nil {
			if s.logger != nil {
				s.logger.Error("failed to register srv", zap.Error(err))
			}
		}

	}
}

func (s *Scheduler) srvUnregister(records []SRVRecord) {
	for _, record := range records {
		if s.logger != nil {
			s.logger.Info("unregister srv record", zap.String("domain", record.Domain))
		}
		err := s.register.SRVUnregister(record)
		if err != nil {
			if s.logger != nil {
				s.logger.Error("failed to unregister srv", zap.Error(err))
			}
		}

	}
}
