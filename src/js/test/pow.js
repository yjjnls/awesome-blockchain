'use strict';

var crypto = require('crypto');
var ed = require('ed25519');
var BlockChain = require("../blockchain");
var Consensus = require("../consensus/pow");
var Promise = require("bluebird");

let blockchain = new BlockChain(Consensus);

// console.log(blockchain.get_last_block().hash);

var password = 'I am tester!';

var hash = crypto.createHash('sha256').update(password).digest();
var keypair = ed.MakeKeypair(hash);

async function create_block(prev_time) {
    return new Promise((resolve, reject) => {
        blockchain.generate_block(keypair, () => {
            console.log(`|${blockchain.get_last_block().hash}|${blockchain.get_last_block().timestamp - prev_time}|`);
            resolve();
        });
    });
}

(async () => {
    for (var i = 0; i < 20; ++i) {
        let prev_time = blockchain.get_last_block().timestamp;
        await create_block(prev_time);
    }
})();

