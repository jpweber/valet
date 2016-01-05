FROM golang:latest
MAINTAINER Jim Weber

RUN mkdir -p /go/src/valet/conf

COPY *.go /go/src/valet/
COPY ./conf/* /go/src/valet/conf/

RUN cd /go/src/valet ; go build

CMD /go/src/valet/valet -P 8000

EXPOSE 8000

