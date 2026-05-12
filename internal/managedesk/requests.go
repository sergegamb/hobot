package managedesk

type Request struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Status  struct {
		Name string `json:"name"`
	} `json:"status"`
}

type RequestsResponse struct {
	Requests []Request `json:"requests"`
}

func (c *Client) GetRequests() (
	[]Request,
	error,
) {

	var response RequestsResponse

	_, err := c.HTTP.R().
		SetResult(&response).
		Get("/api/v3/requests")

	if err != nil {
		return nil, err
	}

	return response.Requests, nil
}
