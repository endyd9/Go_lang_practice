package main

import (
	"fmt"

	"github.com/endyd9/learngo/accounts/accounts"
)

func main() {
	account := accounts.NewAccount("dooyong")
	account.Deposit(100)
	fmt.Println(account)
}