export CA_HOST=localhost:30500;
export PEER_HOST=localhost:30501;
export EVENT_HOST=localhost:30503;
export ORDERER_HOST=localhost:32001;
export ORDERER_DOMAIN=orderer1.orgorderer;
export PEER_DOMAIN=peer0.org1;
export CA_SERVER_NAME=ca;
export MSPID=Org1MSP;
export TLS_ENABLED=true;


rm latency.txt
rm throughput.txt

ID=$1
LOOP=$2
for i in `seq 1 $LOOP`; do
    node query.js -u user97 --channel mychannel --chaincode mycc1 -m getResultByID  -a "$ID"
    ID=$ID$i 
done