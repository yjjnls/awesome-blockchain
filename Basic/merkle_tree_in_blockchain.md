# Merkle tree in blockchain
- [Merkle tree in blockchain](#merkle-tree-in-blockchain)
    - [区块链节点中的数据都存在哪里](#区块链节点中的数据都存在哪里)
    - [比特币中的区块结构是怎样的](#比特币中的区块结构是怎样的)
    - [什么是 Merkle Tree 和 Merkle Proof](#什么是-merkle-tree-和-merkle-proof)
    - [以太坊中的 merkle tree](#以太坊中的-merkle-tree)
    - [Merkle Patricia tree](#merkle-patricia-tree)
    - [更深入的 Merkle Patricia tree](#更深入的-merkle-patricia-tree)
        - [Transaction trie](#transaction-trie)
        - [State Trie](#state-trie)
        - [Storage trie](#storage-trie)


## 区块链节点中的数据都存在哪里

在持久化方面，区块链数据可以直接存储在一个扁平的文件中，也可以存储在简单的数据库系统中，比特币和以太坊都区块链数据存储在 google的 LevelDb中。

## 比特币中的区块结构是怎样的

![](http://upload-images.jianshu.io/upload_images/11336404-713157d70bdf8b6a?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

Version: 用于区分软件版本号Previous Block Hash：是指向前一个区块头的 hash。在比特币中，区块头的 hash一般都是临时算出，并没有包含在本区块头或者区块中，但在持久化的时候可以作为索引存储以提高性能

Nonce、Difficulty Target和 Timestamp : 用在 pow共识算法中。

Merkle Root： 是区块中所有交易的指纹，merkle tree的树根。交易在区块链节点之间传播，所有节点都按相同的算法（merkle tree）将交易组合起来，如此可以判断交易是否完全一致,此外也用于轻量钱包中快速验证一个交易是否存在于一个区块中。

## 什么是 Merkle Tree 和 Merkle Proof

![](http://upload-images.jianshu.io/upload_images/11336404-e83df55bc68ea283?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

如上图，merkle Tree是一颗平衡树，树根也就是 Merkle Root存在区块头中。树的构建过程是递归地的计算 Hash的过程，如：先是 Hash交易 a得到 Ha，Hash交易 b得到 Hb，再 Hash前两个 Hash（也就是 Ha和 Hb）得到 Hab，其他节点也是同理递归，最终得到 Merkle Root。

Merkle tree在区块链中有两个作用：

1.  仅仅看 merkle root就可以知道区块中的所有交易是不是一样的

2.  对于轻量节点来说（不存储所有的交易信息，只同步区块头）提供了快速验证一个交易是否存在交易中的方法。

merkle proof从某处出发向上遍历，算出 merkle Root的所需要经过的路径节点。在上图的例子中，如果轻量钱包要验证 Txb（红色方框）是否已经包含在区块中，可以向全量节点请求 merkle Proof，用于证明 Txb的存在，过程为：

1.  全量节点只要返回黄色部分的节点信息（Ha与 Hcd）

2.  轻量节点执行计算 Hash(Txb)=Hb à Hash(Ha + Hb)=Hab à Hash(Hab + Hcd)=Habcd，计算出来的 merkleRoot(也就是 Habcd)跟已知区块头中的 merkleRoot比较，如果一样则认为交易确实已经入块。

在上图的区块中，仅仅存在少量的区块。如果区块所包含的交易很多，merkle proof仅仅需要带 log2(N)个节点，此时 merkle proof的优势就会变得非常明显。

## 以太坊中的 merkle tree

在比特币中，系统底层不维护每个账户的余额，只有 UTXO（Unspent Transaction Outputs）。账户之间的转账通过交易完成，确切地说，比特币用户将 UTXO作为交易的输入，可以花掉一个或者多个 UTXO。

一个 UTXO像一张现金纸币，要么不使用，要么全部使用，而不能只花一部分。举个例子来说，一个用户有一个价值 1比特币的 UTXO，如果他想转账 0.5给某人，那他可以创建一个交易，以这个价值 1比特币的 UTXO为输入，另外产生 0.5比特币的 OTXO作为这个交易的输出（找零给自己）。

比特币这个公开底层系统本身不单独维护每个账户的余额，不过比特币钱包可以记录每个用户所拥有的 UTXO，这样计算出用户的余额。

以太坊相比比特币，额外引入了账号状态数据，比如 nonce、余额 balance和合约数据，这些是区块链的关键数据，具有以下特性：

随着交易的入块需要不断高效地更新，所有的这些数据在不同节点之间能够高效地验证是一致的，状态数据不断更新的过程中，历史版本的数据数据需要保留。

系统中的每个节点执行完相同区块和交易后，那么这些节点中对应的所有账户数据都是一样的，账户列表相同，账户对应的余额等数据也相同。总的来说，这些账户数据就像状态机的状态，每个节点执行相同区块后，达到的状态应该是完全一致的。但是，这个状态并不是直接写到区块里面，因为这些数据都是可以由区块和交易重新产生的，如果写到区块里面会增加区块的大小，加重区块同步的负担。

![](http://upload-images.jianshu.io/upload_images/11336404-080b75ab1e7c9261?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

如上所示，区块头中保存了三个 merkle tree的 root:

tansaction root: 跟比特币中的 Merkle Root作用相同，相当于区块中交易的指纹，用于快速验   证交易是否相同以及证明某个交易的存在。

state root: 这颗树是账户状态（余额和 nonce等）存放的地方，除此之外，还保存着 storage root，也就是合约数据保存的地方。receipts root:区块中合约相关的交易输出的事件。

## Merkle Patricia tree

在 Transaction Root中，用类似比特币的二进制 merkle tree是能够解决问题的，因为它更适用于处理队列数据，一旦建好就不再修改。但是对于 state tree，情况就复杂多了，本质上来说，状态数据更像一个 map，包含着账号和账号状态的映射关系。除此之外，state tree还需要经常更新，经常插入或者删除，这样重新计算 Root的性能就显得尤其重要。

Trie是一种字典树，用于存储文本字符，并利用了单词之间共享前缀的特点，所以也叫做前缀树。Trie树在有些时候是比较浪费空间的，如下所示，即使这颗树只有两个词，如果这两个词很长，那么这颗树的节点也会变得非常多，无论是对于存储还是对于 cpu来说都是不可接受的。如下所示：

![](http://upload-images.jianshu.io/upload_images/11336404-d1e60c3f5bab8680?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

相比 Trie树，Patricia Trie将那些公共的的路径压缩以节省空间和提高效率，如下所示：

![](http://upload-images.jianshu.io/upload_images/11336404-e70174b908356ec8?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

以太坊中的 Merkle Patricia trie，顾名思义，它是 Patricia trie和 Merkle Tree的结合，即具有 merkle tree的特性，也具有 Patricia Trie的特征：

1.密码学安全，每个节点都都是按 hash引用，hash用来在 LevelDb中找对应的存储数据；

2.像 Patricia trie树一样，这些可以根据 Path来找对应的节点以找到 value；

3.引入了多种节点类型：

a.空节点 (比如说当一颗树刚刚创建为空的时候)

b.叶子节点，最普通的 [key, value]

c.扩展节点，跟叶子节点类似，不过值变成了指向别的节点的 hash,[key, hash]

d.分支节点，是一个长度为 17的列表，前 16元素为可能的十六进制字符，最后一个元素为 value(如果这是 path的终点的话)

举个例子：

![](http://upload-images.jianshu.io/upload_images/11336404-18ee531a20343340?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

在上图中的 trie包含了 4对 key value，需要注意的是，key是按照 16进制来显示的，也就是 a7占用一个字节，11占用一个字节等等

1.第一层的 Node是扩展节点，4个 Key都有公有的前缀 a7,next node指向一个分支节点

2.第二层是一个分支节点，由于 key转换成了十六进制，每个 branch最多有 16个分支。下标也作为 path的一部分用于 path查找。比如说下标为 1的元素中指向最左边的叶子节点（key-end为 1355）,到叶子节点就完成了 key搜索：扩展节点中 a7 + 分支节点下标 1 + 叶子节点 1355 = a711355

3.叶子节点和扩展节点的区分。正如上面提到的，叶子节点和扩展节点都是两个字段的节点，也就是 [key，value]，存储中没有专门字段用来标识类型。为了区分这两种节点类型并节省空间，在 key中加入了 4bits（1 nibble）的 flags的前缀，用 flags的倒数第二低的位指示是叶子节点还是扩展节点。此外，加入了 4bits之后，key的长度还有可能不是偶数个 nibble（存储中只能按字节存储），为此，如果 key是奇数个 nibble，在 flags nibble之后再添加一个空的 nibble，并且用 flags的最低位表示是否有添加，详见上图左下角。

## 更深入的 Merkle Patricia tree

更详细的字段关系如下图所示：

![](http://upload-images.jianshu.io/upload_images/11336404-eddd4da7b467b169?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

下面将通过代码片段的形式，逐一验证各个 trie的结构（前提条件是先在本地搭建起以太坊私有链）。

### Transaction trie
如下所示，在本地环境发送交易并使之入块，查看区块的交易列表，TransactionsRoot和 RawTransaction：

```
> eth.getBlock(49).transactions
["0xdf648e4ce9bed9d3b0b35d969056ac496207692f96bd13327807e920e97a1b2f"]
> eth.getBlock(49).transactionsRoot
"0x1a65885367afcc561411fe68ed870e4952b11477ad5314de1ae8f26d48a03724"
>eth.getRawTransaction("0xdf648e4ce9bed9d3b0b35d969056ac496207692f96bd13327807e920e97a1b2f")
"0xf86505850430e2340083015f90947b04e3fe46e1cd9939bf572307fdc076478b5252018042a0e9893deacc678345ea700e714b84ce31ffe4a50267c324436fab2c48906871ada03704497c029452a1b19b1f4876e958ec7e873600408d89a8bf46e53c6e5f921e"
```

在 trie包中写单测函数，key为交易在区块中的 index，RLP编码，value为签名过的原始交易 RawTransaction：

```
func TestMyTrieCalculateTxTree(t *testing.T) {
    var trie Trie
    keybuf := new(bytes.Buffer)
    rlp.Encode(keybuf, uint(0))
    valueBytes, _ :=
        hexutil.Decode("0xf86505850430e2340083015f90947b04e3fe46e1cd9939bf572307fdc076478b5252018042a0e9893deacc678345ea700e714b84ce31ffe4a50267c324436fab2c48906871ada03704497c029452a1b19b1f4876e958ec7e873600408d89a8bf46e53c6e5f921e")
   trie.Update(keybuf.Bytes(), valueBytes)
    t.Logf("Got Root:%s", trie.Hash().String())
}
```

运行输出得到的 Hash，也即 transactionsRoot为：

```
0x1a65885367afcc561411fe68ed870e4952b11477ad5314de1ae8f26d48a03724，跟 eth.getBlock(49).transactionsRoot得到的是一致的。
$ go test -v -run TestMyTrieCalculateTxTree
=== RUN  TestMyTrieCalculateTxTree
--- PASS: TestMyTrieCalculateTxTree (0.00s)
    my_trie_test.go:18: Got Root:0x1a65885367afcc561411fe68ed870e4952b11477ad5314de1ae8f26d48a03724
PASS
ok     github.com/ethereum/go-ethereum/trie    0.036s
```

### State Trie

获取最新的区块的 stateRoot，以及打印出账号 0x08e5f4cc4d1b04c450d00693c95ae58825f6a307的余额

```
> eth.getBlock(eth.blockNumber).stateRoot
"0xccc450ac770b0a644b81a8c0729733cf06d19f177e04fe664e1562dc3a620d60"
> eth.getBalance("0x08e5f4cc4d1b04c450d00693c95ae58825f6a307")
2.3229575729235784806170618e+25
```

在 state包中写单测函数，state trie的数据以 trie节点 hash为 key存在 leveldb中，所以整个 state trie的入口 key就是 stateRoot。state tree中存储数据的 path为 account的 hash，value为 RLP编码过的结构体数据。为了简单起见和节省篇幅，这里省去了错误检查。

```
func TestMyTrieCalculateStateTree(t *testing.T) {
    ldb, _ := ethdb.NewLDBDatabase("/Users/peace/ethereum/geth/chaindata", 0, 0)
    tr, _ := trie.New(common.HexToHash("0xccc450ac770b0a644b81a8c0729733cf06d19f177e04fe664e1562dc3a620d60"),
        trie.NewDatabase(ldb))

    accBytes, _ := hexutil.Decode("0x08e5f4cc4d1b04c450d00693c95ae58825f6a307")
    keyBytes := crypto.Keccak256Hash(accBytes).Bytes()
    valueBytes, _ := tr.TryGet(keyBytes)

    var acc Account
    rlp.DecodeBytes(valueBytes, &acc)
    t.Logf("balance:%d", acc.Balance)
}
```

运行输出得到 0x08e5f4cc4d1b04c450d00693c95ae58825f6a307的余额，跟 eth.getBalance接口得到的结果是一致的。

```
peaces-MacBook-Air:state peace$ go test -v -run TestMyTrieCalculateStateTree
=== RUN  TestMyTrieCalculateStateTree
--- PASS: TestMyTrieCalculateStateTree (0.01s)
    my_state_test.go:25: balance:23229575729235784806170618
PASS
ok     github.com/ethereum/go-ethereum/core/state  0.051s
```

### Storage trie

如下合约，为了简单起见，合约中省去了构造函数等不相关的内容，部署后地址为：

```
0x9ea9b9eeac924fd784b064dabf174a55113c4064。 
pragma solidity ^0.4.0;
contract testStorage {
   uint storeduint = 2018;
   string storedstring = 'Onething, OneWorld!';
}
```

获取到当前最新块的 stateRoot为 0x86bce3794034cddb3126ec488d38cb3ee840eeff4e64c3afe0ec5a5ca8b5f6ed。

```sh
> eth.getBlock(eth.blockNumber).stateRoot
"0x86bce3794034cddb3126ec488d38cb3ee840eeff4e64c3afe0ec5a5ca8b5f6ed"
```

在 state包中写单测函数，首先获以 0x86bce3794034cddb3126ec488d38cb3ee840eeff4e64c3afe0ec5a5ca8b5f6ed创建 trie，取获取合约账号 0x9ea9b9eeac924fd784b064dabf174a55113c4064的 storageRoot,之后再以这个 storageRoot创建 trie。在取合约内部数据时，key为 hash过的 32字节 index，value为 RLP编码过的值。

```
func TestMyTrieGetStorageData(t *testing.T) {
    ldb, _ := ethdb.NewLDBDatabase("/Users/peace/ethereum/geth/chaindata", 0, 0)
    statTr, _ :=
        trie.New(common.HexToHash("0x86bce3794034cddb3126ec488d38cb3ee840eeff4e64c3afe0ec5a5ca8b5f6ed"),
            trie.NewDatabase(ldb))

    accBytes, _ := hexutil.Decode("0x9ea9b9eeac924fd784b064dabf174a55113c4064")
    accKeyBytes := crypto.Keccak256Hash(accBytes).Bytes()
    accValueBytes, _ := statTr.TryGet(accKeyBytes)

    var acc Account
    rlp.DecodeBytes(accValueBytes, &acc)
    t.Logf("storageRoot:%s", acc.Root.String())

    storageTr, _ := trie.New(common.HexToHash(acc.Root.String()),
        trie.NewDatabase(ldb))
    index0KeyBytes, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000000")
    index0ValuesBytes, _ := storageTr.TryGet(crypto.Keccak256Hash(index0KeyBytes).Bytes())
    var storedUint uint
    rlp.DecodeBytes(index0ValuesBytes, &storedUint)
    t.Logf("storedUint: %d", storedUint)

    index1KeyBytes, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000001")
    index1ValuesBytes, _ := storageTr.TryGet(crypto.Keccak256Hash(index1KeyBytes).Bytes())
    t.Logf("raw bytes: %s", hexutil.Encode(index1ValuesBytes))
    var storedString string
    rlp.DecodeBytes(index1ValuesBytes, &storedString)
    t.Logf("storedString: %s", storedString)
}
```

运行输出以下数据 storedUint为 2018，跟合约里的数据是一致的。值得注意的是 storedString的数据后面多了一个十六进制的 26(十进制为 38)，是字符串长度 (19)的两倍，更多的细节请参见 http://solidity.readthedocs.io/en/latest/miscellaneous.html#layout-of-state-variables-in-storage。

同时，更复杂的数据结构如变长数组、map等规则会更加复杂，同时这里也忽略了一些字段打包存储等细节，但是都围绕着 storageTrie，基本原理没有改变。

```
go test -v -run TestMyTrieGetStorageData
=== RUN  TestMyTrieGetStorageData
--- PASS: TestMyTrieGetStorageData (0.01s)
    my_state_test.go:41: storageRoot:0x3fa426aa67fff5c38788fe04e4f9815652d0b259a44efed794c309577ddc2057
    my_state_test.go:49: storedUint: 2018
    my_state_test.go:53: raw bytes: 0xa04f6e657468696e672c204f6e65576f726c642100000000000000000000000026
    my_state_test.go:56: storedString: Onething, OneWorld!&
PASS
ok     github.com/ethereum/go-ethereum/core/state  0.047s
```
