FROM node:10.4.1-alpine as builder
WORKDIR /app

RUN npm i yarn -g
RUN apk add --no-cache git
COPY ./frontend/ /app

RUN rm -rf ./node_modules && yarn install && yarn replace && yarn dev && yarn build
RUN cat ./node_modules/irisnet-crypto/chains/iris/conf.json

FROM golang:1.10.3-alpine3.7 as go-builder
ENV GOPATH       /root/go
ENV REPO_PATH    $GOPATH/src/github.com/irisnet/explorer/backend
ENV PATH         $GOPATH/bin:$PATH

RUN mkdir -p GOPATH REPO_PATH

COPY ./backend/ $REPO_PATH
WORKDIR $REPO_PATH

RUN apk add --no-cache make git && go get github.com/golang/dep/cmd/dep && dep ensure && make build


FROM alpine:3.7
WORKDIR /app/backend
COPY --from=builder /app/dist/ /app/frontend/dist
COPY --from=go-builder /root/go/src/github.com/irisnet/explorer/backend/build/ /app/backend/
CMD ./irisplorer