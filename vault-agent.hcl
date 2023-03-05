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
