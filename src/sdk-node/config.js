var defaultVars = require("./defaultVars");
module.exports = {
    peerHost: process.env.PEER_HOST || "localhost:7051",
    eventHost: process.env.EVENT_HOST || "localhost:7053",
    ordererHost: process.env.ORDERER_HOST || "localhost:7050",
    ordererDomain: process.env.ORDERER_DOMAIN || defaultVars.ORDERER_DOMAIN,
    peerDomain: process.env.PEER_DOMAIN || defaultVars.PEER_DOMAIN,
    caServer:
    process.env.CA_HOST ||
    (process.env.NAMESPACE
        ? "ca." + process.env.NAMESPACE + ":7054"
        : "localhost:7054"),
    caServerName: process.env.CA_SERVER_NAME || defaultVars.CA_SERVER_NAME,
    mspID: process.env.MSPID || defaultVars.MSPID,
    anotherUserSecret: "adminpw",
    user: "admin",
    // convert to boolean
    tlsEnabled: process.env.TLS_ENABLED == "true",
    // tlsEnabled: defaultVars.TLS_ENABLED == "true",
    // we use \r\n to put PEM string into process.env, so we have to replace it to newline
    peerPem: (process.env.PEER_PEM || "").replace(/\\r\\n/g, "\r\n") || defaultVars.PEER_PEM,
    ordererPem: (process.env.ORDERER_PEM || "").replace(/\\r\\n/g, "\r\n") || defaultVars.ORDERER_PEM
};

/*

    export CA_HOST=localhost:30500;
    export PEER_HOST=localhost:30501;
    export EVENT_HOST=localhost:30503;
    export ORDERER_HOST=localhost:32001;
    export ORDERER_DOMAIN=orderer1.orgorderer;
    export PEER_DOMAIN=peer0.org1;
    export CA_SERVER_NAME=ca;
    export MSPID=Org1MSP;
    export TLS_ENABLED=true;

node enrollAdmin.js
node registerUser.js -u user9

node query.js -u user9 --channel mychannel --chaincode mycc -m query -a a
node invoke.js -u user9 --channel mychannel --chaincode mycc -m invoke -a a -a b -a 10

 */
