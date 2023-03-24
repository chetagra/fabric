# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG UBUNTU_VER
FROM golang:${GO_VER} as golang

RUN apt update && apt install -y \
	bash \
	gcc \
	git \
	make \
	musl-dev;

ADD . /home/chetan/fabric
WORKDIR /home/chetan/fabric

FROM golang as tools
RUN make configtxgen configtxlator cryptogen peer discover idemixgen

FROM golang:${GO_VER}
# git is required to support `go list -m`
RUN apt upgrade && apt update
RUN apt install -y \
	bash \
	git \
	jq \
	tzdata;
ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
VOLUME /etc/hyperledger/fabric
COPY --from=tools /home/chetan/fabric/build/bin /usr/local/bin
COPY --from=tools /home/chetan/fabric/sampleconfig ${FABRIC_CFG_PATH}
