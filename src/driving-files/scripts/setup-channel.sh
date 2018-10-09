#!/bin/bash

echo
echo " ============================================== "
echo " ==========initialize mychannel========== "
echo " ============================================== "
echo

source scripts/header.sh

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="60"}
COUNTER=1
MAX_RETRY=5
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer/msp/tlscacerts/tlsca.orgorderer-cert.pem

echo_b "Channel name : "$CHANNEL_NAME


ORDERER0_IP=orderer0.orgorderer
ORDERER1_IP=orderer1.orgorderer


verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
                echo_r "================== ERROR !!! FAILED to execute End-2-End Scenario =================="
		echo
   		exit 1
	fi
}

setGlobals () {
    PEER=$1
    ORG=$2

	if [ $ORG -eq 1 ] ; then
        CORE_PEER_LOCALMSPID="Org1MSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1/peers/peer0.org1/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1/users/Admin@org1/msp
        if [ $PEER -eq 0 ]; then
            CORE_PEER_ID=peer0.org1
            CORE_PEER_ADDRESS=peer0.org1:7051
            CORE_PEER_CHAINCODELISTENADDRESS=peer0.org1:7052
        else
            CORE_PEER_ADDRESS=peer1.org1:7051
            CORE_PEER_ID=peer1.org1
            CORE_PEER_CHAINCODELISTENADDRESS=peer1.org1:7052
        fi
    elif [ $ORG -eq 2 ] ; then
        CORE_PEER_LOCALMSPID="Org2MSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2/peers/peer0.org2/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2/users/Admin@org2/msp
        if [ $PEER -eq 0 ]; then
            CORE_PEER_ADDRESS=peer0.org2:7051
            CORE_PEER_ID=peer0.org2
            CORE_PEER_CHAINCODELISTENADDRESS=peer0.org2:7052
        else
            CORE_PEER_ADDRESS=peer1.org2:7051
            CORE_PEER_ID=peer1.org2
            CORE_PEER_CHAINCODELISTENADDRESS=peer1.org2:7052
        fi
    fi

	env |grep CORE
}

createChannel() {
    PEER=$1
    ORG=$2
	setGlobals $PEER $ORG

	local ORDERER=
    if [ $ORG -eq 1 ]; then
        ORDERER=$ORDERER0_IP:7050
    elif [ $ORG -eq 2 ]; then
        ORDERER=$ORDERER1_IP:7050
    fi


    if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    		peer channel create -o $ORDERER -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx -t $TIMEOUT >&log.txt
    	else
    		peer channel create -o $ORDERER -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -t $TIMEOUT >&log.txt
    	fi
	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo_g "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
}

updateAnchorPeers() {
    PEER=$1
    ORG=$2
    setGlobals $PEER $ORG

    local ORDERER=
    if [ $ORG -eq 1 ]; then
        ORDERER=$ORDERER0_IP:7050
    elif [ $ORG -eq 2 ]; then
        ORDERER=$ORDERER1_IP:7050
    fi

    if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer channel update -o $ORDERER -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
	else
		peer channel update -o $ORDERER -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Anchor peer update failed"
	echo_g "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
	echo
}

## Sometimes Join takes time hence RETRY atleast for 5 times
joinChannelWithRetry () {
	PEER=$1
    ORG=$2
    setGlobals $PEER $ORG

        set -x
    peer channel join -b $CHANNEL_NAME.block  >&log.txt
    res=$?
        set +x
    cat log.txt
    if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
        COUNTER=` expr $COUNTER + 1`
        echo "peer${PEER}.org${ORG} failed to join the channel, Retry after $DELAY seconds"
        sleep $DELAY
        joinChannelWithRetry $PEER $ORG
    else
        COUNTER=1
    fi
    verifyResult $res "After $MAX_RETRY attempts, peer${PEER}.org${ORG} has failed to Join the Channel"
}

joinChannel () {
	for org in 1 2; do
        for peer in 0; do
            joinChannelWithRetry $peer $org
            echo "===================== peer${peer}.org${org} joined on the channel \"$CHANNEL_NAME\" ===================== "
            sleep $DELAY
            echo
        done
    done
}




## Create channel
echo_b "Creating channel..."
createChannel 0 1
sleep 3

## Join all the peers to the channel
echo_b "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo_b "Updating anchor peers for org1..."
updateAnchorPeers 0 1
sleep 3
echo_b "Updating anchor peers for org2..."
updateAnchorPeers 0 2




echo
echo_g "===================== All GOOD, initialization completed ===================== "
echo

echo
echo " _____   _   _   ____  "
echo "| ____| | \ | | |  _ \ "
echo "|  _|   |  \| | | | | |"
echo "| |___  | |\  | | |_| |"
echo "|_____| |_| \_| |____/ "
echo

exit 0
