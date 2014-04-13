package detectright

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	VERSION_MAJOR  int    = 0
	VERSION_MINOR  int    = 1
	VERSION_PATCH  int    = 0
	VERSION_SUFFIX string = "beta"
)

type DRClient struct {
	baseUrl           string
	actionDetect      string
	actionTestHeaders string
	apiKey            string
	properties        map[string]interface{}
	headers           map[string]interface{}
	debugMode         int
}

/************************************
	       Utility functions
*************************************/
func getVersion() string {
	return strconv.Itoa(VERSION_MAJOR) + "." + strconv.Itoa(VERSION_MINOR) + "." + strconv.Itoa(VERSION_PATCH) + "-" + VERSION_SUFFIX
}

func parseConfigFile(filePath string) map[string]string {

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error! Could not open config file: %v\n", err)
		fmt.Println("")
		os.Exit(0)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	params := map[string]string{}

	for err == nil {
		s, err := readln(r)
		if err != nil {
			break
		}
		if err == nil && s != "" {
			parts := strings.SplitN(s, "=", 2)
			params[parts[0]] = strings.Trim(parts[1], " ")
		}
	}

	return params
}

func getUrlData(url string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Requested-With", "DetecRight Go Client v"+getVersion())

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func urlEncode(domain string, qsParams map[string]string) string {

	var Url *url.URL
	Url, err := url.Parse(domain) // E.g.:  http://www.yourdomain.com
	if err != nil {
		panic("Invalid domain")
	}

	Url.Path += ""
	parameters := url.Values{}

	for k, v := range qsParams {
		parameters.Add(k, v)
	}

	Url.RawQuery = parameters.Encode()
	return Url.String()
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

/************************************
				DetectRight client functions
*************************************/

func InitClient() *DRClient {

	/* Initialize an instance of the client */
	drc := &DRClient{
		baseUrl:           "",
		actionDetect:      "",
		actionTestHeaders: "",
		apiKey:            "",
		properties:        map[string]interface{}{},
		headers:           map[string]interface{}{},
		debugMode:         false,
	}
	/* Then load the config and return the DRClient instance */
	drc.LoadConf()
	return drc
}

func (drc *DRClient) LoadConf() {
	conf := parseConfigFile("detectright.conf")
	drc.apiKey = conf["api_key"]
	drc.baseUrl = conf["base_url"]
	drc.actionDetect = conf["action_detect"]
	drc.actionTestHeaders = conf["action_test_headers"]
	drc.debugMode = strconfig["debug"]
}

func (drc *DRClient) IsFilled() bool {
	if len(drc.properties) >= 1 {
		return true
	}
	return false
}

func (drc *DRClient) IsEmpty() bool {
	return len(drc.properties) == 0
}

func (drc *DRClient) IsReady() bool {
	return len(drc.headers) > 0
}

func (drc *DRClient) Prepare() bool {
	if drc.IsEmpty() == true {
		if drc.IsReady() == false {
			return false
		}
		drc.GetProfileFromHeaders()
	}
	return false
}

func (drc *DRClient) IsMobile() bool {
	if drc.GetProperty("mobile") == "1" || drc.GetProperty("mobile") == "yes" {
		return true
	}
	return false
}

func (drc *DRClient) GetTestHeaders() bool {

	if drc.apiKey == "" {
		return false
	}

	payload := map[string]string{
		"of": "JSON", // output format
		"k":  drc.apiKey,
	}

	jsonContent, _ := json.Marshal(payload)

	if drc.debugMode == 1 {
		fmt.Println("GetTestHeaders URL:\n", url)
		fmt.Println("GetTestHeaders JSON Payload:\n", string(jsonContent), "\n------------------\n")
	}
	url := urlEncode(drc.baseUrl+drc.actionTestHeaders, payload)

	drc.headers = drc.GetProfile(url)

	return true
}

func (drc *DRClient) GetProperty(propname string) interface{} {
	prop := drc.properties[propname]
	switch prop.(type) {
	case string:
		retVal, _ := prop.(string)
		return retVal
	case int:
		retVal, _ := prop.(string)
		return retVal
	case []interface{}:
		retVal, _ := prop.(string)
		return retVal
	}
	retVal, _ := prop.(string)
	return retVal
}

func (drc *DRClient) GetProperties() map[string]interface{} {
	return drc.properties
}

func (drc *DRClient) GetHeaders() map[string]interface{} {
	return drc.headers
}

func (drc *DRClient) SetHeaders(headers map[string]interface{}) {
	drc.headers = headers
}

func (drc *DRClient) SetHeadersFromUA(userAgent string) {
	drc.headers = map[string]interface{}{"HTTP_USER_AGENT": userAgent}
}

func (drc *DRClient) GetProfileFromHeaders() bool {

	if drc.apiKey == "" {
		return false
	}

	payload := map[string]string{
		"of":  "JSON",
		"if":  "JSON",
		"k":   drc.apiKey,
		"raw": "0",
		"h":   "",
	}

	headers, _ := json.Marshal(drc.headers)
	payload["h"] = base64.StdEncoding.EncodeToString(headers)

	jsonContent, _ := json.Marshal(payload)

	if drc.debugMode == 1 {
		fmt.Println("GetProfileFromHeaders URL:\n", url)
		fmt.Println("GetProfileFromHeaders JSON Payload:\n", string(jsonContent), "\n------------------\n")
	}
	url := urlEncode(drc.baseUrl+drc.actionDetect, payload)

	drc.properties = drc.GetProfile(url)

	return true

}

func (drc *DRClient) GetProfile(url string) map[string]interface{} {

	res := map[string]interface{}{}

	if url == "" {
		return res
	}

	properties := drc.GetContentResult(url)

	if drc.debugMode == 1 {
		fmt.Println("GetProfile Result:\n", properties, "\n------------------\n")
	}

	err := json.Unmarshal([]byte(properties), &res)
	if err != nil {
		fmt.Println("JSON parsing error:", err)
	}

	return res
}

func (drc *DRClient) GetContentResult(url string) string {
	return getUrlData(url)
}
