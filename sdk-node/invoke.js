'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Chaincode Invoke
 */
var program = require("commander");
var defaultConfig = require("./config");
var path = require('path');

program
    .version("0.1.0")
    .option("-u, --user []", "User id", "user1")
    .option("--name, --channel []", "A channel", "mychannel")
    .option("--chaincode, --chaincode []", "A chaincode name", "chaincode_example02")
    .option("--host, --host []", "Host", defaultConfig.peerHost)
    .option("--ehost, --event-host []", "Host", defaultConfig.eventHost)
    .option("--ohost, --orderer-host []", "Host", defaultConfig.ordererHost)
    .option("-m, --method []", "A method", "invoke")
    .option(
        "-a, --arguments [value]",
        "A repeatable value",
        (val, memo) => memo.push(val) && memo,
        []
    )
    .option("-l, --loop []", "Loop", "8")
    .parse(process.argv);

// node invoke.js -u user9 --channel mychannel --chaincode mycc -m invoke -a a -a b -a 10
var store_path = path.join(__dirname, 'hfc-key-store')
const config = Object.assign({}, defaultConfig, {
    channelName: program.channel,
    user: program.user,
    storePath: store_path
});

var controller = require("./controller")(config);
var numLoop = program.loop;
var request = {
    //targets: let default to the peer assigned to the client
    chaincodeId: program.chaincode,
    fcn: program.method,
    args: program.arguments
};
var timeWait = 300 ;
invoke();
async function invoke() {
    for (var i = 0; i < numLoop; i++) {
        getTimer(request);
        await wait(timeWait);
    }
}
function wait(ms) {
    program.arguments[0] = program.arguments[0] + "a";
    return new Promise(r => setTimeout(r, ms))
}

async function getTimer(request) {
    var start = Date.now();
    await getTimeInvoke(request, start);
}

// each method require different certificate of user
function getTimeInvoke(request, start) {
    controller
        .invoke(program.user, request, start)
        .then(results => {
            // console.log(
            //     "Send transaction promise and event listener promise have completed",
            //     results
            // );
        })
        .catch(err => {
            console.error(err);
        });
}

