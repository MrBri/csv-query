FROM golang:1.21

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go types.go upload_template.html ./

RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 3333
CMD ["app"]
