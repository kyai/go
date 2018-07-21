# OpenWallet应用开发教程

## 框架库安装

```shell

go get github.com/blocktree/OpenWallet

```

## 命令行工具openw

```shell

# 创建项目
openw new MyWalletApp

# 运行项目
cd MyWalletApp
openw run

# 访问示例网站
http://127.0.0.1:8080

# 打包项目
openw pack

# 部署项目到服务器，启动后台执行
nohup ./openwApp &

```