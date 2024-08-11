// provider.go
package provider

import (
	"context"
	"terraform-provider-idcloudhost/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type: schema.TypeString,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_server": resources.ResourceServer(),
		},
		ConfigureContextFunc: contextConfig,
	}
}

func contextConfig(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {

	return nil, nil
}
