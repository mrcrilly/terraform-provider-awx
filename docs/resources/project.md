---
layout: "awx"
page_title: "AWX: awx_project"
sidebar_current: "docs-awx-resource-project"
description: |-
  *TBD*
---

# awx_project

*TBD*

## Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_project" "base_service_config" {
  name                 = "base-service-configuration"
  scm_type             = "git"
  scm_url              = "https://github.com/nolte/ansible_playbook-baseline-online-server"
  scm_branch           = "feature/centos8-v2"
  scm_update_on_launch = true
  organisation_id      = data.awx_organization.default.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of this project
* `organisation_id` - (Required) Numeric ID of the project organization
* `scm_type` - (Required) One of "" (manual), git, hg, svn
* `description` - (Optional) Optional description of this project.
* `local_path` - (Optional) Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.
* `scm_branch` - (Optional) Specific branch, tag or commit to checkout.
* `scm_clean` - (Optional) 
* `scm_credential_id` - (Optional) Numeric ID of the scm used credential
* `scm_delete_on_update` - (Optional) 
* `scm_update_cache_timeout` - (Optional) 
* `scm_update_on_launch` - (Optional) 
* `scm_url` - (Optional) 

