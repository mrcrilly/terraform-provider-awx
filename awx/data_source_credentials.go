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

func dataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialsRead,
		Schema: map[string]*schema.Schema{
			"credentials": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)

	creds, _, err := client.CredentialsService.ListCredentials(map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   "Unable to fetch credentials from AWX API",
		})
		return diags
	}

	parsedCreds := make([]map[string]interface{}, 0)
	for _, c := range creds {
		parsedCreds = append(parsedCreds, map[string]interface{}{
			"id":       c.ID,
			"username": c.Inputs["username"],
			"kind":     c.Kind,
		})
	}

	err = d.Set("credentials", parsedCreds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
