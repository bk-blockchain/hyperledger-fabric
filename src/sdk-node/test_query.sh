rm latency.txt
rm throughput.txt

ID=$1
LOOP=$2
for i in `seq 1 $LOOP`; do
    node query.js -u user97 --channel mychannel --chaincode mycc1 -m getResultByID  -a "$ID"
    ID=$ID$i 
done