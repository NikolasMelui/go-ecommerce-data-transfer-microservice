package entity

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/cconfig"
)

// Client ...
type Client struct {
	BaseURL    string
	httpClient *http.Client
}

// NewClient ...
func NewClient(BaseURL string) *Client {
	return &Client{
		BaseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Get ...
// TODO: Add the interface for src
func (c *Client) Get(url string, src interface{}) error {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", cconfig.Config.SourceURL, url), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")
	req.SetBasicAuth(cconfig.Config.SourceBasicAuthLogin, cconfig.Config.SourceBasicAuthPassword)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("Unknown error, status code: %d", res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	err = xml.Unmarshal(body, &src)
	if err != nil {
		return err
	}

	return nil
}

// Set ...
// TODO: Add the interface for data
func (c *Client) Set(url string, data *map[string]interface{}) error {

	requestBody, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", cconfig.Config.TargetURL, url), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("Unknown error, status code: %d", res.StatusCode)
	}

	fmt.Println(res)

	// if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
	// 	return err
	// }
	return nil
}
