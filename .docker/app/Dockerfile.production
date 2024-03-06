FROM golang:1.22 AS builder
WORKDIR /artifaces
COPY ./app .
ENV GO111MODULE=on
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./main.go && \
    chmod +x /artifaces/main

FROM alpine:latest
ENV TZ=Asia/Tokyo
RUN apk --no-cache add tzdata
WORKDIR /artifaces
COPY --from=builder /artifaces/main .
CMD [ "./main" ]
