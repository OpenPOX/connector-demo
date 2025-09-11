package confluencev2

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const (
	APIResourceBase     = "https://api.atlassian.com"
	APIResourceGetPages = "/wiki/rest/api/content"
)

// GetResourceURL 获取资源URL
func GetResourceURL(cloudID, api string) string {
	return APIResourceBase + "/ex/confluence/" + cloudID + api
}

// API represents the Confluence V2 API client
type API struct {
	client        *http.Client
	CloudID       string
	Authorization string
}

// NewAPI creates a new Confluence V2 API client
func NewAPI(auth, cloudID string) *API {
	return &API{
		client:        &http.Client{},
		CloudID:       cloudID,
		Authorization: auth,
	}
}

// GetPagesResponse represents the response from the GetPages API
type GetPagesResponse struct {
	Results []struct {
		ParentID *string `json:"parentId"`
		SpaceID  string  `json:"spaceId"`
		OwnerID  string  `json:"ownerId"`
		Title    string  `json:"title"`
		ID       string  `json:"id"`
		Status   string  `json:"status"`
		Body     struct {
			Storage struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			} `json:"storage"`
			AtlasDocFormat struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			} `json:"atlas_doc_format"`
		} `json:"body"`
	}
}

// GetPages retrieves pages from Confluence
func (a *API) GetPages(params url.Values) (GetPagesResponse, error) {
	api := GetResourceURL(a.CloudID, APIResourceGetPages)
	httpReq, err := http.NewRequest("GET", api, strings.NewReader(params.Encode()))
	if err != nil {
		return GetPagesResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+a.Authorization)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return GetPagesResponse{}, err
	}
	defer resp.Body.Close()

	var out GetPagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return GetPagesResponse{}, err
	}
	return out, nil
}
