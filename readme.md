# README

## Introduction
This is a Go command-line application to edit a template block from stdin. It searches for a template block in stdin with the same 'destination' attribute and edits it. If no template block with the same 'destination' attribute is found, a new one is created. 

## Usage
To use this application, run the following command:

`go run main.go templateblock upsert -d <destination> [-s <source>|-c <contents>]`

* `-d` or `--destination`: Required flag. It specifies the destination of the template block.
* `-s` or `--source`: Required if `-c` is not provided. It specifies the source of the template block.
* `-c` or `--contents`: Required if `-s` is not provided. It specifies the contents of the template block.


## Example
- vault-agent.hcl
```
template {
    destination = ".secrets"
    source      = ".secrets.ctmpl"
}
```

- edit vault-agent.hcl
```
cat vault-agent.hcl | edithcl terraformblock upsert -d ".secrets" -c "{{ .Data }} > vault-agent.hcl
```

- vault-agent.hcl 
```
template {
    destination = ".secrets"
    content     = "{{ .Data }}"
}
```

- add block to vault-agent.hcl
```
cat vault-agent.hcl | edithcl terraformblock upsert -d ".octopus" -s ".octopus.ctmpl" > vault-agent.hcl
```

- vault-agent.hcl 
```
template {
    destination = ".secrets"
    content     = "{{ .Data }}"
}

template {
    destination = ".octopus"
    source      = ".octopus.ctmpl"
}
```

## Additional Information
* If you edit a template block with a destination and contents attribute, the source attribute will be removed, and vice versa. This is because the vault agent only supports one of the two attributes. [More info here](https://developer.hashicorp.com/vault/docs/agent/template#templating-configuration-example).
