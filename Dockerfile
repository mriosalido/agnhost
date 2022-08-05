FROM golang:1.18-alpine AS builder
WORKDIR /go/src/agnhost
COPY . .
RUN go mod tidy
RUN go build
FROM alpine:3
COPY --from=builder /go/src/agnhost/agnhost /usr/local/bin

# from dnsutils image
# install necessary packages:
# - bind-tools: contains dig, which can used in DNS tests.
# - CoreDNS: used in some DNS tests.
# from hostexec image
# install necessary packages:
# - curl, nc: used by a lot of e2e tests
# - iproute2: includes ss used in NodePort tests
# from iperf image
# install necessary packages: iperf, bash
RUN apk --update add bind-tools curl netcat-openbsd iproute2 iperf bash && rm -rf /var/cache/apk/* \
  && ln -s /usr/bin/iperf /usr/local/bin/iperf \
  && ls -altrh /usr/local/bin/iperf

#ADD https://github.com/coredns/coredns/releases/download/v1.6.2/coredns_1.6.2_linux_BASEARCH.tgz /coredns.tgz
#RUN tar -xzvf /coredns.tgz && rm -f /coredns.tgz

# PORT 80 needed by: test-webserver
# PORT 8080 needed by: netexec, nettest, resource-consumer, resource-consumer-controller
# PORT 8081 needed by: netexec
# PORT 9376 needed by: serve-hostname
# PORT 5000 needed by: grpc-health-checking
EXPOSE 80 8080 8081 9376 5000

# from netexec
RUN mkdir /uploads

# from porter
#ADD porter/localhost.crt localhost.crt
#ADD porter/localhost.key localhost.key

ENTRYPOINT ["/usr/local/bin/agnhost"]
CMD ["pause"]
