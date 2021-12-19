package vaultpp

import (
	"errors"
	"fmt"
	pathx "path"
	"path/filepath"
	"strings"

	"github.com/apex/log"
)

type KVExport struct {
	Mount   string           `yaml:"mount"`
	Vault   string           `yaml:"vault"`
	Secrets []KVExportSecret `yaml:"secrets"`
}

type KVExportSecret struct {
	DataPath string                 `yaml:"data_path"`
	Path     string                 `yaml:"path"`
	Data     map[string]interface{} `yaml:"data"`
}

func (v *VaultPP) WalkTree(filePath string) ([]KVExportSecret, error) {
	mounts, err := v.GetKVMounts()
	if err != nil {
		return nil, err
	}
	m, err := MatchKVMountInfo(mounts, filePath)
	if err != nil {
		return nil, err
	}
	items, err := v.WalkTreeWithMount(m, strings.TrimPrefix(filePath, m.Path))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (v *VaultPP) WalkTreeWithMount(mount *KVMountInfo, path string) ([]KVExportSecret, error) {
	var items []KVExportSecret
	keys, err := v.ListKeys(mount, path)
	if err != nil {
		log.WithError(err).WithField("path", path).Error("failed to list keys")
		return nil, errors.New("Issue walking dir")
	}

	for _, k := range keys {
		if strings.HasSuffix(k, "/") {
			newItems, err := v.WalkTreeWithMount(mount, fmt.Sprintf("%s%s", path, k))
			if err != nil {
				return nil, err
			}
			items = append(items, newItems...)
		} else {
			// items = append(items, fmt.Sprintf("%s%s", path, k))
			items = append(items, KVExportSecret{
				Path:     pathx.Join(mount.Path, path, k),
				DataPath: pathx.Join(mount.Path, "data/", fmt.Sprintf("%s%s", path, k)),
			})
		}
	}
	return items, nil
}

func (v *VaultPP) ListKeys(mount *KVMountInfo, path string) ([]string, error) {
	var fullPath string
	if mount.Version > 1 {
		fullPath = filepath.Join(mount.Path, "metadata", path)
	} else {
		fullPath = filepath.Join(mount.Path, path)
	}

	data, err := v.Client.Logical().List(fullPath)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets found")
	}

	keys := []string{}

	if data.Data != nil {
		for _, k := range data.Data["keys"].([]interface{}) {
			keys = append(keys, k.(string))
		}
	}

	return keys, nil
}
