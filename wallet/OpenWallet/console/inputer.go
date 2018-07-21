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

package console

import (
	"fmt"
	"github.com/blocktree/OpenWallet/common"
	"github.com/blocktree/OpenWallet/logger"
)

//PasswordPrompt 提示输入密码
//@param 是否二次确认
func InputPassword(isConfirm bool, minLen int) (string, error) {

	var (
		confirm  string
		password string
		err      error
	)

	for {

		// 等待用户输入密码
		password, err = Stdin.PromptPassword("Enter wallet password: ")
		if err != nil {
			openwLogger.Log.Errorf("unexpected error: %v\n", err)
			return "", err
		}

		if len(password) < minLen {
			fmt.Printf("The length of the password is less than %d chars. Please re-enter it.\n", minLen)
			continue
		}

		// 二次确认密码
		if isConfirm {

			confirm, err = Stdin.PromptPassword("Confirm wallet password: ")

			if password != confirm {
				fmt.Printf("The two password is not equal, please rre-enter it.\n")
				continue
			}

		}

		break
	}

	return password, nil
}

//InputText 输入文本
func InputText(prompt string, required bool) (string, error) {

	var (
		text string
		err  error
	)

	for {

		// 等待用户输入
		text, err = Stdin.PromptInput(prompt)
		if err != nil {
			openwLogger.Log.Errorf("unexpected error: %v\n", err)
			return "", err
		}

		if len(text) == 0 && required {
			fmt.Printf("Input can not be empty!\n")
			continue
		}

		break
	}

	return text, nil
}

//InputNumber 输入数值
func InputNumber(prompt string, zero bool) (uint64, error) {

	var (
		num uint64
	)

	for {
		// 等待用户输入参数
		line, err := Stdin.PromptInput(prompt)
		if err != nil {
			openwLogger.Log.Errorf("unexpected error: %v\n", err)
			return 0, err
		}
		num = common.NewString(line).UInt64()

		if !zero && num <= 0 {
			fmt.Printf("Input can not be empty, and must be greater than 0.\n")
			continue
		}

		break
	}

	return num, nil
}

//InputRealNumber 输入实数值
//@param p 是否正数
func InputRealNumber(prompt string, p bool) (string, error) {

	var (
		num     float64
		realNum string
	)

	for {
		// 等待用户输入参数
		line, err := Stdin.PromptInput(prompt)
		if err != nil {
			openwLogger.Log.Errorf("unexpected error: %v\n", err)
			return "0", err
		}
		num = common.NewString(line).Float64()

		if p && num <= 0 {
			fmt.Printf("Input can not be empty, and must be greater than 0.\n")
			continue
		}

		realNum = line

		break
	}

	return realNum, nil
}
