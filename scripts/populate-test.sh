#!/bin/bash

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=$(cat ~/.vault-token)

ST=$(vault status -format=json | jq -r .storage_type)
if [ $ST != "inmem" ]; then
    echo "Only do these tests on inmem stuff" 1>&2
    exit 2
fi

vault kv put secret/foo/bar foo=bar
vault kv put secret/bar/bar foo=bar
vault kv put secret/thiing/var/bar foo=bar
vault kv put secret/usr/bin/bar foo=bar
