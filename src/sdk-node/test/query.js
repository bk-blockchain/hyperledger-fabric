'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Chaincode query
 */
var program = require("commander");
var defaultConfig = require("./config");
var path = require('path');

program
    .version("0.1.0")
    .option("-u, --user []", "User id", "user1")
    .option("--name, --channel []", "A channel", "mychannel")
    .option("--chaincode, --chaincode []", "A chaincode", "mycc")
    .option("-m, --method []", "A method", "query")
    .option(
        "-a, --arguments [value]",
        "A repeatable value",
        (val, memo) => memo.push(val) && memo,
        []
    )
    .parse(process.argv);


var store_path = path.join(__dirname, 'hfc-key-store');
const config = Object.assign({}, defaultConfig, {
    channelName: program.channel,
    user: program.user,
    storePath: store_path
});

// console.log("Config:", config);

var controller = require("./controller")(config);
var numLoop = 8;
var request = {
    //targets: let default to the peer assigned to the client
    chaincodeId: program.chaincode,
    fcn: program.method,
    args: program.arguments
};
var timeWait = 1000 / numLoop;
invoke();
async function invoke() {
    for (var i = 0; i < numLoop; i++) {
        await setTimeout(function () {
            program.arguments[0] = program.arguments[0] + "a";
            getTimer(request);
        },timeWait);
    }
}
// function wait(ms) {
//     program.arguments[0] = program.arguments[0] + "a";
//     return new Promise(r => setTimeout(r, ms))
// }

async function getTimer(request) {
    var start = Date.now();
    await getTimeInvoke(request, start);
}

// each method require different certificate of user
function getTimeInvoke(request, start) {
    controller
        .query(program.user, request, start)
        .then(ret => {
            console.log(
        	    "Query results: ",
        	    ret.toString()
		    );
        })
        .catch(err => {
            console.error(err);
        });
}