package main

import (
	"fmt"
	"sync"
)

func SingleHashWorker(out chan interface{}, data string, mainWg *sync.WaitGroup, lock *sync.Mutex) {
	defer mainWg.Done()
	firstPartChannel := make(chan string)
	secondPartChannel := make(chan string)

	go func() {
		firstPartChannel <- DataSignerCrc32(data)
		close(firstPartChannel)
	}()

	go func() {
		lock.Lock()
		dataMd5 := DataSignerMd5(data)
		lock.Unlock()

		secondPartChannel <- DataSignerCrc32(dataMd5)
		close(secondPartChannel)
	}()

	out <- (<-firstPartChannel + "~" + <-secondPartChannel)
}

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~), где data - то что пришло на вход (по сути - числа из первой функции)
// * DataSignerMd5 может одновременно вызываться только 1 раз, считается 10 мс. Если одновременно запустится несколько - будет перегрев на 1 сек
// * DataSignerCrc32, считается 1 сек
func SingleHash(in, out chan interface{}) {
	var lock sync.Mutex

	mainWg := new(sync.WaitGroup)
	for iface := range in {
		mainWg.Add(1)
		go SingleHashWorker(out, fmt.Sprintf("%d", iface), mainWg, &lock)
	}

	mainWg.Wait()
}
