FROM golang:alpine as builder
RUN apk --no-cache add git ca-certificates
RUN mkdir /src
ADD . /src/
WORKDIR /src/cmd/order_service
RUN CGO_ENABLED=0 GOOS=linux go build -o ../order_service -ldflags '-s -w -extldflags "-static"' ./

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/cmd/order_service/.env /app/.env
COPY --from=builder /src/cmd/order_service/order_service /app/order_service
WORKDIR /app
CMD ["./order_service"]
