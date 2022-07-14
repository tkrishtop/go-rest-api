# builder
FROM golang:1.18 as builder

COPY . /go/src/webapp
WORKDIR /go/src/webapp

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/webapp main.go

# application construction
# FROM registry.access.redhat.com/ubi8/ubi:latest
FROM gcr.io/distroless/static-debian11

# Metadata
ARG release_tag=0.0.0
LABEL name="Webapp" \
      maintainer="tkrishtop" \
      summary="Test webapp on Golang." \
      description="Test webapp on Golang description." \
      release=${release_tag}
LABEL quay.expires-after=never

COPY --from=builder /go/bin/webapp /
CMD ["/webapp"]
