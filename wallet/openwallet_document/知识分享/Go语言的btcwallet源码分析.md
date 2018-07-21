# Go语言的btcwallet源码分析

## 项目地址

https://github.com/btcsuite/btcwallet

## 前言

目前开源社区上有很多不同语言领域的区块链项目和程序，尽管代码开源，但每个项目的架构和设计都各有特色，对初入门的开发者去学习区块链会有较大门槛和难度。为此原因，我开始对一些主流的开源项目进行代码分析，帮助入门开发者了解区块链是如何实现。而且在这个过程中，有可能分析源代码的漏洞哦！！！对源代码的开发者也是一种帮助。

因为主流的项目代码会不断更新，而且我也有可能理解错误的地方，大家可以在评论中指出。我会完善文档。

## 程序入口分析

> 按照官方安装教程，我们最终编译出一个btcwallet的命令行工具。我们通过btcwallet命令完成钱包的创建，所以btcwallet是钱包的入口程序，先从它的源文件进行分析。

### btcwallet.go

```go

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 可以看到walletMain是处理，钱包启动的主函数
	if err := walletMain(); err != nil {
		os.Exit(1)
	}
}

func walletMain() error {
	//loadConfig解析命令，里面完成钱包创建
	tcfg, _, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg
	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

        //下面的代码启动服务监听
        ...
}

```

### config.go

> 这个文件钱包的用户配置管理，提供初始默认配置，解析用户命令更在配置功能。

```go

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
//      1) Start with a default config with sane settings
//      2) Pre-parse the command line to check for an alternative config file
//      3) Load configuration file overwriting defaults with any specified options
//      4) Parse CLI options and overwrite/add any specified options
//
// The above results in btcwallet functioning properly without any config
// settings while still allowing the user to override settings with config files
// and command line options.  Command line options always take precedence.
func loadConfig() (*config, []string, error) {
        //初始化默认配置，解析命令
        ...
        
        //如果命令为create，创建钱包。
	if cfg.Create {
		// 钱包已经创建过了，就不能再创建
		if dbFileExists {
			err := fmt.Errorf("The wallet database file `%v` "+
				"already exists.", dbPath)
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}

		// 确保钱包对应的区块链网络数据文件夹存在
		if err := checkCreateDir(netDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}

		// 执行初始化钱包.
		if err := createWallet(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
			return nil, nil, err
		}

		// Created successfully, so exit now with success.
		os.Exit(0)
	} 
}

```

### walletsetup.go

> 这个文件负责根据用户配置创建钱包目录及钱包模型。

```go

// 创建钱包过程同时提示用户输入密码
func createWallet(cfg *config) error {
	dbDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	loader := wallet.NewLoader(activeNet.Params, dbDir)

	// 如果本地已经有一个keystore文件，尝试从keystore中恢复密钥
	netDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	keystorePath := filepath.Join(netDir, keystore.Filename)
	var legacyKeyStore *keystore.Store
	_, err := os.Stat(keystorePath)
	...

	// 命令行提示用户设置私钥密码
	reader := bufio.NewReader(os.Stdin)
	privPass, err := prompt.PrivatePass(reader, legacyKeyStore)

	// 如果keystore存在，且密码通过，解密keystore
	if legacyKeyStore != nil {
		err = legacyKeyStore.Unlock(privPass)
        //处理keystore的地址迁移
		...
	}

    // 命令行提示用户是否需要设置一个密码来加密公钥，如果不设置
    // 默认使用"public"作为密码
	pubPass, err := prompt.PublicPass(reader, privPass,
		[]byte(wallet.InsecurePubPassphrase), []byte(cfg.WalletPass))

	// 命令行提示用户是否需要设置已有的种子，用户不提供，则随机生成新的
	seed, err := prompt.Seed(reader)

    //得到公钥密码，私钥密码，种子后，创建钱包模型
	fmt.Println("Creating the wallet...")
	w, err := loader.CreateNewWallet(pubPass, privPass, seed)

    ...
}

```

## wallet包分析

> 这个包负责分层确定性模型下的钱包抽象实现。实现诸如，创建账户，查询账户余额，查询账户交易记录等等。

### loader.go

Loader是一个加载工具，主要用途为：创建新钱包，打开新钱包，并回调给其他子线程加载钱包。

```go

// 创建新钱包
func (l *Loader) CreateNewWallet(pubPassphrase, privPassphrase, seed []byte) (*Wallet, error) {
    
    ...

    // 拼接一个钱包数据路径，btcwallet用的数据是boltdb，是一个go原生写的数据库程序。
	dbPath := filepath.Join(l.dbDirPath, walletDbName)
	
    ...
    // 在该路径上创建文件数据库
	db, err := walletdb.Create("bdb", dbPath)

	// 创建钱包
	err = Create(db, pubPassphrase, privPassphrase, seed, l.chainParams)

	// 打开钱包
	w, err := Open(db, pubPassphrase, nil, l.chainParams)
	w.Start()

	l.onLoaded(w, db)
	return w, nil
}

```

### wallet.go

这个文件比较负责，创建钱包流程分开了几个函数

```go

// 创建钱包
func Create(db walletdb.DB, pubPass, privPass, seed []byte, params *chaincfg.Params) error {
    // 如果用户没有提供种子，程序会随机创建种子。
    // 这个种子是钱包的一切，只要被被黑客拿到，无论你的私钥密码或公钥密码没有泄露.
    // 黑客都可以按分层确定性模型，用户暴力算法找到有utxo的私钥。所以钱包程序尽量地少与seed交互。
	if seed == nil {
		hdSeed, err := hdkeychain.GenerateSeed(
			hdkeychain.RecommendedSeedLen)
		if err != nil {
			return err
		}
		seed = hdSeed
    }
    ...
	// boltdb是一款KV数据库，所以每个钱包都有自己的索引键。
	addrMgrNamespace, err := db.Namespace(waddrmgrNamespaceKey)
	
    //这里才是正真创建钱包的主逻辑，我们跳转到waddrmgr进行分析
	err = waddrmgr.Create(addrMgrNamespace, seed, pubPass, privPass,
		params, nil)

	// 创建交易单表
	txMgrNamespace, err := db.Namespace(wtxmgrNamespaceKey)
	
	return wtxmgr.Create(txMgrNamespace)
}

```

## waddrmgr包分析

> 这个包负责按照BIP32，BIP44，即分层确定性模型管理钱包。

### manager.go

```go

// 创建一个地址管理器，使用种子创建根私有和根公钥，
// 所有私钥包含“根”和“后代”都由用户设置的私钥密码加密保护，
// 所有公钥包含“根”和“后代”都由用户设置的公钥密码加密保护，
func Create(namespace walletdb.Namespace, seed, pubPassphrase, privPassphrase []byte, chainParams *chaincfg.Params, config *ScryptOptions) error {
    
	// 通过种子生成根私钥
	root, err := hdkeychain.NewMaster(seed, chainParams)

	// 根，衍生后代作为BTC币种，BIP44硬性定BTC主网在m/44'/0'，更新查看BIP32和44
	coinTypeKeyPriv, err := deriveCoinTypeKey(root, chainParams.HDCoinType)
	
	// 衍生索引位0作为默认第一个账户
	acctKeyPriv, err := deriveAccountKey(coinTypeKeyPriv, 0)

	// 导出公钥
	acctKeyPub, err := acctKeyPriv.Neuter()


	/*
	btcwallet采用这样的加密方案：
	1. 随机生成PUB密钥，加密种子生成的根公钥，生成PUBENC密文。
	2. 随机生成PRV密钥，加密种子生成的根私钥，生成PRVENC密文。
	3. 用户设置公钥密码计算出UPUB密钥，加密PUB密钥，生成UPUBENC密文。
	4. 用户设置私有密码计算出UPRV密钥，加密PRV密钥，生成UPRVENC密文。
	5. 所有密文都保存到数据库。

	这样多了一层加密的作用是可以避免黑客暴力破解了用户密码后，仍然不能直接解密出私钥。

	解密过程就是：
	1. 用户输入密码，计算出UPUB/UPRV密钥，解密UPUBENC/UPRVENC得出出PUB/PRV密钥。
	2. PUB/PRV密钥解密PUBENC/PRVENC密文，得出根公钥/根私钥。
	3. 得到根公钥/根私钥后，可以按分层确定性模型计算后代。
	*/

	// 用户设置公钥密码计算出UPUB密钥
	masterKeyPub, err := newSecretKey(&pubPassphrase, config)
    
    // 用户设置私有密码计算出UPRV密钥
	masterKeyPriv, err := newSecretKey(&privPassphrase, config)

	// 随机生成“盐"，好像没有用到这个。
	var privPassphraseSalt [saltSize]byte
	_, err = rand.Read(privPassphraseSalt[:])

    //随机生成PUB密钥
	cryptoKeyPub, err := newCryptoKey()
    
    //随机生成PRV密钥
	cryptoKeyPriv, err := newCryptoKey()
	
	// 随机创建赎回脚本密码
	cryptoKeyScript, err := newCryptoKey()

	// 加密PUB密钥，生成UPUBENC密文
	cryptoKeyPubEnc, err := masterKeyPub.Encrypt(cryptoKeyPub.Bytes())
	
	// 加密PRV密钥，生成UPRVENC密文
	cryptoKeyPrivEnc, err := masterKeyPriv.Encrypt(cryptoKeyPriv.Bytes())

	// 加密赎回脚本，生成密文
	cryptoKeyScriptEnc, err := masterKeyPriv.Encrypt(cryptoKeyScript.Bytes())
	
	// 根私钥导出根公钥
	coinTypeKeyPub, err := coinTypeKeyPriv.Neuter()

	//加密种子生成的根公钥，生成PUBENC密文
	coinTypePubEnc, err := cryptoKeyPub.Encrypt([]byte(coinTypeKeyPub.String()))
	
	// 加密种子生成的根私钥，生成PRVENC密文
	coinTypePrivEnc, err := cryptoKeyPriv.Encrypt([]byte(coinTypeKeyPriv.String()))
	
	// 加密钱包首个账户
	acctPubEnc, err := cryptoKeyPub.Encrypt([]byte(acctKeyPub.String()))

	// 加密钱包首个账户
	acctPrivEnc, err := cryptoKeyPriv.Encrypt([]byte(acctKeyPriv.String()))

	...

	// 在一个事务中，保存所有密文
	err = namespace.Update(func(tx walletdb.Tx) error {
		// 只保存密文，记住不保存明文
		pubParams := masterKeyPub.Marshal()
		privParams := masterKeyPriv.Marshal()
		err = putMasterKeyParams(tx, pubParams, privParams)
		err = putCryptoKeys(tx, cryptoKeyPubEnc, cryptoKeyPrivEnc,
			cryptoKeyScriptEnc)
		err = putCoinTypeKeys(tx, coinTypePubEnc, coinTypePrivEnc)
		err = putWatchingOnly(tx, false)
		// Save the initial synced to state.
		err = putSyncedTo(tx, &syncInfo.syncedTo)
		err = putStartBlock(tx, &syncInfo.startBlock)
		// Save the initial recent blocks state.
		err = putRecentBlocks(tx, recentHeight, recentHashes)
		// Save the information for the imported account to the database.
		err = putAccountInfo(tx, ImportedAddrAccount, nil,
			nil, 0, 0, ImportedAddrAccountName)
		err = putAccountInfo(tx, DefaultAccountNum, acctPubEnc,
			acctPrivEnc, 0, 0, defaultAccountName)
		return err
	})
	return nil
}

```

> 到此为止，我们的钱包已经创建成功了。我们的默认第一个账户在m/44'/0'。我们在看看创建新账户的代码，看看的规则。

```go

// 创建新账户，先检查账户表仲是否存在同名账户，因为我们之前创建钱包时已通过用户密码加密了一大堆密文。
// 所以创建新账户需要解锁根公钥。
func (m *Manager) NewAccount(name string) (uint32, error) {
	
	...

	// 检查账户名是否唯一
	_, err := m.lookupAccount(name)
	if err == nil {
		str := fmt.Sprintf("account with the same name already exists")
		return 0, managerError(ErrDuplicateAccount, str, err)
	}

	var account uint32
	var coinTypePrivEnc []byte

	// 查找最后一个账户索引位，在事务中创建新用户
	err = m.namespace.Update(func(tx walletdb.Tx) error {
		var err error
		// 获取当前最后的账户位置
		account, err = fetchLastAccount(tx)
		if err != nil {
			return err
		}
		// 递增一位
		account++

		// 获取根私钥的密文
		_, coinTypePrivEnc, err = fetchCoinTypeKeys(tx)

		// 解密出根私钥
		serializedKeyPriv, err := m.cryptoKeyPriv.Decrypt(coinTypePrivEnc)
		coinTypeKeyPriv, err := hdkeychain.NewKeyFromString(string(serializedKeyPriv))
		zero.Bytes(serializedKeyPriv)
		
		// 衍生新账户的私钥，假如我们初始钱包后，再创建新账户，新账户在m/44'/1'。
		acctKeyPriv, err := deriveAccountKey(coinTypeKeyPriv, account)
		coinTypeKeyPriv.Zero()
		
		// 导出新账户公钥
		acctKeyPub, err := acctKeyPriv.Neuter()
		
		// 加密密钥对
		acctPubEnc, err := m.cryptoKeyPub.Encrypt([]byte(acctKeyPub.String()))
		acctPrivEnc, err := m.cryptoKeyPriv.Encrypt([]byte(acctKeyPriv.String()))
		
		// 保存新账户
		err = putAccountInfo(tx, account, acctPubEnc, acctPrivEnc, 0, 0, name)

		// 记录新账户的索引位为最后位置
		if err := putLastAccount(tx, account); err != nil {
			return err
		}
		return nil
	})
	return account, err
}

```

> 还没写完。。。To be continue