package p

import "fmt"

type Error struct {
}

func f() string {
	return ""
}

func g() *Error {
	return nil
}

func h() (int32, *Error) {
	return 1, nil
}

func test() {
	f()
	g()
	h()
	x, _ := h()
	fmt.Printf("%d", x)
	_, err := h()
	fmt.Printf("%s", err)
	go f()
	go g()
	defer h()
}
