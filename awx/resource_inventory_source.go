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

func resourceInventorySource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInventorySourceCreate,
		ReadContext:   resourceInventorySourceRead,
		UpdateContext: resourceInventorySourceUpdate,
		DeleteContext: resourceInventorySourceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"overwrite_vars": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"verbosity": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},
			"update_cache_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  0,
				Optional: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "scm",
				Optional: true,
			},
			"source_project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_path": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventorySourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService

	result, err := awxService.CreateInventorySource(map[string]interface{}{
		"name":           d.Get("name").(string),
		"inventory":      d.Get("inventory_id").(int),
		"overwrite_vars": d.Get("overwrite_vars").(bool),
		"verbosity":      d.Get("verbosity").(int),
		"source":         d.Get("source").(string),
		"source_project": d.Get("source_project_id").(int),
		"source_path":    d.Get("source_path").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventorySourceTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventorySourceRead(ctx, d, m)

}

func resourceInventorySourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNummeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}

	_, err := awxService.UpdateInventorySource(id, map[string]interface{}{
		"name":           d.Get("name").(string),
		"overwrite_vars": d.Get("overwrite_vars").(bool),
		"verbosity":      d.Get("verbosity").(int),
		"source":         d.Get("source").(int),
		"source_project": d.Get("source_project_id").(int),
		"source_path":    d.Get("source_path").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventorySourceTitle, id, err)
	}

	return resourceInventorySourceRead(ctx, d, m)
}

func resourceInventorySourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNummeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := awxService.DeleteInventorySource(id); err != nil {
		return buildDiagDeleteFail(
			"inventroy source",
			fmt.Sprintf("inventroy source %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func resourceInventorySourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNummeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := awxService.GetInventorySourceByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventorySourceTitle, id, err)
	}
	d = setInventorySourceResourceData(d, res)
	return nil
}

func setInventorySourceResourceData(d *schema.ResourceData, r *awx.InventorySource) *schema.ResourceData {
	d.Set("name", r.Name)

	d.Set("inventory_id", r.Inventory)
	d.Set("overwrite_vars", r.OverwriteVars)
	d.Set("verbosity", r.Verbosity)
	d.Set("source", r.Source)
	d.Set("source_project_id", r.SourceProject)
	d.Set("source_path", r.SourcePath)

	return d
}
