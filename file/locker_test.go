package file

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestLock(t *testing.T) {
	file := "/test-app.lock"
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lerr := Lock(f, true)
			_, werr := f.WriteString("001")
			fmt.Printf("Lock ok?: %v, Write ok?: %v\n", lerr, werr)
		}()
	}
	wg.Wait()
	if err = Unlock(f); err != nil {
		t.Fatal(err)
	}
}

func TestRLock(t *testing.T) {
	file := "/test-app.lock"
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lerr := RLock(f, true)
			_, werr := f.WriteString("001")
			fmt.Printf("Rlock ok?: %v, Write ok?: %v\n", lerr, werr)

		}()
	}
	wg.Wait()
	if err = Unlock(f); err != nil {
		t.Fatal(err)
	}
}
