'use strict';
var Crypto = require("./crypto");

class TxOutput {
    constructor(amount, ScriptPubKey) {
        this.amount_ = amount;
        this.script_pubkey_ = ScriptPubKey;
    }
    toObject() {
        let output = {
            "amount": this.amount_,
            "ScriptPubKey": this.script_pubkey_
        };
        return output;
    }
}

class TxInput {
    constructor(id, index, ScriptSig) {
        this.id_ = id;
        this.index_ = index;
        this.script_sig_ = ScriptSig;
    }
    toObject() {
        let input = {
            "id": this.id_,
            "index": this.index_,
            "ScriptSig": this.script_sig_
        };
        return input;
    }
}

class Transaction {
    constructor(input, output) {
        this.input_ = [];
        for (i = 0; i < input.length; ++i) {
            this.input_.push(input[i].toObject());
        }
        this.output_ = [];
        for (var i = 0; i < output.length; ++i) {
            this.output_.push(output[i].toObject());
        }
        this.id_ = Crypto.calc_hash(JSON.stringify(this.input_) + JSON.stringify(this.output_));
        return this.toObject();
    }
    get_id() { return this.id_; }
    get_input() { return this.input_; }
    get_output() { return this.output_; }
    toObject() {
        let tx = {
            "id": this.id_,
            "input": this.input_,
            "output": this.output_
        };
        return tx;
    }
}

module.exports = {
    TxOutput,
    TxInput,
    Transaction
};