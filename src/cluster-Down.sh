
#!/usr/bin/env bash

# Copyright IBM Corp., All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

#This is a one step script to delete all the deployment and
#services executed during the execution of the cluster


echo "Delete ca"
kubectl delete -f setup-cluster/render/org1-ca.yaml
kubectl delete -f setup-cluster/render/org2-ca.yaml
sleep 5

echo "Delete zookeeper"
kubectl delete -f setup-cluster/render/zookeeper0-kafka.yaml
kubectl delete -f setup-cluster/render/zookeeper1-kafka.yaml
kubectl delete -f setup-cluster/render/zookeeper2-kafka.yaml
sleep 5

echo "Delete Kafka"
kubectl delete -f setup-cluster/render/kafka0-kafka.yaml
kubectl delete -f setup-cluster/render/kafka1-kafka.yaml
kubectl delete -f setup-cluster/render/kafka2-kafka.yaml
kubectl delete -f setup-cluster/render/kafka3-kafka.yaml
sleep 5

echo "Delete orderer"
kubectl delete -f setup-cluster/render/orderer0.orgorderer.yaml
kubectl delete -f setup-cluster/render/orderer1.orgorderer.yaml
sleep 5

echo "Delete Peer0 Org1"
kubectl delete -f setup-cluster/render/peer0.org1.yaml
sleep 5

echo "Delete rest of the Peers"
kubectl delete -f setup-cluster/render/peer0.org2.yaml
#kubectl delete -f setup-cluster/render/peer1.org1.yaml -f setup-cluster/render/peer1.org2.yaml
sleep 5

echo "Delete Cli"
kubectl delete -f setup-cluster/render/org1-cli.yaml -f setup-cluster/render/org2-cli.yaml

echo "Delete old configs"

sudo rm driving-files/channel-artifacts/*
sudo rm -rf driving-files/crypto-config/
sudo rm -rf kafka/*
rm setup-cluster/render/*

PEERS=(172.31.33.147   172.31.40.239)
for PEER in "${PEERS[@]}"; do
    ssh -i ~/Hyperledger.pem ubuntu@$PEER 'sudo rm /data/driving-files/channel-artifacts/*'
    ssh -i ~/Hyperledger.pem ubuntu@$PEER 'sudo rm -rf /data/driving-files/crypto-config/'
    ssh -i ~/Hyperledger.pem ubuntu@$PEER 'sudo rm -rf /data/kafka/*'
done

echo "Remove chaincode images"
echo "Remove chaincode containers"

./rm-chaincode.sh

PEERS=(172.31.33.147   172.31.40.239)
for PEER in "${PEERS[@]}"; do
    ssh -i ~/Hyperledger.pem ubuntu@$PEER < ./rm-chaincode.sh
done

echo "CLUSTER  Down Completed"

exit 0

