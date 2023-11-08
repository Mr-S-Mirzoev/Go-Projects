// Требуется реализовать функцию uniqRandn, 
// которая генерирует слайс уникальных, рандомных чисел длины n.

import (
    "fmt""math/rand"
)

func main() {
    fmt.Println(uniqRandn(10))
}

func uniqRandn(n int) []int { // O(N)
    hashTable := make(map[int]bool, 0, n)
    var rand int
    for i := 0; i < n; i++ { // O(N)
        found := true

        for (found) { // O(1)
            rand := rand.IntRandom()
            found = hashTable[rand]
        }
        
        hashTable[rand] = true
    }
    
    // O(N)
    //...
}