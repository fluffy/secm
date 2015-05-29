FROM ubuntu:14.04
MAINTAINER Cullen Jennings <fluffy@iii.ca>

RUN apt-get update 
RUN apt-get install -y git 
RUN apt-get install -y emacs24-nox

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN apt-get install -y golang
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux

WORKDIR /usr/local/ 

#RUN git clone https://github.com/fluffy/secm.git
COPY . /usr/local/ks

WORKDIR /usr/local/ks

RUN go build ks.go

EXPOSE 8080

# the IP of the DB link gets put as a host in /etc/hosts with name db
CMD ["./ks","db"]




