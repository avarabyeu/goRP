FROM ubuntu:16.04

MAINTAINER Andrei Varabyeu <andrei_varabyeu@epam.com>

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
#export DEBIAN_FRONTEND=noninteractive

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 0C49F3730359A14518585931BC711F9BA15703C6 && \
    echo "deb http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.4.list

RUN apt-get update && \
    apt-get -y install software-properties-common && \
    apt-add-repository -y ppa:webupd8team/java && \
    apt-get update && \
    apt-get install -y debconf-utils && \
    echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections && \
    apt-get -y install unzip oracle-java8-installer redis-server fonts-droid-fallback supervisor mongodb-org-server mongodb-org-tools && \
    rm -rf /var/lib/apt/lists/*


RUN sed 's/^daemonize yes/daemonize no/' -i /etc/redis/redis.conf

ENV CONSUL_VERSION=0.7.5
ENV FABIO_VERSION=1.3.8
ENV CONSUL_ARCHIVE_FILENAME=consul_${CONSUL_VERSION}_linux_amd64.zip

RUN wget https://releases.hashicorp.com/consul/$CONSUL_VERSION/$CONSUL_ARCHIVE_FILENAME && \
    unzip $CONSUL_ARCHIVE_FILENAME -d /usr/local/bin/ && \
    rm -rf $CONSUL_ARCHIVE_FILENAME && \
    wget -O /usr/local/bin/fabio https://github.com/eBay/fabio/releases/download/v${FABIO_VERSION}/fabio-${FABIO_VERSION}-go1.8-linux_amd64 && \
    chmod +x /usr/local/bin/fabio

## OK. We're all set. Let's install RP binaries
ENV REPO_URL=https://dl.bintray.com/epam/reportportal/com/epam/reportportal
ENV SERVICE_API_VERSION=2.7.1
ENV SERVICE_UI_VERSION=2.7.0
ENV SERVICE_AUTHORIZATION_VERSION=2.7.0
ENV SERVICE_REGISTRY_VERSION=2.6.0
ENV SERVICE_GATEWAY_VERSION=2.6.0

VOLUME /tmp

RUN mkdir /reportportal
RUN mkdir -p /data/db/

RUN wget -O /reportportal/service-api.jar $REPO_URL/service-api/$SERVICE_API_VERSION/service-api-$SERVICE_API_VERSION.jar && \
    wget -O /reportportal/service-authorization.zip $REPO_URL/service-authorization/$SERVICE_AUTHORIZATION_VERSION/service-authorization-$SERVICE_AUTHORIZATION_VERSION.zip && \
    mkdir /reportportal/service-authorization/ && \
    unzip /reportportal/service-authorization.zip -d /reportportal/service-authorization/ && \
    rm -rf /reportportal/service-authorization.zip && \
    mv /reportportal/service-authorization/service-authorization-$SERVICE_AUTHORIZATION_VERSION.jar /reportportal/service-authorization/service-authorization.jar && \
    wget -O /reportportal/service-ui.jar $REPO_URL/service-ui/$SERVICE_UI_VERSION/service-ui-$SERVICE_UI_VERSION.jar

RUN unzip /reportportal/service-ui.jar -d /reportportal/tmp/ && \
    mv /reportportal/tmp/BOOT-INF/classes/public/ /reportportal/public && \
    rm -f /reportportal/service-ui.jar && \
    rm -rf /reportportal/tmp


COPY docker/supervisor.conf /etc/supervisor.conf

COPY bin/gorpRoot /reportportal/gorpRoot
COPY bin/gorpUI /reportportal/gorpUI

EXPOSE 8080
EXPOSE 9001
EXPOSE 8500

ENTRYPOINT ["/usr/bin/supervisord", "--nodaemon", "--configuration", "/etc/supervisor.conf"]