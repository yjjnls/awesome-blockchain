# Awesome Blockchain 
[![Awesome](https://awesome.re/badge.svg)](https://github.com/yjjnls/awesome-blockchain)
>Curated list of resources for the development and applications of block chain.

The blockchain is an incorruptible digital ledger of economic transactions that can be programmed to record not just financial transactions but virtually everything of value (by [Don Tapscott](https://www.linkedin.com/pulse/whats-next-generation-internet-surprise-its-all-don-tapscott)).

<font color=#0099ff size=3>**This is not a simple collection of Internet resources, but verified and organized data ensuring it's really suitable for your learning process and uesful for your development and application.**</font> 


## Contents
- [Awesome Blockchain](#awesome-blockchain)
    - [Contents](#contents)
    - [Frequently Asked Questions (F.A.Q.s) & Answers](#frequently-asked-questions-faqs--answers)
    - [Basic Introduction](#basic-introduction)
    - [Further Extesnsion](#further-extesnsion)
        - [Books](#books)
    - [Development Tutorial](#development-tutorial)
        - [BitCoin](#bitcoin)
        - [Ethereum](#ethereum)
        - [Fabric](#fabric)
    - [Releated Tools](#releated-tools)
        - [Solidity](#solidity)
        - [truffle](#truffle)
        - [web3.js](#web3js)
    - [Projects and Applications](#projects-and-applications)
        - [Monero](#monero)
        - [IOTA](#iota)
        - [EOS](#eos)
        - [IFPS](#ifps)
    - [Contribute](#contribute)

## Frequently Asked Questions (F.A.Q.s) & Answers


**Q: What's a Blockchain?**

A: A blockchain is a distributed database with a list (that is, chain) of records (that is, blocks) linked and secured by
digital fingerprints (that is, cryptho hashes).
Example from [`blockchain.rb`](https://github.com/openblockchains/awesome-blockchains/blob/master/blockchain.rb/blockchain.rb):

```
[#<Block:0x1eed2a0
  @timestamp     = 1637-09-15 20:52:38,
  @data          = "Genesis",
  @previous_hash = "0000000000000000000000000000000000000000000000000000000000000000",
  @hash          = "edbd4e11e69bc399a9ccd8faaea44fb27410fe8e3023bb9462450a0a9c4caa1b">,
 #<Block:0x1eec9a0
  @timestamp     = 1637-09-15 21:02:38,
  @data          = "Transaction Data...",
  @previous_hash = "edbd4e11e69bc399a9ccd8faaea44fb27410fe8e3023bb9462450a0a9c4caa1b",
  @hash          = "eb8ecbf6d5870763ae246e37539d82e37052cb32f88bb8c59971f9978e437743">,
 #<Block:0x1eec838
  @timestamp     = 1637-09-15 21:12:38,
  @data          = "Transaction Data......",
  @previous_hash = "eb8ecbf6d5870763ae246e37539d82e37052cb32f88bb8c59971f9978e437743",
  @hash          = "be50017ee4bbcb33844b3dc2b7c4e476d46569b5df5762d14ceba9355f0a85f4">,
  ...
```

![](Basic/img/blockchain-jesus.png)


**Q: What's a Hash? What's a (One-Way) Crypto(graphic) Hash Digest Checksum**?

A: A hash e.g. `eb8ecbf6d5870763ae246e37539d82e37052cb32f88bb8c59971f9978e437743`
is a small digest checksum calculated
with a one-way crypto(graphic) hash digest checksum function
e.g. SHA256 (Secure Hash Algorithm 256 Bits)
from the data. Example from [`blockchain.rb`](blockchain.rb/blockchain.rb):

```ruby
def calc_hash
  sha = Digest::SHA256.new
  sha.update( @timestamp.to_s + @previous_hash + @data )
  sha.hexdigest   ## returns "eb8ecbf6d5870763ae246e37539d82e37052cb32f88bb8c59971f9978e437743"
end
```

A blockchain uses

- the block timestamp (e.g. `1637-09-15 20:52:38`) and
- the hash from the previous block (e.g. `edbd4e11e69bc399a9ccd8faaea44fb27410fe8e3023bb9462450a0a9c4caa1b`) and finally
- the block data (e.g. `Transaction Data...`)

to calculate the new hash digest checksum, that is, the hash
e.g. `be50017ee4bbcb33844b3dc2b7c4e476d46569b5df5762d14ceba9355f0a85f4`.


**Q: What's a Merkle Tree?**

A: A Merkle tree is a hash tree named after Ralph Merkle who patented the concept in 1979
(the patent expired in 2002). A hash tree is a generalization of hash lists or hash chains where every leaf node (in the tree) is labelled with a data block and every non-leaf node (in the tree)
is labelled with the crypto(graphic) hash of the labels of its child nodes. For more see the [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree) Wikipedia Article.

Note: By adding crypto(graphic) hash functions you can "merkelize" any data structure.


**Q: What's a Merkelized DAG (Directed Acyclic Graph)?**

A: It's a blockchain secured by crypto(graphic) hashes that uses a directed acyclic graph data structure (instead of linear "classic" linked list).

Note: Git uses merkelized dag (directed acyclic graph)s for its blockchains.


**Q: Is the Git Repo a Blockchain?**

A: Yes, every branch in the git repo is a blockchain.
The "classic" Satoshi-blockchain is like a git repo with a single master branch (only).

## Basic Introduction
*  [Encryption](./Basic/crypto.md)  数字加密相关知识
- [ ] [Consensus]()  共识算法详解
*  [Account and transaction model](./Basic/account.md)  账户与交易模型
- [ ] [Bitcoin basics]()  比特币基础知识
- [ ] [Ethereum basics]()  以太坊基础知识
<!-- - [ ] []()链上治理 -->
- [ ] [Exchange]()  数字交易所基础知识
- [ ] [Application]()  应用与思考：区块链不能做什么？
* [Digital currency ranking](https://coinmarketcap.com/)  数字货币排行

## Further Extesnsion
### Books
*  [区块链技术指南](https://yeasy.gitbooks.io/blockchain_guide/content/)
*  [区块链原理、设计与应用](https://github.com/yjjnls/books/blob/master/block%20chain/%E5%8C%BA%E5%9D%97%E9%93%BE%E5%8E%9F%E7%90%86%E3%80%81%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%BA%94%E7%94%A8.pdf)
*  [区块链 从数字货币到信用社会](https://github.com/yjjnls/books/blob/master/block%20chain/%E5%8C%BA%E5%9D%97%E9%93%BE%20%E4%BB%8E%E6%95%B0%E5%AD%97%E8%B4%A7%E5%B8%81%E5%88%B0%E4%BF%A1%E7%94%A8%E7%A4%BE%E4%BC%9A.pdf)
*   [**Attack of the 50 Foot Blockchain: Bitcoin, Blockchain, Ethereum & Smart Contracts**](https://davidgerard.co.uk/blockchain/table-of-contents/) by David Gerard, London, 2017 --
_What is a bitcoin? ++
The Bitcoin ideology ++
The incredible promises of Bitcoin! ++
Early Bitcoin: the rise to the first bubble ++
How Bitcoin mining centralised ++
Who is Satoshi Nakamoto? ++
Spending bitcoins in 2017 ++
Trading bitcoins in 2017: the second crypto bubble ++
Altcoins ++
Smart contracts, stupid humans ++
Business bafflegab, but on the Blockchain ++
Case study: Why you can’t put the music industry on a blockchain_

*   [**Mastering Bitcoin - Programming the Open Blockchain**](https://github.com/bitcoinbook/bitcoinbook/blob/second_edition/ch09.asciidoc) 2nd Edition,
by Andreas M. Antonopoulos, 2017 - FREE (Online Source Version) --
_What Is Bitcoin? ++
How Bitcoin Works ++
Bitcoin Core: The Reference Implementation ++
Keys, Addresses ++
Wallets ++
Transactions ++
Advanced Transactions and Scripting ++
The Bitcoin Network ++
The Blockchain ++
Mining and Consensus ++
Bitcoin Security ++
Blockchain Applications_

*   [**Programming Blockchains in Ruby from Scratch Step-by-Step Starting w/ Crypto Hashes... ( Beta / Rough Draft )**](https://github.com/yukimotopress/programming-blockchains-step-by-step)
by Gerald Bauer et al, 2018 - FREE (Online Version) --
_(Crypto) Hash ++
(Crypto) Block ++
(Crypto) Block with Proof-of-Work ++
Blockchain! Blockchain! Blockchain! ++
Blockchain Broken? ++
Timestamping ++
Mining, Mining, Mining - What's Your Hash Rate? ++
Bitcoin, Bitcoin, Bitcoin ++
(Crypto) Block with Transactions (Tx)_


*   [**Programming Cryptocurrencies and Blockchains in Ruby ( Beta / Rough Draft )**](http://yukimotopress.github.io/blockchains)
by Gerald Bauer et al, 2018 - FREE (Online Version) @ Yuki & Moto Press Bookshelf --
_Digital $$$ Alchemy - What's a Blockchain? -
How-To Turn Digital Bits Into $$$ or €€€? •
Decentralize Payments. Decentralize Transactions. Decentralize Blockchains. •
The Proof of the Pudding is ... The Bitcoin (BTC) Blockchain(s)
++
Building Blockchains from Scratch -
A Blockchain in Ruby in 20 Lines! A Blockchain is a Data Structure  •
What about Proof-of-Work? What about Consensus?   •
Find the Lucky Number - Nonce == Number Used Once
++
Adding Transactions -
The World's Worst Database - Bitcoin Blockchain Mining  •
Tulips on the Blockchain! Adding Transactions
++
Blockchain Lite -
Basic Blocks  •
Proof-of-Work Blocks  •
Transactions
++
Merkle Tree -
Build Your Own Crypto Hash Trees; Grow Your Own Money on Trees  •
What's a Merkle Tree?   •
Transactions
++
Central Bank -
Run Your Own Federated Central Bank Nodes on the Blockchain Peer-to-Peer over HTTP  •
Inside Mining - Printing Cryptos, Cryptos, Cryptos on the Blockchain
++
Awesome Crypto
++
Case Studies - Dutch Gulden  • Shilling  • CryptoKitties (and CryptoCopycats)_

*   [**Blockchain for Dummies, IBM Limited Edition**](https://www.ibm.com/blockchain/what-is-blockchain.html) by Manav Gupta, 2017 - FREE (Digital Download w/ Email) --
_Grasping Blockchain Fundamentals ++
Taking a Look at How Blockchain Works ++
Propelling Business with Blockchains ++
Blockchain in Action: Use Cases ++
Hyperledger, a Linux Foundation Project ++
Ten Steps to Your First Blockchain application_

*   [**Get Rich Quick "Business Blockchain" Bible - The Secrets of Free Easy Money**](https://github.com/bitsblocks/get-rich-quick-bible), 2018 - FREE --
_Step 1: Sell hot air. How? ++
Step 2: Pump up your tokens. How? ++
Step 3: Revolutionize the World. How?_

*   [**Mastering Ethereum - Building Contract Services and Decentralized Apps on the Blockchain**](https://github.com/ethereumbook/ethereumbook) -
by Andreas M. Antonopoulos, Gavin Wood, 2018 - FREE (Online Source Version)
_What is Ethereum ++
Introduction ++
Ethereum Clients ++
Ethereum Testnets ++
Keys and Addresses ++
Wallets	++
Transactions ++
Contract Services ++
Tokens ++
Oracles ++
Accounting & Gas ++
EVM (Ethereum Virtual Machine) ++ 	
Consensus ++		
DevP2P (Peer-To-Peer) Protocol ++
Dev Tools and Frameworks ++
Decentralized Apps ++
Ethereum Standards (EIPs/ERCs)_

*   [**Building Decentralized Apps on the Ethereum Blockchain**](https://www.manning.com/books/building-ethereum-dapps) by Roberto Infante, 2018 - FREE chapter 1 --
_Understanding decentralized applications ++
The Ethereum blockchain ++
Building contract services in (JavaScript-like) Solidity ++
Running contract services on the Ethereum blockchain ++
Developing Ethereum Decentralized apps with Truffle ++
Best design and security practice_



*   [**Best of Bitcoin Maximalist - Scammers, Morons, Clowns, Shills & BagHODLers - Inside The New New Crypto Ponzi Economics**](https://github.com/bitsblocks/bitcoin-maximalist), 2018 - FREE 

*   [**Crypto Facts - Decentralize Payments - Efficient, Low Cost, Fair, Clean - True or False?**](https://github.com/bitsblocks/crypto-facts), 2018 - FREE 

*   [**IslandCoin White Paper - A Pen and Paper Cash System - How to Run a Blockchain on a Deserted Island**](https://github.com/bitsblocks/islandcoin-whitepaper)
by Tal Kol -- 
_Motivation ++
Consensus ++
Transaction and Block Specification -
Transaction format •
Block format •
Genesis block ++
References_


## Development Tutorial
### [BitCoin](https://github.com/bitcoin/bitcoin) 
[<img src="https://bitcoin.org/img/icons/logotop.svg" align="right" width="120">](https://bitcoin.org/zh_CN/)
*  [BitCoin white paper](https://bitcoin.org/bitcoin.pdf) / [比特币白皮书](BitCoin/white%20paper.md)
*  [Mastering BitCoin](https://github.com/bitcoinbook/bitcoinbook) / [精通比特币](http://zhibimo.com/read/wang-miao/mastering-bitcoin/index.html)

### Ethereum
[<img src="https://github.com/yjjnls/Notes/blob/master/img/ethereum.png" align="right" width="80">](https://www.hyperledger.org/projects/fabric)
*  [Ethereum white paper](https://github.com/ethereum/wiki/wiki/White-Paper) / [以太坊白皮书](./Ethereum/white%20paper.md)
*  [Ethereum wiki](https://github.com/ethereum/wiki/wiki)
*  [Ethereum problems](https://github.com/ethereum/wiki/wiki/Problems)


### Fabric 
[<img src="https://www.hyperledger.org/wp-content/uploads/2018/03/Hyperledger_Fabric_Logo_Color.png" align="right" width="120">](https://www.hyperledger.org/projects/fabric)


## Releated Tools
### Solidity
### truffle
### web3.js

## Projects and Applications
### Monero

### IOTA

### EOS

<!-- [<img src="https://avatars2.githubusercontent.com/u/10536621?s=200&v=4" align="right" width="40">](https://github.com/ipfs/ipfs) -->
### IFPS


<!-- [<img src="https://avatars3.githubusercontent.com/u/22163706?s=200&v=4" align="right" width="40">](https://github.com/mvs-org/metaverse)
### Metaverse
元界链


[<img src="https://avatars0.githubusercontent.com/u/28505705?s=200&v=4" align="right" width="40">](https://github.com/Bytom/bytom)
### BYTOM 
比原链 -->




## Contribute

Contributions welcome! 

1. Fork it (https://github.com/yjjnls/awesome-blockchain/fork)
2. Clone it (`git clone https://github.com/yjjnls/awesome-blockchain`)
3. Create your feature branch (`git checkout -b your_branch_name`)
4. Commit your changes (`git commit -m 'Description of a commit'`)
5. Push to the branch (`git push origin your_branch_name`)
6. Create a new Pull Request

If you found this resource helpful, give it a 🌟 otherwise contribute to it and give it a ⭐️.