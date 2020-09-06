---
layout: "awx"
page_title: "AWX: awx_host"
sidebar_current: "docs-awx-resource-host"
description: |-
  *TBD*
---

# awx_host

*TBD*

## Example Usage

```hcl
resource "awx_host" "k3snode1" {
  name         = "k3snode1"
  description  = "pi node 1"
  inventory_id = data.awx_inventory.default.id
  group_ids = [
    data.awx_inventory_group.default.id,
    data.awx_inventory_group.pinodes.id,
  ]
  enabled   = true
  variables = <<YAML
---
ansible_host: 192.168.178.29
YAML
}
```

## Argument Reference

The following arguments are supported:

* `inventory_id` - (Required) 
* `name` - (Required) 
* `description` - (Optional) 
* `enabled` - (Optional) 
* `group_ids` - (Optional) 
* `instance_id` - (Optional) 
* `variables` - (Optional) 

