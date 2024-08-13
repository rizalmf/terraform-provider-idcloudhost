package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"terraform-provider-idcloudhost/provider/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVm() *schema.Resource {
	return &schema.Resource{
		CreateContext: vmCreate,
		ReadContext:   vmRead,
		UpdateContext: vmUpdate,
		DeleteContext: vmDelete,
		Schema:        schemas.VmSchema,
	}
}

func vmCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/user-resource/vm"
	fullUrl := baseUrl + path

	name := d.Get("name").(string)
	billing_account_id := d.Get("billing_account_id").(int)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	private_network_uuid := d.Get("private_network_uuid").(string)
	os_name := d.Get("os_name").(string)
	os_version := d.Get("os_version").(string)
	vcpu := d.Get("vcpu").(int)
	ram := d.Get("ram").(int)
	disks := d.Get("disks").(int)

	client := &http.Client{}
	form := url.Values{}
	form.Add("name", name)
	form.Add("billing_account_id", strconv.Itoa(billing_account_id))
	form.Add("username", username)
	form.Add("password", password)
	form.Add("network_uuid", private_network_uuid)
	form.Add("os_name", os_name)
	form.Add("os_version", os_version)
	form.Add("vcpu", strconv.Itoa(vcpu))
	form.Add("ram", strconv.Itoa(ram))
	form.Add("disks", strconv.Itoa(disks))
	form.Add("reserve_public_ip", "false")
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
		return diag.FromErr(err)
	}
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return diag.FromErr(err)
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))
	}

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return diag.FromErr(err)
	}

	defer resp.Body.Close()

	uuid, ok := result["uuid"].(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("fail to get generated UUID"))
	}

	disks_uuid := ""
	// Navigate the map to get the value of storage[0].uuid
	if storageArray, ok := result["storage"].([]interface{}); ok {
		if firstStorage, ok := storageArray[0].(map[string]interface{}); ok {
			if storageUuid, ok := firstStorage["uuid"].(string); ok {
				disks_uuid = storageUuid
			}
		}
	}
	if disks_uuid == "" {
		return diag.FromErr(fmt.Errorf("fail to get generated storage UUID"))
	}

	d.SetId(uuid)
	d.Set("uuid", uuid)
	d.Set("disks_uuid", disks_uuid)

	return vmRead(ctx, d, m)
}

func vmRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/user-resource/vm"
	uuid := d.Id()

	fullUrl, err := url.Parse(baseUrl + path)
	if err != nil {
		return diag.FromErr(err)

	}

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("uuid", uuid)
	fullUrl.RawQuery = queryParams.Encode()
	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))

	}
	defer resp.Body.Close()

	return nil
}

func vmUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/user-resource/vm"
	fullUrl := baseUrl + path

	uuid := d.Id()
	disks_uuid := d.Get("disks_uuid").(string)
	name := d.Get("name").(string)
	ram := d.Get("ram").(int)
	vcpu := d.Get("vcpu").(int)
	disks := d.Get("disks").(int)
	// float_ip_address := d.Get("float_ip_address").(string)

	if d.HasChanges("name", "ram", "vcpu") {
		client := &http.Client{}
		form := url.Values{}
		form.Add("uuid", uuid)
		form.Add("name", name)
		form.Add("ram", strconv.Itoa(ram))
		form.Add("vcpu", strconv.Itoa(vcpu))
		req, err := http.NewRequest("PATCH", fullUrl, strings.NewReader(form.Encode()))
		if err != nil {
			return diag.FromErr(err)
		}
		req.PostForm = form
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("apikey", apiKey)
		resp, err := client.Do(req)

		if err != nil {
			return diag.FromErr(err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return diag.FromErr(fmt.Errorf(string(bodyBytes)))
		}
		defer resp.Body.Close()
	}

	if d.HasChange("disks") {
		path = "/v1/user-resource/vm/storage"
		fullUrl = baseUrl + path
		client := &http.Client{}
		form := url.Values{}
		form.Add("uuid", uuid)
		form.Add("disk_uuid", disks_uuid)
		form.Add("size_gb", strconv.Itoa(disks))
		req, err := http.NewRequest("PATCH", fullUrl, strings.NewReader(form.Encode()))
		if err != nil {
			return diag.FromErr(err)
		}
		req.PostForm = form
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("apikey", apiKey)
		resp, err := client.Do(req)

		if err != nil {
			return diag.FromErr(err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return diag.FromErr(fmt.Errorf(string(bodyBytes)))
		}
		defer resp.Body.Close()
	}

	if d.HasChange("float_ip_address") {
		oldIntrface, newIntrface := d.GetChange("float_ip_address")
		oldAddr := oldIntrface.(string)
		newAddr := newIntrface.(string)

		// unassign old ip
		path = fmt.Sprintf("v1/network/ip_addresses/%s/unassign", oldAddr)
		fullUrl = baseUrl + path
		client := &http.Client{}
		req, err := http.NewRequest("POST", fullUrl, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		req.Header.Set("apikey", apiKey)
		resp, err := client.Do(req)

		if err != nil {
			return diag.FromErr(err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return diag.FromErr(fmt.Errorf(string(bodyBytes)))
		}
		defer resp.Body.Close()

		if newAddr != "" {
			// assign new address
			path = fmt.Sprintf("v1/network/ip_addresses/%s/assign", newAddr)
			fullUrl = baseUrl + path
			client := &http.Client{}
			req, err := http.NewRequest("POST", fullUrl, nil)
			if err != nil {
				return diag.FromErr(err)
			}
			req.Header.Set("apikey", apiKey)
			resp, err := client.Do(req)

			if err != nil {
				return diag.FromErr(err)
			}

			if resp.StatusCode > 299 || resp.StatusCode < 200 {
				bodyBytes, _ := io.ReadAll(resp.Body)
				return diag.FromErr(fmt.Errorf(string(bodyBytes)))
			}
			defer resp.Body.Close()
		}
	}
	return vmRead(ctx, d, m)
}

func vmDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/user-resource/vm"
	fullUrl := baseUrl + path

	uuid := d.Id()
	client := &http.Client{}
	form := url.Values{}
	form.Add("uuid", uuid)
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return diag.FromErr(err)
	}
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))
	}
	defer resp.Body.Close()

	return nil
}
