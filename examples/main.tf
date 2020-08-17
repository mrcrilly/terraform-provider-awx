terraform {
  required_providers {
    awx = {
      versions = ["0.1"]
      source = "github.com/mrcrilly/awx"
    }
  }
}

provider "awx" {}

resource "awx_credential_ssh" "main" {
  name = "Mikes SSH Creds"
  organisation_id = 1
}

resource "awx_credential_azure_key_vault" "main" {
  name = "Mikes Azure KV Creds"
  organisation_id = 1
  url = "https://something.com"
  client = "1234"
  secret = "my new secret string"
  tenant = "mytenantid"
}