#!/bin/sh

kubectl create -f ./kube/node-discovery-server/node-discovery-server.yaml
kubectl create -f ./kube/node-discovery-server/node-discovery-server-service.yaml

kubectl create -f ././kube/ubuntu/ubuntu.yaml
sleep 30

kubectl create -f ./kube/genesis/genesis.yaml
kubectl create -f ./kube/genesis/genesis-service.yaml
sleep 60

kubectl create -f ./kube/node/benchmark-node.yaml
kubectl create -f ./kube/node/benchmark-node-service.yaml
