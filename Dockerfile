FROM golang:1.21 AS build-stage

ARG ENVIRONMENT

COPY . /usr/src

WORKDIR /usr/src
RUN ./build.sh

FROM scratch AS export-stage
COPY --from=build-stage /usr/src/main/CharlesGo .
COPY --from=build-stage /usr/src/main/LinuxGo .

