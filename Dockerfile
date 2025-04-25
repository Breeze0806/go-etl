FROM debian:stable-backports AS base

RUN sed -i 's@deb.debian.org@mirrors.aliyun.com@g' /etc/apt/sources.list.d/debian.sources

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

RUN wget -q -O - https://golang.google.cn/dl/go1.20.14.linux-amd64.tar.gz | tar -C /usr/local -xzf -

ENV GOROOT="/usr/local/go"
ENV PATH="${GOROOT}/bin:${PATH}"
ENV GOPATH="/goproject"
ENV PATH="${GOPATH}/bin:${PATH}"
ENV GOPROXY="https://goproxy.cn,direct"
ENV GO111MODULE="on"
ENV IGNORE_PACKAGES="db2"
RUN mkdir -p "$GOPATH/src/github.com/Breeze0806/go-etl" "$GOPATH/bin" "$GOPATH/pkg"
WORKDIR $GOPATH/src/github.com/Breeze0806/go-etl
COPY . .

RUN make dependencies \
    && make release \
    && mv datax-$(git describe --abbrev=0 --tags)-linux-x86_64.tar.gz datax-no-db2-linux-x86_64.tar.gz

ENTRYPOINT ["tail", "-f", "/dev/null"]

FROM base AS production
WORKDIR /opt
COPY --from=builder /goproject/src/github.com/Breeze0806/go-etl/datax-no-db2-linux-x86_64.tar.gz .
RUN tar zxvf datax-no-db2-linux-x86_64.tar.gz
ENTRYPOINT ["tail", "-f","/dev/null"]
