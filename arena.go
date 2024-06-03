package iracluster

// #include <stdlib.h>
import "C"
import "unsafe"

type Arena []unsafe.Pointer

func (a *Arena) CString(str string) *C.char {
	ptr := C.CString(str)
	*a = append(*a, unsafe.Pointer(ptr))
	return ptr
}

func (a *Arena) Add(ptr unsafe.Pointer) {
	*a = append(*a, ptr)
}

func (a *Arena) free() {
	for _, ptr := range *a {
		C.free(ptr)
	}
}
