var slot = require("./slot");
var Msg = require("../message");
var MessageType = require("../message").type;

var PBFT_N = slot.delegates;
var PBFT_F = Math.floor((PBFT_N - 1) / 3);

var State = {
    Idle: 0,
    Prepare: 1,
    Commit: 2,
};


class Pbft {
    constructor(blockchain) {
        this.block_chain_ = blockchain;

        this.pending_block_ = {};
        this.prepare_info_ = null;
        this.commit_infos_ = {};
        this.state_ = State.Idle;

        this.prepare_hash_cache_ = {};
        this.commit_hash_cache_ = {};
        this.current_slot_ = 0;
    }

    make_consensus(block) {
        let time_stamp = block.get_timestamp();
        let block_slot = slot.get_slot_number(time_stamp);
        if (block_slot > this.current_slot_) {
            this.clear_state_();
        }
        this.pending_block_[block.get_hash()] = block;

        if (this.state_ === State.Idle) {
            this.current_slot_ = block_slot;
            this.state_ = State.Prepare;
            this.prepare_info_ = {
                "height": block.get_height(),
                "hash": block.get_hash(),
                "votesNumber": 1,
                "votes": {}
            };
            this.prepare_info_.votes[this.block_chain_.get_account_id()] = true;
            var self = this;
            setImmediate(() => {
                // console.log(`node: ${self.block_chain_.get_account_id()} \t[prepared] broadcast prepare msg to peer: ${self.block_chain_.list_peers()}`);
                self.block_chain_.broadcast(Msg.prepare({
                    "height": block.get_height(),
                    "hash": block.get_hash(),
                    "signer": self.block_chain_.get_account_id()
                }));
            });
        }
    }

    clear_state_() {
        this.state_ = State.Idle;
        this.prepare_info_ = null;
        this.commit_infos_ = {};
        this.pending_block_ = {};
    }

    commit(hash) {
        var block = this.pending_block_[hash];
        block.emit('consensus completed', block.toObject());
        this.clear_state_();
    }

    processMessage(msg) {
        var key = msg.data.hash + ':' + msg.data.height + ':' + msg.data.signer;
        switch (msg.type) {
            case MessageType.Prepare:
                // 如果从未收到过这个prepare消息，则转播出去，否则不处理，防止无限广播
                // hash+height其实就是这个消息的ID（n），而signer就是该消息的副本来源 i
                if (!this.prepare_hash_cache_[key]) {
                    this.prepare_hash_cache_[key] = true;
                    this.block_chain_.broadcast(msg);
                } else {
                    return;
                }
                // 如果当前为prepare状态，且收到的prepare消息是同一个id的消息（关于同一个block的）则这个消息的投票数加1
                // 超过2f+1票后就进入commit状态，发送commit消息
                if (this.state_ === State.Prepare &&
                    msg.data.height === this.prepare_info_.height &&
                    msg.data.hash === this.prepare_info_.hash &&
                    !this.prepare_info_.votes[msg.data.signer]) {
                    this.prepare_info_.votes[msg.data.signer] = true;
                    this.prepare_info_.votesNumber++;
                    if (this.prepare_info_.votesNumber > 2 * PBFT_F) {
                        // console.log(`node: ${this.block_chain_.get_account_id()} \t[commit] broadcast commit msg to peer: ${this.block_chain_.list_peers()}`);
                        this.state_ = State.Commit;
                        var commitInfo = {
                            "height": this.prepare_info_.height,
                            "hash": this.prepare_info_.hash,
                            "votesNumber": 1,
                            "votes": {}
                        };
                        commitInfo.votes[this.block_chain_.get_account_id()] = true;
                        this.commit_infos_[commitInfo.hash] = commitInfo;
                        this.block_chain_.broadcast(Msg.commit({
                            "height": this.prepare_info_.height,
                            "hash": this.prepare_info_.hash,
                            "signer": this.block_chain_.get_account_id()
                        }));
                    }
                }
                break;
            case MessageType.Commit:
                if (!this.commit_hash_cache_[key]) {
                    this.commit_hash_cache_[key] = true;
                    this.block_chain_.broadcast(msg);
                } else {
                    return;
                }
                // prepare消息是只能处理一个，但是commit消息却是可以处理多个，哪一个先达成共识就先处理哪个，然后剩下没处理的都清空
                var commit = this.commit_infos_[msg.data.hash];
                if (commit) {
                    if (!commit.votes[msg.data.signer]) {
                        commit.votes[msg.data.signer] = true;
                        commit.votesNumber++;
                        if (commit.votesNumber > 2 * PBFT_F) {
                            // console.log(`node: ${this.block_chain_.get_account_id()} \t[commited] do commit block}`);
                            this.commit(msg.data.hash);
                        }
                    }
                } else {
                    this.commit_infos_[msg.data.hash] = {
                        hash: msg.data.hash,
                        height: msg.data.height,
                        votesNumber: 1,
                        votes: {}
                    };
                    this.commit_infos_[msg.data.hash].votes[msg.data.signer] = true;
                }
                break;
            default:
                break;
        }
    }
}

module.exports = Pbft;
