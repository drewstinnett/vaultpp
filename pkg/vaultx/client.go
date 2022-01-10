package vaultx

import (
	"github.com/hashicorp/vault/api"
)

type VaultPP struct {
	Client      *api.Client
	ContextFile string
}

func NewVaultPP(vctx *Context) (*VaultPP, error) {
	config := &api.Config{
		Address: vctx.Address,
	}
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	c.SetToken(vctx.Token)
	return &VaultPP{
		Client:      c,
		ContextFile: "~/.vaultx/contexts.yaml",
	}, nil
}
