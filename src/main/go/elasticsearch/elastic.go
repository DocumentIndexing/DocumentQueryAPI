package elasticsearch

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"searchQuery/log"
)

var elasticURL string

func Setup(url string) {
	elasticURL = url
}

func MultiSearch(index string, docType string, request []byte, pretty bool) ([]byte, error) {
	action := "_msearch"
	return post(index, docType, action, request, pretty)
}

func Search(index string, docType string, request []byte, pretty bool) ([]byte, error) {
	action := "_search"
	return post(index, docType, action, request, pretty)
}

func buildContext(index string, docType string) string {
	context := ""
	if len(index) > 0 {
		context = index + "/"
		if len(docType) > 0 {
			context += docType + "/"
		}
	}
	return context
}

func GetStatus() ([]byte, error) {
	return get("_cat", "health", "", nil)
}

func post(index string, docType string, action string, request []byte, pretty bool) ([]byte, error) {
	reader := bytes.NewReader(request)
	prettyUrl := ""
	if pretty {
		prettyUrl = "?pretty=true"
	}
	url := elasticURL + buildContext(index, docType) + action + prettyUrl
	log.Debug("URL", log.Data{"url": url})
	req, err := http.NewRequest("POST", url, reader)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	response, err := ioutil.ReadAll(resp.Body)
	log.Debug("Reader", log.Data{"Data": string(response)})
	return response, err
}

func get(index string, docType string, action string, request []byte) ([]byte, error) {
	reader := bytes.NewReader(request)
	req, err := http.NewRequest("GET", elasticURL+buildContext(index, docType)+action, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)

	return response, err
}
