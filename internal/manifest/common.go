package manifest

type Channel struct {
	Version      string `yaml:"version,omitempty"`
	Id           string `yaml:"id,omitempty"`
	Owner        string `yaml:"owner,omitempty"`
	ContractFile string `yaml:"contract-file,omitempty"`
	AbiFile      string `yaml:"abi-file,omitempty"`
}

type GatewayNode struct {
	Address string `yaml:"address,omitempty"`
}

type Tx struct {
	Tx *Transaction `yaml:"tx"`
}

type Transaction struct {
	Function string   `yaml:"function,omitempty"`
	Args     []string `yaml:"args,omitempty"`
	Receipt  *Receipt `yaml:"receipt,omitempty"`
}

type Receipt struct {
	Status int32  `yaml:"status,omitempty"`
	Result string `yaml:"result,omitempty"`
}
