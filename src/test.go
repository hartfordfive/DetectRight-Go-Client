package main

import (
	"detectright"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func dateStampAsString() string {
	t := time.Now()
	return ymdToString() + " " + fmt.Sprintf("%02d", t.Hour()) + ":" + fmt.Sprintf("%02d", t.Minute()) + ":" + fmt.Sprintf("%02d", t.Second())
}

func ymdToString() string {
	t := time.Now()
	y, m, d := t.Date()
	return strconv.Itoa(y) + fmt.Sprintf("%02d", m) + fmt.Sprintf("%02d", d)
}

type RequestHandler struct{}

func (rh *RequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		res.WriteHeader(http.StatusNotFound)
		res.Header().Set("Cache-control", "public, max-age=0")
		res.Header().Set("Content-Type", "text/html")
		res.Header().Set("Server", "GoLang Test Webserver")
		fmt.Fprintf(res, "Invalid path")
		return
	}

	// Initialize the DetectRigh Go client
	drc := detectright.InitClient()

	// Store all the headers from the current request in header map
	drcHeaders := make(map[string]interface{})
	for k, v := range req.Header {
		drcHeaders[k] = v[0]
	}

	// Sets the headers of the current rquest
	drc.SetHeaders(drcHeaders)

	// Fetches the device profile from HQ with the collected headers
	drc.GetProfileFromHeaders()
	drc.SetTestHeaders()

	response := map[string]interface{}{
		"headers":  drc.GetHeaders(),
		"response": drc.GetProperties(),
	}

	output, _ := json.Marshal(response)
	fmt.Fprintf(res, string(output))

}

func main() {

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := http.ListenAndServe("0.0.0.0:8088", &RequestHandler{})
		if err != nil {
			fmt.Println("GoLog Error:", err)
			os.Exit(0)
		}
		wg.Done()
	}()

	fmt.Println("[" + dateStampAsString() + "] Logging server started on 0.0.0.0:8088")

	wg.Wait()

}
