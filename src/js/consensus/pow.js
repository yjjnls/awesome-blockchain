'use strict';

class Pow {
    constructor(blockchain) {
        this.difficulty_ = 10000;
        this.nonce_ = 0;
        this.state_ = "Idle";
        this.block_chain_ = blockchain;
    }

    // calc_difficulty(block) {
    //     this.difficulty_ = parseInt("0001000000000000000000000000000000000000000000000000000000000000", 16);
    //     // this.difficulty_ = parseInt("0000100000000000000000000000000000000000000000000000000000000000", 16);
    // }

    // make_consensus(block_data) {
    //     this.state_ = "Busy";
    //     this.calc_difficulty(block_data);

    //     let self = this;
    //     setImmediate((block) => {
    //         self.nonce_ = 0;
    //         while (self.nonce_ < Number.MAX_SAFE_INTEGER) {
    //             block.set_consensus_data({
    //                 "difficulty": self.difficulty_,
    //                 "nonce": self.nonce_
    //             });

    //             var hash = block.calc_block_hash();
    //             if (parseInt(hash, 16) < self.difficulty_) {
    //                 block.emit('consensus completed');
    //                 self.state_ = "Idle";
    //                 break;
    //             } else
    //                 self.nonce_ += 1;
    //         }
    //         block.emit('consensus failed');
    //         self.state_ = "Idle";
    //     }, block_data);

    // }
    // verify(block) {
    //     let hash = block.calc_block_hash();
    //     if (hash != block.get_hash())
    //         return false;

    //     let difficulty = block.get_consensus_data().difficulty;
    //     if (parseInt(hash, 16) < difficulty) {
    //         return true;
    //     }
    //     return false;
    // }

    prepared() {
        return this.state_ == "Idle";
    }
    calc_difficulty(block) {
        // adaption
        let prev = this.block_chain_.get_last_block();
        if (!prev.difficulty) {
            this.difficulty_ = 10000;
            return;
        }
        var time_period = block.get_timestamp() - prev.timestamp;
        if (time_period < 3000) {
            this.difficulty_ = prev.difficulty + 9000;
        } else if (prev.difficulty < 3000) {
            this.difficulty_ = prev.difficulty + 9000;
        } else {
            this.difficulty_ = prev.difficulty - 3000;
        }
    }
    make_consensus(block_data) {
        this.state_ = "Busy";
        this.calc_difficulty(block_data);

        let self = this;
        setImmediate((block) => {
            self.nonce_ = 0;
            let target = Number.MAX_SAFE_INTEGER / self.difficulty_ * 100;
            while (self.nonce_ < Number.MAX_SAFE_INTEGER) {
                block.set_consensus_data({
                    "difficulty": self.difficulty_,
                    "nonce": self.nonce_
                });

                var hash = block.calc_block_hash();
                hash = hash.toString('utf8').substring(0, 16);
                if (parseInt(hash, 16) < target) {
                    block.emit('consensus completed');
                    self.state_ = "Idle";
                    break;
                } else
                    self.nonce_ += 1;
            }
            block.emit('consensus failed');
            self.state_ = "Idle";
        }, block_data);

    }
    verify(block) {
        let hash = block.calc_block_hash();
        if (hash != block.get_hash())
            return false;

        let difficulty = block.get_consensus_data().difficulty;
        let target = Number.MAX_SAFE_INTEGER / difficulty * 100;
        hash = hash.toString('utf8').substring(0, 16);

        if (parseInt(hash, 16) < target) {
            return true;
        }
        return false;
    }
}

module.exports = Pow;

