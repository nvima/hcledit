#/bin/bash

destination='.secrets'
contents='{{ with secret "secret/my-secret" }}{{ .Data.data.foo }}{{ end }}'
vaultFile="vault-agent.hcl"
vaultFileCat=$(cat $vaultFile 2>&1)
if [ $? -eq 0 ]; then
  output=$(cat $vaultFile | go run main.go templateblock upsert -d "$destination" -c "$contents" 2>&1)
  if [ $? -eq 0 ]; then
    echo "$output" > $vaultFile
  else
    echo "hcledit failed"
    echo "$output"
    exit 3
  fi
else
  echo "cat failed"
  echo "$vaultFileCat"
  exit 3
fi

