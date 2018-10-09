rm latency.txt
rm throughput.txt

ID=$1
LOOP=$2
node query.js -u user97 --channel mychannel --chaincode mycc1 -m getResultByID  -a "$ID" -l "$LOOP"
