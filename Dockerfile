# parent image
FROM registry.access.redhat.com/ubi8/ubi:latest as builder

RUN dnf update -y                                  \
  ; dnf install openssl gcc git ca-certificates -y \
  ; dnf clean all                                  \
  ; dnf install unzip -y ; curl -kLfsSo /tmp/instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip https://download.oracle.com/otn_software/linux/instantclient/215000/instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip ; mkdir -p /opt/oracle  ; cd /tmp ; unzip instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip -d /opt/oracle

# Install Go
RUN curl -kLfsSo /tmp/go-linux.tar.gz https://golang.org/dl/go1.17.1.linux-amd64.tar.gz  \
  ; rm -rf /usr/local/go ; tar -C /usr/local -xzf /tmp/go-linux.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# workspace directory
WORKDIR /go/src/godrortest

# copy source code
COPY . .

# build executable
#RUN go mod init godrortest; go mod tidy ; go build -o ./bin/godrortest .
RUN go mod tidy ; go build -o ./bin/godrortest .

# multistage image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

# oracle config
RUN microdnf install libaio -y
COPY --from=builder /opt/oracle/instantclient_21_5/ /opt/oracle/instantclient_21_5
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient_21_5

# change workdir
WORKDIR /bin

#copy form BUILDER API compile
COPY --from=builder /go/src/godrortest/bin .
COPY --from=builder /go/src/godrortest/set_machine.sh .

EXPOSE 8080

# set entrypoint
ENTRYPOINT [ "godrortest" ]
