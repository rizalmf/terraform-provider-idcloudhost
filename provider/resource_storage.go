package provider

import (
	"context"
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

func ResourceStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: storageCreate,
		ReadContext:   storageRead,
		UpdateContext: storageUpdate,
		DeleteContext: storageDelete,
		Schema:        schemas.StorageSchema,
	}
}

func storageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/storage/bucket"
	fullUrl := baseUrl + path

	name := d.Get("name").(string)
	billing_account_id := d.Get("billing_account_id").(int)

	client := &http.Client{}
	form := url.Values{}
	form.Add("name", name)
	form.Add("billing_account_id", strconv.Itoa(billing_account_id))
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(form.Encode()))
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

	if resp.StatusCode >= 299 || resp.StatusCode <= 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))
	}
	defer resp.Body.Close()

	d.SetId(name)

	return storageRead(ctx, d, m)
}

func storageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/storage/bucket"
	name := d.Id()

	fullUrl, err := url.Parse(baseUrl + path)
	if err != nil {
		return diag.FromErr(err)

	}

	client := &http.Client{}
	queryParams := url.Values{}
	queryParams.Add("name", name)
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
	if resp.StatusCode >= 299 || resp.StatusCode <= 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))

	}
	defer resp.Body.Close()

	return nil
}

func storageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/storage/bucket"
	fullUrl := baseUrl + path

	name := d.Id()
	billing_account_id := d.Get("billing_account_id").(int)

	d.SetId(name)

	if d.HasChange("billing_account_id") {
		client := &http.Client{}
		form := url.Values{}
		form.Add("name", name)
		form.Add("billing_account_id", strconv.Itoa(billing_account_id))
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

		if resp.StatusCode >= 299 || resp.StatusCode <= 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return diag.FromErr(fmt.Errorf(string(bodyBytes)))
		}
		defer resp.Body.Close()
	}

	return storageRead(ctx, d, m)
}

func storageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	apiKey := config.ApiKey
	baseUrl := config.BaseUrl
	path := "/v1/storage/bucket"
	fullUrl := baseUrl + path

	name := d.Id()
	client := &http.Client{}
	form := url.Values{}
	form.Add("name", name)
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

	if resp.StatusCode >= 299 || resp.StatusCode <= 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf(string(bodyBytes)))
	}
	defer resp.Body.Close()

	return nil
}
