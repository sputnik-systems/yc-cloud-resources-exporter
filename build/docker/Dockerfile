FROM golang:1.19.0-buster as build

RUN apt update \
    && apt install -y curl

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -o ./yc-resources-exporter ./cmd/resources
RUN go build -o ./yc-price-exporter ./cmd/pricing


FROM ubuntu

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/yc-resources-exporter /usr/local/bin/
COPY --from=build /app/yc-price-exporter /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/yc-resources-exporter"]
