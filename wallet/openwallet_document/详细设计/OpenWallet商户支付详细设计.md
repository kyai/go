# OpenWallet商户支付详细设计

在OpenWallet框架中，每一个钱包的创建都会建立这样的文件体系。

### 文件目录结构

使用wmd管理区块链资产时会创建以下文件目录，存储一些钱包相关的数据。目录用途说明如下：

| 参数变量                 | 描述                                                                         |
|--------------------------|------------------------------------------------------------------------------|
| ./data/[symbol]/key/     | 钱包根私钥文件目录，文件命名 [alias]-[ID].key                                 |
| ./data/[symbol]/address/ | 钱包批量地址导出目录，文件命名 address-yyyyMMddHHmmss.txt                     |
| ./data/[symbol]/db/      | 钱包辅助数据库缓存目录，文件命名 [alias]-[ID].db                              |
| ./data/[symbol]/backup/  | 钱包备份文件导出目录，以文件夹归档备份，文件夹命名 [alias]-[ID]-yyyyMMddHHmmss |