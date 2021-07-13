# Awesome Blockchain

[![Awesome](https://awesome.re/badge.svg)](https://github.com/yjjnls/awesome-blockchain)

> Curated list of resources for the development and applications of block chain.

The blockchain is an incorruptible digital ledger of economic transactions that can be programmed to record not just financial transactions but virtually everything of value (by [Don Tapscott](https://www.linkedin.com/pulse/whats-next-generation-internet-surprise-its-all-don-tapscott)).

<font color=#0099ff size=3>**This is not a simple collection of Internet resources, but verified and organized data ensuring it's really suitable for your learning process and useful for your development and application.**</font>

## Contents
<details><summary>Click to expand</summary>

- [Awesome Blockchain](#awesome-blockchain)
  - [Contents](#contents)
  - [Frequently Asked Questions (F.A.Q.s) & Answers](#frequently-asked-questions-faqs--answers)
  - [Basic Introduction](#basic-introduction)
  - [Development Tutorial](#development-tutorial)
    - [BitCoin](#bitcoin)
    - [Ethereum](#ethereum)
    - [Consortium Blockchain](#consortium-blockchain)
      - [Hyperledger](#hyperledger)
      - [XuperChain](#xuperchain)
      - [FISCO-BCOS](#fisco-bcos)
  - [Releated Tools](#releated-tools)
    - [Solidity](#solidity)
    - [truffle](#truffle)
    - [web3.js](#web3js)
  - [Implementation of Blockchain](#implementation-of-blockchain)
  - [Projects and Applications](#projects-and-applications)
    - [Quorum](#quorum)
    - [Monero](#monero)
    - [IOTA](#iota)
    - [EOS](#eos)
    - [IPFS](#ipfs)
      - [Filecoin](#filecoin)
      - [BigchainDB](#bigchaindb)
    - [BitShares](#bitshares)
    - [ArcBlock](#arcblock)
  - [Further Extension](#further-extension)
    - [Papers](#papers)
    - [Books](#books)
    - [Applications](#applications)
      - [Identity Applications](#identity-applications)
        - [Public Blockchain Identity](#public-blockchain-identity)
        - [Blockchain as a collateral](#blockchain-as-a-collateral)
        - [Unclear](#unclear)
        - [Guidance](#guidance)
      - [Internet of Things Applications](#internet-of-things-applications)
      - [Energy Applications](#energy-applications)
      - [Media and Journalism](#media-and-journalism)
      - [DeFi (Decentralised Finance)](#defi-decentralised-finance)
  - [Contribute](#contribute)

</details>

## Frequently Asked Questions (F.A.Q.s) & Answers

**Q: What's a Blockchain?**

A: A blockchain is a distributed database with a list (that is, chain) of records (that is, blocks) linked and secured by
digital fingerprints (that is, crypto hashes).
Example from [`genesis_block.json`](https://github.com/yjjnls/awesome-blockchain/tree/master/src/js/genesis_block.json):

```js
{
    "version": 0,
    "height": 1,
    "previous_hash": null,
    "timestamp": 1550049140488,
    "merkle_hash": null,
    "generator_publickey": "18941c80a77f2150107cdde99486ba672b5279ddd469eeefed308540fbd46983",
    "hash": "d611edb9fd86ee234cdc08d9bf382330d6ccc721cd5e59cf2a01b0a2a8decfff",
    "block_signature": "603b61b14348fb7eb087fe3267e28abacadf3932f0e33958fb016ab60f825e3124bfe6c7198d38f8c91b0a3b1f928919190680e44fbe7289a4202039ffbb2109",
    "consensus_data": {},
    "transactions": []
}
```

![](Basic/img/blockchain-jesus.png)

**Q: What's a Hash? What's a (One-Way) Crypto(graphic) Hash Digest Checksum**?

A: A hash e.g. `d611edb9fd86ee234cdc08d9bf382330d6ccc721cd5e59cf2a01b0a2a8decfff`
is a small digest checksum calculated
with a one-way crypto(graphic) hash digest checksum function
e.g. SHA256 (Secure Hash Algorithm 256 Bits)
from the data. Example from [`crypto.js`](https://github.com/yjjnls/awesome-blockchain/blob/master/src/js/crypto.js):

```js
function calc_hash(data) {
    return crypto.createHash('sha256').update(data).digest('hex');
}
```

A blockchain uses

-   the block header (e.g. `Version`, `TimeStamp`, `Previous Hash...` )and
-   the block data (e.g. `Transaction Data...`)

to calculate the new hash digest checksum.

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

**More Q&A**
- [Blockchain Interview Questions](https://mindmajix.com/blockchain-interview-questions)
- [10 Essential Blockchain Interview Questions](https://www.toptal.com/blockchain/interview-questions)
- [Top 36 Blockchain Job Interview Questions & Answers](https://blockchainsfactory.com/blockchain-interview-questions/)

---
## Basic Introduction

<!--    
### Encryption knowledge
   -->

-   **Encryption knowledge**  
    * [Basic concepts](https://www.jianshu.com/p/a044b303f7d5) - Asymmetric encryption, Digital signature, Certificate  
    * [Digital signature extension](https://www.jianshu.com/p/410e77ec23fa)  - Multi-signature, Blind signature, Group signature, Ring signature
    * [Merkle tree](https://www.jianshu.com/p/a044b303f7d5)  
    <!-- * [Merkle tree in blockchain](./Basic/merkle_tree_in_blockchain.md)   -->
    * [Merkle DAG](http://www.sohu.com/a/247540268_100222281)   
    * [**CryptoNote v2.0**](https://cryptonote.org/whitepaper.pdf) - Untraceable Transactions and Egalitarian Proof-of-work
<!--   
### Consensus
    -->
-   **Consensus**  
    * [Proof of Work](https://www.jianshu.com/p/3462f2ed74d7)
    * [Proof of Stake](https://www.jianshu.com/p/2fd3bce523b0)
    * [Proof of Stake FAQs](https://github.com/ethereum/wiki/wiki/Proof-of-Stake-FAQs) / [Chinese version](https://ethfans.org/posts/Proof-of-Stake-FAQ-new-2018-3-15)
    * [Delegated Proof of Stake](https://www.jianshu.com/p/ccc3fff7a60d)
    * [Practical Byzantine Fault Tolerance](https://www.jianshu.com/p/e991c1385f9f)

<!--    
### Account and transaction model
    -->
-   **Account and transaction model**  
    * [UTXO model](https://www.jianshu.com/p/2f4e75dbc2e4)
<!--
### Exchange
    -->
-   **Exchange**  
<!--
### Applications
    -->
-   **Applications**  
    * [Do You Need a Blockchain?](https://spectrum.ieee.org/computing/networks/do-you-need-a-blockchain)  
    * [What can't blockchain do?](https://www.jianshu.com/p/70f6a29a6296)  
    * [More](./Extension/application.md)
<!--     
### Governance
    -->
-   **Governance**
    * [Blockchains should not be democracies](https://haseebq.com/blockchains-should-not-be-democracies/)                                       
<!-- * [](https://github.com/yfeng125/blockchain-tutorial/blob/master/doc/%E2%80%8B25.%E6%AF%94%E7%89%B9%E5%B8%81%EF%BC%9A%E6%89%A9%E5%AE%B9%E4%B9%8B%E4%BA%89%E3%80%81IFO%E4%B8%8E%E9%93%BE%E4%B8%8A%E6%B2%BB%E7%90%86.md)   -->


<!--     
### Digital currency ranking
    -->
-   **[Digital currency ranking](https://coinmarketcap.com/)**   

---
## Development Tutorial

### [BitCoin](https://github.com/bitcoin/bitcoin)

[<img src="https://bitcoin.org/img/icons/logotop.svg" align="right" width="120">](https://bitcoincore.org)

**Bitcoin** is an experimental digital currency that enables instant payments to anyone, anywhere in the world. Bitcoin uses **peer-to-peer** technology to **operate with no central authority**: managing transactions and issuing money are carried out collectively by the network.

-   [BitCoin white paper: A Peer-to-Peer Electronic Cash System](https://bitcoin.org/bitcoin.pdf) / [Chinese version](BitCoin/white%20paper.md) / [Annotated BitCoin white paper](https://fermatslibrary.com/s/bitcoin)
-   [Mastering BitCoin](https://github.com/bitcoinbook/bitcoinbook) / [Chinese version](http://book.8btc.com/books/6/masterbitcoin2cn/_book/) / [pdf download](http://book.8btc.com/master_bitcoin?export=pdf)
-   [Bitcoin Improvement Proposals (BIPs)](https://github.com/bitcoin/bips/)

+   [But how does bitcoin actually work?](https://www.youtube.com/watch?v=bBC-nXj3Ng4)
+   [Mining visualization](http://www.yogh.io/#mine:last)
+   [Wallets](./BitCoin/awesome.md#wallets-api)
+   [Explorers](./BitCoin/awesome.md#blockchain-explorers)
+   [Libraries](./BitCoin/awesome.md#libraries) - C++, JavaScript, PHP, Ruby, Python, Java, .Net
+   [Web services](./BitCoin/awesome.md#blockchain-api-and-web-services)
+   [Full nodes](./BitCoin/awesome.md#full-nodes)
+   [More](./BitCoin/awesome.md)

### [Ethereum](https://github.com/ethereum)

[<img src="https://github.com/yjjnls/Notes/blob/master/img/ethereum.png" align="right" width="80">](https://www.ethereum.org/)

**Ethereum** is a **decentralized platform that runs smart contracts**: applications that run exactly as programmed without any possibility of downtime, censorship, fraud or third-party interference.

These apps run on a custom built **blockchain, an enormously powerful shared global infrastructure that can move value around and represent the ownership of property.**


-   [Ethereum white paper](https://github.com/ethereum/wiki/wiki/White-Paper) / [Chinese version](./Ethereum/white%20paper.md) / [Annotated Ethereum white paper](https://fermatslibrary.com/s/ethereum-a-next-generation-smart-contract-and-decentralized-application-platform)
-   [Mastering Ethereum](https://github.com/ethereumbook/ethereumbook) / [Chinese version](https://github.com/inoutcode/ethereum_book)
-   [Ethereum Yellow Paper](https://ethereum.github.io/yellowpaper/paper.pdf) / [Chinese version](https://github.com/yuange1024/ethereum_yellowpaper)
-   [Ethereum wiki](https://github.com/ethereum/wiki/wiki)
    -   [Ethereum Design Rationale](https://github.com/ethereum/wiki/wiki/Design-Rationale) / [Chinese version](https://ethfans.org/posts/510)
    -   [Ethereum problems](https://github.com/ethereum/wiki/wiki/Problems)
    -   [Sharding roadmap](https://github.com/ethereum/wiki/wiki/Sharding-roadmap)
    -   [**Ethereum flavored WebAssembly (ewasm)**](https://github.com/ewasm)
    -   [ÐΞVp2p Wire Protocol](https://github.com/ethereum/wiki/wiki/%C3%90%CE%9EVp2p-Wire-Protocol)
    -   [EVM-Awesome-List](https://github.com/ethereum/wiki/wiki/Ethereum-Virtual-Machine-(EVM)-Awesome-List)
    -   [Patricia Tree](https://github.com/ethereum/wiki/wiki/Patricia-Tree)
    -   Consensus
        -   [Ethash](https://github.com/ethereum/wiki/wiki/Ethash)
        -   [Ethash-DAG](https://github.com/ethereum/wiki/wiki/Ethash-DAG)
        -   [Ethash Specification](https://github.com/ethereum/wiki/wiki/Ethash)
        -   [Mining Ethash DAG](https://github.com/ethereum/wiki/wiki/Mining#ethash-dag)
        -   [Dagger-Hashimoto Algorithm](https://github.com/ethereum/wiki/blob/master/Dagger-Hashimoto.md)
        -   [DAG Explanation and Images](https://ethereum.stackexchange.com/questions/1993/what-actually-is-a-dag)
        -   [Ethash in Ethereum Yellowpaper](https://ethereum.github.io/yellowpaper/paper.pdf#appendix.J)
        -   [Ethash C API Example Usage](https://github.com/ethereum/wiki/wiki/Ethash-C-API)
-   [Accounts, Transactions, Gas, and Block Gas Limits in Ethereum](https://hudsonjameson.com/2017-06-27-accounts-transactions-gas-ethereum/)
-   [Ethereum Improvement Proposals](https://eips.ethereum.org/)
-   [Important EIPs and ERCs](https://github.com/ethereumbook/ethereumbook/blob/develop/appdx-standards-eip-erc.asciidoc#table-of-most-important-eips-and-ercs) / [EIP list](https://github.com/ethereum/EIPs)
-   Security
    -   [Ethereum Smart Contract Security Best Practices](https://consensys.github.io/smart-contract-best-practices/) / [Chinese version](https://github.com/ConsenSys/smart-contract-best-practices/blob/master/README-zh.md)
    -   [Onward with Ethereum Smart Contract Security](https://blog.zeppelin.solutions/onward-with-ethereum-smart-contract-security-97a827e47702)
    -   [The Hitchhiker's Guide to Smart Contracts in Ethereum](https://blog.zeppelin.solutions/the-hitchhikers-guide-to-smart-contracts-in-ethereum-848f08001f05)
    -   [**OpenZeppelin**](https://docs.openzeppelin.com/openzeppelin/)
    -   [**openzeppelin contracts**](https://github.com/OpenZeppelin/openzeppelin-contracts) / [doc](https://docs.openzeppelin.com/contracts/2.x/)
    -   [openzepplin sdk](https://github.com/OpenZeppelin/openzeppelin-sdk)
-   Token
    -   [ERC20](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md) / [impl](https://github.com/OpenZeppelin/openzeppelin-contracts/tree/master/contracts/token/ERC20)
    -   [ERC721](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md) / [impl](https://github.com/OpenZeppelin/openzeppelin-contracts/tree/master/contracts/token/ERC721)

+   Utils
    +   [Ethereum Blockchain Explorer](https://etherscan.io/)
    +   [Eth Gas Station](https://ethgasstation.info/)
    +   [Eth Network Status](https://ethstats.net/)
    

-   [**EEA** - Enterprise Ethereum: Private Blockchain For Enterprises](https://101blockchains.com/enterprise-ethereum/)
    -   [What Is Enterprise Ethereum?](https://101blockchains.com/enterprise-ethereum/#1)
    -   [What is The Enterprise Ethereum alliance?](https://101blockchains.com/enterprise-ethereum/#2)
    -   [Benefits of Enterprise Ethereum](https://101blockchains.com/enterprise-ethereum/#3)
    -   [Architecture Stack of the Enterprise Ethereum Blockchain](https://101blockchains.com/enterprise-ethereum/#4)
    -   [What Are The Possible Enterprise Ethereum Use Cases?](https://101blockchains.com/enterprise-ethereum/#5)
    -   [Ethereum Blockchain as a Service Providers](https://101blockchains.com/enterprise-ethereum/#6)
    -   [Real-World Companies Using Enterprise Ethereum](https://101blockchains.com/enterprise-ethereum/#7)
    -   [Final Words](https://101blockchains.com/enterprise-ethereum/#8)

### Consortium Blockchain
*   **Theory**
    -   [**The Byzantine Generals Problem**](https://people.eecs.berkeley.edu/~luca/cs174/byzantine.pdf)
    -   [**Practical Byzantine Fault Tolerance**](http://pmg.csail.mit.edu/papers/osdi99.pdf)
    -   [Is consortium blockchain better?](http://www.infoq.com/cn/news/2018/10/is-consortium-blockchain-better)   
    -   [5 consortium blockchain comparison](http://www.infoq.com/cn/articles/5-consortium-blockchain-comparison) / [quick version](https://upload-images.jianshu.io/upload_images/11336404-f753396df0e930c8.jpg?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)    
    -   [FISCO BCOS vs Fabric](http://www.infoq.com/cn/news/2018/09/uncover-consortium-blockchain)   

*   **Implement a consortium blockchain using ethereum**  
    -   [Building a Private Ethereum Consortium](https://www.microsoft.com/developerblog/2018/06/01/creating-private-ethereum-consortium-kubernetes/)
    -   [Deploying a private Ethereum blockchain to Microsoft Azure Cloud](https://www.youtube.com/watch?v=HsConsFaZG8)
    -   [Ethereum Consortium Network Deployments Made Easy](https://github.com/CatalystCode/ibera-ethereum-consortium-blockchain-network)
    -   [How to Set Up a Private Ethereum Blockchain in 20 Minutes](https://arctouch.com/blog/how-to-set-up-ethereum-blockchain/)


#### Hyperledger

[<img src="https://www.hyperledger.org/wp-content/uploads/2018/03/Hyperledger_Fabric_Logo_Color.png" align="right" width="120">](https://www.hyperledger.org/projects/fabric)

-   [Hyperledger Org](https://wiki.hyperledger.org/)
-   Fabric
    -   [Fabric Org](https://wiki.hyperledger.org/display/Fabric)
    -   [Fabric Design Documents](https://wiki.hyperledger.org/display/fabric/Design+Documents)
    -   [Fabric Wiki](https://hyperledger-fabric.readthedocs.io/en/latest/)
        -   1.4 [En](https://hyperledger-fabric.readthedocs.io/en/release-1.4/) / [Zn](https://hyperledger-fabric.readthedocs.io/zh_CN/release-1.4/) / [Release](https://hyperledger-fabric.readthedocs.io/_/downloads/en/release-1.4/pdf/)
        -   2.2 [En](https://hyperledger-fabric.readthedocs.io/en/release-2.2/) / [Zn](https://hyperledger-fabric.readthedocs.io/zh_CN/release-2.2/)
    -   [Fabric Source Code Analyse](https://yeasy.gitbook.io/hyperledger_code_fabric/overview)
    -   [A Kafka-based Ordering Service for Fabric](https://docs.google.com/document/d/19JihmW-8blTzN99lAubOfseLUZqdrB6sBR0HsRgCAnY/edit)

-   Explorer
    -   [Explorer Proposal](https://docs.google.com/document/d/1GuVNHZ5Jqq-gTVKflnZ1YiJfEoozvugqenC6QEQFQj4/edit)
    -   [Explorer doc](https://blockchain-explorer.readthedocs.io/en/master/architecture/index.html)

-   [IBM OpenTech Hyperledger Fabric 1.4 LTS Course](https://space.bilibili.com/102734951/channel/detail?cid=69148)
-   [edx: Introduction to Hyperledger Blockchain Technologies Free Course](https://www.edx.org/course/introduction-to-hyperledger-blockchain-technologie)


#### [XuperChain](https://github.com/xuperchain/xuperchain)
[<img src="https://avatars3.githubusercontent.com/u/43258643?s=200&v=4" align="right" width="80">](https://xchain.baidu.com/)

**XuperChain**, the first open source project of XuperChain Lab, introduces a highly flexible blockchain architecture with great transaction performance.

**XuperChain** is the underlying solution for union networks with following highlight features:

**High Performance**
*   Creative XuperModel technology makes contract execution and verification run parallelly.
*   [TDPoS](https://xuperchain.readthedocs.io/zh/latest/design_documents/xpos.html) ensures quick consensus in a large scale network.
*   WASM VM using AOT technology.

**Solid Security**
*   Contract account protected by multiple private keys ensures assets safety.
*   [Flexible authorization system](https://xuperchain.readthedocs.io/zh/latest/design_documents/permission_model.html) supports weight threshold, AK sets and could be easily extended.

**High Scalability**
*   Robust [P2P](https://xuperchain.readthedocs.io/zh/latest/design_documents/p2p.html) network supports a large scale network with thousands of nodes.
*   Branch management on ledger makes automatic convergence consistency and supports global deployment.

**Multi-Language Support**: Support pluggable multi-language contract VM using [XuperBridge](https://xuperchain.readthedocs.io/zh/latest/design_documents/XuperBridge.html) technology.

**Flexibility**: Modular and pluggable design provides high flexibility for users to build their blockchain solutions for various business scenarios.

-   [Baidu Blockchain Engine](https://cloud.baidu.com/product/bbe.html)
-   [Homepage](https://xchain.baidu.com/)
-   [Doc](https://xuperchain.readthedocs.io/zh/latest/index.html)
-   [Wiki](https://github.com/xuperchain/xuperchain/wiki) / [English version](https://github.com/xuperchain/xuperchain/wiki/Wiki-in-English)

+   [Getting start](https://github.com/xuperchain/xuperchain/wiki/3.-Getting-Started)
    +   [Account operation](https://xuperchain.readthedocs.io/zh/latest/advanced_usage/contract_accounts.html)
    +   [Multiple nodes deployment](https://xuperchain.readthedocs.io/zh/latest/advanced_usage/multi-nodes.html)
    +   [Wasm contract](https://xuperchain.readthedocs.io/zh/latest/advanced_usage/create_contracts.html)
    +   [Proposal](https://xuperchain.readthedocs.io/zh/latest/advanced_usage/initiate_proposals.html)
    +   [Parallel chain](https://xuperchain.readthedocs.io/zh/latest/advanced_usage/parallel_chain.html)
+   SDK
    +   [Go SDK](https://github.com/xuperchain/xuper-java-sdk)
    +   [Javascript SDK](https://github.com/xuperchain/xuper-sdk-js)
    +   [Java SDK](https://github.com/xuperchain/xuper-python-sdk)
    +   [Python SDK](https://github.com/xuperchain/xuper-python-sdk)
+   [Detailed FAQs](https://xuperchain.readthedocs.io/zh/latest/FAQs.html)
+   [Comparation with Fabric and Ethereum](https://github.com/xuperchain/xuperchain/wiki/%E9%99%84-%E8%AF%84%E6%B5%8B%E6%95%B0%E6%8D%AE%E5%AF%B9%E6%AF%94)

#### [FISCO-BCOS](https://github.com/FISCO-BCOS/Wiki)

## Releated Tools

### Solidity
-   [doc](https://solidity.readthedocs.io/en/develop/index.html) / [Chinese version](https://solidity-cn.readthedocs.io/zh/develop/)

### truffle

### web3.js
-   [doc](https://web3js.readthedocs.io/en/1.0/) / [Chinese version](http://web3.tryblockchain.org/Web3.js-api-refrence.html)

## Implementation of Blockchain
-   [**ATS**: _Functional Blockchain_](https://beta.observablehq.com/@galletti94/functional-blockchain)
-   [**C#**: _Programming The Blockchain in C#_](https://programmingblockchain.gitbooks.io/programmingblockchain/)
-   [**Crystal**: _Write your own blockchain and PoW algorithm using Crystal_](https://medium.com/@bradford_hamilton/write-your-own-blockchain-and-pow-algorithm-using-crystal-d53d5d9d0c52)
-   [**C++**: _Blockchain from Scratch_](https://github.com/openblockchains/awesome-blockchains/tree/master/blockchain.cpp)
-   [**Go: _Building Blockchain in Go_**](https://github.com/Jeiwan/blockchain_go) / [Chinese version 1](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/part-1/basic-prototype.md) / [Chinese version 2](https://zhangli1.gitbooks.io/dummies-for-blockchain/content/)
    -   [_Part 1: Basic Prototype_](https://jeiwan.cc/posts/building-blockchain-in-go-part-1/)
    -   [_Part 2: Proof-of-Work_](https://jeiwan.cc/posts/building-blockchain-in-go-part-2/)
    -   [_Part 3: Persistence and CLI_](https://jeiwan.cc/posts/building-blockchain-in-go-part-3/)
    -   [_Part 4: Transactions 1_](https://jeiwan.cc/posts/building-blockchain-in-go-part-4/)
    -   [_Part 5: Addresses_](https://jeiwan.cc/posts/building-blockchain-in-go-part-5/)
    -   [_Part 6: Transactions 2_](https://jeiwan.cc/posts/building-blockchain-in-go-part-6/)
    -   [_Part 7: Network_](https://jeiwan.cc/posts/building-blockchain-in-go-part-7/)
-   [**Go**: _Building A Simple Blockchain with Go_](https://www.codementor.io/codehakase/building-a-simple-blockchain-with-go-k7crur06v)
-   [**Go**: _Code your own blockchain in less than 200 lines of Go_](https://medium.com/@mycoralhealth/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc)
-   [**Go**: _Code your own blockchain mining algorithm in Go_](https://medium.com/@mycoralhealth/code-your-own-blockchain-mining-algorithm-in-go-82c6a71aba1f)
-   [**Go**: _GoCoin - A full Bitcoin solution written in Go language (golang)_](https://github.com/piotrnar/gocoin)
-   [**Go**: _GoChain - A basic implementation of blockchain in go_](https://github.com/crisadamo/gochain)
-   [**Go**: _Having fun implementing a blockchain using Golang_](https://github.com/izqui/blockchain)
-   [**Go**: _NaiveChain - A naive and simple implementation of blockchains_](https://github.com/kofj/naivechain)
-   [**Java**: _Creating Your First Blockchain with Java_](https://medium.com/programmers-blockchain/create-simple-blockchain-java-tutorial-from-scratch-6eeed3cb03fa)
-   [**Java**: _Write a blockchain with java_](https://www.jianshu.com/p/afd8c465c91a)
-   [**JavaScript**: _A cryptocurrency implementation in less than 1500 lines of code_](https://github.com/conradoqg/naivecoin)
-   [**JavaScript**: _A web-based demonstration of blockchain concepts_](https://github.com/anders94/blockchain-demo/)
-   [**JavaScript**: _Build your own Blockchain in JavaScript_](https://github.com/nambrot/blockchain-in-js)
-   [**JavaScript**: _Code for Blockchain Demo_](https://github.com/seanjameshan/blockchain)
-   [**JavaScript**: _Creating a blockchain with JavaScript_](https://github.com/SavjeeTutorials/SavjeeCoin)
-   [**JavaScript**: _How To Launch Your Own Production-Ready Cryptocurrency_](https://hackernoon.com/how-to-launch-your-own-production-ready-cryptocurrency-ab97cb773371)
-   [**JavaScript**: _Learn & Build a JavaScript Blockchain_](https://medium.com/digital-alchemy-holdings/learn-build-a-javascript-blockchain-part-1-ca61c285821e)
-   [**JavaScript**: _Node.js Blockchain Imlementation: BrewChain: Chain+WebSockets+HTTP Server_](http://www.darrenbeck.co.uk/blockchain/nodejs/nodejscrypto/)
-   [**JavaScript**: _Writing a tiny blockchain in JavaScript_](https://www.savjee.be/2017/07/Writing-tiny-blockchain-in-JavaScript/)
    -   [_Part 1: Implementing a basic blockchain_](https://www.savjee.be/2017/07/Writing-tiny-blockchain-in-JavaScript/)
    -   [_Part 2: Implementing proof-of-work_](https://www.savjee.be/2017/09/Implementing-proof-of-work-javascript-blockchain/)
    -   [_Part 3: Transactions & mining rewards_](https://www.savjee.be/2018/02/Transactions-and-mining-rewards/)
    -   [_Part 4: Signing transactions_](https://www.savjee.be/2018/10/Signing-transactions-blockchain-javascript/)
-   [**Kotlin**: _Let’s implement a cryptocurrency in Kotlin_](https://medium.com/@vasilyf/lets-implement-a-cryptocurrency-in-kotlin-part-1-blockchain-8704069f8580)
-   [**Python**: _A Practical Introduction to Blockchain with Python_](http://adilmoujahid.com/posts/2018/03/intro-blockchain-bitcoin-python/)
-   [**Python**: _Build your own blockchain: a Python tutorial_](http://ecomunsing.com/build-your-own-blockchain)
-   [**Python**: _Learn Blockchains by Building One_](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46)
-   [**Python**: _Let’s Build the Tiniest Blockchain_](https://medium.com/crypto-currently/lets-build-the-tiniest-blockchain-e70965a248b)
-   [**Python: _write-your-own-blockchain_**](https://bigishdata.com/2017/10/17/write-your-own-blockchain-part-1-creating-storing-syncing-displaying-mining-and-proving-work/)
    -   [_Part 1 — Creating, Storing, Syncing, Displaying, Mining, and Proving Work_](https://bigishdata.com/2017/10/17/write-your-own-blockchain-part-1-creating-storing-syncing-displaying-mining-and-proving-work/)
    -   [_Part 2 — Syncing Chains From Different Nodes_](https://bigishdata.com/2017/10/27/build-your-own-blockchain-part-2-syncing-chains-from-different-nodes/)
    -   [_Part 3 — Nodes that Mine_](https://bigishdata.com/2017/11/02/build-your-own-blockchain-part-3-writing-nodes-that-mine/)
    -   [_Part 4.1 — Bitcoin Proof of Work Difficulty Explained_](https://bigishdata.com/2017/11/13/how-to-build-a-blockchain-part-4-1-bitcoin-proof-of-work-difficulty-explained/)
    -   [_Part 4.2 — Ethereum Proof of Work Difficulty Explained_](https://bigishdata.com/2017/11/21/how-to-build-your-own-blockchain-part-4-2-ethereum-proof-of-work-difficulty-explained/)
-   [**Ruby**: _lets-build-a-blockchain_](https://github.com/Haseeb-Qureshi/lets-build-a-blockchain)
-   [**Ruby**: _Programming Blockchains Step-by-Step (Manuscripts Book Edition)_](https://github.com/yukimotopress/programming-blockchains-step-by-step)
-   [**Scala**: _How to build a simple actor-based blockchain_](https://medium.freecodecamp.org/how-to-build-a-simple-actor-based-blockchain-aac1e996c177)
-   [**TypeScript**: _Naivecoin: a tutorial for building a cryptocurrency_](https://lhartikk.github.io/)
    -   [_Minimal working blockchain_](https://lhartikk.github.io/jekyll/update/2017/07/14/chapter1.html)
    -   [_Proof of Work_](https://lhartikk.github.io/jekyll/update/2017/07/13/chapter2.html)
    -   [_Transactions_](https://lhartikk.github.io/jekyll/update/2017/07/12/chapter3.html)
    -   [_Wallet_](https://lhartikk.github.io/jekyll/update/2017/07/11/chapter4.html)
    -   [_Transaction relaying_](https://lhartikk.github.io/jekyll/update/2017/07/10/chapter5.html)
    -   [_Wallet UI and blockchain explorer_](https://lhartikk.github.io/jekyll/update/2017/07/09/chapter6.html)
-   [**TypeScript**: _NaivecoinStake: a tutorial for building a cryptocurrency with the Proof of Stake consensus_](https://naivecoinstake.learn.uno/)

---
## Projects and Applications
[<img src="https://raw.githubusercontent.com/jpmorganchase/quorum/master/logo.png" align="right" width="80">](https://github.com/jpmorganchase/quorum)  
### Quorum

**Quorum** is an Ethereum-based distributed ledger protocol with transaction/contract privacy and new consensus mechanisms.

**Quorum** is a fork of [go-ethereum](https://github.com/ethereum/go-ethereum) and is updated in line with go-ethereum releases.

Key enhancements over go-ethereum:

*   **Privacy** - Quorum supports private transactions and private contracts through public/private state separation, and utilises peer-to-peer encrypted message exchanges (see [Constellation](https://github.com/jpmorganchase/constellation) and [Tessera](https://github.com/jpmorganchase/tessera)) for directed transfer of private data to network participants
*   **Alternative** Consensus Mechanisms - with no need for POW/POS in a permissioned network, Quorum instead offers multiple consensus mechanisms that are more appropriate for consortium chains:
    *   **Raft-based Consensus** - a consensus model for faster blocktimes, transaction finality, and on-demand block creation
    *   **Istanbul BFT** - a PBFT-inspired consensus algorithm with transaction finality, by AMIS.
*   **Peer Permissioning** - node/peer permissioning using smart contracts, ensuring only known parties can join the network
*   **Higher Performance** - Quorum offers significantly higher performance than public geth


[<img src="https://avatars3.githubusercontent.com/u/7450663?s=460&v=4" align="right" width="80">](https://github.com/monero-project/monero)  
### Monero
**Monero** is a private, secure, untraceable, decentralised digital currency. You are your bank, you control your funds, and nobody can trace your transfers unless you allow them to do so.

**Privacy**: Monero uses a cryptographically sound system to allow you to send and receive funds without your transactions being easily revealed on the blockchain (the ledger of transactions that everyone has). This ensures that your purchases, receipts, and all transfers remain absolutely private by default.

**Security**: Using the power of a distributed peer-to-peer consensus network, every transaction on the network is cryptographically secured. Individual wallets have a 25 word mnemonic seed that is only displayed once, and can be written down to backup the wallet. Wallet files are encrypted with a passphrase to ensure they are useless if stolen.

**Untraceability**: By taking advantage of ring signatures, a special property of a certain type of cryptography, Monero is able to ensure that transactions are not only untraceable, but have an optional measure of ambiguity that ensures that transactions cannot easily be tied back to an individual user or computer.

- [Getmonero.org](https://getmonero.org) - The official Monero website
- [Lab.getmonero.org](https://lab.getmonero.org) - The official research group of Monero
- [RPC documentation](https://getmonero.org/resources/developer-guides/daemon-rpc.html) - RPC documentation of the Monero daemon
- [Wallet documentation](https://getmonero.org/resources/developer-guides/wallet-rpc.html) - Wallet documentation of the Monero daemon
- [Cryptonote Whitepaper](https://cryptonote.org/whitepaper.pdf) - White paper of cryptonote, the family of crypto-currencies of Monero
- [Review of the Cryptonote White Paper](https://downloads.getmonero.org/whitepaper_review.pdf) - By the research lab of Monero
- [Cryptonote Standards](https://cryptonote.org/cns) - The 10 Cryptonote standards (equivalent to BIPs for Bitcoin)


+ [**How to get started**](https://github.com/monero-project/monero#compiling-monero-from-source)
+ [**Roadmap**](https://www.getmonero.org/resources/roadmap/)
+ [**What is Monero? Most Comprehensive Guide**](https://blockgeeks.com/guides/monero/) / [Chinese version](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/monero/what-is-monero.md)
+ [**More resouces**](./Extension/monero.md)



[<img src="https://avatars0.githubusercontent.com/u/20126597?s=200&v=4" align="right" width="80">](https://github.com/iotaledger)  
### IOTA  


**IOTA** is a revolutionary new transactional settlement and data integrity layer for the Internet of Things. It’s based on a new distributed ledger architecture, the **Tangle**, which overcomes the inefficiencies of current **Blockchain** designs and introduces a new way of reaching consensus in a **decentralized peer-to-peer system**. For the first time ever, through IOTA people can transfer money without any fees. This means that even infinitesimally small nanopayments can be made through IOTA.

**IOTA** is the missing puzzle piece for **the Machine Economy** to fully emerge and reach its desired potential. We envision IOTA to be the public, permissionless backbone for the Internet of Things that enables true interoperability between all devices.

-   [IOTA](https://iota.org) - Next Generation Blockchain
-   [Whitepaper](https://iota.org/IOTA_Whitepaper.pdf) - The Tangle / [Chinese version](http://www.iotachina.com/wp-content/uploads/2016/11/2016112902003453.pdf)
-   [Wikipedia](https://en.wikipedia.org/wiki/IOTA_(Distributed_Ledger_Technology))
-   [A Primer on IOTA](https://blog.iota.org/a-primer-on-iota-with-presentation-e0a6eb2cc621) - A Primer on IOTA (with Presentation)
-   [IOTA China](http://iotachina.com/) - IOTA China 首页
-   [IOTA Italia](http://iotaitalia.com/) - IOTA Italia
-   [IOTA Korea](http://blog.naver.com/iotakorea) - IOTA 한국
-   [IOTA Japan](http://lhj.hatenablog.jp/entry/iota) - IOTA 日本
-   [IOTA on Reddit](https://www.reddit.com/r/Iota/)


+   [**How to get started**](https://github.com/iotaledger/iri#how-to-get-started)   
+   [**Roadmap**](https://www.iota.org/research/roadmap)
+   [**IOTA Transactions, Confirmation and Consensus**](https://github.com/noneymous/iota-consensus-presentation) / [Chinese version](https://github.com/liuchengxu/blockchain-tutorial/blob/master/content/iota/iota_consensus_v1.0.md)
+   [**More resouces**](./Extension/iota.md)  


[<img src="https://static.eos.io/images/Landing/SectionTokenSale/eos_spinning_logo.gif" align="right" width="80">](https://github.com/EOSIO/eos)  
### EOS

**EOSIO** is software that introduces a blockchain architecture designed to enable vertical and horizontal scaling of decentralized applications (the “EOSIO Software”). This is achieved through an operating system-like construct upon which applications can be built. The software provides accounts, authentication, databases, asynchronous communication and the scheduling of applications across multiple CPU cores and/or clusters. The resulting technology is a blockchain architecture that has the potential to scale to **millions of transactions per second**, eliminates user fees and allows for quick and easy deployment of decentralized applications. For more information, please read the [EOS.IO Technical White Paper](https://github.com/EOSIO/Documentation/blob/master/TechnicalWhitePaper.md).

- [EOS Wiki](https://github.com/EOSIO/eos/wiki) - High Level EOS Software Overview
- [Technical White Paper](https://github.com/EOSIO/Documentation/blob/master/TechnicalWhitePaper.md) - EOS.IO Technical White Paper v2
- [EOS: An Introduction - Black Edition](http://iang.org/papers/EOS_An_Introduction-BLACK-EDITION.pdf) - Ian Grigg's Whitepaper
- [EOSIO Developer Portal](https://developers.eos.io/) - Official EOSIO developer portal, with docs, APIs etc.

+ [**How to get started**](https://developers.eos.io/eosio-home)
+ [**Roadmap**](https://github.com/EOSIO/Documentation/blob/master/Roadmap.md)
+ [**Tools**](https://github.com/yjjnls/awesome-blockchain/blob/master/Extension/eos.md#tools)  
+ [**Language Support**](https://github.com/yjjnls/awesome-blockchain/blob/master/Extension/eos.md#language-support)  


[<img src="https://avatars2.githubusercontent.com/u/10536621?s=200&v=4" align="right" width="80">](https://github.com/ipfs)  
### IPFS
**IPFS** ([the InterPlanetary File System](https://github.com/ipfs/faq/issues/76)) is a new hypermedia distribution protocol, addressed by content and identities. IPFS enables the creation of completely distributed applications. It aims to make the web faster, safer, and more open.

**IPFS** is a distributed file system that seeks to connect all computing devices with the same system of files. In some ways, this is similar to the original aims of the Web, but IPFS is actually more similar to a single bittorrent swarm exchanging git objects. You can read more about its origins in the paper [IPFS - Content Addressed, Versioned, P2P File System](https://github.com/ipfs/ipfs/blob/master/papers/ipfs-cap2pfs/ipfs-p2p-file-system.pdf?raw=true).

**IPFS** is becoming a new major subsystem of the internet. If built right, it could complement or replace HTTP. It could complement or replace even more. It sounds crazy. It _is_ crazy.

- [White Paper](https://github.com/ipfs/papers/raw/master/ipfs-cap2pfs/ipfs-p2p-file-system.pdf) - Academic papers on IPFS / [Chinese version](https://gguoss.github.io/2017/05/28/ipfs/)
- [Specs](https://github.com/ipfs/specs) - Specifications on the IPFS protocol
- [Notes](https://github.com/ipfs/notes) - Various relevant notes and discussions (that do not fit elsewhere)
[<img src="https://camo.githubusercontent.com/651f7045071c78042fec7f5b9f015e12589af6d5/68747470733a2f2f697066732e696f2f697066732f516d514a363850464d4464417367435a76413155567a7a6e3138617356636637485676434467706a695343417365" align="right" width="200">](https://github.com/ipfs)  
- [Reading-list](https://github.com/ipfs/reading-list) - Papers to read to understand IPFS
- [Protocol Implementations](https://github.com/ipfs/ipfs#protocol-implementations)
- [HTTP Client Libraries](https://github.com/ipfs/ipfs#http-client-libraries)
![]()   

+ [**Roadmap**](https://github.com/ipfs/roadmap)
+ [**More resouces**](./Extension/ipfs.md)  

#### [Filecoin](https://filecoin.io/)
- [White paper](https://filecoin.io/filecoin.pdf) / [Chinese version](http://chainx.org/paper/index/index/id/13.html)

#### [BigchainDB](https://www.bigchaindb.com/)
- [White paper](https://www.bigchaindb.com/whitepaper) / [Chinese version](http://blog.csdn.net/fengqing79/article/details/70154076)

### BitShares
- [White paper]() / [Chinese version](https://www.8btc.com/article/3369)

### ArcBlock
- [Blockchain Developer Platform](https://www.arcblock.io) / [White Paper](https://www.arcblock.io/en/whitepaper/latest)

---
## Further Extension
### [Papers](https://github.com/decrypto-org/blockchain-papers)

### Books

-   [**Blockchain guide**](https://yeasy.gitbooks.io/blockchain_guide/content/) by Baohua Yang, 2017 --
    Introduce blockchain related technologies, from theory to practice with bitcoin, ethereum and hyperledger.
    <!-- -   [区块链原理、设计与应用](https://github.com/yjjnls/books/blob/master/block%20chain/%E5%8C%BA%E5%9D%97%E9%93%BE%E5%8E%9F%E7%90%86%E3%80%81%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%BA%94%E7%94%A8.pdf) -->
-   [**Blockchain: from Digital Currency to Credit Society**](https://github.com/yjjnls/books/blob/master/block%20chain/%E5%8C%BA%E5%9D%97%E9%93%BE%20%E4%BB%8E%E6%95%B0%E5%AD%97%E8%B4%A7%E5%B8%81%E5%88%B0%E4%BF%A1%E7%94%A8%E7%A4%BE%E4%BC%9A.pdf)
-   [**Attack of the 50 Foot Blockchain: Bitcoin, Blockchain, Ethereum & Smart Contracts**](https://davidgerard.co.uk/blockchain/table-of-contents/) by David Gerard, London, 2017 --
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

-   [**Mastering Bitcoin - Programming the Open Blockchain**](https://github.com/bitcoinbook/bitcoinbook/blob/develop/ch09.asciidoc) 2nd Edition,
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

-   [**Programming Blockchains in Ruby from Scratch Step-by-Step Starting w/ Crypto Hashes... ( Beta / Rough Draft )**](https://github.com/yukimotopress/programming-blockchains-step-by-step)
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


-   [**Programming Cryptocurrencies and Blockchains in Ruby ( Beta / Rough Draft )**](http://yukimotopress.github.io/blockchains)
    by Gerald Bauer et al, 2018 - FREE (Online Version) @ Yuki & Moto Press Bookshelf --
    _Digital $$$ Alchemy - What's a Blockchain? -
    How-To Turn Digital Bits Into $$$ or €€€? •
    Decentralize Payments. Decentralize Transactions. Decentralize Blockchains. •
    The Proof of the Pudding is ... The Bitcoin (BTC) Blockchain(s)
    \++
    Building Blockchains from Scratch -
    A Blockchain in Ruby in 20 Lines! A Blockchain is a Data Structure  •
    What about Proof-of-Work? What about Consensus?   •
    Find the Lucky Number - Nonce == Number Used Once
    \++
    Adding Transactions -
    The World's Worst Database - Bitcoin Blockchain Mining  •
    Tulips on the Blockchain! Adding Transactions
    \++
    Blockchain Lite -
    Basic Blocks  •
    Proof-of-Work Blocks  •
    Transactions
    \++
    Merkle Tree -
    Build Your Own Crypto Hash Trees; Grow Your Own Money on Trees  •
    What's a Merkle Tree?   •
    Transactions
    \++
    Central Bank -
    Run Your Own Federated Central Bank Nodes on the Blockchain Peer-to-Peer over HTTP  •
    Inside Mining - Printing Cryptos, Cryptos, Cryptos on the Blockchain
    \++
    Awesome Crypto
    \++
    Case Studies - Dutch Gulden  • Shilling  • CryptoKitties (and CryptoCopycats)_

-   [**Blockchain for Dummies, IBM Limited Edition**](https://www.ibm.com/blockchain/what-is-blockchain.html) by Manav Gupta, 2017 - FREE (Digital Download w/ Email) --
    _Grasping Blockchain Fundamentals ++
    Taking a Look at How Blockchain Works ++
    Propelling Business with Blockchains ++
    Blockchain in Action: Use Cases ++
    Hyperledger, a Linux Foundation Project ++
    Ten Steps to Your First Blockchain application_

-   [**Get Rich Quick "Business Blockchain" Bible - The Secrets of Free Easy Money**](https://github.com/bitsblocks/get-rich-quick-bible), 2018 - FREE --
    _Step 1: Sell hot air. How? ++
    Step 2: Pump up your tokens. How? ++
    Step 3: Revolutionize the World. How?_

-   [**Mastering Ethereum - Building Contract Services and Decentralized Apps on the Blockchain**](https://github.com/ethereumbook/ethereumbook) -
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

-   [**Building Decentralized Apps on the Ethereum Blockchain**](https://www.manning.com/books/building-ethereum-dapps) by Roberto Infante, 2018 - FREE chapter 1 --
    _Understanding decentralized applications ++
    The Ethereum blockchain ++
    Building contract services in (JavaScript-like) Solidity ++
    Running contract services on the Ethereum blockchain ++
    Developing Ethereum Decentralized apps with Truffle ++
    Best design and security practice_


-   [**Best of Bitcoin Maximalist - Scammers, Morons, Clowns, Shills & BagHODLers - Inside The New New Crypto Ponzi Economics**](https://github.com/bitsblocks/bitcoin-maximalist), 2018 - FREE

-   [**Crypto Facts - Decentralize Payments - Efficient, Low Cost, Fair, Clean - True or False?**](https://github.com/bitsblocks/crypto-facts), 2018 - FREE

-   [**IslandCoin White Paper - A Pen and Paper Cash System - How to Run a Blockchain on a Deserted Island**](https://github.com/bitsblocks/islandcoin-whitepaper)
    by Tal Kol --
    _Motivation ++
    Consensus ++
    Transaction and Block Specification -
    Transaction format •
    Block format •
    Genesis block ++
    References_

-   [**Blockchain in Action**](https://www.manning.com/books/blockchain-in-action) by Bina Ramamurthy, early access --
      _Learn how blockchain differs from other distributed systems ++
    Smart contract development with Ethereum and the Solidity language ++
    Web UI for decentralized apps ++
    Identity, privacy and security techniques ++
    On-chain and off-chain data storage_
    
-   [**Permissioned Blockchains in Action**](https://www.manning.com/books/permissioned-blockchains-in-action) by Mansoor Ahmed-Rengers & Marta Piekarska-Geater, early access --
      _A guide to creating innovative applications using blockchain technology ++
    Writing smart contracts and distributed applications using Solidity ++
    Configuring DLT networks ++
    Designing blockchain solutions for specific use cases ++
    Identity management in permissioned blockchains networks_
    
-   [**Programming Hyperledger Fabric**](https://www.amazon.com/dp/0578802228) by Siddharth Jain, --
      _A guide to developing blockchain applications for enterprise use cases ++
    Where Fabric fits in to the blockchain landscape ++
    The ins and outs of deploying real-world applications ++
    Developing smart contracts and client applications in Node ++
    Debugging and troubleshooting ++
    Securing production applications_



### Applications

#### Identity Applications

##### Public Blockchain Identity

-   [Blockstack](https://blockstack.org) - Platform for decentralized, server-less apps where users control their data. Identity included.
-   [Evernym](http://www.evernym.com) - Self-Sovereign identity built on top of open source permissioned blockchain.
-   [Jolocom](https://jolocom.com) - Self-sovereing identity wallet.
-   [SIN](https://en.bitcoin.it/wiki/Identity_protocol_v1) - Proposed identity protocol for BitCoin.
-   [uPort](https://www.uport.me) - Self-Sovereign identity on [Ethereum](https://ethereum.org) by [ConsenSys](https://consensys.net).

##### Blockchain as a collateral

-   [ShoCard](https://shocard.com) - Proprietary digital identity service, uses blockchain for time-stamping and secure documents exchange.
-   [Tradle](https://tradle.io/) - Makes a bank on blockchain, identity as a collateral.

##### Unclear

-   [KYC Chain](http://kyc-chain.com) - Secure platform for sharing verifiable identity claims, data or documents among financial institutions.
-   [ObjectChain Collab](http://www.objectchain-collab.com) - Cross-industry collaboration over distributed identity.
-   [UniquID](http://uniquid.com) - Identity both for people and devices.
-   [Vida Identity](https://vidaidentity.com) - Enterprise-grade Blockchain Identity Software.

##### Guidance

-   [ID3](https://idcubed.org) - Institute for Data Driven Design, explores issues around self-sovereign identity, and distributed organizations.
-   [OpenCreds](http://opencreds.org) - W3C Credentials Community Group.
-   [TAO Network Identity](http://tao.network/portfolio-item/the-identity-system/) - Description of blockchain identity by Tao.Network.

#### Internet of Things Applications

-   [Chronicled](http://www.chronicled.com) - IoT devices registry on blockchain.
-   [Filament](http://filament.com) - Software and hardware for decentralized Intranet of Things systems
-   [IOTA](http://www.iotatoken.com) - Decentralized Internet of Things token on blockless blockchain.
-   [Machinomy](http://machinomy.com) - Distributed platform for IoT micropayments.
-   [Project Oaken](https://www.projectoaken.com) - IoT blockchain platform.
-   [Slock.it](https://slock.it) - Ethereum-based platform for building Shared Things.

#### Energy Applications

-   [bankymoon](http://bankymoon.co.za/) - Blockchain consultancy. [Presented](http://goo.gl/L6vJBx) bitcoin-topped smart electricity meter. Once topped up, it chooses a plan, and starts moving energy.
-   [Co-Tricity](https://co-tricity.com/) - Decentralised energy marketplace by [Innogy](https://innovationhub.innogy.com/) and [ConsenSys](https://consensys.net).
-   [Electron](http://www.electron.org.uk/) - Reinventing energy on blockchain.
-   [GridSingularity](http://gridsingularity.com) - Blockchain for Smart Grid. Declare three years of work on the technology.
-   [lo3 energy](http://lo3energy.com) - Energy Services, Product Research & Development. Makers of [Brooklyn Microgrid](http://brooklynmicrogrid.com) along with [ConsenSys](https://consensys.net).
-   [lumo](https://lumoenergy.com.au) - Energy provider. Experiment with blockchain.
-   [PowerLedger](https://powerledger.io) - Decentralised energy marketpace.
-   [PowerPeers](https://www.powerpeers.nl/) - Peer-to-peer energy marketplace in the Netherlands.
-   [Solar Change](http://www.solarchange.co/) - Makers of [Solar Coin](http://solarcoin.org/). AltCoin for a MW of solar power.
-   [Terraledger](https://terraledger.com) - Provider of Renewable Energy Certificates.
-   [ImpactPPA](https://impactppa.com) - Reinvesting of power generated under Power Purchase Agreement in more PPAs.

#### Media and Journalism

-   [Steem](https://steem.io) - Decentralized social network which incentivises content creation and curation.
-   [PopChest](https://popchest.com) - Incentivized distributed video platform.
-   [Civil](https://joincivil.com) - Decentralized newsmaking platform.

#### DeFi (Decentralised Finance)

-   [Uniswap](https://uniswap.org) - Decentralized exchange powered by the Automated Market Maker model (AMM).
-   [Compound](https://compound.finance) - Decentralized lending and borrowing.
-   [1inch Exchange](https://1inch.exchange) - Get the best rates among multiple DEXes.
-   [Synthetix](https://synthetix.io/) - Protocol for synthetic assets.

+   Tools
    +   [Defi Dashboard](https://debank.com/): portfolio tracker, project lists, rankings, etc.
    +   [Zapper](https://zapper.fi/): dashboard for viewing and managing your DeFi investments.
    +   [Furucombo](https://furucombo.app/): easily create flashloans without writing a single line of code.
    +   [Covalent](https://www.covalenthq.com/): an unified API bringing visibility to billions of blockchain data points.
---
## Contribute

Contributions welcome!

1.  Fork it (<https://github.com/yjjnls/awesome-blockchain/fork>)
2.  Clone it (`git clone https://github.com/yjjnls/awesome-blockchain`)
3.  Create your feature branch (`git checkout -b your_branch_name`)
4.  Commit your changes (`git commit -m 'Description of a commit'`)
5.  Push to the branch (`git push origin your_branch_name`)
6.  Create a new Pull Request

If you found this resource helpful, give it a 🌟 otherwise contribute to it and give it a ⭐️.
