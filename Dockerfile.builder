FROM registry.access.redhat.com/ubi8/ubi:8.6

RUN dnf update -y                                      \
  ; dnf install python3 openssl gcc gcc-c++ ca-certificates -y \
  ; dnf clean all

RUN curl -kLfsSo /tmp/go-linux.tar.gz https://golang.org/dl/go1.18.5.linux-amd64.tar.gz  \
  ; rm -rf /usr/local/go ; tar -C /usr/local -xzf /tmp/go-linux.tar.gz

ENV PATH="/usr/local/go/bin:$PATH" \
    GOPATH="/workspace"

WORKDIR /workspace/source

COPY . .

RUN go mod init app; go mod tidy; go build -o ./bin/app
