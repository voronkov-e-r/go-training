package mutex

import (
	"fmt"
	"sync"
)

type Account struct {
	mu      sync.Mutex
	balance int
}

func (acc *Account) Deposit(amount int) {
	if amount <= 0 {
		return
	}
	acc.mu.Lock()
	defer acc.mu.Unlock()
	acc.balance += amount
}

func (acc *Account) WithDraw(amount int) error {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	if acc.balance < amount {
		return fmt.Errorf("Недостаточно средств на балансе.")
	}
	acc.balance -= amount
	return nil
}

func (acc *Account) Balance() int {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	return acc.balance
}

func main() {
	var wg sync.WaitGroup
	account := Account{}
	errMu := sync.Mutex{}
	countOfErrs := 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				account.Deposit(1)
			}
		}()
		if i%2 == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				drawErr := account.WithDraw(10)
				if drawErr != nil {
					errMu.Lock()
					defer errMu.Unlock()
					countOfErrs += 1
				}
			}()
		}
	}
	wg.Wait()
	fmt.Println(account.Balance(), countOfErrs)
}
