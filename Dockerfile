FROM golang:1.10.1
#add dependencies to current module
RUN go get -d -v github.com/stedc1976/bash-exporter/cmd
WORKDIR /go/src/github.com/stedc1976/bash-exporter/cmd
#compile packages and dependencies present into the module and writes the resulting executable to bash-exporter file
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bash-exporter .

FROM alpine:3.7
WORKDIR /root
RUN mkdir -p scripts
# copy executable file
COPY --from=0 /go/src/github.com/stedc1976/bash-exporter/cmd/bash-exporter .
# copy script files
COPY ./scripts/*.sh ./scripts/
CMD ["./bash-exporter"]
