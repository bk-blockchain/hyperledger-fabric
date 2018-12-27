PEERS=(172.31.39.152   172.31.38.202)
for PEER in "${PEERS[@]}"; do
    scp -i ~/Hyperledger.pem -r ./driving-files/channel-artifacts/* ubuntu@$PEER:/data/driving-files/channel-artifacts/
    scp -i ~/Hyperledger.pem -r ./driving-files/crypto-config ubuntu@$PEER:/data/driving-files/
done

exit 0
