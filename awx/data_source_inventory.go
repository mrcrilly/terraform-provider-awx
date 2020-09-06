/*
*TBD*

Example Usage

```hcl
data "awx_inventory" "default" {
  name            = "private_services"
  organisation_id = data.awx_organization.default.id
}
```

*/
package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func dataSourceInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoriesRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organisation_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceInventoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if organisationID, okIOrgID := d.GetOk("organisation_id"); okIOrgID {
		params["organization"] = strconv.Itoa(organisationID.(int))
	}
	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or group_id)",
		)
		return diags
	}
	inventories, _, err := client.InventoriesService.ListInventories(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Inventory Group",
			"Fail to find the group got: %s",
			err.Error(),
		)
	}
	if len(inventories) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(inventories),
		)
		return diags
	}

	inventory := inventories[0]
	d = setInventoryResourceData(d, inventory)
	return diags
}
