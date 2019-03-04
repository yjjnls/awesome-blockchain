'use strict';
var slot = require("./slot.js");
class DPos {
    constructor(blockchain) {
        this.last_slot_ = -1;
        this.block_chain_ = blockchain;
    }

    prepared() {
        let current_slot = slot.get_slot_number();
        let current_id = current_slot % slot.delegates;
        // console.log(current_slot + ' ' + current_id + ' ' + this.block_chain_.get_account_id());

        if (current_id != this.block_chain_.get_account_id())
            return false;

        if (current_slot == this.last_slot_)
            return false;

        this.last_slot_ = current_slot;
        return true;
    }
    make_consensus(block_data) {
        let self = this;
        setImmediate((block) => {
            let time_stamp = block.get_timestamp();
            let block_slot = slot.get_slot_number(time_stamp);
            let target_id = block_slot % slot.delegates;

            let current_slot = slot.get_slot_number();
            let current_id = current_slot % slot.delegates;

            if (target_id != current_id) {
                block.emit('consensus failed');
                return;
            }

            if (target_id != self.block_chain_.get_account_id()) {
                block.emit('consensus failed');
                return;
            }
            block.set_consensus_data({
                "generator_id": self.block_chain_.get_account_id()
            });
            block.emit('consensus completed');
        }, block_data);

    }
    verify(block) {
        let time_stamp = block.timestamp;
        let block_slot = slot.get_slot_number(time_stamp);
        let id = block_slot % slot.delegates;
        if (id != block.consensus_data.generator_id)
            return false;

        return true;
    }
}

module.exports = DPos;
