#!/bin/bash

source .env

ldflags="-s -w"
ldflags+=" -X 'common.VERSION=$VERSION'"
ldflags+=" -X 'common.ENVIRONMENT=$ENVIRONMENT'"
ldflags+=" -X 'common.MAGIC_KEY=$MAGIC_KEY'"
ldflags+=" -X 'common.MQTT_BROKER=$MQTT_BROKER'"
ldflags+=" -X 'common.MQTT_PORT=$MQTT_PORT'"
ldflags+=" -X 'common.SFTP_SERVER=$SFTP_SERVER'"
ldflags+=" -X 'common.SFTP_PORT=$SFTP_PORT'"
ldflags+=" -X 'common.SFTP_USER=$SFTP_USER'"
ldflags+=" -X 'common.SFTP_PASS=$SFTP_PASS'"
ldflags+=" -X 'common.DATADOG_HOST=$DATADOG_HOST'"
ldflags+=" -X 'common.DATADOG_API_KEY=$DATADOG_API_KEY'"
ldflags+=" -X 'common.API_PORT=$API_PORT'"

cd main

env GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -tags ${ENVIRONMENT} -trimpath -ldflags="$ldflags" -o ../bin/CharlesGo || exit -1
go build -trimpath -tags staging -ldflags="$ldflags" -o ../bin/LinuxGo || exit -1

unset VERSION ENVIRONMENT MAGIC_KEY MQTT_BROKER MQTT_PORT SFTP_SERVER SFTP_PORT SFTP_USER SFTP_PASS DATADOG_HOST DATADOG_API_KEY API_PORT