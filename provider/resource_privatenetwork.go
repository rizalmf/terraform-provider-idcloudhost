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

func ResourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: privateNetworkCreate,
		ReadContext:   privateNetworkRead,
		UpdateContext: privateNetworkUpdate,
		DeleteContext: privateNetworkDelete,
		Schema:        schemas.PrivateNteworkSchema,
	}
}

func privateNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/network/network"

	name := d.Get("name").(string)

	fullUrl, err := url.Parse(baseUrl + path)
	if err != nil {
		return diag.FromErr(err)

	}

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("name", name)
	fullUrl.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("POST", fullUrl.String(), nil)
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

	d.SetId(uuid)

	return privateNetworkRead(ctx, d, m)
}

func privateNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/network/network/"
	uuid := d.Id()

	fullUrl, err := url.Parse(baseUrl + path + uuid)
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

func privateNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/network/network/"
	uuid := d.Id()
	fullUrl := baseUrl + path + uuid

	if d.HasChange("name") {
		data := map[string]interface{}{
			"name": d.Get("name").(string),
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

	return privateNetworkRead(ctx, d, m)
}

func privateNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/network/network/"
	uuid := d.Id()
	fullUrl := baseUrl + path + uuid

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fullUrl, nil)
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

	return nil
}
