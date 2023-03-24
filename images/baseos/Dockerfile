# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG UBUNTU_VER

FROM ubuntu:${UBUNTU_VER} as base
RUN apt update && apt install -y tzdata
RUN addgroup --gid 500 chaincode && adduser --disabled-password --gecos "" --uid 500 --gid 500 --home /home/chaincode chaincode
