# 数字加密相关知识
- [数字加密相关知识](#数字加密相关知识)
    - [非对称加密](#非对称加密)
        - [椭圆曲线加密](#椭圆曲线加密)
    - [公钥与私钥](#公钥与私钥)
    - [数字签名](#数字签名)
    - [数字证书](#数字证书)
    - [Merkle Tree](#merkle-tree)
    - [数字签名扩展](#数字签名扩展)
        - [数字签名算法](#数字签名算法)
            - [RSA数字签名算法](#rsa数字签名算法)
            - [DSA数字签名算法](#dsa数字签名算法)
            - [ECDSA椭圆曲线数字签名算法](#ecdsa椭圆曲线数字签名算法)
        - [盲签名](#盲签名)
        - [多重签名](#多重签名)
        - [群签名](#群签名)
        - [环签名](#环签名)
            - [环签名满足的性质](#环签名满足的性质)
            - [环签名实现](#环签名实现)
            - [环签名和群签名的比较](#环签名和群签名的比较)
            - [算法实现](#算法实现)
- [Reference](#reference)


## 非对称加密

`对称加密`指加密和解密使用`相同密钥`的加密算法。它要求发送方和接收方在安全通信之前，商定一个密钥，泄漏密钥就意味着任何人都可以对他们发送或接收的消息解密。

每对用户每次使用对称加密算法时，都需要使用其他人不知道的惟一钥匙，这会使得发收信双方所拥有的钥匙数量呈几何级数增长，密钥管理成为用户的负担。同时，对称加密算法能够提供加密和认证却缺乏了签名功能，使得使用范围有所缩小。

具体算法：DES算法，3DES算法，TDEA算法，Blowfish算法，RC5算法，IDEA算法。

而`非对称加密`指加解密密钥不相关，典型如RSA、EIGamal、椭圆曲线算法。

### 椭圆曲线加密

公开密钥算法总是要基于一个数学上的难题。比如RSA 依据的是：给定两个素数p、q 很容易相乘得到n，而对n进行因式分解却相对困难。那椭圆曲线上有什么难题呢？

![Trapdoor function](https://upload.wikimedia.org/wikipedia/commons/thumb/8/8f/Trapdoor_permutation.svg/300px-Trapdoor_permutation.svg.png)

考虑如下等式：  

    K=kG  [其中 K,G为Ep(a,b)上的点，k为小于n（n是点G的阶）的整数]  

`不难发现，给定k和G，根据加法法则，计算K很容易；但给定K和G，求k就相对困难了。`  
这就是椭圆曲线加密算法采用的难题，我们把点G称为基点（base point）。  

现在我们描述一个利用椭圆曲线进行加密通信的过程：  

1.  用户A选定一条椭圆曲线Ep(a,b)，并取椭圆曲线上一点，作为基点G。
2.  用户A选择一个私有密钥k，并生成`公开密钥K=kG`。
3.  用户A将Ep(a,b)和点K，G传给用户B。
4.  用户B接到信息后 ，将待传输的明文编码到Ep(a,b)上一点M（编码方法很多，这里不作讨论），并产生一个随机整数r
5.  用户B计算点C1=M+rK；C2=rG。
6.  用户B将C1、C2传给用户A。
7.  用户A接到信息后，计算C1-kC2，结果就是点M。因为
    C1-kC2=M+rK-k(rG)=M+rK-r(kG)=M  
    再对点M进行解码就可以得到明文。

在这个加密通信中，如果有一个偷窥者H，他只能看到Ep(a,b)、K、G、C1、C2，`而通过K、G 求k 或通过C2、G求r 都是相对困难的`。因此，H无法得到A、B间传送的明文信息。

## 公钥与私钥

公开密钥与私有密钥是一对，如果用公开密钥对数据进行加密，只有用对应的私有密钥才能解密;如果用私有密钥对数据进行加密，那么只有用对应的公开密钥才能解密。  
通过公钥是无法（或极其困难）推算出私钥的。**注意这里公钥与私钥都可以用来加密，只是私钥是自己保存，而公钥是公开的。**

比如A发信息给B，就用B的公钥加密信息，然后发给B，B用自己的私钥解密，就可以看到信息内容。

## 数字签名

A把要加密的内容用hash函数生成摘要digest，再用自己的私钥对digest加密，生成数字签名signature。连同加密的内容一起发给B。  

B收到后，对**摘要digest（用A的公钥解密）**和**内容（用自己的私钥解密）** 都解密，再对内容使用相同hash，看得到的digest是否相同，相同，说明发送的内容没有被修改。  

>但是如果这里用B的公钥来加密摘要digest，B收到后用自己的私钥解密，这不是也可以验证内容是否被篡改么？

如果仅仅从防篡改的角度来讲确实可以，但是这样无法验证这些内容是谁发来的！**所以用A的私钥来加密摘要digest，相当于A用自己的私钥给这个摘要digest进行签名，B收到后用A的公钥对digest进行解密还能验证这是不是A发来的内容**。

但是这里有个`潜在的问题`。  

>**如果B存储的A的公钥被C替换成了C的公钥，那么C就可以冒充A和B进行通信，而B却完全不知道。**

## 数字证书

证书中心用自己的私钥对`A的公钥和一些相关信息`一起加密，形成`数字证书`。  
A在发送内容的同时，在数字签名后再附上数字证书。  
B收到后，先用CA的公钥解密数字证书，得到A真正的公钥，再用A的公钥来验证签名是否是A的签名。  

B可以每次都到CA的网站上（或者什么别的官方途径）获得CA的公钥。

这么做的目的是为了验证：
1. 确认该信息确实是A所发；
2. 确认A发出的信息是完整的。

*  **公钥防泄漏,私钥防篡改、防假冒**  
    *  B收到后，只有用B自己的私钥才能解密内容，别人是无法解密的。`防泄漏`  
    *  再用上述数字证书来验证数字签名是否来自A，发送内容有没有被篡改。`防篡改`  
数字证书一般挂靠在可信任的机构，无法篡改和伪造。

## Merkle Tree

默克尔树，又叫哈希树，由**一个**root节点，**一组**中间节点和**一组**叶节点组成。  
叶节点包含存储数据或者其哈希值，中间节点和root节点都是其孩子的hash值。  

![markle tree](https://github.com/yjjnls/Notes/blob/master/block%20chain/Basic/img/markle%20tree.jpg)

应用：  
1\.  快速比较数据，两个默克尔树的根节点相同，那么其所代表的数据必然相同  
2\.  快速定位修改，比如上面D1数据被修改，可通过root->N4->N1，快速定位到发生改变的D1  
3\.  零知识证明，比如要证明某个数据中包含D0，那就构造一个默克尔树，公开root、N4、N1、N0，D0拥有者可以检测到D0存在，但不知道其他内容。（D0拥有者可以看到hash值，但看不到完整的数据内容）（比如用户可以查找自己的money是否在交易所的总备用金中，而不必知道其余用户的money信息；或者p2p下载中，文件切片成小块，下载一个分支后就可以验证该分支的数据是否正确，定位错误数据块重新下载或者继续下载下一个分支数据。）

## 数字签名扩展
### 数字签名算法
常见的数字签名算法主要有RSA、DSA、ECDSA三种。

#### RSA数字签名算法
RSA是目前计算机密码学中最经典算法，也是目前为止使用最广泛的数字签名算法，RSA数字签名算法的密钥实现与RSA的加密算法是一样的，算法的名称都叫RSA。密钥的产生和转换都是一样的，包括在售的所有SSL数字证书、代码签名证书、文档签名以及邮件签名大多都采用RSA算法进行加密。

RSA数字签名算法主要包括MD和SHA两种算法，例如我们熟知的MD5和SHA-256即是这两种算法中的一类。

#### DSA数字签名算法
DSA全称Digital Signature Algorithm，DSA只是一种算法，和RSA不同之处在于它不能用作加密和解密，也不能进行密钥交换，只用于签名，所以它比RSA要快很多，其安全性与RSA相比差不多。DSA的一个重要特点是两个素数公开，这样，当使用别人的p和q时，即使不知道私钥，你也能确认它们是否是随机产生的，还是作了手脚。RSA算法却做不到。

DSA的整个签名算法流程如下：

1.  发送方使用SHA-1和SHA-2编码将发送内容加密产生数字摘要；
2.  发送方用自己的专用密钥对摘要进行再次加密得到数字签名；
3.  发送方将原文和加密后的摘要传给接收方；
4.  接收方使用发送方提供的密钥对进行解密 ，同时对收到的内容用SHA-1/SHA-2编码加密产生同样的摘要；
5.  接收方再将解密后的摘要和4步骤中加密产生的摘要进行比对，如果两者一至，则说明传输过程的信息没有被破坏和篡改，否则传输信息则不安全。

#### ECDSA椭圆曲线数字签名算法
ECDSA是用于数字签名，是ECC与DSA的结合，整个签名过程与DSA类似，**所不一样的是签名中采取的算法为ECC**，最后签名出来的值也是分为r,s。而ECC（全称Elliptic Curves Cryptography）是一种椭圆曲线密码编码学。

ECDH每次用一个固定的DH key，导致不能向前保密（forward secrecy），所以一般都是用ECDHE（ephemeral）或其他版本的ECDH算法。ECDH则是基于ECC的DH（Diffie-Hellman）密钥交换算法。

ECC与RSA 相比，有以下的优点：

1.  相同密钥长度下，安全性能更高，如160位ECC已经与1024位RSA、DSA有相同的安全强度。
2.  计算量小，处理速度快，在私钥的处理速度上（解密和签名），ECC远 比RSA、DSA快得多。
3.  存储空间占用小 ECC的密钥尺寸和系统参数与RSA、DSA相比要小得多， 所以占用的存储空间小得多。
4.  带宽要求低使得ECC具有广泛得应用前景。

### 盲签名
前文说到了数字签名，但如果A想要B来对消息签名，但是又不想让B知道该消息的内容，该如何做呢？这里就需要用到盲签名了。

盲签名具有以下特点：
1.  签名者对其所签署的消息是不可见的，即签名者不知道他所签署消息的具体内容。
2.  签名消息不可追踪，即当签名消息被公布后，签名者无法知道这是他哪次的签署的。

传统RSA加解密的结局方案为：

1.  加密: ![](http://upload-images.jianshu.io/upload_images/11336404-be53f4cdd8eff0f1.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

2.  解密:![](http://upload-images.jianshu.io/upload_images/11336404-0c87f6ebf1f4b7b8.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

而盲签名则按照如下步骤来进行：
1.  接收者首先将待签数据进行盲变换，把变换后的盲数据发给签名者。  
![](http://upload-images.jianshu.io/upload_images/11336404-c22f3f330c202a4d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)  
这里一般用对称加密，所以：  
![](http://upload-images.jianshu.io/upload_images/11336404-e2f7a82497d99ebf.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

2.  经签名者签名后再发给接收者。  
![](http://upload-images.jianshu.io/upload_images/11336404-e858b0db3ba78057.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)  

3.  接收者对签名再作去盲变换，得出的便是签名者对原数据的盲签名。   
![](http://upload-images.jianshu.io/upload_images/11336404-717c9c59aaf9e53a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)   
![](http://upload-images.jianshu.io/upload_images/11336404-074b0968d84944d8.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


流程如下：  
![](http://upload-images.jianshu.io/upload_images/11336404-1c5466a4f060888b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

在这个过程中，B无法得知MSG是什么，B也不知道自己什么时候对MSG做了签名，即若B进行了多次签名，当公布出某一具体的MSG-signature时，B并不知道这个签名是自己在那一次进行签署的。


背景：
一般的签名，签名者对自己发出的签名，必须是记得的，比如，在何时何地对谁发的，他自己可以记下来。但是，如果把签名看作是电子现金的话，就涉及到一个匿名性的问题用实际钞票的时候，钞票上有没有写你的名字？当然没有。那我也不希望，银行通过追踪自己发出签名，来获得用户的消费情况。于是就设计出盲签名。

### 多重签名
多重签名技术（multisig）就是多个用户同时对一个数字资产进行签名才有效。

如果一个地址只能由一个私钥签名和支付，表现形式就是1/1；而多重签名的表现形式是m/n，也就是说一共n个私钥可以给一个账户签名，而当m个地址签名时，就可以支付一笔交易。所以，m一定是小于等于n的。

### 群签名
在一个群签名方案中，一个群体中的任意一个成员可以以匿名的方式代表整个群体对消息进行签名。与其他数字签名一样，群签名是可以公开验证的，而且可以只用单个群公钥来验证。

流程：
1. 初始化：群管理者建立群资源，生成对应的群公钥（Group Public Key）和群私钥（Group Private Key）群公钥对整个系统中的所有用户公开，比如群成员、验证者等。
2. 成员加入：在用户加入群的时候，群管理者颁发群证书（Group Certificate）给群成员。
3. 签名：群成员利用获得的群证书签署文件，生成群签名
4. 验证：同时验证者利用**群公钥**仅可以验证所得群签名的正确性，但不能确定群中的正式签署者。
5. 打开：**群管理者利用群私钥可以对群用户生成的群签名进行追踪，并暴露签署者身份**。


### 环签名
是一种简化的群签名，只有环成员没有管理者，不需要环成员间的合作。环签名方案中签名者首先选定一个临时的签名者集合,集合中包括签名者。然后签名者利用自己的私钥和签名集合中其他人的公钥就可以独立的产生签名,而无需他人的帮助。签名者集合中的成员可能并不知道自己被包含在其中。

#### 环签名满足的性质
1.  无条件匿名性: 攻击者无法确定签名是由环中哪个成员生成,即使在获得环成员私钥的情况下,概率也不超过1/n。
2.  正确性: 签名必需能被所有其他人验证。
3.  不可伪造性: 环中其他成员不能伪造真实签名者签名,外部攻击者即使在获得某个有效环签名的基础上,也不能为消息m伪造一个签名。

#### 环签名实现
1.  密钥生成： 为环中每个成员产生一个密钥对(公钥PKi,私钥SKi)。
2.  签名： 签名者用自己的私钥和任意n个环成员(包括自己)的公钥为消息m生成签名a。
3.  签名验证： 验证者根据环签名和消息m,验证签名是否为环中成员所签,如果有效就接收,否则丢弃。

#### 环签名和群签名的比较
1.  匿名性：都是一种个体代表群体签名的体制，验证者能验证签名为群体中某个成员所签，但并不能知道为哪个成员，以达到签名者匿名的作用。
2.  可追踪性：**群签名中，群管理员的存在保证了签名的可追踪性**。群管理员可以撤销签名，揭露真正的签名者。环签名本身无法揭示签名者,除非签名者本身想暴露或者在签名中添加额外的信息。提出了一个可验证的环签名方案,方案中真实签名者希望验证者知道自己的身份,此时真实签名者可以通过透露自己掌握的秘密信息来证实自己的身份。
3.  管理系统：群签名由群管理员管理，环签名不需要管理，签名者只有选择一个可能的签名者集合，获得其公钥，然后公布这个集合即可，所有成员平等。

#### 算法实现
```python
import os, hashlib, random, Crypto.PublicKey.RSA

class ring:
    def __init__(self, k, L=1024):
        self.k = k
        self.l = L
        self.n = len(k)
        self.q = 1 << (L - 1)

    def sign(self, m, z):
        self.permut(m)
        s = [None] * self.n
        u = random.randint(0, self.q)
        c = v = self.E(u) 
        for i in (range(z+1, self.n) + range(z)):
            s[i] = random.randint(0, self.q)
            e = self.g(s[i], self.k[i].e, self.k[i].n)
            v = self.E(v^e) 
            if (i+1) % self.n == 0:
                c = v
        s[z] = self.g(v^u, self.k[z].d, self.k[z].n)
        return [c] + s

    def verify(self, m, X):
        self.permut(m)
        def _f(i):
            return self.g(X[i+1], self.k[i].e, self.k[i].n)
        y = map(_f, range(len(X)-1))
        def _g(x, i):
            return self.E(x^y[i])
        r = reduce(_g, range(self.n), X[0])
        return r == X[0]

    def permut(self, m):
        self.p = int(hashlib.sha1('%s' % m).hexdigest(),16)

    def E(self, x): 
        msg = '%s%s' % (x, self.p)
        return int(hashlib.sha1(msg).hexdigest(), 16)

    def g(self, x, e, n):
        q, r = divmod(x, n)
        if ((q + 1) * n) <= ((1 << self.l) - 1):
            rslt = q * n + pow(r, e, n)
        else:
            rslt = x
        return rslt
```

签名并验证两个由4个用户组成的环签名消息：

```python
size = 4
msg1, msg2 = 'hello', 'world!'

def _rn(_):
  return Crypto.PublicKey.RSA.generate(1024, os.urandom)

key = map(_rn, range(size))
r = ring(key)
for i in range(size):
    s1 = r.sign(msg1, i)
    s2 = r.sign(msg2, i)
    assert r.verify(msg1, s1) and r.verify(msg2, s2) and not r.verify(msg1, s2)
```


# Reference
1. [数字签名是什么？](http://www.ruanyifeng.com/blog/2011/08/what_is_a_digital_signature.html)
2. [Digital Signatures](http://learnmeabitcoin.com/guide/digital_signatures)
3. [比特币背后的密码学原理](https://www.jianshu.com/p/225ff9439132)
4. [《区块链原理设计与应用》第5章：密码学与安全技术](https://github.com/yjjnls/books/blob/master/block%20chain/%E5%8C%BA%E5%9D%97%E9%93%BE%E5%8E%9F%E7%90%86%E3%80%81%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%BA%94%E7%94%A8.pdf)
5. [Secure Hash Algorithms](https://en.wikipedia.org/wiki/Secure_Hash_Algorithms)
6. [Digital Signature Algorithm](https://en.wikipedia.org/wiki/Digital_Signature_Algorithm)
7. [用实例给新手讲解RSA加密算法](https://www.cnblogs.com/jiftle/p/7903762.html)
8. [用Go语言实现环签名的签名和验证（一）](https://blog.csdn.net/qq_38200023/article/details/79902049)
9. [用Go语言实现环签名的签名和验证（二）](https://blog.csdn.net/qq_38200023/article/details/79912530)