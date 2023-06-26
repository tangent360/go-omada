package omada

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type GetNetworksResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int            `json:"totalRows"`
		CurrentPage int            `json:"currentPage"`
		CurrentSize int            `json:"currentSize"`
		Data        []OmadaNetwork `json:"data"`
	} `json:"result"`
}

type OmadaNetwork struct {
	Id     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Domain string `json:"domain,omitempty"`
	Subnet string `json:"gatewaySubnet"`
}

func (c *Controller) GetNetworks() ([]OmadaNetwork, error) {

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/lan/networks?currentPage=1&currentPageSize=999", c.baseURL, c.controllerId, c.siteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return nil, err
	}

	// respBody, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(respBody))

	var networkResponse GetNetworksResponse
	if err := json.NewDecoder(res.Body).Decode(&networkResponse); err != nil {
		return nil, err
	}

	networks := networkResponse.Result.Data
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})

	return networks, nil

}

func (c *Controller) GetAllNetworks() ([]OmadaNetwork, error) {

	var allNetworks []OmadaNetwork

	for _, v := range c.allSiteIds {
		url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/lan/networks?currentPage=1&currentPageSize=999", c.baseURL, c.controllerId, v)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Add("Csrf-Token", c.token)

		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("status code: %d", res.StatusCode)
			return nil, err
		}

		// respBody, _ := ioutil.ReadAll(res.Body)
		// fmt.Println(string(respBody))

		var networkResponse GetNetworksResponse
		if err := json.NewDecoder(res.Body).Decode(&networkResponse); err != nil {
			return nil, err
		}

		networks := networkResponse.Result.Data
		allNetworks = append(allNetworks, networks...)

		sort.Slice(allNetworks, func(i, j int) bool {
			return allNetworks[i].Name < allNetworks[j].Name
		})

	}
	return allNetworks, nil
}
