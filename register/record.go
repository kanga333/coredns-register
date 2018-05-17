package register

// Records represents a collection of DNS records.
type Records struct {
	SRV []SRVRecord `yaml:"srv"`
}

// Add adds Records to Records.
func (r *Records) Add(a *Records) {
	if len(a.SRV) != 0 {
		r.SRV = append(r.SRV, a.SRV...)
	}
}

// InitAddress sets an initial value for a DNS record whose address is empty.
func (r *Records) InitAddress(a string) {
	for i := range r.SRV {
		r.SRV[i].InitAddress(a)
	}
}

// SRVRecord represents SRV Record.
type SRVRecord struct {
	Domain  string `yaml:"domain"`
	Address string `yaml:"address,omitempty"`
	Port    int    `yaml:"port"`
}

// InitAddress sets an initial value for a DNS record whose address is empty.
func (r *SRVRecord) InitAddress(a string) {
	if r.Address == "" {
		r.Address = a
	}
}
