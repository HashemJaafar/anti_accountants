package main

import (
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	equity := 1000.0
	totalWithdrawals := 0.0
	const growthRate = 0.5
	const charity = 0.1

	for k1 := 0; k1 < 60; k1++ {

		income := equity * growthRate
		withdrawals := income * charity

		equity += income - withdrawals
		totalWithdrawals += withdrawals

		fmt.Println(equity, "\t", totalWithdrawals)
	}
}

func Test2(t *testing.T) {
	fmt.Println(time.Now().Sub(time.Date(2018, time.February, 4, 0, 0, 0, 0, time.UTC)).Hours() / 24)
}
