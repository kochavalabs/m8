package cfg

import (
	"errors"
	"io/ioutil"

	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"gopkg.in/yaml.v2"
)

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

// FromFile loads the cli Configuration at a given path, returns and error if the file does not exists
// or is malformed
func FromFile(path string) (*Configuration, error) {
	cfgfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Configuration{}
	if err := yaml.Unmarshal(cfgfile, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ActiveChannelId returns an error if a active channel is not found.
func (c *Configuration) ActiveChannelId() (xdr.ID, error) {
	if c.User == nil {
		return xdr.ID{}, errors.New("missing user cfg")
	}

	for _, channel := range c.Channels {
		if channel.Channel.ChannelAlias == c.User.ActiveChannel {
			id, err := xdr.IDFromHexString(channel.Channel.ChannelID)
			if err != nil {
				return xdr.ID{}, err
			}
			return id, nil
		}
	}

	return xdr.ID{}, errors.New("no active channel found in cfg")
}
