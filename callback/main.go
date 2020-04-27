package main

//#include <stdlib.h>
//typedef void (*callback)(void* user_data, int i);
//void cgoFunc(char* filename, callback start, void* user_data);
//extern void startCgo(void*, int);
import "C"
import (
	"fmt"
	"strconv"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

type GoStartCallback func(int)

func GoTraverse(filename string, v GoStartCallback) {

	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	p := gopointer.Save(v)
	defer gopointer.Unref(p)

	C.cgoFunc(cfilename, C.callback(C.startCgo), p)
}

//export goFn
func goFn(user_data unsafe.Pointer, i C.int) {
	time.Sleep(1000 * time.Millisecond)
	v := gopointer.Restore(user_data).(GoStartCallback)
	v(int(i))
}

func main() {
	for i := 0; i < 50; i++ {
		time.Sleep(20 * time.Millisecond)
		go func(in int) {
			GoTraverse("no."+strconv.FormatInt(int64(i), 10), func(par int) {
				fmt.Println(par + in)
			})
		}(i)
	}
	time.Sleep(10 * time.Second)
}
