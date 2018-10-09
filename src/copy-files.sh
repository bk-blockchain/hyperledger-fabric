PEERS=(172.31.33.147   172.31.40.239)
for PEER in "${PEERS[@]}"; do
    scp -i ~/Hyperledger.pem -r ./driving-files/channel-artifacts/* ubuntu@$PEER:/data/driving-files/channel-artifacts/
    scp -i ~/Hyperledger.pem -r ./driving-files/crypto-config ubuntu@$PEER:/data/driving-files/
done

exit 0

MASTERs=(13.229.151.77   13.229.151.77)
for MASTER in "${MASTERs[@]}"; do
    scp -i ~/Hyperledger.pem -r /data/setup-cluster/render/* ubuntu@$MASTER:/data/setup-cluster/render/
done
