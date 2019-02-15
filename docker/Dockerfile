FROM golang:alpine as builder

RUN apk add --update make git gcc musl-dev

ADD . /src
RUN cd /src && \
    go build -i -o /build/linux-amd64/searchQuery cmd/searchQuery/main.go


FROM alpine

RUN apk add --update curl
RUN adduser -S -D -H -h /app searchUser
USER searchUser
COPY --from=builder /build/linux-amd64 /app
COPY --from=builder /src/templates /app/templates
EXPOSE 10001
ENV ELASTIC_URL=http://elasticSearch:9200/
WORKDIR  /app
CMD ["./searchQuery"]