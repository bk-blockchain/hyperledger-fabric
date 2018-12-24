PEERS=(172.31.41.42   172.31.40.157)
for PEER in "${PEERS[@]}"; do
    scp -i ~/Hyperledger.pem -r ./driving-files/channel-artifacts/* ubuntu@$PEER:/data/driving-files/channel-artifacts/
    scp -i ~/Hyperledger.pem -r ./driving-files/crypto-config ubuntu@$PEER:/data/driving-files/
done

exit 0
