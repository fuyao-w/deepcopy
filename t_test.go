package deepcopy

import (
	"encoding/json"
	"os"
	"testing"
	"time"
	"unsafe"
)

type S struct {
	A int
	B int64
	C string
	inter
}

func TestS(t *testing.T) {
	t.Log(Copy(S{
		A: 1,
		B: 2,
		C: "3",
		inter: inter{
			loc: "beijing",
		},
	}))
}

type inter struct {
	loc string
}
type K struct {
	X *int
	Y *string
	Z *struct {
		H []*int
	}
	n struct {
		int
	}
	string
	F func()
}

func PtrInt(i int) *int {
	return &i
}
func PtrStr(i string) *string {
	return &i
}

func TestName(t *testing.T) {
	k := K{
		X:      PtrInt(1),
		Y:      PtrStr("s"),
		Z:      &struct{ H []*int }{H: []*int{PtrInt(1), PtrInt(2)}},
		n:      struct{ int }{int: 1},
		string: "Name",
		F: func() {
			t.Log("func F")
		},
	}
	s := S{
		A:     1,
		B:     2,
		C:     "3",
		inter: inter{loc: "fsdfsdf"},
	}

	var list = []interface{}{
		s, []int{1, 2, 3}, []int8{1, 2, 3}, []string{"1", "2"},
		map[string]string{
			"1": "a",
		},
		[3]int{1, 2, 3},
		make(chan int, 3),
		make(chan *K, 1),
		[]interface{}{
			1, k, PtrStr("dd"), &k, func(c <-chan int) {
				t.Log("recv")
			},
		},
		unsafe.Pointer(&k),
		(*unsafe.Pointer)(unsafe.Pointer(&k)),
		complex(1, 2),
		real(complex(1, 3)),
		k,
	}

	assertF := func(i interface{}) {
		if sli, ok := Copy(i).([]interface{}); ok {
			if f, ok := sli[4].(func(c <-chan int)); ok {
				f(make(chan int))
			}
		}
	}
	assertK := func(i interface{}) {
		if newK, ok := Copy(i).(K); ok {
			k.F = nil
			newK.F()
			t.Log("assertK", newK.F == nil, k.F == nil)

		}
	}
	assertUnsafe := func(i interface{}) {
		if h, ok := i.(unsafe.Pointer); ok {
			t.Log("unsafe", h)
		}
		if h, ok := i.(*unsafe.Pointer); ok {
			t.Log("*unsafe", h)
		}
	}

	for _, i2 := range list {
		assertF(i2)
		assertUnsafe(i2)
		assertK(i2)
		t.Log(Copy(i2))
		t.Log(str(Copy(i2)))
	}

	t.Log(Copy(str).(func(interface{}) string)(k))
	chh := Copy(make(ch, 3)).(ch)
	chh <- &k
	t.Log("recv", <-chh)
	Copy(nil)
	tm := Copy(time.Now()).(time.Time)
	t.Log(tm.Unix())

}

type Inner struct {
	a int
}

func TestPointer(t *testing.T) {
	t.Log(Copy(GetErr()), Copy(GetErr()) == nil)
	innter := Inner{a: 1}
	t.Log(Copy(innter))

}

func GetErr() error {
	var res *os.PathError
	return res
}

type ch chan *K

func str(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}

type dp struct {
	a int
}

func (d dp) DeepCopy() interface{} {
	return dp{d.a}
}
func TestDpInter(t *testing.T) {
	d := dp{4}
	t.Log(Copy(d))
}
