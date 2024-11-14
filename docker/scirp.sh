#!/bin/bash

pip3 install dfss -i https://pypi.tuna.tsinghua.edu.cn/simple --upgrade
scripts_dir=$(dirname $(readlink -f "$0"))


target=${1^^}

pushd $scripts_dir

if [ "$target" == "BM1684X" ]; then
    python3 -m dfss --url=open@sophgo.com:sophon-stream/tools/stream-agent/application-stream-1684x.tar.gz
    image_file="application-stream-1684x.tar.gz"
    image_name="application-stream-1684x"
elif [ "$target" == "BM1684" ]; then
    python3 -m dfss --url=open@sophgo.com:sophon-stream/tools/stream-agent/application-stream-1684.tar.gz
    image_file="application-stream-1684.tar.gz"
    image_name="application-stream-1684"
elif [ "$target" == "BM1688" ]; then
    python3 -m dfss --url=open@sophgo.com:sophon-stream/tools/stream-agent/application-stream-1688.tar.gz
    image_file="application-stream-1688.tar.gz"
    image_name="application-stream-1688"
elif [ "$target" == "CV186AH" ]; then
    python3 -m dfss --url=open@sophgo.com:sophon-stream/tools/stream-agent/application-stream-cv186ah.tar.gz
    image_file="application-stream-cv186ah.tar.gz"
    image_name="application-stream-cv186ah"
else
    echo "Unknown TARGET type: $target; Please input BM1684X, BM1684, BM1688 or CV186AH"
    exit 0
fi

docker load -i $image_file

docker run -itd -p 8088:8089 --privileged -v /opt/sophon:/opt/sophon -v /dev:/dev $image_name

popd