package register

type Records struct {
	SRV []SRVRecord `yaml:"srv"`
}

func (r *Records) Add(a *Records) {
	if len(a.SRV) != 0 {
		r.SRV = append(r.SRV, a.SRV...)
	}
}

func (r *Records) InitAddress(a string) {
	for i := range r.SRV {
		r.SRV[i].InitAddress(a)
	}
}

type SRVRecord struct {
	Domain  string `yaml:"domain"`
	Address string `yaml:"address,omitempty"`
	Port    int    `yaml:"port"`
}

func (r *SRVRecord) InitAddress(a string) {
	if r.Address == "" {
		r.Address = a
	}
}
