package main

import (
	"detectright"
	"fmt"
	"lib"
	"net/http"
	"os"
	"sync"
)

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

	//req.Header.Get("User-Agent")
	fmt.Println(req.Header)

	//var drc = detectright.DRClient{}
	//fmt.Println("DRClient:", drc)

	drc := detectright.DRClient{
		BaseUrl:           "",
		ActionDetect:      "detect.jsp",
		ActionTestHeaders: "getTestHeader.jsp?",
		ApiKey:            "",
		Properties:        map[string]string{},
		Headers:           map[string]string{},
		Debug:             true,
	}

	fmt.Println("DRClient:", drc)
	drc.LoadConf()

	fmt.Println(req.Header)
	drc.SetHeaders(req.Header)

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

	fmt.Println("[" + tools.DateStampAsString() + "] Logging server started on 0.0.0.0:8088")

	wg.Wait()

}
