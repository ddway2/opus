package opus

import (
	"fmt"
)

/*
#cgo windows CFLAGS: -Ithird_party/opus/include/opus
#cgo windows LDFLAGS: -Lthird_party/opus/lib/windows -lopus
#include <opus.h>
#include <opus_custom.h>

int
bridge_decoder_custom_get_last_packet_duration(OpusCustomDecoder *st, opus_int32 *samples)
{
	return opus_custom_decoder_ctl(st, OPUS_GET_LAST_PACKET_DURATION(samples));
}
*/
import "C"

var errDecCustomUninitialized = fmt.Errorf("opus decoder uninitialized")

type DecoderCustom struct {
	p *C.struct_OpusCustomDecoder

	channels int
	mode     *OpusMode
}

func NewDecoderCustom(channels int, mode *OpusMode) (*DecoderCustom, error) {
	var dec DecoderCustom
	err := dec.Init(channels, mode)
	if err != nil {
		return nil, err
	}
	return &dec, nil
}

func (dec *DecoderCustom) Init(channels int, mode *OpusMode) error {
	if dec.p != nil {
		return fmt.Errorf("opus decoder already initialized")
	}
	if channels != 1 && channels != 2 {
		return fmt.Errorf("Number of channels must be 1 or 2: %d", channels)
	}

	dec.channels = channels
	dec.mode = mode
	var errno C.int

	dec.p = C.opus_custom_decoder_create(
		dec.mode.P,
		C.int(channels),
		&errno)
	if int(errno) != 0 {
		return Error(errno)
	}
	return nil
}

func (dec *DecoderCustom) Close() {
	if dec.p != nil {
		C.opus_custom_decoder_destroy(dec.p)
	}
	dec.p = nil
}

// Decode encoded Opus data into the supplied buffer. On success, returns the
// number of samples correctly written to the target buffer.
func (dec *DecoderCustom) Decode(data []byte, pcm []int16) (int, error) {
	if dec.p == nil {
		return 0, errDecUninitialized
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}
	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: target buffer empty")
	}
	if cap(pcm)%dec.channels != 0 {
		return 0, fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}
	n := int(C.opus_custom_decode(
		dec.p,
		(*C.uchar)(&data[0]),
		C.int(len(data)),
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)/dec.channels)))
	if n < 0 {
		return 0, Error(n)
	}
	return n, nil
}

// Decode encoded Opus data into the supplied buffer. On success, returns the
// number of samples correctly written to the target buffer.
func (dec *DecoderCustom) DecodeFloat32(data []byte, pcm []float32) (int, error) {
	if dec.p == nil {
		return 0, errDecUninitialized
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}
	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: target buffer empty")
	}
	if cap(pcm)%dec.channels != 0 {
		return 0, fmt.Errorf("opus: target buffer capacity must be multiple of channels")
	}
	n := int(C.opus_custom_decode_float(
		dec.p,
		(*C.uchar)(&data[0]),
		C.int(len(data)),
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)/dec.channels)))
	if n < 0 {
		return 0, Error(n)
	}
	return n, nil
}

// LastPacketDuration gets the duration (in samples)
// of the last packet successfully decoded or concealed.
func (dec *DecoderCustom) LastPacketDuration() (int, error) {
	var samples C.opus_int32
	res := C.bridge_decoder_custom_get_last_packet_duration(dec.p, &samples)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return int(samples), nil
}
