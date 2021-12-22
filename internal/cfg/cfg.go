package cfg

type Configuration struct {
	Version  string       `yaml:"version"`
	User     *UserCfg     `yaml:"user"`
	Channels []ChannelCfg `yaml:"channels"`
}

type UserCfg struct {
	PrivateKey    string `yaml:"private-key"`
	PublicKey     string `yaml:"public-key"`
	ActiveChannel string `yaml:"active-channel"`
}

type ChannelCfg struct {
	Channel Channel `yaml:"channel"`
}

type Channel struct {
	ChannelURL   string `yaml:"channel-url"`
	ChannelID    string `yaml:"channel-id"`
	ChannelAlias string `yaml:"channel-alias"`
}

func FromFile(path string) (*Configuration, error) {
	return nil, nil
}

func ToFile(path string, cfg *Configuration) error {
	return nil
}
