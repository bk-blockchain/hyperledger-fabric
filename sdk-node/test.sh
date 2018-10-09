rm latency.txt
rm throughput.txt

ID=$1
LOOP=$2
for i in `seq 1 $LOOP`; do
    node invoke.js -u user97 --channel mychannel --chaincode mycc1 -m initResult  -a "$ID" -a "2222" -a "Nam" -a "3333" -a "toan" -a "4444" -a "hung" -a "Gioi" -a "20172" -l "8" &
    ID=$ID$i 
done