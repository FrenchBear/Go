// g47_priority_queue.go
// Learning go, Implementation of a PriorityQueue
//
// 2025-09-10	PV		First version

package main

import (
	"container/heap"
	"fmt"
)

func main() {
    // Create a priority queue, put some items in it.
    items := map[string]int{
        "banana": 3, "apple": 2, "pear": 4,
    }
    pq := make(PriorityQueue[string], len(items))
    i := 0
    for value, priority := range items {
        pq[i] = &Item[string]{
            value:    value,
            priority: priority,
            index:    i,
        }
        i++
    }
    heap.Init(&pq)

    // Insert a new item and then modify its priority.
    item := &Item[string]{
        value:    "orange",
        priority: 1,
    }
    heap.Push(&pq, item)
    pq.update(item, item.value, 5) // Change orange's priority to 5

    // Take the items out; they arrive in decreasing priority order.
    for pq.Len() > 0 {
        item := heap.Pop(&pq).(*Item[string])
        fmt.Printf("%.2d:%s ", item.priority, item.value)
    }
	fmt.Println()
    // Output: 05:orange 04:pear 03:banana 02:apple
}