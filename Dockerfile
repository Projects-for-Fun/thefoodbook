FROM golang:1.20.5-alpine AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy

RUN go build -o /thefoodbook ./cmd/thefoodbook

FROM scratch

COPY --from=builder /thefoodbook /thefoodbook

EXPOSE 3000

ENTRYPOINT ["/thefoodbook"]
CMD ["webservice"]

LABEL name="thefoodbook" \
      version="0.0.1" \
      summary="thefoodbook webservice" \
      description="The Food Book webservice"