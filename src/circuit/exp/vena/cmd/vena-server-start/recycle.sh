#!/bin/sh
. ../../../env-local.sh
export CGO_CFLAGS="-I$GOCIRCUIT/misc/starter-kit-osx/zookeeper/include"
export CGO_LDFLAGS="$GOCIRCUIT/misc/starter-kit-osx/zookeeper/lib/libzookeeper_mt.a"
killall 4r-sumr
4crossbuild && \
echo localhost | 4deploy && \
echo localhost | 4clear && \
echo "-----------------------------------------------------------------------------------------------" && \
go build -a && \
echo "-----------------------------------------------------------------------------------------------" && \
cd ../sumr-spawn-api && go build -a && cd ../sumr-spawn-shard && \
echo "–––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––" && \
./sumr-spawn-shard -sumr ../../../sumr-local.config -durable /tutorial/sumr && \
echo "–––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––" && \
cd ../sumr-spawn-api && \
./sumr-spawn-api -api ../../../api-local.config -durable /tutorial/sumr
