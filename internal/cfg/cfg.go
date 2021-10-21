package cfg

type Configuration struct {
	PrivateKey    string                 `yaml:"private-key"`
	PublicKey     string                 `yaml:"public-key"`
	ActiveChannel string                 `yaml:"active-channel"`
	Channels      []ChannelConfiguration `yaml:"channels"`
}

type ChannelConfiguration struct {
	ChannelURL   string `yaml:"channel-url"`
	ChannelID    string `yaml:"channel-id"`
	ChannelAlias string `yaml:"channel-alias"`
}
