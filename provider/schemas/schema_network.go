package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var PrivateNteworkSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"network_uuid": {
		Type:     schema.TypeString,
		Computed: true,
	},
}
