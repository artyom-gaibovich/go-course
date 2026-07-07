package main

import "fmt"

type ringBuffer struct {
	buf    []int
	sendx  int
	recvx  int
	qcount int
}

func newRing(size int) *ringBuffer {
	return &ringBuffer{buf: make([]int, size)}
}

func (r *ringBuffer) push(v int) bool {
	if r.qcount == len(r.buf) {
		return false
	}
	r.buf[r.sendx] = v
	r.sendx = (r.sendx + 1) % len(r.buf)
	r.qcount++
	return true
}

func (r *ringBuffer) pop() (int, bool) {
	if r.qcount == 0 {
		return 0, false
	}
	v := r.buf[r.recvx]
	r.recvx = (r.recvx + 1) % len(r.buf)
	r.qcount--
	return v, true
}

func main() {
	r := newRing(3)
	r.push(1)
	r.push(2)
	r.push(3)
	fmt.Println("push 4 into full ring:", r.push(4))

	v, _ := r.pop()
	fmt.Println("popped:", v)
	fmt.Println("now there is room, push 4:", r.push(4))
	fmt.Printf("sendx=%d recvx=%d qcount=%d (wrapped around)\n", r.sendx, r.recvx, r.qcount)
}
