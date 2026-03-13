package paypal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	sandboxBaseURL    = "https://api-m.sandbox.paypal.com"
	productionBaseURL = "https://api-m.paypal.com"
)

type Client struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *http.Client
}

func New(clientID, clientSecret string, sandbox bool) *Client {
	base := productionBaseURL
	if sandbox {
		base = sandboxBaseURL
	}
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      base,
		httpClient:   &http.Client{},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	data := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	creds := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+creds)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: status %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	return tr.AccessToken, nil
}

type createOrderRequest struct {
	Intent        string         `json:"intent"`
	PurchaseUnits []purchaseUnit `json:"purchase_units"`
	AppContext    appContext      `json:"application_context"`
}

type purchaseUnit struct {
	Amount amount `json:"amount"`
}

type amount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type appContext struct {
	ReturnURL   string `json:"return_url"`
	CancelURL   string `json:"cancel_url"`
	BrandName   string `json:"brand_name"`
	UserAction  string `json:"user_action"`
}

type OrderResponse struct {
	ID    string `json:"id"`
	Links []link `json:"links"`
}

type link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

func (o *OrderResponse) ApprovalURL() string {
	for _, l := range o.Links {
		if l.Rel == "approve" {
			return l.Href
		}
	}
	return ""
}

func (c *Client) CreateOrder(ctx context.Context, amountUSD float64, returnURL, cancelURL string) (*OrderResponse, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	body := createOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []purchaseUnit{
			{Amount: amount{CurrencyCode: "USD", Value: fmt.Sprintf("%.2f", amountUSD)}},
		},
		AppContext: appContext{
			ReturnURL:  returnURL,
			CancelURL:  cancelURL,
			BrandName:  "Shop01",
			UserAction: "PAY_NOW",
		},
	}

	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/checkout/orders", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create order failed: status %d: %s", resp.StatusCode, b)
	}

	var order OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (c *Client) CaptureOrder(ctx context.Context, orderID string) error {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/v2/checkout/orders/"+orderID+"/capture",
		strings.NewReader("{}"))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("capture failed: status %d: %s", resp.StatusCode, b)
	}
	return nil
}
