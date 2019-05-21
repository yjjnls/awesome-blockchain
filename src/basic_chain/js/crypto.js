'use strict';

var crypto = require("crypto");
var ed = require("ed25519");

function calc_hash(data) {
    return crypto.createHash('sha256').update(data).digest('hex');
}

function sign(keypair, data) {
    return ed.Sign(Buffer.from(data, 'utf-8'), keypair).toString('hex');
}

function verify_signature(data, signature, publickey) {
    var res = ed.Verify(Buffer.from(data, 'utf-8'), Buffer.from(signature, 'hex'), Buffer.from(publickey, 'hex'));
    return res;
}

module.exports = {
    calc_hash,
    sign,
    verify_signature
};