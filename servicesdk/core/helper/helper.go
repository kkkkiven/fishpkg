package helper

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"unsafe"
)

// HostAddrCheck check the host address is valid
func HostAddrCheck(addr string) bool {
	items := strings.Split(addr, ":")
	if items == nil || len(items) != 2 {
		return false
	}

	a := net.ParseIP(items[0])
	if a == nil {
		return false
	}

	if m, err := regexp.MatchString("^[0-9]*$", items[1]); err != nil || m == false {
		return false
	}

	p, err := AsciiToInt(StringToBytes(items[1]))
	if err != nil {
		return false
	}
	if p < 0 || p > 65535 {
		return false
	}

	return true
}

// AsciiToInt converts bytes to int.
func AsciiToInt(bts []byte) (ret int, err error) {
	// ASCII numbers all start with the high-order bits 0011.
	// If you see that, and the next bits are 0-9 (0000 - 1001) you can grab those
	// bits and interpret them directly as an integer.
	var n int
	if n = len(bts); n < 1 {
		return 0, fmt.Errorf("converting empty bytes to int")
	}
	for i := 0; i < n; i++ {
		if bts[i]&0xf0 != 0x30 {
			return 0, fmt.Errorf("%s is not a numeric character", string(bts[i]))
		}
		ret += int(bts[i]&0xf) * pow(10, n-i-1)
	}
	return ret, nil
}

// pow for integers implementation.
// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

type WaitGroupWrapper struct {
	wg sync.WaitGroup
}

func (w *WaitGroupWrapper) AddAndRun(cb func()) {
	w.wg.Add(1)
	go func() {
		cb()
		w.wg.Done()
	}()
}

func (w *WaitGroupWrapper) Wait() {
	w.wg.Wait()
}

// GetLocalAddress 获取外部接口ip
func GetLocalAddress(prober string) string {
	if prober == "" {
		prober = "8.8.8.8:80"
	}

	conn, err := net.Dial("udp", prober)
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[:idx]
}
