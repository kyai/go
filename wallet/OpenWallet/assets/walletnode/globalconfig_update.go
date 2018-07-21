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

package walletnode

import (
	// "encoding/json"
	"fmt"
	bconfig "github.com/astaxie/beego/config"
	"github.com/pkg/errors"
	"path/filepath"
	s "strings"
)

// Update <Symbol>.ini file
func updateConfig(symbol string) error {
	// defer func() { return recover() }()
	// if err := loadConfig(symbol); err != nil {
	// 	return err
	// }
	configFilePath, _ := filepath.Abs("conf")
	configFileName := s.ToUpper(symbol) + ".ini"
	absFile := filepath.Join(configFilePath, configFileName)

	c, err := bconfig.NewConfig("ini", absFile)
	if err != nil {
		return (errors.New(fmt.Sprintf("Load Config Failed: %s", err)))
	}

	err = c.Set("serverAddr", serverAddr)
	err = c.Set("serverPort", serverPort)
	err = c.Set("isTestNet", isTestNet)
	err = c.Set("rpcuser", rpcUser)
	err = c.Set("rpcpassword", rpcPassword)
	err = c.Set("apiurl", apiURL)
	err = c.Set("mainnetdatapath", mainNetDataPath)
	err = c.Set("testnetdatapath", testNetDataPath)
	err = c.Set("threshold", isTestNet)

	if err := c.SaveConfigFile(absFile); err != nil {
		return err
	}

	return nil
}
