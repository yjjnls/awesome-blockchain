'use strict';

var crypto = require('crypto');
var ed = require('ed25519');
var BlockChain = require("../blockchain");
var Consensus = require("../consensus/dpos");
var Promise = require("bluebird");

var password = 'I am tester!';

var hash = crypto.createHash('sha256').update(password).digest();
var keypair = ed.MakeKeypair(hash);

let blockchains = [];
for (var i = 0; i < 20; ++i) {
    let blockchain;
    if (i == 10)
        blockchain = new BlockChain(Consensus, keypair, i, true);
    else
        blockchain = new BlockChain(Consensus, keypair, i);
    blockchain.start();
    blockchains.push(blockchain);
}
// // test1
// setTimeout(() => {
//     for (var i = 0; i < 20; ++i) {
//         console.log(`${i} --> ${blockchains[i].list_peers()}`);
//     }
// }, 3000);

// test2
function print_blockchian() {
    console.log("--------------------------------");
    for (i = 0; i < 20; ++i) {
        blockchains[i].print();
    }
}

setInterval(print_blockchian, 10000);
