package detectright

import (
	"./lib"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	DEBUG = true
)

type DRClient struct {
	baseUrl           string
	actionDetect      string
	actionTestHeaders string
	apiKey            string
	properties        map[string]string
	headers           map[string]string
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

func (drc *DRClient) loadConf() {
	conf := tools.ParseConfigFile("detectright.conf")
	drc.apiKey = conf["api_key"]
	drc.baseUrl = conf["base_url"]
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
		conf := tools.ParseConfigFile("detectright.conf")
		drc.apiKey = conf["api_key"]
		drc.baseUrl = conf["base_url"]
	}

	payload := map[string]string{
		"of": "JSON", // output format
		"k":  drc.apiKey,
	}

	jsonContent, _ := json.Marshal(payload)

	if DEBUG {
		fmt.Println("DEBUG: JSON Payload = ", string(jsonContent))
	}
	url := tools.UrlEncode(drc.baseUrl+drc.actionTestHeaders, payload)

	if DEBUG {
		fmt.Println("DEBUG: URL = ", url)
	}

	drc.headers = drc.GetProfile(url)

	return true
}

func (drc *DRClient) GetProperty(propname string) string {
	return "property"
}

func (drc *DRClient) GetProperties() map[string]string {
	return drc.properties
}

func (drc *DRClient) GetHeaders() map[string]string {
	return drc.headers
}


func (drc *DRClient) SetHeaders(headers map[string]string) {
     drc.headers = headers
}


/********** TODO ***************/
func (drc *DRClient) SetHeadersFromUA(userAgent string) bool {

	return true
}

func (drc *DRClient) GetProfileFromHeaders() bool {

	if drc.apiKey == "" {
		conf := tools.ParseConfigFile("detectright.conf")
		drc.apiKey = conf["api_key"]
		drc.baseUrl = conf["base_url"]
	}

	payload := map[string]string{
		"of":  "JSON",
		"if":  "JSON",
		"k":   drc.apiKey,
		"raw": "0",
		"h":   "",
	}

	headers,_ := json.Marshal(drc.headers)
	payload["h"] = base64.StdEncoding.EncodeToString(headers)

	jsonContent, _ := json.Marshal(payload)

	if DEBUG {
		fmt.Println("GetProfileFromHeaders DEBUG: JSON Payload = ", string(jsonContent))
	}
	url := tools.UrlEncode(drc.baseUrl+drc.actionTestHeaders, payload)

	if DEBUG {
		fmt.Println("GetProfileFromHeaders DEBUG: URL = ", url)
	}

	drc.headers = drc.GetProfile(url)

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

func main() {

	fmt.Println(drc.GetTestHeaders())
	drc.loadConf()

}
