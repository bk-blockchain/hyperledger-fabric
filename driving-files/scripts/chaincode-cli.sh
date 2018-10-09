CHANNEL_NAME="$1"
DELAY="$2"
LANGUAGE="$3"
TIMEOUT="$4"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer/msp/tlscacerts/tlsca.orgorderer-cert.pem

echo "Channel name : "$CHANNEL_NAME

CC_SRC_PATH="github.com/hyperledger/fabric/peer/chaincode/chaincode_example02/"


ORDERER0_IP=orderer0.orgorderer
ORDERER1_IP=orderer1.orgorderer

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

installChaincode () {
	PEER=$1
	ORG=$2
	setGlobals $PEER $ORG
	VERSION=$3
        set -x
	peer chaincode install -n $CHAINCODE_NAME -v ${VERSION} -l ${LANGUAGE} -p ${CC_SRC_PATH} >&log.txt
	res=$?
        set +x
	cat log.txt
	verifyResult $res "Chaincode installation on peer${PEER}.org${ORG} has Failed"
	echo "===================== Chaincode is installed on peer${PEER}.org${ORG} ===================== "
	echo
}

instantiateChaincode () {
	PEER=$1
	ORG=$2
	setGlobals $PEER $ORG
	VERSION=$3

	local ORDERER=
    if [ $ORG -eq 1 ]; then
        ORDERER=$ORDERER0_IP:7050
    elif [ $ORG -eq 2 ]; then
        ORDERER=$ORDERER1_IP:7050
    fi

	# while 'peer chaincode' command can get the orderer endpoint from the peer (if join was successful),
	# lets supply it directly as we know it using the "-o" option
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
                set -x
		peer chaincode instantiate -o $ORDERER -C $CHANNEL_NAME -n $CHAINCODE_NAME -l ${LANGUAGE} -v ${VERSION} -c '{"Args":["init","a","100","b","200"]}' -P "OR	('Org1MSP.member','Org2MSP.member')" >&log.txt
		res=$?
                set +x
	else
                set -x
		peer chaincode instantiate -o $ORDERER --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE_NAME -l ${LANGUAGE} -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "OR	('Org1MSP.member','Org2MSP.member')" >&log.txt
		res=$?
                set +x
	fi
	cat log.txt
	verifyResult $res "Chaincode instantiation on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Instantiation on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' ===================== "
	echo
}


upgradeChaincode () {
	PEER=$1
    ORG=$2
    setGlobals $PEER $ORG
    VERSION=$3

	local ORDERER=
    if [ $ORG -eq 1 ]; then
        ORDERER=$ORDERER0_IP:7050
    elif [ $ORG -eq 2 ]; then
        ORDERER=$ORDERER1_IP:7050
    fi

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer chaincode upgrade -o $ORDERER -C $CHANNEL_NAME -n $CHAINCODE_NAME -v $VERSION  -l ${LANGUAGE} -p ${CC_SRC_PATH} -c  '{"Args":["init","a","100","b","200"]}' -P "OR	('Org1MSP.member','Org2MSP.member')" >&log.txt
    else
		peer chaincode upgrade -o $ORDERER --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE_NAME -v $VERSION  -l ${LANGUAGE} -p ${CC_SRC_PATH} -c '{"Args":["init","a","100","b","200"]}' -P "OR	('Org1MSP.member','Org2MSP.member')"  >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Chaincode upgradation on PEER$PEER on channel '$CHANNEL_NAME' failed"
	echo "===================== Chaincode Upgradation on PEER$PEER on channel '$CHANNEL_NAME' ===================== "
	echo
}

chaincodeInvoke () {
	PEER=$1
    ORG=$2
    setGlobals $PEER $ORG
    DATE_WITH_TIME=`date "+%Y-%m-%d %H:%M:%S.%3N"`
    echo "Started time: " $DATE_WITH_TIME

	local ORDERER=
    if [ $ORG -eq 1 ]; then
        ORDERER=$ORDERER0_IP:7050
    elif [ $ORG -eq 2 ]; then
        ORDERER=$ORDERER1_IP:7050
    fi

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		peer chaincode invoke -o $ORDERER -C $CHANNEL_NAME -n $CHAINCODE_NAME -c $CONSTRUCTOR  >&log.txt
	else
		peer chaincode invoke -o $ORDERER  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE_NAME -c $CONSTRUCTOR >&log.txt
	fi
	res=$?
	cat log.txt
	verifyResult $res "Invoke execution on PEER$PEER failed "
	echo "===================== Invoke transaction on PEER$PEER on channel '$CHANNEL_NAME' ===================== "
	echo
}

chaincodeQuery () {
    PEER=$1
    ORG=$2
    setGlobals $PEER $ORG

    echo "===================== Querying on PEER$PEER on channel '$CHANNEL_NAME'... ===================== "

    peer chaincode query -C $CHANNEL_NAME -n $CHAINCODE_NAME -c $CONSTRUCTOR >&log.txt
    echo
    cat log.txt
}



CHAINCODE_NAME=mycc
VERSION=1.0

installChaincode 0 1 $VERSION

installChaincode 0 2 $VERSION

instantiateChaincode 0 1 $VERSION

sleep 3

CONSTRUCTOR='{"Args":["query","a"]}'
chaincodeQuery 0 1
sleep 3

CONSTRUCTOR='{"Args":["invoke","a","b","10"]}'
chaincodeInvoke 0 1
sleep 3

CONSTRUCTOR='{"Args":["query","a"]}'
chaincodeQuery 0 1

exit 0

#upgradeChaincode 0 1 $VERSION

#CONSTRUCTOR='{"Args":["query","a"]}'
#chaincodeQuery 0 1
sleep 3

#CONSTRUCTOR='{"Args":["invoke","a","b","10"]}'
#chaincodeInvoke 0 1


