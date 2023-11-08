// problems? if there are, fix

/*
func main() {
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i)
        }()
    }
}
*/

func main() {
    var wg sync.WaitGroup // init
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(int i) {
            defer wg.Done()
            fmt.Println(i)
        }(i)
    }
    wg.Wait()
}

