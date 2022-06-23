FROM ubuntu:jammy

RUN apt-get update && apt-get -y install npm unzip

RUN npm install -g sass

ADD v5.2.0-beta1.zip /src/

RUN cd /src; unzip v5.2.0-beta1.zip

WORKDIR /src/bootstrap-5.2.0-beta1/scss/

COPY input.scss .
