FROM golang:1.19-alpine
WORKDIR /go/src/be-kredit-pintar
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o be-kredit-pintar app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /app/
COPY --from=0 /go/src/be-kredit-pintar/be-kredit-pintar .
COPY --from=0 /go/src/be-kredit-pintar/db/migration ./db/migration
CMD ["./be-service-insurance-auth"]
