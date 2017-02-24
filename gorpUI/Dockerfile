FROM alpine:latest

MAINTAINER Andrei Varabyeu <andrei_varabyeu@epam.com>

ADD ./bin/gorpUI /
ADD /gorpUI/public.tar.gz /public/

ENV RP_STATICSPATH=/public

EXPOSE 8080
ENTRYPOINT ["/gorpUI"]