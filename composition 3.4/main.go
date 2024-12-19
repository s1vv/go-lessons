package main1

import (
	"fmt"
	"math/rand"
	"sync"
)

// генерирует случайные слова из 5 букв
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	randomWord := func(n int) string {
		const letters = "aeiourtnsl"
		chars := make([]byte, n)
		for i := range chars {
			chars[i] = letters[rand.Intn(len(letters))]
		}
		return string(chars)
	}

	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				// fmt.Println("[generate] Cancel received, exiting.")
				return
			case out <- randomWord(5):
				// fmt.Println("[generate] Generated word sent.")
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)

	filterWord := func(w string) string {
		runeU := make(map[rune]struct{}, 5)
		for _, r := range w {
			runeU[r] = struct{}{}
		}
		if len(w) != len(runeU) {
			return ""
		}
		return w
	}

	go func() {
		defer close(out)
		for w := range in {
			select {
			case out <- filterWord(w):
				// fmt.Println("[takeUnique] Sent unique word.")
			case <-cancel:
				// fmt.Println("[takeUnique] Cancel received, exiting.")
				return
			}
		}
	}()
	return out
}

// переворачивает слова
func reverse(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)

	rev := func(w string) (result string) {
		for _, v := range w {
			result = string(v) + result
		}
		return
	}

	go func() {
		defer close(out)
		for w := range in {
			select {
			case out <- rev(w):
				// fmt.Println("[reverse] Reversed word sent.")
			case <-cancel:
				// fmt.Println("[reverse] Cancel received, exiting.")
				return
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	mergeChannel := func(ch <-chan string) {
		defer wg.Done()
		for val := range ch {
			select {
			case out <- val:
				// fmt.Println("[merge] Value sent to merged channel.")
			case <-cancel:
				// fmt.Println("[merge] Cancel received, exiting goroutine.")
				return
			}
		}
	}

	wg.Add(2)
	go mergeChannel(c1)
	go mergeChannel(c2)

	go func() {
		wg.Wait()
		close(out)
		// fmt.Println("[merge] Merged channel closed.")
	}()

	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan string, n int) {
	rev := func(w string) (result string) {
		for _, v := range w {
			result = string(v) + result
		}
		return
	}
	for range n {
		select {
		case rw := <-in:
			if rw != "" {
				fmt.Println(rev(rw) + " -> " + rw)
			}
		case <-cancel:
			// fmt.Println("[print] Cancel received, stopping.")
			return
		}
	}
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 51)
}
