---
layout: "awx"
page_title: "AWX: awx_inventory"
sidebar_current: "docs-awx-resource-inventory"
description: |-
  *TBD*
---

# awx_inventory

*TBD*

## Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_inventory" "default" {
  name            = "acc-test"
  organisation_id = data.awx_organization.default.id
  variables       = <<YAML
---
system_supporters:
  - pi
YAML
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `organisation_id` - (Required) 
* `description` - (Optional) 
* `host_filter` - (Optional) 
* `kind` - (Optional) 
* `variables` - (Optional) 

