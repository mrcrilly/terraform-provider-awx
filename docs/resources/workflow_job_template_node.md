---
layout: "awx"
page_title: "AWX: awx_workflow_job_template_node"
sidebar_current: "docs-awx-resource-workflow_job_template_node"
description: |-
  *TBD*
---

# awx_workflow_job_template_node

*TBD*

## Example Usage

```hcl
resource "random_uuid" "workflow_node_base_uuid" {}

resource "awx_workflow_job_template_node" "default" {
  workflow_job_template_id = awx_workflow_job_template.default.id
  unified_job_template_id  = awx_job_template.baseconfig.id
  inventory_id             = awx_inventory.default.id
  identifier               = random_uuid.workflow_node_base_uuid.result
}
```

## Argument Reference

The following arguments are supported:

* `identifier` - (Required) 
* `unified_job_template_id` - (Required) 
* `workflow_job_template_id` - (Required) 
* `all_parents_must_converge` - (Optional) 
* `diff_mode` - (Optional) 
* `extra_data` - (Optional) 
* `inventory_id` - (Optional) Inventory applied as a prompt, assuming job template prompts for inventory.
* `job_tags` - (Optional) 
* `job_type` - (Optional) 
* `limit` - (Optional) 
* `scm_branch` - (Optional) 
* `skip_tags` - (Optional) 
* `verbosity` - (Optional) 

