package provider

import (
	"context"
	// "terraform-provider-idcloudhost/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	ApiKey  string
	BaseUrl string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type:     schema.TypeString,
				Required: true,
			},
			"baseurl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.idcloudhost.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"idcloudhost_s3":              ResourceStorage(),
			"idcloudhost_private_network": ResourcePrivateNetwork(),
			"idcloudhost_float_ip":        ResourceFloatIp(),
		},
		ConfigureContextFunc: contextConfig,
	}
}

func contextConfig(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {

	config := &Config{
		ApiKey:  rd.Get("apikey").(string),
		BaseUrl: rd.Get("baseurl").(string),
	}

	return config, nil
}
