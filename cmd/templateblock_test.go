package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTemplateBlockUpsertCmd(t *testing.T) {
	// Test case where input contains a template block with the same destination
	t.Run("add new block without stdin", func(t *testing.T) {
		input := ``
		expectedOutput := `template {
  destination = "path/to/destination"
  contents    = "new contents"
}
`
		cmd := templateBlockUpsertCmd()
		cmd.SetArgs([]string{"upsert", "--destination", "path/to/destination", "--contents", "new contents"})
		cmd.SetIn(strings.NewReader(input))
		output := bytes.Buffer{}
		cmd.SetOut(&output)
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output.String())
	})
	// Test case where input contains a template block with the same destination
	t.Run("edit templateblock with stdin", func(t *testing.T) {
		input := `
auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      type = "iam"
      role = "dev-role-iam"
    }
  }
  sink "file" {
    config = {
      path = "./agent-token"
    }
  }
}

vault {
  address = "https://172.16.148.3:8200"
}

pid_file = "./vault.pid"

template {
  source      = "/etc/vault.d/secret.ctmpl"
  destination = "/opt/vault/secret"
}

template {
  destination = ".secrets"
  source      = ".secrets.ctmpl"
}
`
		expectedOutput := `
auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      type = "iam"
      role = "dev-role-iam"
    }
  }
  sink "file" {
    config = {
      path = "./agent-token"
    }
  }
}

vault {
  address = "https://172.16.148.3:8200"
}

pid_file = "./vault.pid"

template {
  source      = "/etc/vault.d/secret.ctmpl"
  destination = "/opt/vault/secret"
}

template {
  destination = ".secrets"
  contents    = "{{ with secret \"secret/my-secret\" }}{{ .Data.data.foo }}{{ end }}"
}
`
		cmd := templateBlockUpsertCmd()
		cmd.SetArgs([]string{"upsert", "--destination", ".secrets", "--contents", "{{ with secret \"secret/my-secret\" }}{{ .Data.data.foo }}{{ end }}"})
		cmd.SetIn(strings.NewReader(input))
		output := bytes.Buffer{}
		cmd.SetOut(&output)
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output.String())
	})

	t.Run("add templateblock with stdin", func(t *testing.T) {
		input := `
auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      type = "iam"
      role = "dev-role-iam"
    }
  }
  sink "file" {
    config = {
      path = "./agent-token"
    }
  }
}

vault {
  address = "https://172.16.148.3:8200"
}

pid_file = "./vault.pid"

template {
  source      = "/etc/vault.d/secret.ctmpl"
  destination = "/opt/vault/secret"
}

template {
  source      = "/etc/.octopus.ctmpl"
  destination = "/etc/.octopus"
}
`
		expectedOutput := `
auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      type = "iam"
      role = "dev-role-iam"
    }
  }
  sink "file" {
    config = {
      path = "./agent-token"
    }
  }
}

vault {
  address = "https://172.16.148.3:8200"
}

pid_file = "./vault.pid"

template {
  source      = "/etc/vault.d/secret.ctmpl"
  destination = "/opt/vault/secret"
}

template {
  source      = "/etc/.octopus.ctmpl"
  destination = "/etc/.octopus"
}

template {
  destination = ".secrets"
  contents    = "{{ with secret \"secret/my-secret\" }}{{ .Data.data.foo }}{{ end }}"
}
`
		cmd := templateBlockUpsertCmd()
		cmd.SetArgs([]string{"upsert", "--destination", ".secrets", "--contents", "{{ with secret \"secret/my-secret\" }}{{ .Data.data.foo }}{{ end }}"})
		cmd.SetIn(strings.NewReader(input))
		output := bytes.Buffer{}
		cmd.SetOut(&output)
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output.String())
	})

	// Test case where required flag 'destination' is not provided
	t.Run("missing destination flag", func(t *testing.T) {
		cmd := templateBlockUpsertCmd()
		cmd.SetArgs([]string{"upsert", "--source", "new contents"})
		cmd.SetIn(strings.NewReader(""))
		output := bytes.Buffer{}
		cmd.SetOut(&output)
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "destination is required")
	})

	// Test case where required flag 'destination' is not provided
	t.Run("missing source or contents flag", func(t *testing.T) {
		cmd := templateBlockUpsertCmd()
		cmd.SetArgs([]string{"upsert", "--destination", ".secrets"})
		cmd.SetIn(strings.NewReader(""))
		output := bytes.Buffer{}
		cmd.SetOut(&output)
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "either source or contents is required")
	})

}
