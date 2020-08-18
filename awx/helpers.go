package awx

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/mrcrilly/goawx/client"
	"strconv"
)

func CredentialsServiceDeleteByID(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())
	client := m.(*awx.AWX)
	err := client.CredentialsService.DeleteCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete existing credentials",
			Detail:   fmt.Sprintf("Unable to delete existing credentials with id %d: %s", id, err.Error()),
		})
	}

	return diags
}
