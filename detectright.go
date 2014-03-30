package main

import (
	"fmt"
	"encoding/json"
	"./lib"
)

type DRClient struct {
	baseUrl      string
	actionDetect string
	actionTestHeaders string
	apiKey       string
	Properties   map[string]string
	Headers      map[string]string
}

var drc = DRClient{
	baseUrl:      "",
	actionDetect: "detect.jsp",
	actionTestHeaders: "getTestHeader.jsp?",
	apiKey:       "",
	Properties:   map[string]string{},
	Headers:      map[string]string{},
}


func (drc *DRClient) loadConf() {
     conf:= tools.ParseConfigFile("detectright.conf")
     drc.apiKey = conf["api_key"]
     drc.baseUrl = conf["base_url"]
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

     	if drc.apiKey == "" {
	   conf:= tools.ParseConfigFile("detectright.conf")
	   drc.apiKey = conf["api_key"]
	}

	payload := map[string]string{
		"of": "JSON", // output format
		"k":  drc.apiKey,
	}

	jsonContent,_ := json.Marshal(payload)
	fmt.Println(string(jsonContent))
	url := tools.UrlEncode(drc.baseUrl+drc.actionTestHeaders, payload)
	fmt.Println(url)

	drc.Headers = drc.GetProfile(url)

	return true
}

func (drc *DRClient) GetProperty(propname string) string {

     return "placeholder"
}

func (drc *DRClient) GetProperties() map[string]string {

     res := map[string]string{
     	 "property": "place_holder",
     }
     return res
}

func (drc *DRClient) GetHeaders() map[string]string {

     res := map[string]string{
         "property": "place_holder",
     }
     return res

}

func (drc *DRClient) SetHeadersFromUA(userAgent string) bool {

     return true
}

func (drc *DRClient) GetProfileFromHeaders() map[string]string {

     res := map[string]string{
         "property": "place_holder",
     }
     return res
}

func (drc *DRClient) GetProfile(url string) map[string]string {

     res := map[string]string{}

     if url == "" {
     	return res
     }
     
     properties := drc.GetContentResult(url)
     err := json.Unmarshal(jsonBlob, &animals)
     if err != nil {
     	fmt.Println("error:", err)
	}
     
     fmt.Println(properties)
     

     return res
}


func (drc *DRClient) GetContentResult(url string) string {
     return tools.GetUrlData(url)
}


func main() {

	fmt.Println(drc.GetTestHeaders())
	drc.loadConf()

}
