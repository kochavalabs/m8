package manifest

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/xdr"
	"gopkg.in/yaml.v2"
)

const (
	maxRetry = 10
)

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

func loadAbi(path string) (*xdr.Abi, error) {
	abiFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	abi := &xdr.Abi{}
	if err := json.Unmarshal(abiFile, abi); err != nil {
		return nil, err
	}
	return abi, nil
}

func loadContract(path string) ([]byte, error) {
	contractFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return contractFile, nil
}

// TODO must replace with WS Connection, P2P, or sync tx execution to prevent polling
func pollForReceipt(channelId string, transactionId string, client mazzaroth.Client) (*xdr.Receipt, error) {
	retry := 0
	for {
		receipt, err := client.ReceiptLookup(context.Background(), channelId, transactionId)
		if err != nil {
			if retry != maxRetry {
				retry++
				time.Sleep(time.Second * 5)
				continue
			}
			return nil, err
		}
		return receipt, nil
	}
}

func FromFile(path string, manifestType string) ([]*Manifest, error) {
	manifestFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifests := make([]*Manifest, 0, 0)
	r := bytes.NewReader(manifestFile)

	decoder := yaml.NewDecoder(r)
	for {
		manifest := &Manifest{}
		if err := decoder.Decode(manifest); err != nil {
			if err != io.EOF {
				return nil, os.ErrClosed
			}
			break
		}
		if manifest.Type == manifestType {
			manifests = append(manifests, manifest)
		}
	}
	return manifests, nil
}

func ExecuteDeployments(ctx context.Context, manifests []*Manifest, client mazzaroth.Client, sender string, privKey ed25519.PrivateKey) error {
	for _, m := range manifests {
		if m.Type != "deployment" {
			continue
		}

		senderId, err := xdr.IDFromHexString(sender)
		if err != nil {
			return err
		}

		channelId, err := xdr.IDFromHexString(m.Channel.Id)
		if err != nil {
			return err
		}

		abi, err := loadAbi(m.Channel.AbiFile)
		if err != nil {
			return err
		}

		contract, err := loadContract(m.Channel.ContractFile)
		if err != nil {
			return err
		}

		tx, err := mazzaroth.Transaction(&senderId, &channelId).
			Contract(0, 0).Abi(abi).ContractBytes(contract).Version(m.Channel.Version).Sign(privKey)
		if err != nil {
			return err
		}

		id, _, err := client.TransactionSubmit(ctx, tx)
		if err != nil {
			return err
		}

		fmt.Println("Contract Deployed:transaction id:", id)
		receipt, err := pollForReceipt(m.Channel.Id, fmt.Sprintf("%b", id), client)
		if err != nil {
			return err
		}
		receiptJson, err := json.MarshalIndent(receipt, "", "\t")
		fmt.Println("Contract Deployment Complete:receipt:", string(receiptJson))

		for _, t := range m.Transactions {
			args := make([]xdr.Argument, 0, 0)
			if len(t.Tx.Args) > 0 {
				for _, a := range t.Tx.Args {
					args = append(args, xdr.Argument(a))
				}
			}
			tx, err := mazzaroth.Transaction(&senderId, &channelId).
				Call(0, 0).Function(t.Tx.Function).Arguments(args...).Sign(privKey)
			if err != nil {
				return err
			}
			id, _, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}
			fmt.Println("transaction submitted:id:", fmt.Sprintf("%s", id))
			receipt, err := pollForReceipt(m.Channel.Id, fmt.Sprintf("%s", id), client)
			if err != nil {
				return err
			}
			receiptJson, err := json.MarshalIndent(receipt, "", "\t")
			fmt.Println("transaction complete:receipt:", string(receiptJson))
		}
	}
	return nil
}

func ExecuteTests(ctx context.Context, manifests []*Manifest, client mazzaroth.Client, sender string, privKey ed25519.PrivateKey) error {
	for _, m := range manifests {
		if m.Type != "test" {
			continue
		}

		senderId, err := xdr.IDFromHexString(sender)
		if err != nil {
			return err
		}

		channelId, err := xdr.IDFromHexString(m.Channel.Id)
		if err != nil {
			return err
		}

		for _, t := range m.Transactions {
			args := make([]xdr.Argument, 0, 0)
			if len(t.Tx.Args) > 0 {
				for _, a := range t.Tx.Args {
					args = append(args, xdr.Argument(a))
				}
			}
			tx, err := mazzaroth.Transaction(&senderId, &channelId).
				Call(0, 0).Function(t.Tx.Function).Arguments(args...).Sign(privKey)
			if err != nil {
				return err
			}
			id, _, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}
			fmt.Println("transaction submitted:id:", fmt.Sprintf("%s", id))
			receipt, err := pollForReceipt(m.Channel.Id, fmt.Sprintf("%s", id), client)
			if err != nil {
				return err
			}
			receiptJson, err := json.MarshalIndent(receipt, "", "\t")
			fmt.Println("transaction complete:receipt:", string(receiptJson))
		}
	}
	return nil
}
