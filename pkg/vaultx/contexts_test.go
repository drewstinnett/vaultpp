package vaultx_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drewstinnett/vaultx/pkg/vaultx"
)

func TestInitContextFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		log.Fatal(err)
	}
	cfile := filepath.Join(tmpdir, "subdir", "contexts.yaml")
	err = vaultx.InitContextFile(cfile)
	require.NoError(t, err)
}

func TestContextFreshen(t *testing.T) {
	ts := []struct {
		env       map[string]string
		tokenHash string
		shouldErr bool
	}{
		{
			env: map[string]string{
				"VAULT_ADDR":  "http//localhost:8200",
				"VAULT_TOKEN": "foo",
			},
			tokenHash: "683716d9d7f82eed174c6caebe086ee93376c79d7c61dd670ea00f7f8d6eb0a8",
		},
	}
	for _, tc := range ts {
		os.Clearenv()
		for k, v := range tc.env {
			os.Setenv(k, v)
		}
		c, err := vaultx.NewContextWithEnv()
		c.Freshen()

		if tc.shouldErr {
			require.Error(t, err)
		} else {
			require.Equal(t, tc.tokenHash, c.TokenHash)
		}

	}
}

func TestNewContextWithEnv(t *testing.T) {
	ts := []struct {
		env       map[string]string
		shouldErr bool
	}{
		{
			env: map[string]string{
				"VAULT_ADDR":  "http//localhost:8200",
				"VAULT_TOKEN": "foo",
			},
		},
		{
			env:       map[string]string{},
			shouldErr: true,
		},
		{env: map[string]string{"VAULT_ADDR": "foo", "VAULT_TOKEN": ""}, shouldErr: true},
		{env: map[string]string{"VAULT_ADDR": "", "VAULT_TOKEN": "foo"}, shouldErr: true},
		{env: map[string]string{"VAULT_ADDR": "foo", "VAULT_TOKEN": "foo", "VAULT_NAMESPACE": "test"}},
	}
	for _, tc := range ts {
		os.Clearenv()
		for k, v := range tc.env {
			os.Setenv(k, v)
		}
		c, err := vaultx.NewContextWithEnv()

		if tc.shouldErr {
			require.Error(t, err)
		}
		if c != nil {
			if _, ok := tc.env["VAULT_ADDR"]; ok {
				require.Equal(t, tc.env["VAULT_ADDR"], c.Address)
			}
			if _, ok := tc.env["VAULT_TOKEN"]; ok {
				require.Equal(t, tc.env["VAULT_TOKEN"], c.Token)
			}
			if _, ok := tc.env["VAULT_NAMESPACE"]; !ok {
				require.Equal(t, "root", c.Namespace)
			} else {
				require.Equal(t, tc.env["VAULT_NAMESPACE"], c.Namespace)
			}
		}
	}
}
