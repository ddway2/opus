package opus

/*
#cgo windows CFLAGS: -Ithird_party/opus/include/opus
#cgo windows LDFLAGS: -Lthird_party/opus/lib/windows -lopus
#include <opus.h>

*/

import "C"

type OpusMode struct {
	P *C.struct_OpusCustomMode
}

func NewOpusMode(sample_rate int32, frame_size int) (*OpusMode, error) {
	var mode OpusMode
	err := mode.Init(sample_rate, frame_size)
	if err != nil {
		return nil, err
	}
	return &mode, nil
}

func (mode *OpusMode) Init(sample_rate int32, frame_size int) error {
	var err C.int
	mode.p = C.opus_custom_mode_create(C.opus_int32(sample_rate), C.int(frame_size), &err)
	if int(err) != 0 {
		return Error(int(err))
	}

	return nil
}

func (mode *OpusMode) Close() {
	if mode.p != nil {
		C.opus_custom_mode_destroy(mode.p)
	}
	mode.p = nil
}
