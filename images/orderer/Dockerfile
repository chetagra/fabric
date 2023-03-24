# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG UBUNTU_VER
FROM ubuntu:${UBUNTU_VER} as base
RUN apt update && apt install -y tzdata
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

FROM golang:${GO_VER} as golang
RUN apt update && apt install -y \
	gcc \
	musl-dev \
	git \
	bash \
	make;
ADD . /home/chetan/fabric
WORKDIR /home/chetan/fabric

FROM golang as orderer
ARG GO_TAGS
RUN make orderer GO_TAGS=${GO_TAGS}

FROM base
ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger
COPY --from=orderer /home/chetan/fabric/build/bin /usr/local/bin
COPY --from=orderer /home/chetan/fabric/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=orderer /home/chetan/fabric/sampleconfig/orderer.yaml ${FABRIC_CFG_PATH}
COPY --from=orderer /home/chetan/fabric/sampleconfig/configtx.yaml ${FABRIC_CFG_PATH}
EXPOSE 7050
CMD ["orderer"]
