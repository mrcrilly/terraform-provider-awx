---
layout: "awx"
page_title: "AWX: awx_job_template_credential"
sidebar_current: "docs-awx-resource-job_template_credential"
description: |-
  *TBD*
---

# awx_job_template_credential

*TBD*

## Example Usage

```hcl
resource "awx_job_template_credentials" "baseconfig" {
  job_template_id = awx_job_template.baseconfig.id
  credential_id   = awx_credential_machine.pi_connection.id
}
```

## Argument Reference

The following arguments are supported:

* `credential_id` - (Required, ForceNew) 
* `job_template_id` - (Required, ForceNew) 

