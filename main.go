package main

import (
	"Lesson19/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}
	_, err := svc.RegisterAccount("+992000001")
	if err != nil {
		log.Print(err)
	}
	_, err = svc.Deposit(1, 10000_000_000)
	if err != nil {
		log.Print()
	}
	svc.Pay(1, "auto", 20_000_000)
	svc.Pay(1, "auto", 30_000_000)
	svc.Pay(1, "auto", 40_000_000)
	svc.Pay(1, "auto", 50_000_000)
	svc.Pay(1, "auto", 60_000_000)
	svc.Pay(1, "auto", 10_000_000)
	for v := range svc.SumPaymentsWithProgress() {
		log.Print("Part: ", v.Part, " Result: ", v.Result)
	}
}
