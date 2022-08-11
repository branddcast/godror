FROM registry.access.redhat.com/ubi8/ubi-minimal:8.6

RUN microdnf clean all \
  ; microdnf install libaio -y \
  ; microdnf install unzip -y \
  ; curl -kLfsSo /tmp/instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip https://download.oracle.com/otn_software/linux/instantclient/215000/instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip \
  ; mkdir -p /opt/oracle \
  ; cd /tmp \
  ; unzip instantclient-basiclite-linux.x64-21.5.0.0.0dbru.zip -d /opt/oracle

COPY /instantclient_21_5 /opt/oracle/instantclient_21_5
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient_21_5

WORKDIR /opt/app
RUN chown 1001 /opt/app; \
    chmod "g+rwX" /opt/app; \
    chown 1001:root /opt/app;
	
USER 1001

EXPOSE 8080

ENTRYPOINT [ "/bin/bash" ]
