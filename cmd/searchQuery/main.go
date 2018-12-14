package main

import (
	"github.com/gorilla/pat"
	"net/http"
	"os"
	"searchQuery/elasticsearch"
	"searchQuery/handlers"
	"searchQuery/log"
	"time"

)

var server *http.Server
var elasticURL = getEnv("ELASTIC_URL", "http://localhost:9200/")
var bindAddr = getEnv("PORT", "10001")

func getEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	if len(envValue) == 0 {
		envValue = defaultValue
	}
	return envValue
}

func main() {
	log.Namespace = "searchQuery"
	healthCheckEndpoint := 	getEnvironmentVariable("HEALTHCHECK_ENDPOINT", "/healthcheck")
	log.Debug("Starting server", log.Data{"Port": bindAddr, "ElasticSearchUrl": elasticURL})

	// Setup libraries and handlers
	elasticsearch.Setup(elasticURL)
	errSearch := handlers.SetupSearch()
	if errSearch != nil {
		log.ErrorC("Failed to setup search templates", errSearch, log.Data{})
	}

	// Setup web handlers for the search query services
	router := pat.New()
	router.Get("/search", handlers.SearchHandler)
	router.Get(healthCheckEndpoint, handlers.HealthCheckHandlerCreator())
	server = &http.Server{
		Addr:         ":" + bindAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.ErrorC("Failed to bind to port address", err, log.Data{"Port": bindAddr})
	}


}
func getEnvironmentVariable(name string, defaultValue string) string {
	environmentValue := os.Getenv(name)
	if environmentValue != "" {
		return environmentValue
	}
	return defaultValue
}
