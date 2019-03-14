'use strict';

var MessageType = {
    Connection: 0,
    Block: 1,
    Transaction: 2,
    PrePrepare: 3,
    Prepare: 4,
    Commit: 5,
    Sync: 6,
    SyncBlock: 7
};

module.exports = {
    type: MessageType,
    connection: (data) => { return { type: MessageType.Connection, data: data }; },
    block: (data) => { return { type: MessageType.Block, data: data }; },
    preprepare: (data) => { return { type: MessageType.PrePrepare, data: data }; },
    prepare: (data) => { return { type: MessageType.Prepare, data: data }; },
    commit: (data) => { return { type: MessageType.Commit, data: data }; },
    transaction: (data) => { return { type: MessageType.Transaction, data: data }; },
    sync: (data) => { return { type: MessageType.Sync, data: data }; },
    sync_block: (data) => { return { type: MessageType.SyncBlock, data: data }; }
};