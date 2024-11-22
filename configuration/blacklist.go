package configuration

import (
	"net/http"
	"net/url"
)

type BlacklistGroup []string

const TOPTAL_API_URL string = "https://www.toptal.com"
const TOPTAL_TYPE_PATH string = "developers/gitignore/api"

type Client struct {
	client    *http.Client
	UserAgent string
	ApiUrl    *url.URL
}

func NewClient() (*Client, error) {
	c := &Client{
		client:    http.DefaultClient,
		UserAgent: "todor",
	}

	api_url, err := url.Parse(TOPTAL_API_URL)
	if err != nil {
		return &Client{}, err
	}
	c.ApiUrl = api_url

	return c, nil
}
