package vaultpp_test

import (
	"testing"

	"github.com/drewstinnett/vaultpp/pkg/vaultpp"
	"github.com/stretchr/testify/require"
)

func TestMatchKVMountInfo(t *testing.T) {
	tt := []struct {
		mounts     []vaultpp.KVMountInfo
		match      string
		matchMount string
		shouldErr  bool
	}{
		{
			mounts: []vaultpp.KVMountInfo{
				{Path: "kv/"},
				{Path: "another/secret/"},
				{Path: "secret/"},
			},
			match:      "secret/foo/bar",
			matchMount: "secret/",
		},
		{
			mounts: []vaultpp.KVMountInfo{
				{Version: 2, Path: "kv/"},
				{Version: 1, Path: "secret/another/secret/"},
				{Version: 2, Path: "secret/"},
			},
			match:      "secret/foo/bar",
			matchMount: "secret/",
		},
		{
			mounts: []vaultpp.KVMountInfo{
				{Version: 2, Path: "kv/foo/"},
				{Version: 2, Path: "foo/kv/foo/"},
				{Version: 2, Path: "kv/"},
			},
			match:      "kv/baz/bar",
			matchMount: "kv/",
		},
		{
			mounts: []vaultpp.KVMountInfo{
				{Version: 2, Path: "kv/foo/"},
				{Version: 2, Path: "foo/kv/foo/"},
				{Version: 2, Path: "kv/"},
			},
			match:     "never-exist/",
			shouldErr: true,
		},
	}
	for _, tt := range tt {
		info, err := vaultpp.MatchKVMountInfo(tt.mounts, tt.match)
		if tt.shouldErr {
			require.Error(t, err)
			require.Nil(t, info)
		} else {
			require.Equal(t, tt.matchMount, info.Path)
		}
	}
}
