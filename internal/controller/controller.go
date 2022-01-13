package controller

import (
	"crypto/ed25519"

	"github.com/kochavalabs/mazzaroth-go"
)

type Controller struct {
	client mazzaroth.Client
	pk     *ed25519.PrivateKey
}

func NewController(gatewateAddress string, pk *ed25519.PrivateKey) (*Controller, error) {
	return &Controller{}, nil
}

func (c *Controller) BlockLookup(channelID string, blockid string, headers bool) ([]byte, error) {

	return nil, nil
}
