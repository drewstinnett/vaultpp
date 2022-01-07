package unsealers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	hvault "github.com/hashicorp/vault/api"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type OPUnsealer struct{}

func (o *OPUnsealer) Prerequisites() error {
	_, err := exec.LookPath("op")
	if err != nil {
		return errors.New("No op binary found")
	}
	if os.Getenv("OP_UNSEAL_PATH") == "" {
		return errors.New("No OP_UNSEAL_PATH environment variable set")
	}
	opToken := os.Getenv("OP_SESSION_my")
	if opToken == "" {
		return errors.New("No op token found, be sure to set OP_SESSION_my")
	}
	_, err = exec.Command("op", "get", "account").Output()
	if err != nil {
		return errors.New("Could not do 'op get account', is your token still valid?")
	}
	return nil
}

func (o *OPUnsealer) FetchUnsealData(params map[string]interface{}) (*UnsealData, error) {
	if params == nil {
		return nil, errors.New("No params provided")
	}
	if params["path"] == nil {
		return nil, errors.New("No path provided")
	}
	out, err := exec.Command("op", "get", "document", params["path"].(string)).Output()
	if err != nil {
		return nil, err
	}
	unsealData := &UnsealData{}
	err = yaml.Unmarshal(out, &unsealData)
	if err != nil {
		return nil, err
	}
	return unsealData, nil
}

func (o *OPUnsealer) Unseal(ud UnsealData) error {
	for _, vaultName := range ud.Vaults {

		cf := hvault.Config{
			Address: fmt.Sprintf("https://%v:8200", vaultName),
		}
		client, err := hvault.NewClient(&cf)
		if err != nil {
			return err
		}
		sealStatus, err := client.Sys().SealStatus()
		if err != nil {
			return err
		}
		if sealStatus.Sealed {
			// Reset Unseal status
			_, err = client.Sys().ResetUnsealProcess()
			if err != nil {
				return err
			}
			for x, shard := range ud.Keys {
				log.WithFields(log.Fields{
					"vault": vaultName,
					"shard": x,
				}).Info("Unsealing")
				status, err := client.Sys().Unseal(shard)
				if err != nil {
					return err
				}
				if !status.Sealed {
					return nil
				}
			}

			log.WithFields(log.Fields{"instance": vaultName}).Warn("Unseal failed")
			return errors.New("UnsealFailed")
		} else {
			log.WithFields(log.Fields{"instance": vaultName}).Info("Vault is already unsealed")
		}
	}
	return nil
}
