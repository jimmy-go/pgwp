FROM golang:1.9.3-alpine
# DeGOps 0.0.4

# NOTE: added apk for CGO too.
RUN apk --update --no-cache add curl bash git alpine-sdk util-linux gcc musl-dev
# Install glide
RUN curl https://glide.sh/get | sh

WORKDIR /go/src
