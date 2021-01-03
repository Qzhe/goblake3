// +build darwin

package goblake3
/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -lblake3
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
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
uint8_t* Blake3HasherFinalize(blake3_hasher* hasher, size_t outlen) {
    uint8_t* output = (uint8_t*)malloc(outlen);
    memset(output, 0, outlen);
    blake3_hasher_finalize(hasher, output, BLAKE3_OUT_LEN);
    return output;
}

uint8_t* Blake3HasherFinalizeSeek(blake3_hasher* hasher, uint64_t seek, size_t outlen) {
    uint8_t* output = (uint8_t*)malloc(outlen);
    memset(output, 0, outlen);
    blake3_hasher_finalize_seek(hasher, seek, output, outlen);
    return output;
}
*/
import "C"
import (
	"unsafe"
)

type Blake3Hasher struct {
	hasher *C.struct_blake3_hasher
	seek   uint64
	outlen int32
}

func New() *Blake3Hasher {
	hasher := C.NewBlake3Hasher()
	tmp := (*C.struct_blake3_hasher)((unsafe.Pointer)(hasher))
	outlen := int32(C.blake3_out_len)
	return &Blake3Hasher{hasher: tmp, seek: 0, outlen: outlen}
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
	output := C.Blake3HasherFinalize(tmp, C.size_t(b3h.outlen))
	// unsigned char => []byte
	return C.GoBytes(unsafe.Pointer(output), C.int(C.blake3_out_len))
}

func (b3h *Blake3Hasher) FinalizeSeek() []byte {
	hasher := (*C.blake3_hasher)(unsafe.Pointer(b3h.hasher))
	output := C.Blake3HasherFinalizeSeek(hasher, C.ulonglong(b3h.seek), C.size_t(b3h.outlen))
	outPtr := unsafe.Pointer(output)
	if outPtr == nil {
		return nil
	}
	b3h.seek += uint64(b3h.outlen)
	return C.GoBytes(outPtr, C.int(b3h.outlen))
}
