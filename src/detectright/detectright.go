package detectright

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/karlseguin/ccache"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	VERSION_MAJOR          int    = 0
	VERSION_MINOR          int    = 3
	VERSION_PATCH          int    = 0
	VERSION_SUFFIX         string = "beta"
	DEFAULT_CACHE_TTL_MINS int    = 5
)

type DRClient struct {
	config            map[string]string
	baseUrl           string
	actionDetect      string
	actionTestHeaders string
	actionAnalytics   string
	apiKey            string
	properties        map[string]interface{}
	headers           map[string]interface{}
	debugMode         int
	localCache        *ccache.Cache
	profileHits       []*PageVisit
	profileHitsBuffer int
	request           *http.Request
	mutex             *sync.Mutex
}

type PageVisit struct {
	TS          int32
	DrBrowserId string
	PageUrl     string
	DrUdid      string
	UA          string
	Referrer    string
}

func getVersion() string {
	return strconv.Itoa(VERSION_MAJOR) + "." + strconv.Itoa(VERSION_MINOR) + "." + strconv.Itoa(VERSION_PATCH) + "-" + VERSION_SUFFIX
}

func (drc *DRClient) parseConfigFile(filePath string) map[string]string {

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
		s, err := drc.readln(r)
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

// Fethches a given url
func (drc *DRClient) getUrlData(url string) string {

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

// Encodes the given query string map as a url-encoded string
func (drc *DRClient) urlEncode(domain string, qsParams map[string]string) string {

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

func (drc *DRClient) readln(r *bufio.Reader) (string, error) {
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

// Initializes the DRClient and loads the config file
func InitClient() *DRClient {

	/* Initialize an instance of the client */
	drc := &DRClient{
		config:            map[string]string{},
		baseUrl:           "",
		actionDetect:      "",
		actionTestHeaders: "",
		actionAnalytics:   "",
		apiKey:            "",
		properties:        map[string]interface{}{},
		headers:           map[string]interface{}{},
		debugMode:         0,
		localCache:        ccache.New(ccache.Configure().MaxItems(16777216).ItemsToPrune(100)),
		profileHits:       []*PageVisit{},
		profileHitsBuffer: 1,
		request:           nil,
		mutex:             &sync.Mutex{},
	}
	/* Then load the config and return the DRClient instance */
	drc.loadConf()
	return drc
}

// Loads the application configuration file
func (drc *DRClient) loadConf() {
	conf := drc.parseConfigFile("detectright.conf")
	drc.apiKey = conf["api_key"]
	drc.baseUrl = conf["base_url"]
	drc.actionDetect = conf["action_detect"]
	drc.actionTestHeaders = conf["action_test_headers"]
	drc.config = conf

	dm, err := strconv.Atoi(conf["debug"])
	if err == nil {
		drc.debugMode = dm
	} else {
		drc.debugMode = 0
	}
}

func (drc *DRClient) GetUA() {
	//return drc.properties["User-Agent"].(string)
}

// Determines if the client is mobile or not
func (drc *DRClient) IsMobile() bool {
	if drc.GetProperty("mobile") == "1" || drc.GetProperty("mobile") == "yes" {
		return true
	}
	return false
}

// Retreives test headers for the DetectRight server and sets them
// as the headers for the current device profile request
func (drc *DRClient) SetTestHeaders() bool {

	if drc.apiKey == "" {
		return false
	}

	payload := map[string]string{
		"of": "JSON", // output format
		"k":  drc.apiKey,
	}

	jsonContent, _ := json.Marshal(payload)
	url := drc.urlEncode(drc.baseUrl+drc.actionTestHeaders, payload)

	if drc.debugMode == 1 {
		fmt.Println("SetTestHeaders URL:\n", url)
		fmt.Println("SetTestHeaders JSON Payload:\n", string(jsonContent))
	}

	drc.headers = drc.getProfile(url)
	if drc.debugMode == 1 {
		fmt.Println("SetTestHeaders Result:\n", drc.GetHeaders(), "\n------------------\n")
	}

	return true
}

// Attempts to retreive the given property from the current
// map of properties.  Will fail and return an empty string in the case
// the property doesn't exist or if the properties haven't been fetched
// from the DetectRight API yet.
func (drc *DRClient) GetProperty(propname string) interface{} {
	prop := drc.properties[propname]
	if prop == nil {
		return ""
	}
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

// Returns all the properties from the current property map
func (drc *DRClient) GetProperties() map[string]interface{} {
	return drc.properties
}

// Returns all headers in the current headers map
func (drc *DRClient) GetHeaders() map[string]interface{} {
	return drc.headers
}

// Sets the header map for the device profile API request
func (drc *DRClient) SetHeaders(headers map[string]interface{}) {
	drc.headers = headers
}

// Sets the headers for the device profile API request to contain only
// the specified user agent
func (drc *DRClient) SetHeadersFromUA(userAgent string) {
	drc.headers = map[string]interface{}{"HTTP_USER_AGENT": userAgent}
}

func (drc *DRClient) ReportAnalyticsToHQ() {

	apiUrl := drc.config["analytics_reporting_url"]

	var pv []string
	var pvt []*PageVisit

	drc.mutex.Lock()
	for _, v := range drc.profileHits {
		ov := *v
		pvt = append(pvt, &ov)
		val, _ := json.Marshal(ov)
		fmt.Println("\tVal:", string(val))
		pv = append(pv, string(val))
	}
	drc.mutex.Unlock()

	// Reset the pagve visits
	if drc.debugMode == 1 {
		fmt.Println("Profile hits:", len(drc.profileHits))
		fmt.Println("Resetting profile hits!")
	}
	drc.profileHits = []*PageVisit{}

	jsonEncoded, _ := json.Marshal(pv)

	data := url.Values{}
	data.Set("data", string(jsonEncoded))

	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("X-DR-API-KEY", drc.apiKey)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if drc.debugMode == 1 {
		fmt.Println("Response Code:", resp.Status)
		fmt.Println("Response Body:\n", string(responseBody))
	}

}

func (drc *DRClient) SetHttpRequest(req *http.Request) {
	drc.request = req
}

func (drc *DRClient) getCookie(cookieName string) string {

	cookie := drc.request.Header.Get("Cookie")
	if cookie != "" {
		cookies := strings.Split(cookie, "; ")
		for i := 0; i < len(cookies); i++ {
			parts := strings.Split(cookies[i], "=")
			if parts[0] == cookieName {
				return parts[1]
			}
		}
	}
	return ""
}

// Attempts to retreive a device profile based on the current headers
func (drc *DRClient) GetProfileFromHeaders() bool {

	// Check if the item is in the local cache, and if so, fetch it
	// and return it immediately
	ua := drc.headers["User-Agent"]
	var u string
	if v, ok := ua.(string); ok {
		u = v
	}

	h := sha1.New()
	h.Write([]byte(ua.(string)))
	uaHash := base64.URLEncoding.EncodeToString(h.Sum(nil)[:])

	cachedProfile := drc.localCache.Get(uaHash)

	if drc.debugMode == 1 {
		fmt.Println("Device UA:", ua)
		fmt.Println("Device Hash:", uaHash)
	}

	// If we obtained a valid cache object, then decode it and return it
	if cachedProfile != nil {

		if drc.debugMode == 1 {
			fmt.Println("Profile object found in cache!")
		}

		cp, _ := cachedProfile.(string)
		cachedOjb := map[string]interface{}{}
		err := json.Unmarshal([]byte(cp), &cachedOjb)
		if err != nil {
			fmt.Println("JSON parsing error:", err)
		}

		drc.properties = cachedOjb
		drc.properties["profile_source"] = "cache"

		// Now check if the buffer is full, if so, then send
		// the analytics requests to HQ and empty the buffer
		intval, _ := strconv.Atoi(drc.config["profile_hits_buffer"])
		if len(drc.profileHits) >= intval {
			go func() {
				drc.ReportAnalyticsToHQ()
			}()
		}

		// Finally, make sure to increment the profile hit count
		// We will use the DR browser ID (id.browser)

		referrer := ""
		_, hasReferrer := drc.headers["Referrer"]
		if hasReferrer {
			referrer = drc.headers["Referrer"].(string)
		}

		browser_id, _ := drc.properties["id.browser"].(string)
		dr_udid := drc.getCookie("udid")

		drc.profileHits = append(drc.profileHits, &PageVisit{
			TS:          int32(time.Now().Unix()),
			DrBrowserId: browser_id,
			PageUrl:     drc.headers["X-Requested-Page"].(string),
			DrUdid:      dr_udid,
			UA:          u,
			Referrer:    referrer,
		})

		return true
	} else if cachedProfile == nil && drc.debugMode == 1 {
		fmt.Println("Profile object not found in cache!")
	}

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
	url := drc.urlEncode(drc.baseUrl+drc.actionDetect, payload)

	if drc.debugMode == 1 {
		fmt.Println("GetProfileFromHeaders URL:\n", url)
		fmt.Println("GetProfileFromHeaders JSON Payload:\n", string(jsonContent), "\n------------------\n")
	}

	profile := drc.getProfile(url)

	// Store in local cache
	jsonObject, _ := json.Marshal(profile)

	cache_ttl, _ := strconv.Atoi(drc.config["cache_ttl"])
	drc.localCache.Set(uaHash, string(jsonObject), time.Minute*time.Duration(cache_ttl))

	drc.properties = profile
	drc.properties["profile_source"] = "api"

	return true

}

func (drc *DRClient) getProfile(url string) map[string]interface{} {

	res := map[string]interface{}{}

	if url == "" {
		return res
	}

	properties := drc.getUrlData(url)

	if drc.debugMode == 1 {
		fmt.Println("getProfile Result:\n", properties, "\n------------------\n")
	}

	err := json.Unmarshal([]byte(properties), &res)
	if err != nil {
		fmt.Println("JSON parsing error:", err)
	}

	return res
}
