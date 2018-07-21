# OpenWallet钱包地址规范设计

## 前言

现在每一种主链加密货币都自己的公钥哈希地址格式，

> 例如：
> bitcoin:16XyEUNgwF3KX6UyMEWtpQDWWXkmqgH7V8，
> ethereum:0x4fffdaa5dbba850ae41aa6d031a6dffd91614608。

如果你拥有多种加密货币，你就需要安装多个加密货币钱包，向别人提供充值地址收币，及仔细查看清楚别人地址，向对方支付加密货币。如果使用OpenWallet框架创建钱包应用。钱包用户就可以通过统一地址格式，接收不同币种的充值。地址如何解析，如何路由到真实的链上资产账户，如何完成最后的账户过程则交由OpenWallet框架完成处理。

## 地址编码

通常secp256k1算法生成的公钥长度33字节（压缩）或65字节（未压缩）。我们采用bitcoin的地址生成规制，对公钥进行sha256压缩，再进行ripemd160压缩，在首字节添加0x48作为标识区分，字节串末位添加4字节的checksum，最后转换为base58，生成结果：openw:WDcC8pHh9RhL3MoUJSv4fpNzZTzmVrmw97。

为了更好地方便扫描器识别，我们使用【openw:】作为协议头，组成地址：WDcC8pHh9RhL3MoUJSv4fpNzZTzmVrmw97。

关键点：

- 使用0x48在首字节作为版本标识，转为base58后地址以W开头。
- 使用openw:作为外部使用的协议。

## 地址解码

OpenWallet框架解析openw:开头的地址，在OpenWallet私链中根据地址查找到钱包账户，使用币种标识和地址在资产路由器中找到币种的区块链协议地址。通过OpenWallet框架提供功能的协议驱动完成币种公链上的转账。这样钱包使用这就不再为地址账户管理而烦恼。