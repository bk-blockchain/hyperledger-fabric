
var fs = require('fs');
let peerPem = fs.readFileSync('/data/driving-files/crypto-config/peerOrganizations/org1/peers/peer0.org1/tls/ca.crt');
let ordererPem = fs.readFileSync('/data/driving-files/crypto-config/ordererOrganizations/orgorderer/orderers/orderer0.orgorderer/tls/ca.crt');
module.exports = {
    PEER_PEM: peerPem,
    ORDERER_PEM: ordererPem,
    ORDERER_DOMAIN: "orderer0.orgorderer",
    PEER_DOMAIN: "peer0.org1",
    TLS_ENABLED: "true",
    MSPID: "Org1MSP",
    CA_SERVER_NAME: "ca.org1"
};
