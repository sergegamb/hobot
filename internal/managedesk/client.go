package managedesk

import (
    "crypto/tls"
    "log"

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
    log.Printf("[ManageDeskAPI] NewClient: Initializing client with baseURL=%s\n", baseURL)

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

    log.Println("[ManageDeskAPI] NewClient: Client initialized successfully")

    return &Client{
        BaseURL: baseURL,
        APIKey: apiKey,
        HTTP: httpClient,
    }
}
