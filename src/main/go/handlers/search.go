package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"searchQuery/elasticsearch"
	"searchQuery/log"
	"strconv"
	"text/template"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

type searchRequest struct {
	Term      string
	From      int
	Size      int
	PrintType string
	Highlight bool
	Now       string
}

var searchTemplates *template.Template

func SetupSearch() error {
	//Load the templates once, the main entry point for the templates is search.tmpl. The search.tmpl takes
	//the SearchRequest struct and uses the Request to build up the multi-query queries that is used to query elastic.
	templates, err := template.ParseFiles(
		"templates/search/search.tmpl",
	)

	searchTemplates = templates
	return err
}

func formatMultiQuery(rawQuery []byte) ([]byte, error) {
	//Is minify thread Safe? can I put this as a global?
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]js$"), js.Minify)

	linearQuery, err := m.Bytes("application/js", rawQuery)

	if err != nil {
		return nil, err
	}

	//Put new lines in for ElasticSearch to determine the headers and the queries are detected
	return bytes.Replace(linearQuery, []byte("$$"), []byte("\n"), -1), nil

}

func paramGet(params url.Values, key, defaultValue string) string {
	value := params.Get(key)
	if len(value) < 1 {
		value = defaultValue
	}
	return value
}

func paramGetBool(params url.Values, key string, defaultValue bool) bool {
	val, ok := params[key]

	if !ok {
		return defaultValue
	}

	return val[0] == "true" || val[0] == ""
}

func SearchHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	size, err := strconv.Atoi(paramGet(params, "size", "10"))
	if err != nil {
		log.Debug("Expected number only characters for paramater 'size'",
			log.Data{"Size": paramGet(params, "size", "10"), "Error": err.Error()})
		http.Error(w, "Invalid size paramater", http.StatusBadRequest)
		return
	}
	from, err := strconv.Atoi(paramGet(params, "from", "0"))
	if err != nil {
		log.Debug("Expected number only characters for parameter 'from'",
			log.Data{"From": paramGet(params, "from", "0"), "Error": err.Error()})
		http.Error(w, "Invalid from paramater", http.StatusBadRequest)
		return
	}

	pretty := paramGetBool(params, "pretty", false)

	reqParams := searchRequest{
		Term:      params.Get("term"),
		From:      from,
		Size:      size,
		Highlight: paramGetBool(params, "highlight", true),
		Now:       time.Now().UTC().Format(time.RFC3339),
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Debug("Expected number only characters for parameter 'from'",
			log.Data{"From": paramGet(params, "from", "0"), "Error": err.Error()})
		http.Error(w, "Invalid from paramater", http.StatusBadRequest)
		return
	}
	json.Unmarshal(body, &reqParams)

	var doc bytes.Buffer
	err = searchTemplates.Execute(&doc, reqParams)

	if err != nil {
		log.Debug("Failed to create search from template", log.Data{"Error": err.Error(), "Params": reqParams})
		http.Error(w, "Failed to create query", http.StatusInternalServerError)
		return
	}

	////Put new lines in for ElasticSearch to determine the headers and the queries are detected
	//formattedQuery, err := formatMultiQuery(doc.Bytes())
	//if err != nil {
	//	log.Debug("Failed to format query for elasticsearch", log.Data{"Error": err.Error()})
	//	http.Error(w, "Failed to create query", http.StatusInternalServerError)
	//	return
	//}

	responseData, err := elasticsearch.Search(params.Get("index"), "", doc.Bytes(), pretty)
	if err != nil {
		log.Debug("Failed to query elasticsearch", log.Data{"Error": err.Error()})
		http.Error(w, "Failed to run search query", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(responseData)
}
