---
layout: "awx"
page_title: "AWX: awx_inventory"
sidebar_current: "docs-awx-datasource-inventory"
description: |-
  *TBD*
---

# awx_inventory

*TBD*

## Example Usage

```hcl
data "awx_inventory" "default" {
  name            = "private_services"
  organisation_id = data.awx_organization.default.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) 
* `name` - (Optional) 
* `organisation_id` - (Optional) 

