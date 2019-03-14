'use strict';

var EventEmitter = require('events').EventEmitter;
var Crypto = require("./crypto");


class Block extends EventEmitter {
    constructor(data, consensus) {
        super();
        // body
        this.transactions_ = data ? data.transactions : [];
        // header
        this.version_ = 0;
        this.height_ = data ? data.previous_block.height + 1 : -1;
        this.previous_hash_ = data ? data.previous_block.hash : null;
        this.timestamp_ = (new Date()).getTime();
        this.merkle_hash_ = data ? this.calc_merkle_hash(data.transactions) : null;
        this.generator_publickey_ = data ? data.keypair.publicKey.toString('hex') : null;
        this.hash_ = null;
        this.block_signature_ = null;
        // header extension
        this.consensus_data_ = {};

        if (consensus) {
            let self = this;
            setImmediate(() => {
                self.make_proof(consensus, data.keypair);
            });
        }

    }

    get_version() { return this.version_; }
    get_height() { return this.height_; }
    get_hash() { return this.hash_; }
    get_previous_hash() { return this.previous_hash_; }
    get_timestamp() { return this.timestamp_; }
    get_signature() { return this.block_signature_; }
    get_publickey() { return this.generator_publickey_; }
    get_transactions() { return this.transactions_; }
    get_consensus_data() { return this.consensus_data_; }
    set_consensus_data(data) { this.consensus_data_ = data; }
    toObject() {
        let block = {
            "version": this.version_,
            "height": this.height_,
            "previous_hash": this.previous_hash_,
            "timestamp": this.timestamp_,
            "merkle_hash": this.merkle_hash_,
            "generator_publickey": this.generator_publickey_,
            "hash": this.hash_,
            "block_signature": this.block_signature_,
            "consensus_data": this.consensus_data_,
            "transactions": this.transactions_
        };
        return block;
    }
    set_data(data) {
        this.version_ = data.version;
        this.height_ = data.height;
        this.previous_hash_ = data.previous_hash;
        this.timestamp_ = data.timestamp;
        this.merkle_hash_ = data.merkle_hash;
        this.generator_publickey_ = data.generator_publickey;
        this.hash_ = data.hash;
        this.block_signature_ = data.block_signature;
        this.consensus_data_ = data.consensus_data;
        this.transactions_ = data.transactions;
    }

    calc_merkle_hash() {
        // calc merkle root hash according to the transactions in the block
        var hashes = [];
        for (var i = 0; i < this.transactions_.length; ++i) {
            hashes.push(Crypto.calc_hash(this.transactions_.toString('utf-8')));
        }
        while (hashes.length > 1) {
            var tmp = [];
            for (i = 0; i < hashes.length / 2; ++i) {
                let data = hashes[i * 2] + hashes[i * 2 + 1];
                tmp.push(Crypto.calc_hash(data));
            }
            if (hashes.length % 2 === 1) {
                tmp.push(hashes[hashes.length - 1]);
            }
            hashes = tmp;
        }
        return hashes[0] ? hashes[0] : null;
    }

    prepare_data() {
        let tx = "";
        for (var i = 0; i < this.transactions_.length; ++i) {
            tx += this.transactions_[i].toString('utf-8');
        }
        let data = this.version_.toString()
            + this.height_.toString()
            + this.previous_hash_
            + this.timestamp_.toString()
            + this.merkle_hash_
            + this.generator_publickey_
            + JSON.stringify(this.consensus_data_)
            + tx;

        return data;
    }
    // calc the hash of the block
    calc_block_hash() {
        return Crypto.calc_hash(this.prepare_data());
    }
    sign(keypair) {
        var hash = this.calc_block_hash();
        return Crypto.sign(keypair, hash);
    }
    make_proof(consensus, keypair) {
        let self = this;
        this.on('consensus completed', () => {
            self.hash_ = self.calc_block_hash();
            self.block_signature_ = self.sign(keypair);
            self.emit('block completed', self.toObject());
        });

        consensus.make_consensus(this);
    }

    static verify_signature(block) {
        var hash = block.hash;
        return Crypto.verify_signature(hash, block.block_signature, block.generator_publickey);
    }

}

module.exports = Block;