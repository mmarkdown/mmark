# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY lang/ lang/
COPY mast/ mast/
COPY mparser/ mparser/
COPY render/ render/
COPY rfc/ rfc/

RUN CGO_ENABLED=0 GOOS=linux go build -o /mmark

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian13 AS build-release-stage

WORKDIR /

COPY --from=build-stage /mmark /mmark

WORKDIR /data

USER nonroot:nonroot

ENTRYPOINT ["/mmark"]
