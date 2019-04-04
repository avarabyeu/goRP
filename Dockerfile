FROM alpine:3.9

RUN apk add --update --no-cache \
      ca-certificates curl

## Take version from Docker Hub during build
ARG dockerTag

RUN echo $dockerTag | awk -F"v" '{ print $2 }' > version
RUN cat ./version

# Install Hugo
RUN version=$(cat version) && curl -L >gorp.tar.gz https://github.com/avarabyeu/goRP/releases/download/$dockerTag/goRP_${version}_linux_amd64.tar.gz \
 && tar -xzvf gorp.tar.gz -C /usr/bin \
 && rm gorp.tar.gz

ENTRYPOINT ["gorp"]
