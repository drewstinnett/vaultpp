package vaultpp

import (
	"errors"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultPP struct {
	Client *api.Client
}

func NewVaultPP(addr, token string) (*VaultPP, error) {
	if addr == "" {
		if os.Getenv("VAULT_ADDR") == "" {
			return nil, errors.New("Must set VAULT_ADDR, or provide a vault address")
		} else {
			addr = os.Getenv("VAULT_ADDR")
		}
	}

	if token == "" {
		if os.Getenv("VAULT_TOKEN") == "" {
			return nil, errors.New("Must set VAULT_TOKEN, or provide a vault token")
		} else {
			token = os.Getenv("VAULT_TOKEN")
		}
	}
	config := &api.Config{
		Address: addr,
	}
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	c.SetToken(token)
	return &VaultPP{
		Client: c,
	}, nil
}
