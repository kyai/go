/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package bitcoin

import (
	"github.com/asdine/storm"
	"path/filepath"
)

/*
	步骤：
	1.添加需要扫块的钱包，及传入初始高度，-1为本地高度。
	2.获取已扫描的本地高度。
	3.获取高度+1的区块hash，通过区块链hash获取区块链数据，获取mempool数据。
	4.判断区块链的父区块hash是否与本地上一区块hash一致。
	5.解析新区块链的交易单数组。
	6.遍历交易单结构，检查每个接收地址是否存在钱包的地址表中
	7.检查地址是否合法，存在地址表，生成充值记录。
	8.定时程推送充值记录到钱包的充值通道。先检查交易hash是否存在区块中。
	9.接口返回确认，标记充值记录已确认。
*/

type BTCBlockScanner struct {
	addressInScanning  map[string]string //加入扫描的钱包资产账户
	CurrentBlockHeight uint64            //当前区块高度
	isScanning bool
}

func NewBTCBlockScanner() *BTCBlockScanner {
	bs := BTCBlockScanner{}
	bs.addressInScanning = make(map[string]string)
	return &bs
}

func (bs *BTCBlockScanner) addAddress(address string, accountID string) {
	bs.addressInScanning[address] = accountID
}

func (bs *BTCBlockScanner) start() {

}

func (bs *BTCBlockScanner) scanning() {
	//获取本地区块高度

}

func (bs *BTCBlockScanner) currentBlockHeight() uint64 {

	var (
		blockHeight uint64 = 0
	)
	//获取本地区块高度
	db, err := storm.Open(filepath.Join(dbPath, blockchainFile))
	if err != nil {
		return 0
	}
	defer db.Close()

	db.Get("blockchain","blockHeight", &blockHeight)

	//如果本地没有记录，查询接口的高度
	if blockHeight == 0 {

	}

	return blockHeight
}

func GetBlockHeight() {

}