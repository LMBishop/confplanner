FROM golang:1.23 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o confplanner

FROM scratch AS production
WORKDIR /app
COPY --from=builder /build/confplanner ./
EXPOSE 4000
CMD ["/app/confplanner"]
