'use strict';

var Block = require("./block");
const genesis_block = require("./genesis_block.json");
var Node = require("./network");
var Account = require("./account");
var Transcation = require("./transcation");
var Msg = require("./message");
var MessageType = require("./message").type;
var Promise = require("bluebird");

var Pbft = require("./consensus/pbft");
let pbft = true;
class BlockChain {
    constructor(Consensus, keypair, id, is_bad = false) {
        // todo
        this.pending_block_ = {};
        this.chain_ = [];

        this.is_bad_ = is_bad;
        this.pbft_ = new Pbft(this);

        // ///////////////////////////////////////
        this.genesis_block_ = genesis_block;
        this.last_block_ = genesis_block;
        this.save_last_block();

        this.account_ = new Account(keypair, id);
        this.consensus_ = new Consensus(this);
        this.node_ = null;
    }
    start() {
        this.node_ = new Node(this.get_account_id());
        this.node_.on("message", this.on_data.bind(this));
        this.node_.start();
        // start loop
        var self = this;
        setTimeout(function next_loop() {
            self.loop(function () {
                setTimeout(next_loop, 1000);
            });
        }, 5000);
    }
    loop(cb) {
        let self = this;
        if (this.consensus_.prepared()) {
            if (!self.is_bad_) {
                this.generate_block(this.get_account_keypair(), () => {
                    // broadcast block
                    let block = self.get_last_block();
                    console.log(`node: ${self.get_account_id()} generate block! block height: ${block.height} hash: ${block.hash}`);
                });
            } else {
                self.fork();
            }
        }
        cb();
    }

    save_last_block() {
        // query from db via hash
        // if not exist, write into db, else do nothing
        // todo（tx is also need to store?）
        if (this.pending_block_[this.last_block_.hash]) {
            delete this.pending_block_[this.last_block_.hash];
        }
        this.chain_.push(this.last_block_);
    }
    generate_block(keypair, cb) {
        // load transcations
        var tx = [];
        // create block
        let block = new Block({
            "keypair": keypair,
            "previous_block": this.last_block_,
            "transactions": tx
        }, this.consensus_);
        // make proof of the block/mine
        let self = this;
        block.on('block completed', (data) => {
            if (data.height == self.last_block_.height + 1) {
                console.log("block completed");
                self.commit_block(data);

                self.broadcast(Msg.block(data));

                if (cb) cb();
            } else {
                // fork or store into tmp
                console.log('fork');
                // todo
                self.pending_block_[data.hash] = data;

            }
        });
    }
    commit_block(block_data) {
        if (pbft && !this.is_bad_) {
            var block = new Block();
            block.set_data(block_data);
            let self = this;
            block.on('consensus completed', (data) => {
                self.last_block_ = data;
                self.save_last_block();
            });
            this.pbft_.make_consensus(block);

        } else {
            this.last_block_ = block_data;
            this.save_last_block();
        }
    }
    get_height() {
        return this.last_block_.height;
    }
    get_block(hash) {
        // query block with hash value
        // todo
        for (var i = 0; i < this.chain_.length; ++i) {
            if (this.chain_[i] == hash) {
                return this.chain_[i];
            }
        }
        return null;
    }
    get_last_block() {
        return this.last_block_;
    }
    get_genesis_block() {
        return this.generate_block_;
    }
    get_amount() {
        // get the amount of the account
        return this.account_.get_amount();
    }
    get_account_id() {
        // get the node id
        return this.account_.get_id();
    }
    get_account_keypair() {
        return this.account_.get_key();
    }
    broadcast(data) {
        this.node_.broadcast(data);
    }
    list_peers() {
        return this.node_.list_peers();
    }
    // verify the block is valid
    verify(block) {
        // verify the block signature
        if (!Block.verify_signature(block))
            return false;
        // verify consensus
        if (!this.consensus_.verify(block))
            return false;
        // verify transcations
        let tx = block.transcations;
        for (var i = 0; i < tx.length; ++i) {
            // todo (check tx is exist and valid)
            if (!Transcation.verify(tx[i]))
                return false;
        }
        return true;
    }
    on_data(msg) {
        switch (msg.type) {
            case MessageType.Block:
                {
                    let block = msg.data;
                    // console.log(`node: ${this.get_account_id()} receive block: height ${block.height}`);
                    // check if exist
                    if (this.pending_block_[block.hash] || this.get_block(block.hash))
                        return;
                    // verify
                    if (!this.verify(block))
                        return;

                    this.pending_block_[block.hash] = block;

                    // add to chain
                    if (block.height == this.last_block_.height + 1) {
                        // console.log("on block data");
                        this.commit_block(block);
                        // console.log("----------add block");
                    } else {
                        // fork or store into tmp
                        // console.log('fork');
                        // todo
                    }
                    // broadcast
                    this.broadcast(msg);
                }
                break;
            case MessageType.Transcation:
                {
                    // check if exist(pending or in chain) verify, store(into pending) and broadcast
                }
                break;
            default:
                if (pbft && !this.is_bad_) {
                    this.pbft_.processMessage(msg);
                }
                break;
        }
    }
    print() {
        // todo chain_
        let output = '';
        for (var i = 0; i < this.chain_.length; ++i) {
            let height = this.chain_[i].height;
            let hash = this.chain_[i].hash.substr(0, 6);
            let generator_id = this.chain_[i].consensus_data.generator_id;
            if (generator_id == undefined) generator_id = null;
            output += `(${height}:${hash}:${generator_id}) -> `;
        }
        console.log(`node: ${this.get_account_id()} ${output}`);
    }
    async fork() {
        console.log('----------fork----------');
        // load transcations
        var tx1 = [{
            amount: 1000,
            recipient: 'bob',
            sender: 'alice'
        }];
        // create block
        let block1 = new Block({
            "keypair": this.get_account_keypair(),
            "previous_block": this.last_block_,
            "transactions": tx1
        }, this.consensus_);
        // make proof of the block/mine
        let self = this;
        let block_data1 = await new Promise((resolve, reject) => {
            block1.on('block completed', (data) => {
                if (data.height == self.last_block_.height + 1) {
                    resolve(data);
                } else {
                    reject('block1 failed');
                }
            });
        });

        // load transcations
        var tx2 = [{
            amount: 1000,
            recipient: 'cracker',
            sender: 'alice'
        }];
        // create block
        let block2 = new Block({
            "keypair": this.get_account_keypair(),
            "previous_block": this.last_block_,
            "transactions": tx2
        }, this.consensus_);
        let block_data2 = await new Promise((resolve, reject) => {
            block2.on('block completed', (data) => {
                if (data.height == self.last_block_.height + 1) {
                    resolve(data);
                } else {
                    reject('block2 failed');
                }
            });
        });

        var i = 0;
        for (var id in this.node_.peers_) {
            let socket = this.node_.peers_[id];
            if (i % 2 == 0) {
                var msg1 = Msg.block(block_data1);
                this.node_.send(socket, msg1);
            } else {
                var msg2 = Msg.block(block_data2);
                this.node_.send(socket, msg2);
            }
            i++;
        }
        console.log("fork");
        this.commit_block(block_data1);
    }
}

module.exports = BlockChain;