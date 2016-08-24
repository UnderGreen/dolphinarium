FROM debian:jessie
MAINTAINER Sergei Antipov <greendayonfire@gmail.com>

ENV DEBIAN_FRONTEND=noninteractive MYSQL_VERSION=5.6

# Installation of Galera Cluster software
RUN \
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
            ca-certificates \
            --no-install-recommends && \
    sed -i 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    locale-gen && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /var/lib/mysql/*

ADD rootfs/ /

COPY entrypoint.sh /

# Define mountable directories.
VOLUME ["/var/lib/mysql"]
EXPOSE 3306 4567 4568 4444

ENTRYPOINT ["/entrypoint.sh"]
