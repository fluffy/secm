FROM ubuntu:14.04
MAINTAINER Cullen Jennings <fluffy@iii.ca>

RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y emacs24-nox

RUN apt-get install -y autotools-dev
RUN apt-get install -y autoconf
RUN apt-get install -y apache2-dev
RUN apt-get install -y curl
RUN apt-get install -y curl libc6 libcurl3 zlib1g
RUN apt-get install -y libjansson-dev
RUN apt-get install -y libpcre3 libpcre3-dev
RUN apt-get install -y pkg-config
RUN apt-get install -y libcurl4-openssl-dev

RUN apt-get install -y apache2

RUN apt-get install -y tcsh

WORKDIR /home
RUN git clone https://github.com/pingidentity/mod_auth_openidc.git
WORKDIR /home/mod_auth_openidc
RUN ./autogen.sh
RUN ./configure
RUN make
RUN make install

RUN apt-get install -y openssh-server supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir /var/lock/apache2
RUN mkdir /var/run/sshd
RUN chmod 0755 /var/run/sshd
RUN mkdir /var/run/apache2

RUN mkdir /data
RUN cp /etc/apache2/sites-available/000-default.conf /data/site.conf
RUN touch /data/site.cert
RUN touch /data/site-chain.cert
RUN touch /data/site.key
RUN ln -s /data/site.conf /etc/apache2/sites-available/site.conf
RUN a2ensite site
RUN a2dissite 000-default.conf
RUN a2enmod proxy 
RUN a2enmod proxy_http
RUN a2enmod ssl
        
EXPOSE 22 80 443
CMD ["/usr/bin/supervisord"]
