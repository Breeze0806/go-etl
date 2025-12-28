FROM ubuntu:22.04 AS base

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        gnupg \
        software-properties-common \
        git \
        build-essential \
        unzip \
        libaio1 \
        && \
    rm -rf /var/lib/apt/lists/*

FROM base AS builder
ENV GOROOT="/usr/local/go"
ENV PATH="${GOROOT}/bin:${PATH}"
ENV GOPATH="/goproject"
ENV PATH="${GOPATH}/bin:${PATH}"
ENV GO111MODULE="on"

# if you located in China, you can use this proxy to speed up
#ENV GOPROXY="https://goproxy.cn,direct"

# if you located in China and do not want to use DB2 as the data source, you can use this to disable download db2 odbc
#ENV IGNORE_PACKAGES="db2"

ENV GO_WEB_SITE="go.dev"
# if you located in China, you can use golang mirror to speed up
#ENV GO_WEB_SITE="golang.google.cn"

RUN wget -q -O - https://${GO_WEB_SITE}/dl/go1.20.14.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
    && mkdir -p "$GOPATH/src/github.com/Breeze0806/go-etl" "$GOPATH/bin" "$GOPATH/pkg"

WORKDIR $GOPATH/src/github.com/Breeze0806/go-etl
COPY . .

RUN make dependencies \
    && make release \
    && mv go-etl-$(git describe --abbrev=0 --tags)-linux-x86_64.tar.gz go-etl-linux-x86_64.tar.gz

ENTRYPOINT ["tail", "-f", "/dev/null"]

FROM base AS production
ENV LD_LIBRARY_PATH="/usr/local/go-etl/ibmdb/clidriver/lib:/usr/local/go-etl/oracle/instantclient_21_18:/usr/local/go-etl/sqlite3/sqlite-tools-linux-x64-3500200"
RUN mkdir -p /usr/local/go-etl/ibmdb /usr/local/go-etl/oracle /usr/local/go-etl/sqlite3
WORKDIR /usr/local/go-etl
COPY --from=builder /goproject/src/github.com/ibmdb/clidriver ./ibmdb/clidriver
RUN wget https://download.oracle.com/otn_software/linux/instantclient/2118000/instantclient-basiclite-linux.x64-21.18.0.0.0dbru.zip && \
    unzip instantclient-basiclite-linux.x64-21.18.0.0.0dbru.zip -d oracle && \
    rm instantclient-basiclite-linux.x64-21.18.0.0.0dbru.zip
RUN wget https://www.sqlite.org/2025/sqlite-tools-linux-x64-3500200.zip && \
    unzip sqlite-tools-linux-x64-3500200.zip -d sqlite3 &&  \
    rm  sqlite-tools-linux-x64-3500200.zip
COPY --from=builder /goproject/src/github.com/Breeze0806/go-etl/go-etl-linux-x86_64.tar.gz .
RUN tar zxvf go-etl-linux-x86_64.tar.gz \
    && rm -f go-etl-linux-x86_64.tar.gz
ENTRYPOINT ["tail", "-f","/dev/null"]