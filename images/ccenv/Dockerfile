# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG UBUNTU_VER
FROM golang:${GO_VER}
RUN apt update && apt install -y \
	g++ \
	gcc \
	git \
	musl-dev

RUN mkdir -p /chaincode/output /chaincode/input
RUN addgroup --gid 500 chaincode && adduser --disabled-password --gecos "" --uid 500 --gid 500 --home /home/chaincode chaincode
RUN chown -R chaincode:chaincode /chaincode
USER chaincode
