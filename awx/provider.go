package awx

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_HOSTNAME", "http://localhost"),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_USERNAME", "admin"),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_PASSWORD", "password"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_credential_ssh": resourceCredentialSSH(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"awx_credentials": dataSourceCredentials(),
			"awx_credential": dataSourceCredentialByID(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := awx.NewAWX(hostname, username, password, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create AWX client",
			Detail:   "Unable to auth user against AWX API: check the hostname, username and password",
		})
		return nil, diags
	}

	return c, diags
}
