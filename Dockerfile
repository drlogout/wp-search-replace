FROM wordpress:php7.2-apache

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -y \
    iputils-ping \
    mysql-client \
    vim && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /tmp

# install go
RUN curl -O https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz && \
    tar xf go1.11.1.linux-amd64.tar.gz && \
    mv go /usr/local/

# install wp cli
RUN curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x wp-cli.phar && \
    mv wp-cli.phar /usr/local/bin/wp

# install wp-search-replace
ENV GOPATH=/tmp/go
ENV GOROOT=/usr/local/go
RUN mkdir -p /tmp/go/src/

ADD wp-search-replace /tmp/go/src/wp-search-replace

WORKDIR /tmp/go/src/wp-search-replace

RUN /usr/local/go/bin/go build -o /usr/local/bin/wp-search-replace && \
    rm -rf /tmp/go && \
    rm -rf /usr/local/go

COPY docker-entrypoint.sh /usr/local/bin/
COPY export /usr/local/bin/

WORKDIR /var/www/html

VOLUME [ "/db" ]

ENTRYPOINT ["docker-entrypoint.sh"]
