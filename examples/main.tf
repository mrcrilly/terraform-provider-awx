terraform {
  required_providers {
    awx = {
      source = "github.com/mrcrilly/awx"
      version = "0.1"
    }
  }
}

provider "awx" {}
provider "aws" {
  region = "ap-southeast-2"
}

provider "azurerm" {
  features {}
}

resource "aws_instance" "windows" {
  ami = "ami-0343f4bc3607550d2"
  instance_type = "t2.small"
  get_password_data = true
  key_name = "<example>"

  tags = {
    Name = "deleteme"
  }
}

data "azurerm_key_vault" "example" {
  name = "<example>"
  resource_group_name = "<example>"
}

data "azurerm_key_vault_secret" "example" {
  name = "windows-server-pem-file"
  key_vault_id = data.azurerm_key_vault.example.id
}

locals {
  decrypted_password = rsadecrypt(aws_instance.windows.password_data, base64decode(data.azurerm_key_vault_secret.example.value))
}

resource "azurerm_key_vault_secret" "example" {
  name = "deleteme-windows-server"
  value = local.decrypted_password
  key_vault_id = data.azurerm_key_vault.example.id
}

resource "awx_credential_machine" "windows" {
  organisation_id = 1
  name = "Windows Server 2018"
  username = "administrator"
}

resource "awx_credential_input_source" "azure-to-windows" {
  description = "link azure key vault secret to windows server"
  input_field_name = "password"
  target = awx_credential_machine.windows.id
  source = awx_credential_azure_key_vault.windows.id
  metadata = {
    secret_field = azurerm_key_vault_secret.example.name
  }
}

resource "awx_credential_azure_key_vault" "windows" {
  name = "Primary Key Vault"
  organisation_id = 1
  url = data.azurerm_key_vault.example.vault_uri
  client = "<example>"
  secret = "<example>"
  tenant = data.azurerm_key_vault.example.tenant_id
}

