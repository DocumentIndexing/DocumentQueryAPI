SearchQuery
================

A Go application microservice to provide query functionality on to the Document Search

## /search?

Search taks a combination of URL and Payload Json parameters

### URL Parameters

1. `term` _[0..1]_
2. `size` _[0..1]_ _default_:10
3. `from` _[0..1]_ _default_:0
4. `index` _[0..1]_
5. `highlight`

### Json Payload (_Optional_)
1.	`term`                string
2.	`from`                int
3.	`size`                int
4.	`printType`           string
5.	`highlight`           bool


### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable | Default | Description
| -------------------- | ------- | -----------
| BIND_ADDR            | 10001  | The host and port to bind to
| ELASTIC_URL	       | "http://127.0.0.1.xip.io/" | Http url of the ElasticSearch server
| HEALTHCHECK_ENDPOINT | /healthcheck             | endpoint that reports health status


## Releasing
To package up the API uses `make package` command, this builds the project inside a docker alpine image.
The docker image is defined in two sections in the `Dockerfile` 
The builder section pull the **golang** alpine image and add the projects source code and compiles the 
binary `searchQuery`.

The second section pulls a fresh base image, the latest **alpine** image and copies the `searchQuery` binary and templates 
from the builder image to the new image.

It exposes the port 10001 and defines the default ELASTIC_URL env as `elasticSearch:9200` as this is what it would 
logically be defined as if configured using **Docker Swarm**



## Deploying
Export the following variables;
* export `ELASTIC_SEARCH_URL` to elastic search url.

