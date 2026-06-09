package main

import (
	"fmt"
	"slices"
	"unsafe"
)

func main() {
	fmt.Println("=== slices.Clip: уменьшает cap, но не освобождает память ===")

	big := make([]int, 10, 1000)
	for i := range big {
		big[i] = i + 1
	}
	fmt.Println("big: len:", len(big), "cap:", cap(big))
	fmt.Printf("big data ptr: %p\n", unsafe.SliceData(big))

	clipped := slices.Clip(big)
	fmt.Println("clipped: len:", len(clipped), "cap:", cap(clipped))
	fmt.Printf("clipped data ptr: %p (тот же массив!)\n", unsafe.SliceData(clipped))

	fmt.Println()
	fmt.Println("=== Реальное освобождение: copy в новый слайс ===")

	src := make([]int, 5, 10000)
	for i := range src {
		src[i] = i * 10
	}
	fmt.Printf("src ptr: %p, cap: %d\n", unsafe.SliceData(src), cap(src))

	freed := make([]int, len(src))
	copy(freed, src)
	src = nil
	fmt.Printf("freed ptr: %p (новый массив), cap: %d\n", unsafe.SliceData(freed), cap(freed))

	fmt.Println()
	fmt.Println("=== Утечка через маленький подслайс от большого массива ===")

	huge := make([]byte, 1<<20)
	for i := range huge {
		huge[i] = byte(i)
	}
	fmt.Printf("huge: len=%d, cap=%d\n", len(huge), cap(huge))

	// Опасно: huge не будет собран GC, пока живёт leak
	leak := huge[100:110]
	fmt.Printf("leak (опасный подслайс): len=%d, cap=%d\n", len(leak), cap(leak))

	// Безопасно: tiny не держит огромный backing array
	tiny := make([]byte, 10)
	copy(tiny, huge[100:110])
	huge = nil
	fmt.Printf("tiny (безопасная копия): len=%d, cap=%d\n", len(tiny), cap(tiny))
}
