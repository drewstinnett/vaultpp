package vaultx_test

import (
	"testing"

	"github.com/drewstinnett/vaultx/pkg/vaultx"
	"github.com/stretchr/testify/require"
)

func TestMatchKVMountInfo(t *testing.T) {
	tt := []struct {
		mounts     []vaultx.KVMountInfo
		match      string
		matchMount string
		shouldErr  bool
	}{
		{
			mounts: []vaultx.KVMountInfo{
				{Path: "kv/"},
				{Path: "another/secret/"},
				{Path: "secret/"},
			},
			match:      "secret/foo/bar",
			matchMount: "secret/",
		},
		{
			mounts: []vaultx.KVMountInfo{
				{Version: 2, Path: "kv/"},
				{Version: 1, Path: "secret/another/secret/"},
				{Version: 2, Path: "secret/"},
			},
			match:      "secret/foo/bar",
			matchMount: "secret/",
		},
		{
			mounts: []vaultx.KVMountInfo{
				{Version: 2, Path: "kv/foo/"},
				{Version: 2, Path: "foo/kv/foo/"},
				{Version: 2, Path: "kv/"},
			},
			match:      "kv/baz/bar",
			matchMount: "kv/",
		},
		{
			mounts: []vaultx.KVMountInfo{
				{Version: 2, Path: "kv/foo/"},
				{Version: 2, Path: "foo/kv/foo/"},
				{Version: 2, Path: "kv/"},
			},
			match:     "never-exist/",
			shouldErr: true,
		},
	}
	for _, tt := range tt {
		info, err := vaultx.MatchKVMountInfo(tt.mounts, tt.match)
		if tt.shouldErr {
			require.Error(t, err)
			require.Nil(t, info)
		} else {
			require.Equal(t, tt.matchMount, info.Path)
		}
	}
}
