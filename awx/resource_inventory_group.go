/*
*TBD*

Example Usage

```hcl
*TBD*
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

func resourceInventoryGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInventoryGroupCreate,
		ReadContext:   resourceInventoryGroupRead,
		UpdateContext: resourceInventoryGroupUpdate,
		DeleteContext: resourceInventoryGroupDelete,

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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceInventoryGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.GroupService

	result, err := awxService.CreateGroup(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventoryGroupTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.GroupService
	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	_, err := awxService.UpdateGroup(id, map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventoryGroupTitle, id, err)
	}

	return resourceInventoryGroupRead(ctx, d, m)

}

func resourceInventoryGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.GroupService

	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteGroup(id); err != nil {
		return buildDiagDeleteFail(
			diagElementInventoryGroupTitle,
			fmt.Sprintf("ID: %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func resourceInventoryGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.GroupService

	id, diags := convertStateIDToNummeric(diagElementInventoryGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetGroupByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventoryGroupTitle, id, err)
	}
	d = setInventoryGroupResourceData(d, res)
	return diags
}

func setInventoryGroupResourceData(d *schema.ResourceData, r *awx.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("variables", normalizeJsonYaml(r.Variables))

	d.SetId(strconv.Itoa(r.ID))
	return d
}
