// Ограничитель скорости
package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled error = errors.New("canceled")

// начало решения

// throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
// Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	canceled := make(chan struct{}, 1)
	handle = func() error {

		delay := time.Duration(100/limit) * time.Millisecond
		time.Sleep(delay)
		fn()

		select {

		case <-canceled:
			return nil
		default:
			return nil

		}

	}

	cancel = func() {
		canceled <- struct{}{}
	}

	return handle, cancel
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := throttle(5, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for i := 0; i < n; i++ {
		handle()

	}
	cancel()
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))

}
