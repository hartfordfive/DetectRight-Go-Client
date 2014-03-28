package main

import (
	"fmt"
	//"net/http"
	//"net/url"
	//"os"
	//"strconv"
	//"strings"
	//"time"
	"encoding/json"
	//"io"
	"./lib"
)

type DRClient struct {
	baseUrl      string
	actionDetect string
	apiKey       string
	Properties   map[string]string
	Headers      map[string]string
}

var drc = DRClient{
	baseUrl:      "http://5.44.238.171:7070/",
	actionDetect: "detect.jsp",
	apiKey:       "",
	Properties:   map[string]string{},
	Headers:      map[string]string{},
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

func (drc *DRClient) GetTestHeaders() map[string]string {

	payload := map[string]string{
		"of": "SERP",
		"k":  drc.apiKey,
	}

	jsonContent,_ := json.Marshal(payload)
	fmt.Println(string(jsonContent))
	qs := tools.UrlEncode( string(jsonContent) )
	fmt.Println(qs)
	return payload
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

     res := map[string]string{
         "property": "place_holder",
     }
     return res
}

func (drc *DRClient) GetContentResult(url string) string {

     return "placeholder"
}


func main() {

	fmt.Println(drc.GetTestHeaders())

}
