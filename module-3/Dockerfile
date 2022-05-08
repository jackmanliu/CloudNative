FROM golang:1.17 as buildup
COPY ./httpserver.go /httpserver/
WORKDIR /httpserver/
ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go build httpserver.go

FROM busybox
COPY --from=buildup /httpserver/httpserver /httpserver/httpserver
EXPOSE 80
ENTRYPOINT ["/httpserver/httpserver"]