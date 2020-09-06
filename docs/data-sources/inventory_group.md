---
layout: "awx"
page_title: "AWX: awx_inventory_group"
sidebar_current: "docs-awx-datasource-inventory_group"
description: |-
  *TBD*
---

# awx_inventory_group

*TBD*

## Example Usage

```hcl
data "awx_inventory_group" "default" {
  name         = "k3sPrimary"
  inventory_id = data.awx_inventory.default.id
}
```

## Argument Reference

The following arguments are supported:

* `inventory_id` - (Required) 
* `id` - (Optional) 
* `name` - (Optional) 

