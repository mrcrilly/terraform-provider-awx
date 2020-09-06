/*
*TBD*

Example Usage

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

*/
package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		DeleteContext: resourceHostDelete,
		UpdateContext: resourceHostUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "",
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variables": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJsonYaml,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.HostService

	result, err := awxService.CreateHost(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementHostTitle, err)
	}

	hostID := result.ID
	if d.IsNewResource() {
		rawGroups := d.Get("group_ids").([]interface{})
		for _, v := range rawGroups {

			_, err := awxService.AssociateGroup(hostID, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{})
			if err != nil {
				return buildDiagnosticsMessage(
					diagElementHostTitle,
					"Assign Group Id %v to hostid %v fail, got  %s",
					v, hostID, err.Error(),
				)
			}
		}
	}
	d.SetId(strconv.Itoa(result.ID))
	return resourceHostRead(ctx, d, m)
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	_, err := awxService.UpdateHost(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementHostTitle, id, err)
	}

	if d.HasChange("group_ids") {
		// TODO Check whats happen with removin groups ....
		rawGroups := d.Get("group_ids").([]interface{})
		for _, v := range rawGroups {
			_, err := awxService.AssociateGroup(id, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{})
			if err != nil {
				return buildDiagnosticsMessage(
					diagElementHostTitle,
					"Assign Group Id %v to hostid %v fail, got  %s",
					v, id, err.Error(),
				)
			}
		}
	}
	return resourceHostRead(ctx, d, m)

}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := awxService.GetHostByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementHostTitle, id, err)
	}
	d = setHostResourceData(d, res)
	return nil
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteHost(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setHostResourceData(d *schema.ResourceData, r *awx.Host) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("enabled", r.Enabled)
	d.Set("instance_id", r.InstanceID)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	d.Set("group_ids", d.Get("group_ids").([]interface{}))
	return d
}
