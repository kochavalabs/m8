package cfg

import (
	"errors"
)

type Configuration struct {
	Version  string        `yaml:"version"`
	User     *UserCfg      `yaml:"user"`
	Channels []*ChannelCfg `yaml:"channels"`
}

type UserCfg struct {
	PrivateKey    string `yaml:"private-key"`
	PublicKey     string `yaml:"public-key"`
	ActiveChannel string `yaml:"active-channel"`
}

type ChannelCfg struct {
	Channel *Channel `yaml:"channel"`
}

type Channel struct {
	ChannelAddress string `yaml:"channel-address"`
	ChannelID      string `yaml:"channel-id"`
	ChannelAlias   string `yaml:"channel-alias"`
}

// ActiveChannelId returns an error if a active channel is not found.
func (c *Configuration) ActiveChannel() (*Channel, error) {
	if c.User == nil {
		return nil, errors.New("missing user cfg")
	}

	for _, channel := range c.Channels {
		if channel.Channel.ChannelAlias == c.User.ActiveChannel {
			return channel.Channel, nil
		}
	}

	return nil, errors.New("no active channel found in cfg")
}

func (c *Configuration) ContainsChannel(channelId string, channelAlias string) bool {
	for _, channel := range c.Channels {
		if channel.Channel.ChannelID == channelId ||
			channel.Channel.ChannelAlias == channelAlias {
			return true
		}
	}
	return false
}
