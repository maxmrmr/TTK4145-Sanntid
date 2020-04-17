package hardware

import "C"

func io_init() bool {
	return int(C.io_init())!=0
}

func io_setBit(channel int) { 
	C.io_set_bit(C.int(channel))
}

func io_clearBit(channel int) {
	C.io_clear_bit(C.int(channel), C.int(value))
}

func io_writeAnalog(channel int, value int) {
	C.io_writeAnalog(C.int(channel), C.int(value))
}

func io_readBit(channel int) bool {
	return int(C.io_read_bit(C.int(channel))) !=0
}

func io_readAnalog(channel int) int {
	return int(C.io_readAnalog(C.int(channel)))
}