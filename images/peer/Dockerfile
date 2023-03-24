# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG UBUNTU_VER

FROM ubuntu:${UBUNTU_VER} as peer-base
RUN apt update && apt install -y tzdata
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

FROM golang:${GO_VER} as golang
RUN apt update && apt install -y \
	bash \
	gcc \
	git \
	make \
	musl-dev
ADD . /home/chetan/fabric
WORKDIR /home/chetan/fabric

FROM golang as peer
ARG GO_TAGS
RUN make peer GO_TAGS=${GO_TAGS}

FROM peer-base
RUN apt update && apt install -y \
	build-essential \ 
	libltdl-dev \
	libssl-dev \  
	automake \
	autoconf \
	libtool \
	git \
	make \
	wget \
	opensc
RUN wget https://dist.opendnssec.org/source/softhsm-2.6.1.tar.gz
RUN tar -xzf softhsm-2.6.1.tar.gz
RUN cd softhsm-2.6.1
RUN ./softhsm-2.6.1/configure
RUN make
RUN make install
ENV SOFTHSM2_CONF /etc/softhsm2.conf
RUN softhsm2-util --init-token --slot 0 --label "fabric" --pin 1234 --so-pin 1234
RUN pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so --label fabric --pin 1234 --login --keypairgen --id 1 --key-type EC:prime256v1
RUN pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so -l --pin 1234 -y privkey -d 01 --set-id 018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5
RUN pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so -l --pin 1234 -y pubkey -d 01 --set-id 018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5
ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger
COPY --from=peer /home/chetan/fabric/build/bin /usr/local/bin
COPY --from=peer /home/chetan/fabric/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=peer /home/chetan/fabric/sampleconfig/core.yaml ${FABRIC_CFG_PATH}
EXPOSE 7051
CMD ["peer","node","start"]
