SearchQuery
================

A Go application microservice to provide query functionality on to the Document Search

## /search?

URL Parameters

1. `term` _[0..1]_
2. `size` _[0..1]_ _default_:10
3. `from` _[0..1]_ _default_:0
4. `index` _[0..1]_
5. `type` _[0..*]_
6. `sort` _[0..1]_ _default_: `relevance`, _options_: `relevance`,`release_date`,`release_date_asc`,`first_letter`,`title`
7. `queries` _[0..4]_ _default_: `search` , _options_: `search`,`counts`,`departments`,`featured`

### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable | Default | Description
| -------------------- | ------- | -----------
| BIND_ADDR            | 10001  | The host and port to bind to
| ELASTIC_URL	       | "http://localhost:9200/" | Http url of the ElasticSearch server
| HEALTHCHECK_ENDPOINT | /healthcheck             | endpoint that reports health status

## Releasing
To package up the API uses `make package`

## Deploying
Export the following variables;
* export `ELASTIC_SEARCH_URL` to elastic search url.

