'use strict';

class Account {
    constructor(keypair, id) {
        this.keypair_ = keypair;
        this.id_ = id;
        this.amount_ = 0;
        if (!this.keypair_) {
            // load from file
        }
        if (!this.id_) {
            // load from file
        }
    }

    get_id() { return this.id_; }
    get_key() { return this.keypair_; }
    get_amount() { return this.amount_; }
    set_amount(amount) {this.amount_ = amount;}
}

module.exports = Account;