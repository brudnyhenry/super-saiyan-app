FROM golang:1.11 as builder
LABEL maintainer="Adam Placzek <brudny.henry@gmail.com>"

RUN go version

COPY . /go/src/github.com/brudnyhenry/super-saiyan-app
WORKDIR /go/src/github.com/brudnyhenry/super-saiyan-app
RUN  go get github.com/golang/dep/cmd/dep
RUN dep init  ||    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/app .

WORKDIR /go/src/github.com/callicoder/go-docker


FROM scratch
WORKDIR /root/



COPY --from=builder /go/bin/app .
COPY template template
COPY static static

EXPOSE 8080
ENTRYPOINT ["./app"]
