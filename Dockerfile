
FROM golang:1.20-alpine as build

COPY . /tmp/build
WORKDIR /tmp/build
RUN apk update \
    && apk add --no-cache \
    git
RUN version=$(git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
RUN echo $version
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'k8f/cmd.tversion=${version}'" .

FROM alpine:3.18.2
RUN apk update \
    && apk add --no-cache \
        git \
        python3 \
        py3-pip \
        bash \
        py-pip \
        gcc libffi-dev musl-dev openssl-dev python3-dev \
    && pip3 install --upgrade pip \
    && pip3 install --no-cache-dir \
        awscli \
        azure-cli \
    && rm -rf /var/cache/apk/*

COPY --from=build /tmp/build/k8f .
ENTRYPOINT ["./k8f" ]