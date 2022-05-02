package easyredir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type EasyRedir struct {
	Client *Client
	Rules  *Rules
}

type Options struct {
	apiKey    string
	apiSecret string
}

type Client struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
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

func New(options ...func(*Options)) (e EasyRedir) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}

	e.Client = &Client{
		baseURL:    baseURLV1,
		apiKey:     opts.apiKey,
		apiSecret:  opts.apiSecret,
		httpClient: &http.Client{},
	}

	e.Rules = &Rules{}

	return e
}

func WithAPIKey(key string) func(*Options) {
	return func(o *Options) {
		o.apiKey = key
	}
}

func WithAPISecret(secret string) func(*Options) {
	return func(o *Options) {
		o.apiSecret = secret
	}
}

func (e *EasyRedir) GetRules() (rules *Rules, err error) {
	url := fmt.Sprintf("%s/rules", e.Client.baseURL)

	if err := e.Client.getJSON(url, e.Rules); err != nil {
		return e.Rules, fmt.Errorf("unable to get rules: %w", err)
	}

	return e.Rules, nil
}

func (cl *Client) getJSON(url string, v interface{}) (err error) {
	body, err := cl.sendRequest(url)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("unable to parse json: %w", err)
	}

	return nil
}

func (cl *Client) sendRequest(url string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new request: %w", err)
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("received status code: %d", res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	return body, nil
}

func (rs *Rules) String() (str string) {
	var sb strings.Builder

	for _, r := range rs.Data {
		fmt.Fprintf(&sb, "%s: %s --> %s\n", r.ID, r.Attributes.SourceURLs, r.Attributes.TargetURL)
	}

	str = sb.String()

	return str
}
