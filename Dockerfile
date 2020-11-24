FROM golang:1.15.5
#add dependencies to current module
RUN go get -d -v -insecure github.com/stedc1976/bash-exporter/cmd
WORKDIR /go/src/github.com/stedc1976/bash-exporter/cmd
#compile packages and dependencies present into the module and write the resulting executable to the bash-exporter file
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bash-exporter .

FROM alpine:3.12
WORKDIR /root
RUN mkdir -p scripts
# copy executable file
COPY --from=0 /go/src/github.com/stedc1976/bash-exporter/cmd/bash-exporter .
# copy script files
COPY ./scripts/*.sh ./scripts/

CMD ["./bash-exporter"]
