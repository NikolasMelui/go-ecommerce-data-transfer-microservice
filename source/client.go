package source

import (
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
	*http.Client
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		&http.Client{
			Timeout: time.Minute,
		},
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetData ...
func (c *Client) GetData(url string, src interface{}) error {

	req, err := http.NewRequest("GET", cconfig.Config.SourceURL+url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")
	req.SetBasicAuth(cconfig.Config.SourceBasicAuthLogin, cconfig.Config.SourceBasicAuthPassword)

	res, err := c.Do(req)
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
