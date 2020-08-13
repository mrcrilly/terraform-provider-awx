terraform {
  required_providers {
    awx = {
      versions = [
        "0.1"]
      source = "github.com/mrcrilly/awx"
    }
  }
}

provider "awx" {}

resource "awx_credential_ssh" "main" {
  name = "Mikes SSH Creds"
  organisation_id = 1
}
