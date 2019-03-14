'use strict';

var Block = require("./block");
const genesis_block = require("./genesis_block.json");
var Node = require("./network");
var Account = require("./account");
var Transaction = require("./transaction").Transaction;
var TxInput = require("./transaction").TxInput;
var TxOutput = require("./transaction").TxOutput;
var Msg = require("./message");
var MessageType = require("./message").type;
var Promise = require("bluebird");
var level = require("level");
var Crypto = require("./crypto");

var Pbft = require("./consensus/pbft");
let pbft = false;
class BlockChain {
    constructor(Consensus, keypair, id, is_bad = false) {
        // todo
        this.pending_block_ = {};
        this.tx_pool = {};
        // this.chain_ = [];

        this.is_bad_ = is_bad;
        this.pbft_ = new Pbft(this);

        // ///////////////////////////////////////
        this.genesis_block_ = genesis_block;


        this.account_ = new Account(keypair, id);
        this.consensus_ = new Consensus(this);
        this.node_ = null;
    }
    async start() {
        this.db_ = level(`/tmp/data_${this.get_account_id()}`);
        try {
            // load blocks
            let last = await this.db_.get("last_block");
            this.last_block_ = JSON.parse(last);
            console.log(`node: ${this.get_account_id()} last block: ${this.last_block_.height}`);
        } catch (err) {
            // empty chain
            this.last_block_ = genesis_block;
            this.save_last_block();
            console.log(`node: ${this.get_account_id()} empty`);
        }

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

    async save_block(block) {
        if (!block)
            block = this.last_block_;
        // query from db via hash
        // if not exist, write into db, else do nothing
        if (this.pending_block_[block.hash]) {
            delete this.pending_block_[block.hash];
        }
        await this.db_.put(block.hash, JSON.stringify(block));
        await this.db_.put("last_block", JSON.stringify(block));
        // console.log(`save block: ${block.hash} to db`);

        // tx
        if (!block.transactions) {
            return;
        }
        for (var i = 0; i < block.transactions.length; ++i) {
            let tx = block.transactions[i];
            if (this.tx_pool[tx.id]) {
                delete this.tx_pool[tx.id];
                // console.log(`node ${this.get_account_id()} delete tx ${tx.id}`);
            }
            await this.db_.put(tx.id, JSON.stringify(tx));
        }
    }
    async save_last_block() {
        await this.save_block();
    }

    generate_block(keypair, cb) {
        // load transactions
        var tx = [this.create_coinbase()];
        var i = 0;
        for (let key in this.tx_pool) {
            if (i == 10)
                break;
            tx.push(this.tx_pool[key]);
            i++;
            console.log(`node ${this.get_account_id()} load tx ${key}`);
        }
        // create block
        let block = new Block({
            "keypair": keypair,
            "previous_block": this.last_block_,
            "transactions": tx
        }, this.consensus_);
        // make proof of the block/mine
        let self = this;
        block.on('block completed', (data) => {
            if (data.previous_hash == self.last_block_.hash &&
                data.height == self.last_block_.height + 1) {
                // console.log("block completed");
                self.commit_block(data);

                self.broadcast(Msg.block(data));

                if (cb) cb();
            } else {
                // [fork]
                self.process_fork(data);
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
    async get_from_db(hash) {
        // query block with hash value
        try {
            let block_data = await this.db_.get(hash);
            let block = JSON.parse(block_data);
            return block;
        } catch (err) {
            return null;
        }
    }
    async iterator_back(cb, hash) {
        if (!hash) {
            return;
        }
        let block = await this.get_from_db(hash);
        let res = cb(block);
        if (res)
            await this.iterator_back(cb, block.previous_hash);
    }
    async iterator_forward(cb, hash) {
        if (!hash) {
            return;
        }
        let block = await this.get_from_db(hash);
        await this.iterator_forward(cb, block.previous_hash);
        cb(block);
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
    get_public_key() {
        return this.get_account_keypair().publicKey.toString('hex');
    }
    send_msg(socket, data) {
        this.node_.send(socket, data);
    }
    broadcast(data) {
        this.node_.broadcast(data);
    }
    list_peers() {
        return this.node_.list_peers();
    }
    sync() {
        let peers = this.list_peers();
        let index = Math.floor(Math.random() * peers.length);
        let id = peers[index];
        this.send_msg(parseInt(id), Msg.sync({ "id": this.get_account_id() }));
    }
    async verify_transaction(tx) {
        let input_amount = 0;
        for (var i = 0; i < tx.input.length; ++i) {
            let input = tx.input[i];
            // coinbase
            if (input.id == null) {
                // check milestone
                if (tx.output[0].amount == 50) {
                    return true;
                } else {
                    return false;
                }
            }
            let vout = null;
            if (this.tx_pool[input.id]) {
                vout = this.tx.tx_pool[input.id];
            } else {
                vout = await this.get_from_db(input.id);
            }
            if (!vout) {
                // invalid vout
                return false;
            }
            vout = vout.output[input.index];
            let res = Crypto.verify_signature(JSON.stringify(vout), input.ScriptSig, vout.ScriptPubKey);
            if (!res) {
                return false;
            }
            input_amount += vout.amount;
        }
        let output_amount = 0;
        for (i = 0; i < tx.output.length; ++i) {
            output_amount += tx.output[i].amount;
        }
        if (input_amount < output_amount) {
            return false;
        }
        return true;
    }
    // verify the block is valid
    async verify(block) {
        // verify the block signature
        if (!Block.verify_signature(block))
            return false;
        // verify consensus
        if (!this.consensus_.verify(block)) {
            // [fork] slot
            this.save_block(block);
            return false;
        }
        // verify transactions
        let tx = block.transactions;
        if (tx) {
            for (var i = 0; i < tx.length; ++i) {
                try {
                    if (await this.db_.get(tx[i].id)) {
                        // [fork] transaction exists
                        return false;
                    }
                } catch (err) {
                    // nothing
                }
                if (!await this.verify_transaction(tx[i]))
                    return false;
            }
        }
        return true;
    }
    process_fork(block) {
        if (block.previous_hash != this.last_block_.hash &&
            block.height == this.last_block_.height + 1) {
            // [fork] right height and different previous block
            this.save_block(block);

        } else if (block.previous_hash == this.last_block_.hash &&
            block.height == this.last_block_.height &&
            block.hash != this.last_block_.hash) {
            // [fork] same height and same previous block, but different block id
            this.save_block(block);
        }
    }
    async on_data(msg) {
        switch (msg.type) {
            case MessageType.Block:
                {
                    let block = msg.data;
                    // console.log(`node: ${this.get_account_id()} receive block: height ${block.height}`);
                    // check if exist
                    let query = await this.get_from_db(block.hash);
                    if (this.pending_block_[block.hash] || query) {
                        // console.log("block already exists");
                        return;
                    }
                    // verify
                    if (!await this.verify(block)) {
                        // console.log("verify failed");
                        return;
                    }

                    this.pending_block_[block.hash] = block;

                    // add to chain
                    if (block.previous_hash == this.last_block_.hash &&
                        block.height == this.last_block_.height + 1) {
                        // console.log("on block data");
                        this.commit_block(block);
                        // console.log("----------add block");
                    } else {
                        // [fork]
                        this.process_fork(block);
                    }
                    // broadcast
                    this.broadcast(msg);
                }
                break;
            case MessageType.Transaction:
                {
                    // check if exist(pending or in chain) verify, store(into pending) and broadcast
                    let tx = msg.data;
                    if (this.tx_pool[tx.id]) {
                        // already exists
                        return;
                    }
                    this.tx_pool[tx.id] = tx;
                    // verify transaction
                    let res = await this.verify_transaction(tx);
                    if (!res) {
                        delete this.tx_pool[tx.id];
                    } else {
                        // console.log(`node ${this.get_account_id()} store tx ${tx.id}`);
                    }

                    // broadcast
                    this.broadcast(msg);
                }
                break;
            case MessageType.Sync:
                {
                    console.log(`${this.get_account_id()} receive sync info`);
                    let data = msg.data;
                    let id = data.id;
                    if (data.hash) {
                        let block = await this.get_from_db(data.hash);
                        this.send_msg(id, Msg.sync_block({ "id": this.get_account_id(), "block": block }));
                        console.log(`---> ${this.get_account_id()} send sync block: ${block.height}`);

                    } else {
                        this.send_msg(id, Msg.sync_block({ "id": this.get_account_id(), "last_block": this.last_block_ }));
                        console.log(`---> ${this.get_account_id()} send sync last block: ${this.last_block_.height}`);
                    }
                }
                break;
            case MessageType.SyncBlock:
                {
                    let data = msg.data;
                    let id = data.id;
                    let block = null;
                    if (data.hasOwnProperty("last_block")) {
                        block = data.last_block;
                        this.last_block_ = block;
                        console.log(`++++ ${this.get_account_id()} change last block: ${block.height}`);
                    } else {
                        block = data.block;
                    }
                    console.log(`<--- ${this.get_account_id()} receive sync block: ${block.height}\n`);

                    this.save_block(block);
                    let hash = block.previous_hash;
                    let res = null;
                    if (hash) {
                        res = await this.get_from_db(hash);
                    }
                    if (!res) {
                        console.log(`---> ${this.get_account_id()} continue sync hash: ${hash}`);
                        this.send_msg(id, Msg.sync({ "id": this.get_account_id(), "hash": hash }));
                    } else {
                        console.log(`==== ${this.get_account_id()} complete syning!`);
                    }
                }
                break;
            default:
                if (pbft && !this.is_bad_) {
                    this.pbft_.processMessage(msg);
                } else {
                    console.log("unkown msg");
                    console.log(msg);
                }
                break;
        }
    }
    // print() {
    //     // todo chain_
    //     let output = '';
    //     for (var i = 0; i < this.chain_.length; ++i) {
    //         let height = this.chain_[i].height;
    //         let hash = this.chain_[i].hash.substr(0, 6);
    //         let generator_id = this.chain_[i].consensus_data.generator_id;
    //         if (generator_id == undefined) generator_id = null;
    //         output += `(${height}:${hash}:${generator_id}) -> `;
    //     }
    //     console.log(`node: ${this.get_account_id()} ${output}`);
    // }
    // async fork() {
    //     console.log('----------fork----------');
    //     // load transactions
    //     var tx1 = [{
    //         amount: 1000,
    //         recipient: 'bob',
    //         sender: 'alice'
    //     }];
    //     // create block
    //     let block1 = new Block({
    //         "keypair": this.get_account_keypair(),
    //         "previous_block": this.last_block_,
    //         "transactions": tx1
    //     }, this.consensus_);
    //     // make proof of the block/mine
    //     let self = this;
    //     let block_data1 = await new Promise((resolve, reject) => {
    //         block1.on('block completed', (data) => {
    //             if (data.height == self.last_block_.height + 1) {
    //                 resolve(data);
    //             } else {
    //                 reject('block1 failed');
    //             }
    //         });
    //     });

    //     // load transactions
    //     var tx2 = [{
    //         amount: 1000,
    //         recipient: 'cracker',
    //         sender: 'alice'
    //     }];
    //     // create block
    //     let block2 = new Block({
    //         "keypair": this.get_account_keypair(),
    //         "previous_block": this.last_block_,
    //         "transactions": tx2
    //     }, this.consensus_);
    //     let block_data2 = await new Promise((resolve, reject) => {
    //         block2.on('block completed', (data) => {
    //             if (data.height == self.last_block_.height + 1) {
    //                 resolve(data);
    //             } else {
    //                 reject('block2 failed');
    //             }
    //         });
    //     });

    //     var i = 0;
    //     for (var id in this.node_.peers_) {
    //         let socket = this.node_.peers_[id];
    //         if (i % 2 == 0) {
    //             var msg1 = Msg.block(block_data1);
    //             this.node_.send(socket, msg1);
    //         } else {
    //             var msg2 = Msg.block(block_data2);
    //             this.node_.send(socket, msg2);
    //         }
    //         i++;
    //     }
    //     console.log("fork");
    //     this.commit_block(block_data1);
    // }
    create_coinbase() {
        let input = new TxInput(null, -1, `${new Date()} node: ${this.get_account_id()} coinbase tx`);
        let output = new TxOutput(50, this.get_public_key());
        let tx = new Transaction([input], [output]);
        return tx;
    }

    async get_utxo(cb) {
        let publicKey = this.get_public_key();
        let spentTXOs = {};
        await this.iterator_back((block) => {
            let txs = block.transactions;
            // tx
            for (var i = 0; i < txs.length; ++i) {
                let tx = txs[i];
                let transaction_id = tx.id;
                // output
                for (var j = 0; j < tx.output.length; ++j) {
                    let output = tx.output[j];
                    // owns
                    if (output.ScriptPubKey == publicKey) {
                        // not spent
                        if (spentTXOs.hasOwnProperty(transaction_id) &&
                            spentTXOs[transaction_id].hasOwnProperty(j)) {
                            continue;
                        } else {
                            if (!cb(transaction_id, j, output)) return false;
                        }
                    }
                }
                // input
                for (j = 0; j < tx.input.length; ++j) {
                    let input = tx.input[j];
                    // not coinbase
                    if (input.id != null && input.index != -1) {
                        if (!spentTXOs[input.id]) {
                            spentTXOs[input.id] = [];
                        }
                        spentTXOs[input.id].push(input.index);
                    }
                }
            }
            return true;
        },
        this.get_last_block().hash);
    }
    async get_balance() {
        let value = 0;
        await this.get_utxo((transaction_id, index, vout) => {
            value += vout.amount;
            return true;
        });
        return value;
    }
    async create_transaction(to, amount) {
        let value = 0;
        let input = [];
        let output = [];
        let self = this;
        let tx = null;
        await this.get_utxo((transaction_id, index, vout) => {
            value += vout.amount;
            let signature = Crypto.sign(self.get_account_keypair(), JSON.stringify(vout));
            input.push(new TxInput(transaction_id, index, signature));
            if (value >= amount) {
                output.push(new TxOutput(amount, to));
                if (value > amount)
                    output.push(new TxOutput(value - amount, self.get_public_key()));
                tx = new Transaction(input, output);
                // stop
                return false;
            }
            return true;
        });
        if (value < amount) {
            throw new Error("amount is not enough!");
        }
        if (tx == null) {
            throw new Error("create transaction failed!");
        }
        this.tx_pool[tx.id] = tx;
        this.broadcast(Msg.transaction(tx));

        return tx;
    }
}

module.exports = BlockChain;