package omada

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"time"
)

type Controller struct {
	httpClient   *http.Client
	baseURL      string
	controllerId string
	token        string
	siteId       string
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

	v, _ := os.LookupEnv("OMADA_DISABLE_HTTPS_VERIFICATION")
	disableHttpsVerification, _ := strconv.ParseBool(v)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: disableHttpsVerification},
	}
	httpClient := &http.Client{
		Jar:       jar,
		Timeout:   (30 * time.Second),
		Transport: transport,
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

func (c *Controller) Login(user string, pass string, siteName string) error {

	endpoint := c.baseURL + "/" + c.controllerId + "/api/v2/login"

	loginBody := LoginBody{
		Username: user,
		Password: pass,
	}

	loginJSON, err := json.Marshal(loginBody)
	if err != nil {
		return err
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

	if login.ErrorCode != 0 {
		return fmt.Errorf("omada login error, code: %d, message: %s", login.ErrorCode, login.Msg)
	}

	token := login.Result.Token
	c.token = token

	c.getSiteId(siteName)

	return nil

}

type currentUserResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    struct {
		ID         string `json:"id"`
		Type       int    `json:"type"`
		RoleType   int    `json:"roleType"`
		Name       string `json:"name"`
		OmadacID   string `json:"omadacId"`
		Adopt      bool   `json:"adopt"`
		Manage     bool   `json:"manage"`
		License    bool   `json:"license"`
		SiteManage bool   `json:"siteManage"`
		Privilege  struct {
			Sites       []Sites
			LastVisited string `json:"lastVisited"`
			All         bool   `json:"all"`
		} `json:"privilege"`
		Disaster     int  `json:"disaster"`
		NeedFeedback bool `json:"needFeedback"`
		DefaultSite  bool `json:"defaultSite"`
		ForceModify  bool `json:"forceModify"`
		Dbnormal     bool `json:"dbnormal"`
	} `json:"result"`
}

type Sites struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *Controller) getSiteId(site string) {

	path := "api/v2/users/current"
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.controllerId, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Csrf-Token", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		log.Fatal(err)
	}

	var currentUserResponse currentUserResponse
	if err := json.NewDecoder(res.Body).Decode(&currentUserResponse); err != nil {
		// respBody, _ := ioutil.ReadAll(res.Body)
		// fmt.Println(string(respBody))
		log.Fatal(err)
	}

	var siteId string
	for _, v := range currentUserResponse.Result.Privilege.Sites {
		if v.Name == site {
			siteId = v.Key
		}
	}

	if siteId == "" {
		log.Fatalf("site not found: %s", site)
	}

	c.siteId = siteId

}
