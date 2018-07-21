# HDKeystore详细设计

## 修订信息

| 版本  | 时间       | 修订人 | 修订内容       |
|-------|------------|--------|----------------|
| 1.0.0 | 2018-04-27 | 麦志泉 | 创始           |
| 1.0.1 | 2018-07-09 | 麦志泉 | 调整多账户说明 |

## 前言

keystore这种保存密钥方案，在区块链领域最早用在ethereum。关于keystore的介绍可以查看文章[《什么是以太坊私钥储存（Keystore）文件？》](https://ethfans.org/posts/what-is-an-ethereum-keystore-file)。

因为区块链是建立在去中心化的基础上，你的账户不是由中心服务器创建出来的，而是由你自己随机创造出一个32字节的密钥，随机性取决于熵值。这方面又得讲到芯片和随机算法了，目前不是我们讨论的范畴，我们只需要设计一个适合我们用的密钥保存方案就是。

## keystore特点

`优点`

1. 它可以基于用户自定义密码来加密私钥。使用时，强制要求用户输入密码解锁。
1. 用户通过密码解锁私钥，可以通过mac来检查密码是否正确。
1. 它保存到硬盘，可以给用户自行备份到其它地方，和恢复到其它设备上。

`可能的风险`

1. 黑客获得keystore后，可通过暴力破解简单密码。
1. 一个账户对应一个keystore，对于中心化的托管服务，他有100万用户就需要保存100万个keystore。这样备份起来会相当麻烦，容易丢失。

## HDKeystore介绍

HDKeystore全名Hierarchical Deterministic Keystore，即分层确定性keystore。这个想法来自于BIP32。
HDKeystore继承原有的keystore，增加以下几点优势：

1. 增加多账户模型记录方案。
1. 引入工作量证明机制，防止用户设置的简单密码被暴力破解。
1. 生成助记词短语，包含多账户使用数量的记录。

### 多账户记录方案

根据BIP32的设计思想，我们在keystore中提供rootpath字段记录，有开发者自定义子账户衍生的路径。
rootpub记录生成子账户的根公钥。rootid为通过[!《OpenWallet钱包地址规范设计》](./详细设计/OpenWallet钱包地址规范设计.md)生成的钱包ID。还有alias字段，可由用户自定义钱包别名。

例如采用BIP44的模型设计：

> m / purpose' / coin_type' / account
> 自定义账户路径：0/44'/88'/，钱包别名："Wallet Test"
> 生成结果：
> "rootpath": "m/44'/88'"
> "alias": "Wallet Test"
> "rootpub": "xpub69wpTMqtZgzro3nfYvHa2tiWeE4oaHpNNA2gFoGEatpridLxKMsM2kBknMzS8W3XXXKsyHj9nhtJm3XwzbWrzTmoxMmw9ddv4qCgaQiNrQu"
> "rootid": "W9Kg8MCnREhQiV1BG2b7xzT2gbXiDmCuxu"

### 用户密码强化方案

`加密过程`

1. 我们定义一个工作量算法FN(p string)。这个算法会对密码进行多轮的串行加密，大概需要6秒，才能得到结果。
1. 用户输入密码$p$，对$p$定向生成一个不可逆的扩展密钥$P_{EX}$。
1. 传入$p$运行$FN(p)$，得到结果长度为32字节的$V$。
1. $P_{EX}$对$V$进行RSA加密得到$V_{ENC}$，读取keystore地址公钥$Addr$，把$(Addr, V_{ENC})$存放到安全的KV缓存中，缓存有定时清理机制。
1. $V$作为AES的密钥，加密私钥$K$得到$K_{ENC}$.
1. 创建HDKeystore文件，继承keystore相关的参数记录，增加记录FN函数需要的常数值，如计算轮数。

`解密过程`

1. 用户输入密码$p$，对$p$定向生成一个不可逆的扩展密钥$P_{EX}$。
1. 使用$Addr$在KV缓存中查找$V_{ENC}$。
1. $V_{ENC}$存在，$P_{EX}$解密$V_{ENC}$得到$V$；
1. $V_{ENC}$不存在，传入$p$运行$FN(p)$，得到结果长度为32字节的$V$。
1. 使用$V$解密keystore中的$K_{ENC}$，得到私钥$K$。

强化后的方案有以下几点优势：

1. 用一个32字节的$V$进行加密私钥，而不直接使用用户密码进行加密的原因是，这样可以避免弱的密码被黑客容易暴力破解得到。32字节长度的密码破解难度相等于破解私钥的难度，所以黑客有能力暴力破解$V$，等于私钥也一样有能力破解。
1. 为了得到$V$，密码$p$要通过FN耗时6秒计算。假设黑客的硬件性能和用户一样，黑客得到了用户keystore，但没有用户密码，这样黑客的每次密码尝试都需要CPU满载计算6秒时间。对于一个密码强度为8位的纯数字，CPU全速运行破解，需要3472.22天。（通过异或计算输入值可以减少一半时间）
1. 用户在第一次输入密码正确，记录$(Addr, V_{ENC})$到KV缓存后，用户下一次输入正确密码就不需要做工作量证明。就算输入错误密码，也无法解密正确的私钥，keystore能给出密码错误提示。
1. 若用户迁移keystore到新设备，新设备运算能力更强，程序可以在第一次解锁keystore成功后，评估解锁时间快了，可以更新运算轮数，已保证6秒，更新keystore相关的工作量算法FN常数和$K_{ENC}$。这样无论以后硬件怎样先进，安全性可以跟上。

> 引入工作量证明所产生的耗时不必固定某个值，大概跟手机和邮箱动态验证码收取到的时间差不多。这样不会给用户第一次输入密码时造成不好的用户体验。

### 数据结构

```json

{
	"alias": "Wallet Test",
	"rootpub": "xpub69wpTMqtZgzro3nfYvHa2tiWeE4oaHpNNA2gFoGEatpridLxKMsM2kBknMzS8W3XXXKsyHj9nhtJm3XwzbWrzTmoxMmw9ddv4qCgaQiNrQu",
	"rootid": "W9Kg8MCnREhQiV1BG2b7xzT2gbXiDmCuxu",
	"crypto": {
		"cipher": "aes-128-ctr",
		"ciphertext": "d4875065acfe6c1d2e3241be3d7a19197956993cfaf34552ab213ed399b2a242",
		"cipherparams": {
			"iv": "94d16f1f09da384e2b9b0b0f8f7a552c"
		},
		"kdf": "scrypt",
		"kdfparams": {
			"dklen": 32,
			"n": 262144,
			"p": 1,
			"r": 8,
			"salt": "0ac805a8db95530073c7e62cc183ccf2c1b110e25e791583447adc68a00dc769"
		},
		"mac": "c5117b139ca3d568db2256b3c3236c1f1d42017ef5fade1062158d12f2b433fc"
	},
	"rootpath": "m/44'/88'",
	"version": 1,
    "mining": {
        "cycle": 2000000,
        "algorithm": "sha256",
        "delay": 6
    }
}

```