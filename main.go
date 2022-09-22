package main

import (
  "crypto/rand"
  "fmt"
  "math"
  "math/big"
  "sync"
  "sync/atomic"
  "time"
)

const (
  targetSide = uint32(0) // 0 for heads, 1 for tails; this is stored as a uint32 due to limitations within the atomic package.
  targetConsecutive = 15 // number of times to consecutively reach desired side
  
  threadCount = 50 // number of goroutines to start
)

var (
  wg sync.WaitGroup
  
  last uint32
  consecutive uint64
  
  heads uint64
  tails uint64
  iterations uint64
)

func main() {
  start := time.Now()
  
  // Start threads.
  
  wg.Add(threadCount)
  for i := 0; i < threadCount; i++ {
    go func(min, max int64) {
      for atomic.LoadUint64(&consecutive) < targetConsecutive {
        atomic.AddUint64(&iterations, 1)
        
        random := randInt(min, max)
        fmt.Println(random)
    
        if random.Cmp(big.NewInt(50)) <= 0 {
          atomic.AddUint64(&heads, 1)
          atomic.StoreUint32(&last, 0)
        } else {
          atomic.AddUint64(&tails, 1)
          atomic.StoreUint32(&last, 1)
        }

        if atomic.LoadUint32(&last) == targetSide {
          atomic.AddUint64(&consecutive, 1)
          continue
        }

        atomic.StoreUint64(&consecutive, 0)
      }
      wg.Done()
    }(1, 100)
  }

  wg.Wait()

  // Convert targetSide to a human-readable string, since it's stored as a uint32 (see line 14).
  
  var side string
  if targetSide == 0 {
    side = "heads"
  } else {
    side = "tails"
  }
  
  // Print out useful (or not) stats.
  
  fmt.Println("Heads:", heads, "\nTails:", tails, "\nIterations:", iterations)
  fmt.Pritnf("\nTarget: %d (%s)\n", targetConsecutive, side)
  fmt.Printf("Predicted %%: %.10F%%\nActual %%: %.10F%%\n", 1 / math.Pow(2, targetConsecutive), 1 / float64(iterations))
  fmt.Println("Execution Time:", time.Since(start))
}

func randInt(min, max int64) *big.Int {
  random, _ := rand.Int(rand.Reader, big.NewInt(max))  
  return random.Add(random, big.NewInt(min))
}
