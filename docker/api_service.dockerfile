FROM golang:alpine as builder
RUN apk --no-cache add git ca-certificates
RUN mkdir /src
ADD . /src/
WORKDIR /src/cmd/api_service
RUN CGO_ENABLED=0 GOOS=linux go build -o ../api_service -ldflags '-s -w -extldflags "-static"' ./

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/cmd/api_service/.env /app/.env
COPY --from=builder /src/cmd/api_service/api_service /app/api_service
WORKDIR /app
CMD ["./api_service"]
