FROM golang:1.20.5-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY pkg ./pkg

RUN go mod tidy

RUN go build -o /thefoodbook ./cmd/thefoodbook

FROM alpine:3.18.2

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

COPY --from=builder /thefoodbook /thefoodbook

USER nonroot

EXPOSE 3000

ENTRYPOINT ["/thefoodbook"]
CMD ["webservice"]

LABEL name="thefoodbook" \
      version="0.0.1" \
      summary="thefoodbook webservice" \
      description="The Food Book webservice"