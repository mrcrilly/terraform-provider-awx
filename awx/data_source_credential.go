/*
Use this data source to query Credential by ID.

Example Usage

```hcl
*TBD*
```

*/
package awx

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func dataSourceCredentialByID() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialByIDRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Credential id",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Username from searched id",
			},
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Kind from searched id",
			},
		},
	}
}

func dataSourceCredentialByIDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id := d.Get("id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential",
			Detail:   "The given credential ID is invalid or malformed",
		})
	}

	d.Set("username", cred.Inputs["username"])
	d.Set("kind", cred.Kind)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
