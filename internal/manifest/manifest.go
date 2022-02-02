package manifest

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pterm/pterm"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"gopkg.in/yaml.v2"
)

const (
	maxRetry                = 10
	maxBlockExpirationRange = 100
)

type Manifest struct {
	Version     string      `yaml:"version"`
	Type        string      `yaml:"type"`
	Channel     Channel     `yaml:"channel"`
	GatewayNode GatewayNode `yaml:"gateway-node"`
	Deploy      *Deploy     `yaml:"deploy"`
	Tests       []*Test     `yaml:"tests"`
}

type Deploy struct {
	Name         string `yaml:"name"`
	Transactions []*Tx  `yaml:"transactions,omitempty"`
}

type Test struct {
	Name         string `yaml:"name"`
	Reset        bool   `yaml:"reset"`
	Transactions []*Tx  `yaml:"transactions,omitempty"`
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
				time.Sleep(time.Second * time.Duration(retry))
				retry++
				continue
			}
			// return first error
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

		if m.Deploy == nil {
			return errors.New("missing deploy block for manifest")
		}

		senderId, err := xdr.IDFromHexString(sender)
		if err != nil {
			return err
		}

		channelId, err := xdr.IDFromHexString(m.Channel.Id)
		if err != nil {
			return err
		}

		owner, err := xdr.IDFromHexString(m.Channel.Owner)
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

		tx, err := mazzaroth.Transaction(senderId, channelId).
			Contract(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Deploy(owner, m.Channel.Version, abi, contract).Sign(privKey)
		if err != nil {
			return err
		}

		spinnerSuccess, err := pterm.DefaultSpinner.Start("deploying contract...")
		if err != nil {
			return err
		}

		id, _, err := client.TransactionSubmit(ctx, tx)
		if err != nil {
			spinnerSuccess.Fail()
			return err
		}
		spinnerSuccess.Success("contract deploy transaction submitted with tx id: ", hex.EncodeToString(id[:]))

		receipt, err := pollForReceipt(m.Channel.Id, hex.EncodeToString(id[:]), client)
		if err != nil {
			return err
		}

		receiptJson, err := json.MarshalIndent(receipt, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println("contract deployment complete:receipt:\n", string(receiptJson))

		for _, t := range m.Deploy.Transactions {
			args := make([]xdr.Argument, 0, 0)
			if len(t.Tx.Args) > 0 {
				for _, a := range t.Tx.Args {
					args = append(args, xdr.Argument(a))
				}
			}

			tx, err := mazzaroth.Transaction(senderId, channelId).
				Call(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Function(t.Tx.Function).Arguments(args...).Sign(privKey)
			if err != nil {
				return err
			}

			id, receipt, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}

			fmt.Println("transaction submitted:id:", hex.EncodeToString(id[:]))
			if receipt == nil {
				receipt, err = pollForReceipt(m.Channel.Id, hex.EncodeToString(id[:]), client)
				if err != nil {
					return err
				}
			}

			receiptJson, err := json.MarshalIndent(receipt, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println("transaction complete:receipt:\n", string(receiptJson))
		}
	}
	return nil
}

func ExecuteTests(ctx context.Context, manifests []*Manifest, client mazzaroth.Client, sender string, privKey ed25519.PrivateKey) error {
	for _, m := range manifests {
		if m.Type != "test" {
			continue
		}

		if m.Tests == nil {
			return errors.New("missing tests for test manifest")
		}

		senderId, err := xdr.IDFromHexString(sender)
		if err != nil {
			return err
		}

		channelId, err := xdr.IDFromHexString(m.Channel.Id)
		if err != nil {
			return err
		}

		owner, err := xdr.IDFromHexString(m.Channel.Owner)
		if err != nil {
			return err
		}

		for _, t := range m.Tests {
			if t.Reset {
				tx, err := mazzaroth.Transaction(senderId, channelId).
					Contract(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Delete().Sign(privKey)
				if err != nil {
					return err
				}

				id, _, err := client.TransactionSubmit(ctx, tx)
				if err != nil {
					return err
				}

				fmt.Println("contract delete:transaction id:", hex.EncodeToString(id[:]))
				receipt, err := pollForReceipt(m.Channel.Id, hex.EncodeToString(id[:]), client)
				if err != nil {
					return err
				}
				receiptJson, err := json.MarshalIndent(receipt, "", "\t")
				if err != nil {
					return err
				}
				fmt.Println("contract delete complete:receipt:\n", string(receiptJson))
			}
			abi, err := loadAbi(m.Channel.AbiFile)
			if err != nil {
				return err
			}

			contract, err := loadContract(m.Channel.ContractFile)
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(senderId, channelId).
				Contract(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Deploy(owner, m.Channel.Version, abi, contract).Sign(privKey)
			if err != nil {
				return err
			}

			id, _, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}

			fmt.Println("contract deployed:transaction id:", hex.EncodeToString(id[:]))
			receipt, err := pollForReceipt(m.Channel.Id, hex.EncodeToString(id[:]), client)
			if err != nil {
				return err
			}

			receiptJson, err := json.MarshalIndent(receipt, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println("contract deployment complete:receipt:\n", string(receiptJson))

			for _, t := range t.Transactions {
				args := make([]xdr.Argument, 0, 0)
				if len(t.Tx.Args) > 0 {
					for _, a := range t.Tx.Args {
						args = append(args, xdr.Argument(a))
					}
				}

				tx, err := mazzaroth.Transaction(senderId, channelId).
					Call(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Function(t.Tx.Function).Arguments(args...).Sign(privKey)
				if err != nil {
					return err
				}

				id, receipt, err := client.TransactionSubmit(ctx, tx)
				if err != nil {
					return err
				}

				fmt.Println("transaction submitted:id:", hex.EncodeToString(id[:]))
				if receipt == nil {
					receipt, err = pollForReceipt(m.Channel.Id, hex.EncodeToString(id[:]), client)
					if err != nil {
						return err
					}
				}

				receiptJson, err := json.MarshalIndent(receipt, "", "\t")
				if err != nil {
					return err
				}
				fmt.Println("transaction complete:receipt: \n", string(receiptJson))
				if t.Tx.Receipt != nil {
					if receipt.Status != xdr.Status(t.Tx.Receipt.Status) {
						return fmt.Errorf("expected transaction status : %d does not match %d", t.Tx.Receipt.Status, receipt.Status)
					}
					if receipt.Result != t.Tx.Receipt.Result {
						return fmt.Errorf("expected transaction results : %s does not match %s", t.Tx.Receipt.Result, receipt.Result)
					}
				}
			}
		}
	}
	return nil
}
