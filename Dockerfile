FROM golang:latest AS builder

WORKDIR $GOPATH/src/github.com/lodthe/cpparserbot
COPY . .

RUN go get -d -v .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app ./

ENTRYPOINT ["./app"]