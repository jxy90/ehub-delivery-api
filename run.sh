#!/bin/bash
#获取环境变量
target=$APP_TARGET
if [ ${brand_code} != "" ]
    then ./ehub-delivery-api export -t ${target}
else
    ./ehub-delivery-api api-server
fi