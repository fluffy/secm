FROM ubuntu:14.04
MAINTAINER Cullen Jennings <fluffy@iii.ca>

RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y emacs24-nox

RUN apt-get install -y autotools-dev
RUN apt-get install -y autoconf
RUN apt-get install -y apache2-dev
RUN apt-get install -y curl
#RUN apt-get install -y jansson
#RUN apt-get install -y pcre3
#RUN apt-get install -y libcurl
RUN apt-get install -y curl libc6 libcurl3 zlib1g
RUN apt-get install -y libjansson-dev
RUN apt-get install -y libpcre3 libpcre3-dev
#RUN apt-get install -y pkgconfig
RUN apt-get install -y pkg-config
#RUN apt-get install -y redis
#RUN apt-get install -y hiredis
RUN apt-get install -y libcurl4-openssl-dev

RUN apt-get install -y apache2

WORKDIR /home
RUN git clone https://github.com/pingidentity/mod_auth_openidc.git
WORKDIR /home/mod_auth_openidc
RUN ./autogen.sh
RUN ./configure
RUN make
RUN make install

EXPOSE 80 443

# the IP of the DB link gets put as a host in /etc/hosts with name db
CMD ["/usr/sbin/apache2ctl","-D","FOREGROUND"]
