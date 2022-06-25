package wallet

import (
	"Lesson19/pkg/types"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAccountNotFound = errors.New("account not found")
var ErrAmountMostBePositive = errors.New("amount must be greater than zero")
var ErrNotEnouthBalance = errors.New("not enough balance in account")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, acc := range s.accounts {
		if acc.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountId++
	account := &types.Account{
		ID:      s.nextAccountId,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) FindAccountById(accountId int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountId {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountId int64, amount types.Money) (*types.Account, error) {
	if amount <= 0 {
		return nil, ErrAmountMostBePositive
	}
	account, err := s.FindAccountById(accountId)
	if err != nil {
		return nil, err
	}
	account.Balance += amount
	return account, nil
}

func (s *Service) Pay(acountId int64, category types.PaymentCategory, amount types.Money) (*types.Payment, error) {
	if amount < 0 {
		return nil, ErrAmountMostBePositive
	}
	account, err := s.FindAccountById(acountId)
	if err != nil {
		return nil, err
	}
	if account.Balance < amount {
		return nil, ErrNotEnouthBalance
	}
	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: acountId,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindPaymentById(paymentId string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentId {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}
func (s *Service) SumPaymentsWithProgress() <-chan types.Progress {
	ch := make(chan types.Money)
	defer close(ch)
	parts := 3
	size := len(s.payments) / parts
	for i := 0; i < parts; i++ {
		go func(ch chan<- types.Money, paymetns []*types.Payment) {
			sum := types.Money(0)
			for _, payment := range paymetns {
				sum += payment.Amount
			}
			ch <- sum
		}(ch, s.payments[i*size:(i+1)*size])
	}
	result := types.Money(0)
	for i := 0; i < parts; i++ {
		result += <-ch
	}
	progres := types.Progress{
		Part:   parts,
		Result: result,
	}
	progress := make(chan types.Progress, 1)
	progress <- progres
	go func() {
		close(progress)
	}()
	return progress
}

func merge(channels []<-chan types.Money) <-chan types.Money {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))
	merged := make(chan types.Money)
	for _, ch := range channels {
		go func(ch <-chan types.Money) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}
	go func() {
		defer close(merged)
		wg.Wait()
	}()
	return merged
}
