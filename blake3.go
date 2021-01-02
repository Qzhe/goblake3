package goblake3
/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -lblake3
#include <stdlib.h>
#include <unistd.h>
#include "blake3.h"

const int blake3_out_len = BLAKE3_OUT_LEN;

// return from heap
blake3_hasher* NewBlake3Hasher() {
    blake3_hasher* hasher = (blake3_hasher*)malloc(sizeof(blake3_hasher));
    blake3_hasher_init(hasher);
    return hasher;
}

void Blake3HasherUpdate(blake3_hasher* hasher, unsigned char* buf, size_t len) {
    blake3_hasher_update(hasher, buf, len);
}

// return from heap
uint8_t* Blake3HasherFinalize(blake3_hasher* hasher) {
    uint8_t* output = (uint8_t*)malloc(BLAKE3_OUT_LEN);
    blake3_hasher_finalize(hasher, output, BLAKE3_OUT_LEN);
    return output;
}
*/
import "C"
import (
	"unsafe"
)

type Blake3Hasher struct {
	hasher *C.struct_blake3_hasher
}

func New() *Blake3Hasher {
	hasher := C.NewBlake3Hasher()
	tmp := (*C.struct_blake3_hasher)((unsafe.Pointer)(hasher))
	return &Blake3Hasher{hasher: tmp}
}

func (b3h *Blake3Hasher) Update(data []byte) {
	len := C.size_t(len(data))
	buf := C.CBytes(data)
	tmp := (*C.blake3_hasher)(unsafe.Pointer(b3h.hasher))
	//                      uchar size_t => uint8 uint
	C.Blake3HasherUpdate(tmp, (*C.uchar)(buf), len)
	C.free(buf)
}

func (b3h *Blake3Hasher) Finalize() []byte {
	tmp := (*C.blake3_hasher)(unsafe.Pointer(b3h.hasher))
	output := C.Blake3HasherFinalize(tmp)
	// unsigned char => []byte
	return C.GoBytes(unsafe.Pointer(output), C.int(C.blake3_out_len))
}