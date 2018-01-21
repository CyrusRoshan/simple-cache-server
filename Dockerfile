FROM golang:1.9.2

ARG app_path=$GOPATH/src/github.com/CyrusRoshan/simple-cache-server

RUN mkdir -p ${app_path}
ADD . ${app_path}
WORKDIR ${app_path}

EXPOSE 9000

RUN go build main.go
CMD ["./main"]
