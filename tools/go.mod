module github.com/mrcrilly/terraform-provider-awx/tools

go 1.14

require (
	github.com/hashicorp/terraform v0.13.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/magefile/mage v1.10.0
	github.com/mrcrilly/terraform-provider-awx v0.1.2
	github.com/nolte/plumbing v0.0.1
)

replace github.com/mrcrilly/goawx => github.com/nolte/goawx v0.1.6

replace github.com/mrcrilly/terraform-provider-awx => ../.
