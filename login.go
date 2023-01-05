package omada

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Controller struct {
	httpClient   *http.Client
	baseURL      string
	controllerId string
	token        string
}

type ControllerInfo struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		ControllerVer string `json:"controllerVer"`
		APIVer        string `json:"apiVer"`
		Configured    bool   `json:"configured"`
		Type          int    `json:"type"`
		SupportApp    bool   `json:"supportApp"`
		OmadacID      string `json:"omadacId"`
	} `json:"result"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		RoleType int    `json:"roleType"`
		Token    string `json:"token"`
	} `json:"result"`
}

func New(baseURL string) Controller {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: (30 * time.Second),
	}

	return Controller{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Controller) GetControllerInfo() error {

	url := c.baseURL + "/api/info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	var info ControllerInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return err
	}

	c.controllerId = info.Result.OmadacID
	return nil

}

func (c *Controller) Login(user string, pass string) error {

	endpoint := c.baseURL + "/" + c.controllerId + "/api/v2/login"

	loginBody := LoginBody{
		Username: user,
		Password: pass,
	}

	loginJSON, err := json.Marshal(loginBody)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(loginJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return err
	}

	// todo:
	// - how long is login session valid for
	// - when does it need to be refreshed
	// u, err := url.Parse(c.baseURL)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// domain, _ := url.Parse(u.Hostname())
	// cookies := c.httpClient.Jar.Cookies(domain)
	// if len(cookies) == 0 {
	// 	fmt.Println("No cookies found")
	// }
	// fmt.Println(cookies)

	var login LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
		return err
	}

	token := login.Result.Token
	c.token = token
	return nil

}
