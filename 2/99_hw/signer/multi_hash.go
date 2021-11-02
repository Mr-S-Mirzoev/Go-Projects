package main

import (
	"fmt"
	"strings"
	"sync"
)

const iMax int = 6

// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5
// ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash)
// * DataSignerCrc32, считается 1 сек
func MultiHash(in, out chan interface{}) {
	mainWg := new(sync.WaitGroup)
	for iface := range in {
		mainWg.Add(1)
		go func(data string) {
			defer mainWg.Done()
			subhashes := make([]string, 6)

			wg := new(sync.WaitGroup)
			wg.Add(iMax)

			for i := 0; i < iMax; i++ {
				go func(idx int) {
					idxStringed := fmt.Sprintf("%d", idx)
					subhashes[idx] = DataSignerCrc32(idxStringed + data)
					wg.Done()
				}(i)
			}

			wg.Wait()
			out <- strings.Join(subhashes, "")
		}(iface.(string))
	}

	mainWg.Wait()
}
