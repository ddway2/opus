package opus

import (
	"fmt"
	"unsafe"
)

/*
#cgo windows CFLAGS: -Ithird_party/opus/include/opus
#cgo windows LDFLAGS: -Lthird_party/opus/lib/windows -lopus
#include <opus.h>

int
bridge_decoder_custom_get_last_packet_duration(OpusDecoder *st, opus_int32 *samples)
{
	return opus_custom_decoder_ctl(st, OPUS_GET_LAST_PACKET_DURATION(samples));
}
*/

import "C"

type DecoderCustom struct {
	p *C.struct_OpusCustomDecoder

	mem      []byte
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

	size := C.opus_custom_decoder_get_size(C.int(channels))
	dec.channels = channels
	dec.mode = mode
	dec.mem = make([]byte, size)
	dec.p = (*C.OpusCustomDecoder)(unsafe.Pointer(&dec.mem[0]))
	errno := C.opus_custom_decoder_init(
		dec.p,
		dec.mode.P,
		C.int(channels))
	if errno != 0 {
		return Error(errno)
	}
	return nil
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

// DecodePLC recovers a lost packet using Opus Packet Loss Concealment feature.
//
// The supplied buffer needs to be exactly the duration of audio that is missing.
// When a packet is considered "lost", `DecodePLC` and `DecodePLCFloat32` methods
// can be called in order to obtain something better sounding than just silence.
// The PCM needs to be exactly the duration of audio that is missing.
// `LastPacketDuration()` can be used on the decoder to get the length of the
// last packet.
//
// This option does not require any additional encoder options. Unlike FEC,
// PLC does not introduce additional latency. It is calculated from the previous
// packet, not from the next one.
func (dec *DecoderCustom) DecodePLC(pcm []int16) error {
	if dec.p == nil {
		return errDecUninitialized
	}
	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}
	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: output buffer capacity must be multiple of channels")
	}
	n := int(C.opus_custom_decode(
		dec.p,
		nil,
		0,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))
	if n < 0 {
		return Error(n)
	}
	return nil
}

// DecodePLCFloat32 recovers a lost packet using Opus Packet Loss Concealment feature.
// The supplied buffer needs to be exactly the duration of audio that is missing.
func (dec *DecoderCustom) DecodePLCFloat32(pcm []float32) error {
	if dec.p == nil {
		return errDecUninitialized
	}
	if len(pcm) == 0 {
		return fmt.Errorf("opus: target buffer empty")
	}
	if cap(pcm)%dec.channels != 0 {
		return fmt.Errorf("opus: output buffer capacity must be multiple of channels")
	}
	n := int(C.opus_custom_decode_float(
		dec.p,
		nil,
		0,
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)/dec.channels),
		0))
	if n < 0 {
		return Error(n)
	}
	return nil
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
