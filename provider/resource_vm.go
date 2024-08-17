package provider

import (
	"bytes"
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
		Importer: &schema.ResourceImporter{
			State: vmState,
		},
	}
}

func findPrivateNetwork(c Config, uuid string) (string, error) {
	apiKey := c.ApiKey
	baseUrl := c.BaseUrl
	defaultLocation := c.DefaultLocation
	version := "/v1"
	path := "/network/networks"
	fullUrl := baseUrl + version + path
	if defaultLocation != "" {
		fullUrl = baseUrl + version + "/" + defaultLocation + path
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(string(bodyBytes))
	}
	private_network_uuid := ""
	var result []interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	// Iterate over the networks and extract the VM UUIDs
out:
	for _, network := range result {
		if networkMap, ok := network.(map[string]interface{}); ok {
			if vmUUIDs, ok := networkMap["vm_uuids"].([]interface{}); ok {
				for _, vmUUID := range vmUUIDs {
					if uuidStr, ok := vmUUID.(string); ok {
						if uuidStr == uuid {
							private_network_uuid = networkMap["uuid"].(string)
							break out
						}
					}
				}
			}
		}
	}

	return private_network_uuid, nil
}

func findFloatIp(c Config, uuid string) (string, error) {

	apiKey := c.ApiKey
	baseUrl := c.BaseUrl
	defaultLocation := c.DefaultLocation
	version := "/v1"
	path := "/network/ip_addresses"
	fullUrl := baseUrl + version + path
	if defaultLocation != "" {
		fullUrl = baseUrl + version + "/" + defaultLocation + path
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(string(bodyBytes))
	}
	floatip_address := ""
	var result []interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	// Iterate over the networks and extract the assigned_to
	for _, network := range result {
		if networkMap, ok := network.(map[string]interface{}); ok {
			assigned_to, _ := networkMap["assigned_to"].(string)
			if assigned_to == uuid {
				floatip_address = networkMap["address"].(string)
				break
			}
		}
	}

	return floatip_address, nil
}

// only accept location from provider config
func vmState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	version := "/v1"
	path := "/user-resource/vm"
	generatedUrl := baseUrl + version + path
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}
	uuid := d.Id()

	fullUrl, err := url.Parse(generatedUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("uuid", uuid)
	fullUrl.RawQuery = queryParams.Encode()
	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(bodyBytes))
	}

	var result map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	name, _ := result["name"].(string)
	billing_account_id, _ := result["billing_account"].(float64)
	username, _ := result["username"].(string)
	password := "<Hidden>"
	os_name, _ := result["os_name"].(string)
	os_version, _ := result["os_version"].(string)
	vcpu, _ := result["vcpu"].(float64)
	ram, _ := result["memory"].(float64)
	private_network_uuid, err := findPrivateNetwork(*config, uuid)
	if err != nil {
		return nil, err
	}
	float_ip_address, err := findFloatIp(*config, uuid)
	if err != nil {
		return nil, err
	}

	desired_status, _ := result["status"].(string)

	disks_uuid := ""
	disk_size := 0
	// Navigate the map to get the value of storage[0].uuid
	if storageArray, ok := result["storage"].([]interface{}); ok {
		if firstStorage, ok := storageArray[0].(map[string]interface{}); ok {
			if storageUuid, ok := firstStorage["uuid"].(string); ok {
				disks_uuid = storageUuid
				foundDisk, _ := firstStorage["size"].(float64)
				disk_size = int(foundDisk)
			}
		}
	}
	if disks_uuid == "" {
		return nil, fmt.Errorf("fail to get generated storage UUID")
	}

	d.Set("uuid", uuid)
	d.Set("disks_uuid", disks_uuid)
	d.Set("location", defaultLocation)
	d.Set("name", name)
	d.Set("billing_account_id", billing_account_id)
	d.Set("username", username)
	d.Set("password", password)
	d.Set("os_name", os_name)
	d.Set("os_version", os_version)
	d.Set("vcpu", int(vcpu))
	d.Set("ram", int(ram))
	d.Set("disks", disk_size)
	d.Set("private_network_uuid", private_network_uuid)
	d.Set("float_ip_address", float_ip_address)
	d.Set("desired_status", desired_status)

	return []*schema.ResourceData{d}, nil
}

func vmCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	path := "/user-resource/vm"
	version := "/v1"
	fullUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		fullUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		fullUrl = baseUrl + version + "/" + location + path
	}

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
	defaultLocation := config.DefaultLocation
	path := "/user-resource/vm"
	version := "/v1"
	generatedUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		generatedUrl = baseUrl + version + "/" + location + path
	}
	uuid := d.Id()

	fullUrl, err := url.Parse(generatedUrl)
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
	// path := "/v1/user-resource/vm"
	defaultLocation := config.DefaultLocation
	path := "/user-resource/vm"
	version := "/v1"
	fullUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		fullUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		fullUrl = baseUrl + version + "/" + location + path
	}
	// fullUrl := baseUrl + path

	uuid := d.Id()
	disks_uuid := d.Get("disks_uuid").(string)
	name := d.Get("name").(string)
	ram := d.Get("ram").(int)
	vcpu := d.Get("vcpu").(int)
	disks := d.Get("disks").(int)
	// desired_status := d.Get("desired_status").(string)
	// float_ip_address := d.Get("float_ip_address").(string)

	if d.HasChange("desired_status") {

	}

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
		version = "/v1"
		path = "/user-resource/vm/storage"
		fullUrl = baseUrl + version + path
		if defaultLocation != "" {
			fullUrl = baseUrl + version + "/" + defaultLocation + path
		}
		if location != "" {
			fullUrl = baseUrl + version + "/" + location + path
		}
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

		if oldAddr != "" {
			// unassign old ip
			version = "/v1"
			path = fmt.Sprintf("/network/ip_addresses/%s/unassign", oldAddr)
			fullUrl = baseUrl + version + path
			if defaultLocation != "" {
				fullUrl = baseUrl + version + "/" + defaultLocation + path
			}
			if location != "" {
				fullUrl = baseUrl + version + "/" + location + path
			}
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

		if newAddr != "" {
			// assign new address
			data := map[string]interface{}{
				"vm_uuid": uuid,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				return diag.FromErr(err)
			}
			version = "/v1"
			path = fmt.Sprintf("/network/ip_addresses/%s/assign", newAddr)
			fullUrl = baseUrl + version + path
			if defaultLocation != "" {
				fullUrl = baseUrl + version + "/" + defaultLocation + path
			}
			if location != "" {
				fullUrl = baseUrl + version + "/" + location + path
			}
			client := &http.Client{}
			req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
			if err != nil {
				return diag.FromErr(err)
			}
			req.Header.Set("Content-Type", "application/json")
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
	defaultLocation := config.DefaultLocation
	path := "/user-resource/vm"
	version := "/v1"
	fullUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		fullUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		fullUrl = baseUrl + version + "/" + location + path
	}

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
