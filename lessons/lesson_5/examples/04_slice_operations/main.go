package main

import "fmt"

func main() {

	// ограничения слайсинга
	// 0 ≤ low ≤ high ≤ cap(s) = panic
	// low > high = panic / or complitaion error

	data := make([]int, 6, 10)
	for i := range data {
		data[i] = i + 1
	}
	fmt.Println("data:", data, "len:", len(data), "cap:", cap(data))

	// s[low:high]: len=high-low, cap=cap(data)-low, 10-1=9

	s1 := data[1:4]
	fmt.Println("data[1:4]:", s1, "len:", len(s1), "cap:", cap(s1))

	// s[:n] — первые n элементов, cap остаётся исходным
	s2 := data[:3]
	fmt.Println("data[:3]:", s2, "len:", len(s2), "cap:", cap(s2))

	// s[n:] — с n-го до конца, cap = cap(data) - low = 10-4=6
	s3 := data[4:]
	fmt.Println("data[4:]:", s3, "len:", len(s3), "cap:", cap(s3))

	// Изменение через подслайс видно в оригинале
	s1[0] = 99
	// backing 1 99 3 4 5 6
	// s1 = 1 [99 3 4] 5 6
	fmt.Println("data после s1[0]=99:", data)

	// Расширение через len <= cap — без реаллокации
	data2 := make([]int, 3, 6)
	data2[0], data2[1], data2[2] = 10, 20, 30
	extended := data2[:6]
	fmt.Println("extended:", extended, "len:", len(extended), "cap:", cap(extended))

	// Three-index slicing: s[low:high:max], cap = max-low
	base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	limited := base[2:5:5]
	fmt.Println("base[2:5:5]:", limited, "len:", len(limited), "cap:", cap(limited))

	// Сравним с обычным срезанием
	unlimited := base[2:5]
	fmt.Println("base[2:5]:", unlimited, "len:", len(unlimited), "cap:", cap(unlimited))

	// Three-index защищает от случайной записи в backing array
	appendToSub(limited, base)
	//appendToSub(unlimited, base)

	newSlice := base[0:6]
	fmt.Println(newSlice)
}

func appendToSub(sub, original []int) {
	before := original[5]
	sub = append(sub, 777)
	after := original[5]
	fmt.Printf("sub cap=%d → append изменил original[5] ?: %v (%d → %d)\n",
		cap(sub), before != after, before, after)

}
