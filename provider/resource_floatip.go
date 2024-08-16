package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"terraform-provider-idcloudhost/provider/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFloatIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: floatIpCreate,
		ReadContext:   floatIpRead,
		UpdateContext: floatIpUpdate,
		DeleteContext: floatIpDelete,
		Schema:        schemas.FloatIpSchema,
		Importer: &schema.ResourceImporter{
			State: flaotIpState,
		},
	}
}

// only accept location from provider config
func flaotIpState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	address := d.Id()
	version := "/v1"
	path := "/network/ip_addresses/"
	generatedUrl := baseUrl + version + path
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}

	fullUrl, err := url.Parse(generatedUrl + address)
	if err != nil {
		return nil, err

	}

	client := &http.Client{}
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	name, _ := result["assigned_to_resource_type"].(string)
	billing_account_id, _ := result["billing_account_id"].(float64)

	d.Set("name", name)
	d.Set("billing_account_id", billing_account_id)
	d.Set("address", address)
	d.Set("location", defaultLocation)

	return []*schema.ResourceData{d}, nil
}

func floatIpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	version := "/v1"
	path := "/network/ip_addresses"
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

	data := map[string]interface{}{
		"name":               name,
		"billing_account_id": billing_account_id,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
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
	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	address, ok := result["address"].(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("fail to get float IP address"))
	}

	d.SetId(address)
	d.Set("address", address)

	return floatIpRead(ctx, d, m)
}

func floatIpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	address := d.Id()
	version := "/v1"
	path := "/network/ip_addresses/"
	generatedUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		generatedUrl = baseUrl + version + "/" + location + path
	}

	fullUrl, err := url.Parse(generatedUrl + address)
	if err != nil {
		return diag.FromErr(err)

	}

	client := &http.Client{}
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

func floatIpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	address := d.Id()
	version := "/v1"
	path := "/network/ip_addresses/"
	generatedUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		generatedUrl = baseUrl + version + "/" + location + path
	}
	fullUrl := generatedUrl + address

	if d.HasChanges("name", "billing_account_id") {
		data := map[string]interface{}{
			"name":               d.Get("name").(string),
			"billing_account_id": d.Get("billing_account_id").(int),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return diag.FromErr(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
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

	return floatIpRead(ctx, d, m)
}

func floatIpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	defaultLocation := config.DefaultLocation
	address := d.Id()
	version := "/v1"
	path := "/network/ip_addresses/"
	generatedUrl := baseUrl + version + path
	location := d.Get("location").(string)
	if defaultLocation != "" {
		generatedUrl = baseUrl + version + "/" + defaultLocation + path
	}
	if location != "" {
		generatedUrl = baseUrl + version + "/" + location + path
	}
	fullUrl := generatedUrl + address

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fullUrl, nil)
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
