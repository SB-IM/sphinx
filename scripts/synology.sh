#!/usr/bin/env bash

# List api info
curl -s \
    --data-urlencode "api=SYNO.API.Info" \
    --data-urlencode "version=1" \
    --data-urlencode "method=query" \
    --data-urlencode "query=all" \
    'localhost:5000/webapi/query.cgi' | jq

# Login
curl -s \
    --data-urlencode "api=SYNO.API.Auth" \
    --data-urlencode "version=3" \
    --data-urlencode "method=login" \
    --data-urlencode "account=admin" \
    --data-urlencode 'passwd=SuperD0CK!@#' \
    --data-urlencode "session=FileStation" \
    --data-urlencode "format=sid" \
    'localhost:5000/webapi/auth.cgi'

# sid: VGXQsiXAI_4vIuHq_joiKM8UDRk-jzQ1vA7cFjDd2jYGBEt81sLAbCL5qRcd1eL6kjAq9-GvJqFrAPfuOMkU_g

# File station info
curl -s \
    --data-urlencode "api=SYNO.FileStation.Info" \
    --data-urlencode "version=2" \
    --data-urlencode "method=get" \
    --data-urlencode "_sid=VGXQsiXAI_4vIuHq_joiKM8UDRk-jzQ1vA7cFjDd2jYGBEt81sLAbCL5qRcd1eL6kjAq9-GvJqFrAPfuOMkU_g" \
    'localhost:5000/webapi/entry.cgi' | jq

# List shared folders
curl -s \
    --data-urlencode "api=SYNO.FileStation.List" \
    --data-urlencode "version=2" \
    --data-urlencode "method=list_share" \
    --data-urlencode 'additional=["real_path","owner","time","size"]' \
    --data-urlencode "_sid=VGXQsiXAI_4vIuHq_joiKM8UDRk-jzQ1vA7cFjDd2jYGBEt81sLAbCL5qRcd1eL6kjAq9-GvJqFrAPfuOMkU_g" \
    'localhost:5000/webapi/entry.cgi' | jq

# Donwload file
curl -s \
    --data-urlencode "api=SYNO.FileStation.Download" \
    --data-urlencode "version=2" \
    --data-urlencode "method=download" \
    --data-urlencode 'path=["/photo/Superdock/202104/20210401115140/flightlog/airctl.log"]' \
    --data-urlencode "mode=download" \
    --data-urlencode "_sid=VGXQsiXAI_4vIuHq_joiKM8UDRk-jzQ1vA7cFjDd2jYGBEt81sLAbCL5qRcd1eL6kjAq9-GvJqFrAPfuOMkU_g" \
    'localhost:5000/webapi/entry.cgi' > log
