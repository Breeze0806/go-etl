FROM debian:stable-backports AS base

# if you located in China, you can use aliyun mirrors to speed up
#RUN sed -i 's@deb.debian.org@mirrors.aliyun.com@g' /etc/apt/sources.list.d/debian.sources

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        gnupg \
        software-properties-common \
        git \
        build-essential \
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
RUN mkdir -p /usr/local/go-etl/clidriver
WORKDIR /usr/local/go-etl
COPY --from=builder /goproject/src/github.com/ibmdb/clidriver ./clidriver
ENV LD_LIBRARY_PATH="/usr/local/go-etl/clidriver/lib"
COPY --from=builder /goproject/src/github.com/Breeze0806/go-etl/go-etl-linux-x86_64.tar.gz .
RUN tar zxvf go-etl-linux-x86_64.tar.gz \
    && rm -f go-etl-linux-x86_64.tar.gz
ENTRYPOINT ["tail", "-f","/dev/null"]