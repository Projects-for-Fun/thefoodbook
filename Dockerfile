FROM golang:1.20.5-alpine AS builder

WORKDIR /app

# ENV CGO_ENABLED 0
# ENV GOOS linux

#COPY go.mod ./
#RUN go mod tidy
#RUN go mod download

COPY . ./

RUN go mod tidy
#RUN go mod download

# RUN go build -v -x -o /thefoodbook ./cmd/thefoodbook
RUN go build -o /thefoodbook ./cmd/thefoodbook
#RUN go build -a -installsuffix cgo -o /thefoodbook ./cmd/thefoodbook
#RUN go build -a -installsuffix cgo -o /thefoodbook -mod vendor ./cmd/thefoodbook
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /thefoodbook ./cmd/thefoodbook

FROM scratch

COPY --from=builder /thefoodbook /thefoodbook

EXPOSE 3000

ENTRYPOINT ["/thefoodbook"]
CMD ["webservice"]

LABEL name="thefoodbook" \
      version="0.0.1" \
      summary="thefoodbook webservice" \
      description="The Food Book webservice"