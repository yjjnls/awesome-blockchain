'use strict';

var BlockChain = require("../blockchain");
var Consensus = require("../consensus/pow");

let blockchain = new BlockChain(Consensus);

console.log(blockchain.get_last_block());