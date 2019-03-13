'use strict';

var crypto = require('crypto');
var ed = require('ed25519');

var password = 'I am genesis!';

var hash = crypto.createHash('sha256').update(password).digest();
var keypair = ed.MakeKeypair(hash);

let genesis = {
    "version": 0,
    "height": 1,
    "previous_hash": null,
    "timestamp": 1550049140488,
    "merkle_hash": null,
    "generator_publickey": keypair.publicKey.toString('hex'),
    "hash": null,
    "block_signature": null,
    "consensus_data": {},
    "transactions": []
};

function prepare_data() {
    let tx = "";
    genesis.transactions.forEach(val => {
        tx += val.toString('utf8');
    });
    let data = genesis.version.toString()
        + genesis.height.toString()
        + genesis.previous_hash
        + genesis.timestamp.toString()
        + genesis.merkle_hash
        + genesis.generator_publickey
        + JSON.stringify(genesis.consensus_data)
        + tx;

    return data;
}

function calc_hash(data) {
    return crypto.createHash('sha256').update(data).digest('hex');
}
function calc_block_hash() {
    return calc_hash(prepare_data());
}
function sign(keypair) {
    var hash = calc_block_hash();
    return ed.Sign(Buffer.from(hash, 'utf-8'), keypair).toString('hex');
}


genesis.hash = calc_block_hash();
genesis.block_signature = sign(keypair);
var res = ed.Verify(Buffer.from(genesis.hash, 'utf-8'), Buffer.from(genesis.block_signature, 'hex'), keypair.publicKey);

if (res) {
    console.log(genesis);
    console.log(JSON.stringify(genesis));
}