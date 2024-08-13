package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var VmSchema = map[string]*schema.Schema{
	"uuid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"disks_uuid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"location": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"billing_account_id": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"username": {
		Type:     schema.TypeString,
		Required: true,
	},
	"password": {
		Type:     schema.TypeString,
		Required: true,
	},
	"os_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"os_version": {
		Type:     schema.TypeString,
		Required: true,
	},
	"vcpu": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"ram": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"disks": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"private_network_uuid": {
		Type:     schema.TypeString,
		Required: true,
	},
	"float_ip_address": {
		Type:     schema.TypeString,
		Optional: true,
	},
}
