
#!/usr/bin/env bash

# Copyright IBM Corp., All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

#comment this if you don't want to use the default CHANNEL_NAME "mychannel"
#to create your channel, mannually run:



echo "Deploying ca"
kubectl create -f setup-cluster/render/org1-ca.yaml
kubectl create -f setup-cluster/render/org2-ca.yaml
sleep 5

echo "Deploying zookeeper"
kubectl create -f setup-cluster/render/zookeeper0-kafka.yaml
kubectl create -f setup-cluster/render/zookeeper1-kafka.yaml
kubectl create -f setup-cluster/render/zookeeper2-kafka.yaml
sleep 5

echo "Deploying Kafka"
kubectl create -f setup-cluster/render/kafka0-kafka.yaml
kubectl create -f setup-cluster/render/kafka1-kafka.yaml
kubectl create -f setup-cluster/render/kafka2-kafka.yaml
kubectl create -f setup-cluster/render/kafka3-kafka.yaml
sleep 5

echo "Deploying orderer"
kubectl create -f setup-cluster/render/orderer0.orgorderer.yaml
kubectl create -f setup-cluster/render/orderer1.orgorderer.yaml
sleep 5

echo "Deploying Peer0 Org1"
kubectl create -f setup-cluster/render/peer0.org1.yaml
sleep 5

echo "Deploying rest of the Peers"
kubectl create -f setup-cluster/render/peer0.org2.yaml
#kubectl create -f setup-cluster/render/peer1.org1.yaml -f setup-cluster/render/peer1.org2.yaml
sleep 5

echo "Deploying Cli"
kubectl create -f setup-cluster/render/org1-cli.yaml -f setup-cluster/render/org2-cli.yaml

echo "**********Deployment done successfully**********"


