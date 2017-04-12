FROM alpine:3.5

MAINTAINER Andrei Varabyeu <andrei_varabyeu@epam.com>

ENV JAVA_VERSION 8u111
ENV JAVA_ALPINE_VERSION 8.121.13-r0

RUN apk add --update --no-cache openjdk8-jre="$JAVA_ALPINE_VERSION" \
    unzip redis ttf-droid supervisor ca-certificates openssl && \
    update-ca-certificates


RUN sed 's/^daemonize yes/daemonize no/' -i /etc/redis.conf

VOLUME /tmp

ENV CONSUL_VERSION="0.8.0"
ENV FABIO_VERSION="1.4.2"
ENV CONSUL_ARCHIVE_FILENAME="consul_${CONSUL_VERSION}_linux_amd64.zip"

RUN wget https://releases.hashicorp.com/consul/$CONSUL_VERSION/$CONSUL_ARCHIVE_FILENAME && \
    unzip $CONSUL_ARCHIVE_FILENAME -d /usr/local/bin/ && \
    rm -rf $CONSUL_ARCHIVE_FILENAME && \
    wget -O /usr/local/bin/fabio https://github.com/eBay/fabio/releases/download/v${FABIO_VERSION}/fabio-${FABIO_VERSION}-go1.8.1-linux_amd64 && \
    chmod +x /usr/local/bin/fabio

## OK. We're all set. Let's install RP binaries
ENV REPO_URL="https://dl.bintray.com/epam/reportportal/com/epam/reportportal"
ENV SERVICE_API_VERSION="3.0.0-RC1"
ENV SERVICE_UI_VERSION="3.0.0-RC1"
ENV SERVICE_JIRA_VERSION="3.0.0-RC1"
ENV SERVICE_RALLY_VERSION="3.0.0-RC1"
ENV SERVICE_AUTHORIZATION_VERSION="3.0.0-RC1"

RUN mkdir /reportportal

RUN wget -O /reportportal/service-api.jar $REPO_URL/service-api/$SERVICE_API_VERSION/service-api-$SERVICE_API_VERSION.jar && \
    wget -O /reportportal/service-jira.jar $REPO_URL/service-jira/$SERVICE_JIRA_VERSION/service-jira-$SERVICE_JIRA_VERSION.jar && \
    wget -O /reportportal/service-rally.jar $REPO_URL/service-rally/$SERVICE_RALLY_VERSION/service-rally-$SERVICE_RALLY_VERSION.jar && \
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

COPY bin/gorpRoot /reportportal/gorpRoot
COPY bin/gorpUI /reportportal/gorpUI

ENV GATEWAY_PORT=8080
ENV SERVICE_API_PORT=8081
ENV SERVICE_API_JAVA_OPTS="-Xmx2g"

ENV SERVICE_JIRA_PORT=8083
ENV SERVICE_JIRA_JAVA_OPTS="-Xmx256m"

ENV SERVICE_RALLY_PORT=8084
ENV SERVICE_RALLY_JAVA_OPTS="-Xmx256m"

ENV SERVICE_AUTHORIZATION_PORT=9999
ENV SERVICE_AUTHORIZATION_JAVA_OPTS="-Xmx512m"

ENV SERVICE_UI_PORT=8082
ENV SERVICE_ROOT_PORT=8085


ENV JIRA_ENABLED=false
ENV RALLY_ENABLED=false

EXPOSE 8080
EXPOSE 9001
EXPOSE 8500

COPY docker/supervisor-default.conf /etc/supervisor-default.conf
COPY docker/supervisor-services.conf /etc/supervisor.d/supervisor-services.conf


ENTRYPOINT ["/usr/bin/supervisord", "--nodaemon", "--configuration", "/etc/supervisor-default.conf"]