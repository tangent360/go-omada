package omada

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
)

type clientResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		TotalRows   int      `json:"totalRows"`
		CurrentPage int      `json:"currentPage"`
		CurrentSize int      `json:"currentSize"`
		Data        []Client `json:"data"`
	} `json:"result"`
}

type Client struct {
	Name     string `json:"name"`
	HostName string `json:"hostName,omitempty"`
	Ip       string `json:"ip"`
	MAC      string `json:"mac"`
	DnsName  string
}

func (c *Controller) GetClients() []Client {

	token := c.token
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/clients?currentPage=1&currentPageSize=999", c.baseURL, c.controllerId, c.siteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		log.Fatal(err)
	}

	var clientResponse clientResponse
	if err := json.NewDecoder(res.Body).Decode(&clientResponse); err != nil {
		log.Fatal(err)
	}

	var clients []Client
	for _, client := range clientResponse.Result.Data {
		if client.Ip == "" {
			continue
		}
		client.DnsName = makeDNSSafe(client.Name)
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].DnsName < clients[j].DnsName
	})

	return clients

}
