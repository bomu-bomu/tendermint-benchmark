#!/bin/sh

docker build -f Dockerfile-abci -t watcharaphat/abci ..
docker build -f Dockerfile-tendermint -t watcharaphat/tendermint .
docker build -f Dockerfile-api -t watcharaphat/api ..

