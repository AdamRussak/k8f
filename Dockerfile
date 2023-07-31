
FROM golang:1.20-alpine as build

COPY . /tmp/build
WORKDIR /tmp/build
RUN apk update \
    && apk add --no-cache \
    git
RUN version=$(git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
RUN echo $version
RUN arch=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/) && echo ${arch} && CGO_ENABLED=0 GOOS=linux GOARCH=${arch} go build -ldflags="-X 'k8f/cmd.tversion=${version}'" .

FROM alpine:3.18.2
RUN mkdir -p ~/.aws
COPY --from=build /tmp/build/k8f .
ENTRYPOINT ["./k8f" ]