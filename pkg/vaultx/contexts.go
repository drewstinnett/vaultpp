package vaultx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Context struct {
	Name      string     `yaml:"name"`
	Address   string     `yaml:"addr"`
	Namespace string     `yaml:"namespace"`
	Token     string     `yaml:"token"`
	TokenHash string     `yaml:"token_hash"`
	Expires   *time.Time `yaml:"expires"`
}

type ContextConfig struct {
	Current  string    `yaml:"current"`
	Contexts []Context `yaml:"contexts"`
}

func InitContextFile(filename string) error {
	filename, err := homedir.Expand(filename)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		parent := filepath.Dir(filename)
		if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {
			log.WithFields(log.Fields{
				"directory": parent,
			}).Debug("Creating parent dir for context")
			err = os.MkdirAll(parent, 0o700)
			if err != nil {
				return err
			}
		}

		log.WithFields(log.Fields{
			"filename": filename,
		}).Debug("Creating context file")
		ctxConfig := &ContextConfig{}
		d, err := yaml.Marshal(ctxConfig)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, d, 0o600)
		if err != nil {
			return err
		}

	}
	return nil
}

func SaveContext(cfile string, c *Context) error {
	if cfile == "" {
		cfile, _ = homedir.Expand("~/.vaultx/contexts.yaml")
	}
	existingConfig, err := ReadContextFile(cfile)
	newContexts := []Context{}
	if err != nil {
		return err
	}
	var overwrite bool
	for _, ctx := range existingConfig.Contexts {
		if ctx.Name == c.Name {
			newContexts = append(newContexts, *c)
			overwrite = true
		} else {
			newContexts = append(newContexts, ctx)
		}
	}
	if !overwrite {
		newContexts = append(newContexts, *c)
	}
	existingConfig.Contexts = newContexts
	d, err := yaml.Marshal(existingConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfile, d, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func ReadContextFile(cfile string) (*ContextConfig, error) {
	if cfile == "" {
		defaultPath := "~/.vaultx/contexts.yaml"
		log.WithFields(log.Fields{
			"filename": defaultPath,
		}).Debug("No context file specified, using default")
		cfile = defaultPath
	}
	cfile, err := homedir.Expand(cfile)
	if err != nil {
		return nil, err
	}
	// Init file if it doesn't exist
	err = InitContextFile(cfile)
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadFile(cfile)
	if err != nil {
		return nil, err
	}
	ctxConfig := &ContextConfig{}
	err = yaml.Unmarshal(d, ctxConfig)
	if err != nil {
		return nil, err
	}
	return ctxConfig, nil
}

func ListContexts(cfile string) ([]Context, error) {
	ctxConfig, err := ReadContextFile(cfile)
	if err != nil {
		return nil, err
	}
	envCtx, err := NewContextWithEnv()
	if err == nil {
		ctxConfig.Contexts = append(ctxConfig.Contexts, *envCtx)
	}
	return ctxConfig.Contexts, nil
}

func GetCurrentContext(cfile string) (*Context, error) {
	if cfile == "" {
		defaultPath := "~/.vaultx/contexts.yaml"
		log.WithFields(log.Fields{
			"filename": defaultPath,
		}).Debug("No context file specified, using default")
		cfile = defaultPath
	}
	ctxConfig, err := ReadContextFile(cfile)
	if err != nil {
		return nil, err
	}
	if ctxConfig.Current == "" {
		log.WithFields(log.Fields{
			"filename": cfile,
		}).Debug("Current context is blank, using env")
		ctx, err := NewContextWithEnv()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	} else if ctxConfig.Current == "environment" {
		log.WithFields(log.Fields{
			"filename": cfile,
		}).Debug("Current context is 'environment', using env")
		ctx, err := NewContextWithEnv()
		if err != nil {
			return nil, err
		}
		return ctx, nil

	}
	for _, c := range ctxConfig.Contexts {
		if c.Name == ctxConfig.Current {
			return &c, nil
		}
	}
	return nil, errors.New("No current context")
}

func (c *Context) Freshen() error {
	h := hmac.New(sha256.New, []byte(c.Token))
	sha := hex.EncodeToString(h.Sum(nil))
	c.TokenHash = sha
	return nil
}

func NewContextWithEnv() (*Context, error) {
	c := &Context{
		Name: "environment",
	}

	// Make sure addr is set. Do we wanna default to http://localhost:8200 at some point?
	if os.Getenv("VAULT_ADDR") == "" {
		return nil, errors.New("Must set VAULT_ADDR")
	} else {
		c.Address = os.Getenv("VAULT_ADDR")
	}

	// Need that token
	if os.Getenv("VAULT_TOKEN") == "" {
		return nil, errors.New("Must set VAULT_TOKEN")
	} else {
		c.Token = os.Getenv("VAULT_TOKEN")
	}

	// Default to root namespace
	if os.Getenv("VAULT_NAMESPACE") == "" {
		c.Namespace = "root"
	} else {
		c.Namespace = os.Getenv("VAULT_NAMESPACE")
	}
	return c, nil
}
