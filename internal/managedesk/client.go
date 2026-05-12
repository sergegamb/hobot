package managedesk

import (
    "crypto/tls"

    "github.com/go-resty/resty/v2"
)

type Client struct {
    BaseURL string

    APIKey string

    HTTP *resty.Client
}

func NewClient(
    baseURL string,
    apiKey string,
) *Client {

    httpClient := resty.New()

    httpClient.SetBaseURL(baseURL)

    httpClient.SetHeader(
        "Authtoken",
        apiKey,
    )

    httpClient.SetTLSClientConfig(
        &tls.Config{
            InsecureSkipVerify: true,
        },
    )

    return &Client{
        BaseURL: baseURL,
        APIKey: apiKey,
        HTTP: httpClient,
    }
}
