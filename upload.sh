#!/bin/bash

# 检查用户是否传递了参数
if [ "$#" -ne 1 ]; then
    echo "使用方法: $0 <版本号>"
    exit 1
fi

# 获取传递的版本号参数
VERSION=$1

rm -rf ./bin
rm -rf ./chaosd-linux-amd64
make chaosd
make chaos-tools
mv bin chaosd-linux-amd64
tar -czvf chaosd-v"$VERSION"-linux-amd64.tar.gz chaosd-linux-amd64
protokitgo upload --source_files=chaosd-v"$VERSION"-linux-amd64.tar.gz --dest_dsn=oss@/chaosd/@pmt_setting@OSS_8 --config=setting.yaml --no_append_len --no_append_md5