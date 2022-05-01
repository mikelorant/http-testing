package easyredir

import (
	"encoding/json"
	"fmt"
	_ "github.com/davecgh/go-spew/spew"
	"net/http"
)

type EasyRedir struct {
	Client *Client
	Rules  *Rules
}

type Options struct {
	APIKey    string
	APISecret string
}

type Client struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	HTTPClient *http.Client
}

type Rules struct {
	Data  []RulesData `type:"data"`
	Meta  Meta        `type:"meta"`
	Links Links       `type:"links"`
}

type RulesData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes RulesAttributes `json:"attributes"`
}

type RulesAttributes struct {
	ForwardParams bool     `json:"forward_params"`
	ForwardPath   bool     `json:"forward_path"`
	ResponseType  string   `json:"response_type"`
	SourceURLs    []string `json:"source_urls"`
	TargetURL     string   `json:"target_url"`
}

type Meta struct {
	HasMore bool `json:"has_more"`
}

type Links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

const (
	baseURLV1 = "https://api.easyredir.com/v1"
)

func New(opts *Options) (e EasyRedir) {
	e.Client = &Client{
		baseURL:    baseURLV1,
		apiKey:     opts.APIKey,
		apiSecret:  opts.APISecret,
		HTTPClient: &http.Client{},
	}

	e.Rules = &Rules{}

	return e
}

func (e *EasyRedir) GetRules() (err error) {
	url := fmt.Sprintf("%s/rules", e.Client.baseURL)

	if err := e.Client.getJSON(url, e.Rules); err != nil {
		return fmt.Errorf("unable to get rules: %w", err)
	}

	return nil
}

func (cl *Client) getJSON(url string, v interface{}) (err error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)

	res, err := cl.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return fmt.Errorf("unable to parse json: %w", err)
	}

	return nil
}
