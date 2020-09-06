
resource "awx_organization" "default" {
  name            = "acc-test"
}

resource "awx_inventory" "default" {
  name            = "acc-test"
  organisation_id = awx_organization.default.id
  variables       = <<YAML
---
system_supporters:
  - pi
YAML
}


resource "awx_credential_machine" "credential" {
  organisation_id     = awx_organization.default.id
  name                = "acc-machine-credential"
  username            = "pi"
  ssh_key_data        = file("${path.module}/files/id_rsa")
  ssh_public_key_data = file("${path.module}/files/id_rsa.pub")
  ssh_key_unlock      = "test"
}

resource "awx_credential_scm" "credential" {
  organisation_id     = awx_organization.default.id
  name                = "acc-scm-credential"
  username            = "test"
  ssh_key_data        = file("${path.module}/files/id_rsa")
  ssh_key_unlock      = "test"
}

resource "awx_inventory_group" "default" {
  name         = "common-services"
  inventory_id = awx_inventory.default.id
}
resource "awx_workflow_job_template" "default" {
  name         = "acc-workflow-job"
  organisation_id      = awx_organization.default.id
  inventory_id = awx_inventory.default.id
}
resource "random_uuid" "workflow_node_k3s_uuid" {}

resource "awx_workflow_job_template_node" "default" {

  workflow_job_template_id      = awx_workflow_job_template.default.id
  unified_job_template_id = awx_job_template.template.id
  inventory_id = awx_inventory.default.id
  identifier = random_uuid.workflow_node_k3s_uuid.result
}
resource "random_uuid" "workflow_node_second_uuid" {}

resource "awx_workflow_job_template_node_success" "default" {

  workflow_job_template_node_id = awx_workflow_job_template_node.default.id
  unified_job_template_id = awx_job_template.template.id
  inventory_id = awx_inventory.default.id
  identifier = random_uuid.workflow_node_second_uuid.result
}

resource "awx_host" "k3snode1" {
  name         = "acc-node1"
  inventory_id = awx_inventory.default.id
  enabled   = true
  variables = <<YAML
---
ansible_host: 192.168.178.29
YAML
}

resource "awx_project" "project" {
  name                 = "acc-project"
  scm_type             = "git"
  scm_url              = "https://github.com/nolte/ansible_playbook-baseline-k3s"
  scm_branch           = "feature/controllable-firelld"
  scm_update_on_launch = true
  organisation_id      = awx_organization.default.id
}

## give Certsmanger Time to Work
resource "time_sleep" "wait_seconds" {
  depends_on = [awx_project.project]

  create_duration = "15s"
}

resource "awx_job_template" "template" {
  depends_on = [time_sleep.wait_seconds]
  name           = "acc-job-template"
  job_type       = "run"
  inventory_id   = awx_inventory.default.id
  project_id     = awx_project.project.id
  playbook       = "playbook-install-k3s.yaml"
  become_enabled = true
}

resource "awx_job_template_credential" "template_credentials" {
  job_template_id = awx_job_template.template.id
  credential_id   = awx_credential_machine.credential.id
}

resource "awx_inventory_source" "inventory_surces" {
  name              = "acc-inventory-sources"
  inventory_id      = awx_inventory.default.id
  source_project_id = awx_project.project.id
}


output "inventory_id" {
    value = awx_inventory.default.id
}

