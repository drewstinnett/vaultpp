package vaultx

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/apex/log"
)

func (v *VaultPP) GetMountPaths() ([]string, error) {
	var items []string
	mounts, err := v.Client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}
	for _, mount := range mounts {
		items = append(items, mount.Options["path"])
	}
	return items, nil
}

type KVMountInfo struct {
	Path    string `json:"path"`
	Version int    `json:"version"`
}

func MatchKVMountInfo(mounts []KVMountInfo, path string) (*KVMountInfo, error) {
	// Sort by path length, shortest first
	sort.Slice(mounts, func(i, j int) bool {
		return len(mounts[i].Path) > len(mounts[j].Path)
	})
	for _, mount := range mounts {
		if strings.HasPrefix(path, mount.Path) {
			return &mount, nil
		}
	}
	return nil, errors.New("No matching mount found")
}

func (v *VaultPP) DetectKVVersion(filePath string) (int, error) {
	mounts, err := v.GetKVMounts()
	if err != nil {
		return 0, err
	}
	m, err := MatchKVMountInfo(mounts, filePath)
	if err != nil {
		return 0, err
	}
	return m.Version, nil
}

func (v *VaultPP) GetKVMounts() ([]KVMountInfo, error) {
	items := []KVMountInfo{}
	info, err := v.Client.Logical().Read("sys/mounts")
	if err != nil {
		return nil, err
	}
	ctx := log.WithFields(log.Fields{
		"func": "GetKVMounts",
	})
	for mountPath, d := range info.Data {
		mountType := d.(map[string]interface{})["type"].(string)
		ctx.WithFields(log.Fields{
			"type":      mountType,
			"mountPath": mountPath,
		}).Debugf("Examining mount")
		if mountType != "kv" {
			continue
		}
		log.Debugf("%v", d.(map[string]interface{}))
		options := d.(map[string]interface{})["options"]
		ctx.WithFields(log.Fields{
			"options": options,
		}).Debugf("Mount Options")

		var version string
		if options != nil {
			version = options.(map[string]interface{})["version"].(string)
		} else {
			version = "1"
		}
		versionI, err := strconv.Atoi(version)
		if err != nil {
			return nil, err
		}
		mi := KVMountInfo{
			Path:    mountPath,
			Version: versionI,
		}
		items = append(items, mi)
	}
	return items, nil
}
