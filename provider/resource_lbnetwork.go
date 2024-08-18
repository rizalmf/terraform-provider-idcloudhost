package provider

import (
	"context"
	"terraform-provider-idcloudhost/provider/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: loadBalancerCreate,
		ReadContext:   loadBalancerRead,
		UpdateContext: loadBalancerUpdate,
		DeleteContext: loadBalancerDelete,
		Schema:        schemas.PrivateNteworkSchema,
		Importer: &schema.ResourceImporter{
			State: loadBalancerState,
		},
	}
}

// only accept location from provider config
func loadBalancerState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}

func loadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return loadBalancerRead(ctx, d, m)
}

func loadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}

func loadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return loadBalancerRead(ctx, d, m)
}

func loadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}
