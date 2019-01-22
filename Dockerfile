FROM golang:1.9.4
LABEL MAINTAINER="1156143589@qq.com"

WORKDIR /go/src/cos-storager

COPY Gopkg.lock .
COPY Gopkg.toml .
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep \
    && make install \
    && make

CMD [ "bash", "-c", "go/src/cos-storage" ]