package detectright

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lib"
	"strconv"
)

const (
	DEBUG                 = true
	VERSION_MAJOR  int    = 0
	VERSION_MINOR  int    = 1
	VERSION_PATCH  int    = 0
	VERSION_SUFFIX string = "beta"
)

type DRClient struct {
	BaseUrl           string
	ActionDetect      string
	ActionTestHeaders string
	ApiKey            string
	Properties        map[string]string
	Headers           map[string]string
	Debug             bool
}

/*
var drc = DRClient{
	baseUrl:           "",
	actionDetect:      "detect.jsp",
	actionTestHeaders: "getTestHeader.jsp?",
	apiKey:            "",
	properties:        map[string]string{},
	headers:           map[string]string{},
}
*/

func getVersion() string {
	return strconv.Itoa(VERSION_MAJOR) + "." + strconv.Itoa(VERSION_MINOR) + "." + strconv.Itoa(VERSION_PATCH) + "-" + VERSION_SUFFIX
}

func (drc *DRClient) LoadConf() {
	conf := tools.ParseConfigFile("detectright.conf")
	drc.ApiKey = conf["api_key"]
	drc.BaseUrl = conf["base_url"]
}

func (drc *DRClient) IsFilled() bool {
	if len(drc.Properties) >= 1 {
		return true
	}
	return false
}

func (drc *DRClient) IsEmpty() bool {
	return len(drc.Properties) == 0
}

func (drc *DRClient) IsReady() bool {
	return len(drc.Headers) > 0
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

	if drc.ApiKey == "" {
		conf := tools.ParseConfigFile("detectright.conf")
		drc.ApiKey = conf["api_key"]
		drc.BaseUrl = conf["base_url"]
	}

	payload := map[string]string{
		"of": "JSON", // output format
		"k":  drc.ApiKey,
	}

	jsonContent, _ := json.Marshal(payload)

	if DEBUG {
		fmt.Println("DEBUG: JSON Payload = ", string(jsonContent))
	}
	url := tools.UrlEncode(drc.BaseUrl+drc.ActionTestHeaders, payload)

	if DEBUG {
		fmt.Println("DEBUG: URL = ", url)
	}

	drc.Headers = drc.GetProfile(url)

	return true
}

func (drc *DRClient) GetProperty(propname string) string {
	return "property"
}

func (drc *DRClient) GetProperties() map[string]string {
	return drc.Properties
}

func (drc *DRClient) GetHeaders() map[string]string {
	return drc.Headers
}

func (drc *DRClient) SetHeaders(headers map[string]string) {
	drc.Headers = headers
}

/********** TODO ***************/
func (drc *DRClient) SetHeadersFromUA(userAgent string) bool {

	return true
}

func (drc *DRClient) GetProfileFromHeaders() bool {

	if drc.ApiKey == "" {
		conf := tools.ParseConfigFile("detectright.conf")
		drc.ApiKey = conf["api_key"]
		drc.BaseUrl = conf["base_url"]
	}

	payload := map[string]string{
		"of":  "JSON",
		"if":  "JSON",
		"k":   drc.ApiKey,
		"raw": "0",
		"h":   "",
	}

	headers, _ := json.Marshal(drc.Headers)
	payload["h"] = base64.StdEncoding.EncodeToString(headers)

	jsonContent, _ := json.Marshal(payload)

	if DEBUG {
		fmt.Println("GetProfileFromHeaders DEBUG: JSON Payload = ", string(jsonContent))
	}
	url := tools.UrlEncode(drc.BaseUrl+drc.ActionTestHeaders, payload)

	if DEBUG {
		fmt.Println("GetProfileFromHeaders DEBUG: URL = ", url)
	}

	drc.Headers = drc.GetProfile(url)

	return true

}

func (drc *DRClient) GetProfile(url string) map[string]string {

	res := map[string]string{}

	if url == "" {
		return res
	}

	properties := drc.GetContentResult(url)

	if DEBUG {
		fmt.Println("GetProfile DEBUG: Result from "+url+"\n", properties)
	}

	err := json.Unmarshal([]byte(properties), &res)
	if err != nil {
		fmt.Println("error:", err)
	}

	return res
}

func (drc *DRClient) GetContentResult(url string) string {
	return tools.GetUrlData(url)
}

/*
func main() {

	fmt.Println(drc.GetTestHeaders())
	drc.loadConf()

}
*/
