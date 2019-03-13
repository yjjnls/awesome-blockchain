'use strict';

var crypto = require('crypto');
var ed = require('ed25519');
var BlockChain = require("../blockchain");
var Consensus = require("../consensus/dpos");

var password = 'I am tester!';

var hash = crypto.createHash('sha256').update(password).digest();
var keypair = ed.MakeKeypair(hash);

let blockchains = [];
for (var i = 0; i < 20; ++i) {
    let blockchain = new BlockChain(Consensus, keypair, i);
    blockchain.start();
    blockchains.push(blockchain);
}

// setTimeout(() => {
//     for (var i = 0; i < 20; ++i) {
//         console.log(`${i} --> ${blockchains[i].list_peers()}`);
//     }
// }, 3000);

setTimeout(async () => {
    console.log("=================");
    await blockchains[0].iterator_forward((block) => {
        console.log("-----------------");
        console.log(block.height);
        console.log(block.hash);
        return true;
    }, blockchains[0].get_last_block().hash);
}, 5000);