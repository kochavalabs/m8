package manifest

type Manifest struct {
	Version      string      `yaml:"version"`
	Type         string      `yaml:"type"`
	Channel      Channel     `yaml:"channel"`
	GatewayNode  GatewayNode `yaml:"gateway-node"`
	Transactions []*Tx       `yaml:"transactions,omitempty"`
}

type Tx struct {
	Tx Transaction `yaml:"tx"`
}

func FromFile(path string) (*Manifest, error) {
	return nil, nil
}
