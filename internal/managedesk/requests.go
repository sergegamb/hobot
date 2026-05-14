package managedesk

import (
	"encoding/json"
	"log"
	"strconv"
)

// Filter constants for request queries
const (
	FilterAll     = "all"
	FilterOpen    = "open"
	FilterClosed  = "closed"
	FilterPending = "pending"
)

type Request struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Status  struct {
		Name string `json:"name"`
	} `json:"status"`
	// Priority    string    `json:"priority"`
	Category string `json:"category"`
	// CreatedTime time.Time `json:"created_time"`
	// UpdatedTime time.Time `json:"updated_time"`
	Requester struct {
		DisplayName string `json:"display_name"`
	} `json:"requester"`
	Description string `json:"description"`
}

type RequestsResponse struct {
	Requests []Request `json:"requests"`
}

type PaginatedRequestsResponse struct {
	Requests   []Request `json:"requests"`
	PageIndex  int       `json:"page_index"`
	PageSize   int       `json:"page_size"`
	TotalCount int       `json:"total_count"`
}

// ListInfo contains pagination and sorting configuration
type ListInfo struct {
	RowCount      int    `json:"row_count"`   // Default: 10
	StartIndex    int    `json:"start_index"` // 0-based index
	SortField     string `json:"sort_field"`  // Default: "id"
	SortOrder     string `json:"sort_order"`  // Default: "desc"
	GetTotalCount bool   `json:"get_total_count"`
}

// ListInfoRequest wraps ListInfo for the API request
type ListInfoRequest struct {
	ListInfo *ListInfo `json:"list_info"`
}

// ListInfoResponse represents the response from API with list_info
type ListInfoResponse struct {
	Requests []Request `json:"requests"`
	ListInfo struct {
		RowCount      int `json:"row_count"`
		StartIndex    int `json:"start_index"`
		SortField     string `json:"sort_field"`
		SortOrder     string `json:"sort_order"`
		TotalCount    int `json:"total_count"`
	} `json:"list_info"`
}

// NewListInfo creates a ListInfo with sensible defaults
func NewListInfo() *ListInfo {
	return &ListInfo{
		RowCount:      20,
		SortField:     "subject",
		SortOrder:     "asc",
		StartIndex:    1,
		GetTotalCount: true,
	}
}

func (c *Client) GetRequests() (
	[]Request,
	error,
) {
	log.Println("[ManageDeskAPI] GetRequests: Starting request to /api/v3/requests")

	var response RequestsResponse

	resp, err := c.HTTP.R().
		SetResult(&response).
		Get("/api/v3/requests")

	if err != nil {
		log.Printf("[ManageDeskAPI] GetRequests: Error making request: %v\n", err)
		return nil, err
	}

	log.Printf("[ManageDeskAPI] GetRequests: Success (Status: %d), Got %d requests\n", resp.StatusCode(), len(response.Requests))

	return response.Requests, nil
}

// GetRequestsWithFilters fetches requests with pagination and filter support
func (c *Client) GetRequestsWithFilters(
	filter string,
	page int,
	pageSize int,
) (
	*PaginatedRequestsResponse,
	error,
) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	log.Printf("[ManageDeskAPI] GetRequestsWithFilters: Starting request with filter=%s, page=%d, pageSize=%d\n", filter, page, pageSize)

	var response PaginatedRequestsResponse

	req := c.HTTP.R().
		SetResult(&response).
		SetQueryParam("page_index", strconv.Itoa(page-1)). // API uses 0-based indexing
		SetQueryParam("page_size", strconv.Itoa(pageSize))

	// Add filter to query if not "all"
	if filter != FilterAll && filter != "" {
		req.SetQueryParam("status", filter)
		log.Printf("[ManageDeskAPI] GetRequestsWithFilters: Added status filter: %s\n", filter)
	}

	resp, err := req.Get("/api/v3/requests")

	if err != nil {
		log.Printf("[ManageDeskAPI] GetRequestsWithFilters: Error making request: %v\n", err)
		return nil, err
	}

	log.Printf("[ManageDeskAPI] GetRequestsWithFilters: Success (Status: %d), Got %d requests (Total: %d)\n", 
		resp.StatusCode(), len(response.Requests), response.TotalCount)

	return &response, nil
}

// GetRequestByID fetches a single request by ID
func (c *Client) GetRequestByID(id string) (
	*Request,
	error,
) {
	var request Request

	_, err := c.HTTP.R().
		SetResult(&request).
		Get("/api/v3/requests/" + id)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

// GetRequestsWithListInfo fetches requests with optional list_info parameter
func (c *Client) GetRequestsWithListInfo(listInfo *ListInfo) (
	*ListInfoResponse,
	error,
) {
	// Use defaults if listInfo is nil
	if listInfo == nil {
		listInfo = NewListInfo()
	}

	// Wrap listInfo in the request structure
	request := &ListInfoRequest{
		ListInfo: listInfo,
	}

	// Marshal the wrapped listInfo to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		log.Printf("[ManageDeskAPI] GetRequestsWithListInfo: Error marshaling request: %v\n", err)
		return nil, err
	}

	log.Printf("[ManageDeskAPI] GetRequestsWithListInfo: Sending input_data: %s\n", string(requestJSON))

	var response ListInfoResponse

	resp, err := c.HTTP.R().
		SetResult(&response).
		SetQueryParam("input_data", string(requestJSON)).
		Get("/api/v3/requests")

	if err != nil {
		log.Printf("[ManageDeskAPI] GetRequestsWithListInfo: Error making request: %v\n", err)
		return nil, err
	}

	log.Printf("[ManageDeskAPI] GetRequestsWithListInfo: Success (Status: %d), Got %d requests (Total: %d)\n", 
		resp.StatusCode(), len(response.Requests), response.ListInfo.TotalCount)
	return &response, nil
}