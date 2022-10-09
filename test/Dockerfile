FROM postgres:14-alpine
RUN apk add --no-cache git make musl-dev go busybox-suid

ENV POSTGRES_DB=postgres
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=pwd123

ENV DB_USER postgres
ENV DB_NAME testdb
ENV DB_HOST localhost
ENV DB_PORT 5432
ENV DB_PWD pwd123

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN chown -R postgres /go

WORKDIR testdir
RUN chown -R postgres /testdir

COPY . .

RUN echo -e "\
#!/bin/bash \n \
/bin/sh ../testdir/test/test_init.sh & \n \
COMMANDS_PID=\$! \n \
(while kill -0 \$COMMANDS_PID; do sleep 1; done) && pg_ctl stop & \
" >> /docker-entrypoint-initdb.d/entrypoint.sh

RUN chmod +x /docker-entrypoint-initdb.d/entrypoint.sh