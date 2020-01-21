FROM golang:alpine as build

WORKDIR /app
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get golang.org/x/net/proxy
RUN go get -u "h12.io/socks@0ac3745d74c83be82ab1c6d81ca019af810f80af"
COPY . .

RUN go build .
RUN chmod 777 app
RUN ls -al


FROM alpine
RUN apk --no-cache add ca-certificates
ENV PORT 8080

COPY --from=build /app/app /

CMD ["/app"]
