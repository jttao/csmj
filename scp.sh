#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

#120.77.64.28
IP=wz99qn.hnjinlai.cn
PWD=Jinlai2016!@#

echo "$DIR -> $1"

echo "$PWD"

echo "开始上传>user" 
scp user/main root@120.77.64.28:/game/user/
echo "上传结束>user" 
echo "开始上传>hall" 
scp hall/main root@120.77.64.28:/game/hall/
echo "上传结束>hall" 
echo "开始上传>roommanage" 
scp roommanage/main root@120.77.64.28:/game/roommanage/
echo "上传结束>roommanage" 
echo "开始上传>server" 
scp mahjong/server/main root@120.77.64.28:/game/changsha/
echo "上传结束>server" 
echo "开始上传>gm" 
scp gm/main root@120.77.64.28:/game/gm/
echo "上传结束>gm" 

ssh root@120.77.64.28