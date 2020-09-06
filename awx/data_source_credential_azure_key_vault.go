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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func dataSourceCredentialAzure() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialAzureRead,
		Schema: map[string]*schema.Schema{
			"credential_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"tenant": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialAzureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := d.Get("credential_id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", cred.Name)
	d.Set("description", cred.Description)
	d.Set("organisation_id", cred.OrganizationID)
	d.Set("url", cred.Inputs["url"])
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("tenant", cred.Inputs["tenant"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
