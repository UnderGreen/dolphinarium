FROM debian:jessie
MAINTAINER Sergei Antipov <greendayonfire@gmail.com>

ENV DEBIAN_FRONTEND=noninteractive MYSQL_VERSION=5.6

# Installation of Galera Cluster software
RUN \
    apt-get -y update && \
    apt-get install -y software-properties-common --no-install-recommends && \
    echo 'deb http://releases.galeracluster.com/debian jessie main' > /etc/apt/sources.list.d/galera.list && \
    echo "deb http://repo.percona.com/apt jessie main" > /etc/apt/sources.list.d/percona.list && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv BC19DDBA && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv CD2EFD2A && \
    apt-get -y update && \
    apt-get install -y \
            galera-3 \
            galera-arbitrator-3 \
            mysql-wsrep-$MYSQL_VERSION \
            percona-xtrabackup \
            locales \
            curl \
            --no-install-recommends && \
    locale-gen en_US.UTF-8 && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /var/lib/mysql/*

ADD node.cnf /etc/mysql/my.cnf

COPY jq /usr/bin/jq
RUN chmod a+x /usr/bin/jq

COPY entrypoint.sh /

# Define mountable directories.
VOLUME ["/var/lib/mysql"]
EXPOSE 3306 4567 4568

ENTRYPOINT ["/entrypoint.sh"]
