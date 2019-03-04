'use strict';

class Account {
    constructor(keypair, id) {
        this.keypair_ = keypair;
        this.id_ = id;
        this.amount_ = 0;
        if (!this.keypair_) { }
        if (!this.id_) { }
    }

    get_id() { return this.id_; }
    get_key() { return this.keypair_; }
    get_amount() { return this.amount_; }
}

module.exports = Account;