package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner string
	balance int
}

var errNoMoney = errors.New("잔액부족")

// NewAccount create Account
func NewAccount(owner string) *Account{
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// return balance of account
func (a Account) Balance() int {
	return a.balance
}

// Withdraw x amount on your account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	} 
	a.balance -= amount
	return nil
}

// Change of th account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account \nHas: ", a.Balance())
}