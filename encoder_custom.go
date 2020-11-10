package opus

/*
#cgo windows CFLAGS: -Ithird_party/opus/include/opus
#cgo windows LDFLAGS: -Lthird_party/opus/lib/windows -lopus
#include <opus.h>
#include <opus_custom.h>

int
bridge_encoder_custom_set_dtx(OpusCustomEncoder *st, opus_int32 use_dtx)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_DTX(use_dtx));
}

int
bridge_encoder_custom_get_dtx(OpusCustomEncoder *st, opus_int32 *dtx)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_DTX(dtx));
}

int
bridge_encoder_custom_get_sample_rate(OpusCustomEncoder *st, opus_int32 *sample_rate)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_SAMPLE_RATE(sample_rate));
}


int
bridge_encoder_custom_set_bitrate(OpusCustomEncoder *st, opus_int32 bitrate)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_BITRATE(bitrate));
}

int
bridge_encoder_custom_get_bitrate(OpusCustomEncoder *st, opus_int32 *bitrate)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_BITRATE(bitrate));
}

int
bridge_encoder_custom_set_complexity(OpusCustomEncoder *st, opus_int32 complexity)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_COMPLEXITY(complexity));
}

int
bridge_encoder_custom_get_complexity(OpusCustomEncoder *st, opus_int32 *complexity)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_COMPLEXITY(complexity));
}

int
bridge_encoder_custom_set_max_bandwidth(OpusCustomEncoder *st, opus_int32 max_bw)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_MAX_BANDWIDTH(max_bw));
}

int
bridge_encoder_custom_get_max_bandwidth(OpusCustomEncoder *st, opus_int32 *max_bw)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_MAX_BANDWIDTH(max_bw));
}

int
bridge_encoder_custom_set_inband_fec(OpusCustomEncoder *st, opus_int32 fec)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_INBAND_FEC(fec));
}

int
bridge_encoder_custom_get_inband_fec(OpusCustomEncoder *st, opus_int32 *fec)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_INBAND_FEC(fec));
}

int
bridge_encoder_custom_set_packet_loss_perc(OpusCustomEncoder *st, opus_int32 loss_perc)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_PACKET_LOSS_PERC(loss_perc));
}

int
bridge_encoder_custom_get_packet_loss_perc(OpusCustomEncoder *st, opus_int32 *loss_perc)
{
	return opus_custom_encoder_ctl(st, OPUS_GET_PACKET_LOSS_PERC(loss_perc));
}

int
bridge_encoder_custom_set_vbr(OpusCustomEncoder *st, opus_int32 vbr)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_VBR(vbr));
}

int
bridge_encoder_custom_set_application(OpusCustomEncoder *st, opus_int32 application)
{
	return opus_custom_encoder_ctl(st, OPUS_SET_APPLICATION(application));
}
*/
import "C"
import (
	"fmt"
)

var errEncCustomUninitialized = fmt.Errorf("opus encoder uninitialized")

type EncoderCustom struct {
	p *C.struct_OpusCustomEncoder

	channels int
	mode     *OpusMode
}

func NewEncoderCustom(channels int, mode *OpusMode) (*EncoderCustom, error) {
	var enc EncoderCustom
	err := enc.Init(channels, mode)
	if err != nil {
		return nil, err
	}
	return &enc, nil
}

func (enc *EncoderCustom) Init(channels int, mode *OpusMode) error {
	if enc.p != nil {
		return fmt.Errorf("opus encoder custom already initialized")
	}
	if channels != 1 && channels != 2 {
		return fmt.Errorf("Number of channels must be 1 or 2: %d", channels)
	}

	enc.channels = channels
	enc.mode = mode
	var errno C.int
	enc.p = C.opus_custom_encoder_create(
		enc.mode.P,
		C.int(channels),
		&errno)
	if int(errno) != 0 {
		return Error(int(errno))
	}
	return nil
}

func (enc *EncoderCustom) Close() {
	if enc.p != nil {
		C.opus_custom_encoder_destroy(enc.p)
	}
	enc.p = nil
}

// Encode raw PCM data and store the result in the supplied buffer. On success,
// returns the number of bytes used up by the encoded data.
func (enc *EncoderCustom) Encode(pcm []int16, data []byte) (int, error) {
	if enc.p == nil {
		return 0, errEncCustomUninitialized
	}
	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no target buffer")
	}
	// libopus talks about samples as 1 sample containing multiple channels. So
	// e.g. 20 samples of 2-channel data is actually 40 raw data points.
	if len(pcm)%enc.channels != 0 {
		return 0, fmt.Errorf("opus: input buffer length must be multiple of channels")
	}
	samples := len(pcm) / enc.channels
	n := int(C.opus_custom_encode(
		enc.p,
		(*C.opus_int16)(&pcm[0]),
		C.int(samples),
		(*C.uchar)(&data[0]),
		C.int(cap(data))))
	if n < 0 {
		return 0, Error(n)
	}
	return n, nil
}

// Encode raw PCM data and store the result in the supplied buffer. On success,
// returns the number of bytes used up by the encoded data.
func (enc *EncoderCustom) EncodeFloat32(pcm []float32, data []byte) (int, error) {
	if enc.p == nil {
		return 0, errEncCustomUninitialized
	}
	if len(pcm) == 0 {
		return 0, fmt.Errorf("opus: no data supplied")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("opus: no target buffer")
	}
	if len(pcm)%enc.channels != 0 {
		return 0, fmt.Errorf("opus: input buffer length must be multiple of channels")
	}
	samples := len(pcm) / enc.channels
	n := int(C.opus_custom_encode_float(
		enc.p,
		(*C.float)(&pcm[0]),
		C.int(samples),
		(*C.uchar)(&data[0]),
		C.int(cap(data))))
	if n < 0 {
		return 0, Error(n)
	}
	return n, nil
}

// SetDTX configures the encoder's use of discontinuous transmission (DTX).
func (enc *EncoderCustom) SetDTX(dtx bool) error {
	i := 0
	if dtx {
		i = 1
	}
	res := C.bridge_encoder_custom_set_dtx(enc.p, C.opus_int32(i))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// DTX reports whether this encoder is configured to use discontinuous
// transmission (DTX).
func (enc *EncoderCustom) DTX() (bool, error) {
	var dtx C.opus_int32
	res := C.bridge_encoder_custom_get_dtx(enc.p, &dtx)
	if res != C.OPUS_OK {
		return false, Error(res)
	}
	return dtx != 0, nil
}

// SampleRate returns the encoder sample rate in Hz.
func (enc *EncoderCustom) SampleRate() (int, error) {
	var sr C.opus_int32
	res := C.bridge_encoder_custom_get_sample_rate(enc.p, &sr)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return int(sr), nil
}

// SetBitrate sets the bitrate of the EncoderCustom
func (enc *EncoderCustom) SetBitrate(bitrate int) error {
	res := C.bridge_encoder_custom_set_bitrate(enc.p, C.opus_int32(bitrate))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// SetBitrateToAuto will allow the encoder to automatically set the bitrate
func (enc *EncoderCustom) SetBitrateToAuto() error {
	res := C.bridge_encoder_custom_set_bitrate(enc.p, C.opus_int32(C.OPUS_AUTO))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// SetBitrateToMax causes the encoder to use as much rate as it can. This can be
// useful for controlling the rate by adjusting the output buffer size.
func (enc *EncoderCustom) SetBitrateToMax() error {
	res := C.bridge_encoder_custom_set_bitrate(enc.p, C.opus_int32(C.OPUS_BITRATE_MAX))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// Bitrate returns the bitrate of the EncoderCustom
func (enc *EncoderCustom) Bitrate() (int, error) {
	var bitrate C.opus_int32
	res := C.bridge_encoder_custom_get_bitrate(enc.p, &bitrate)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return int(bitrate), nil
}

// SetComplexity sets the encoder's computational complexity
func (enc *EncoderCustom) SetComplexity(complexity int) error {
	res := C.bridge_encoder_custom_set_complexity(enc.p, C.opus_int32(complexity))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// Complexity returns the computational complexity used by the encoder
func (enc *EncoderCustom) Complexity() (int, error) {
	var complexity C.opus_int32
	res := C.bridge_encoder_custom_get_complexity(enc.p, &complexity)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return int(complexity), nil
}

// SetMaxBandwidth configures the maximum bandpass that the encoder will select
// automatically
func (enc *EncoderCustom) SetMaxBandwidth(maxBw Bandwidth) error {
	res := C.bridge_encoder_custom_set_max_bandwidth(enc.p, C.opus_int32(maxBw))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// MaxBandwidth gets the encoder's configured maximum allowed bandpass.
func (enc *EncoderCustom) MaxBandwidth() (Bandwidth, error) {
	var maxBw C.opus_int32
	res := C.bridge_encoder_custom_get_max_bandwidth(enc.p, &maxBw)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return Bandwidth(maxBw), nil
}

// SetInBandFEC configures the encoder's use of inband forward error
// correction (FEC)
func (enc *EncoderCustom) SetInBandFEC(fec bool) error {
	i := 0
	if fec {
		i = 1
	}
	res := C.bridge_encoder_custom_set_inband_fec(enc.p, C.opus_int32(i))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// InBandFEC gets the encoder's configured inband forward error correction (FEC)
func (enc *EncoderCustom) InBandFEC() (bool, error) {
	var fec C.opus_int32
	res := C.bridge_encoder_custom_get_inband_fec(enc.p, &fec)
	if res != C.OPUS_OK {
		return false, Error(res)
	}
	return fec != 0, nil
}

// SetPacketLossPerc configures the encoder's expected packet loss percentage.
func (enc *EncoderCustom) SetPacketLossPerc(lossPerc int) error {
	res := C.bridge_encoder_custom_set_packet_loss_perc(enc.p, C.opus_int32(lossPerc))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

// PacketLossPerc gets the encoder's configured packet loss percentage.
func (enc *EncoderCustom) PacketLossPerc() (int, error) {
	var lossPerc C.opus_int32
	res := C.bridge_encoder_custom_get_packet_loss_perc(enc.p, &lossPerc)
	if res != C.OPUS_OK {
		return 0, Error(res)
	}
	return int(lossPerc), nil
}

func (enc *EncoderCustom) SetVBR(enable bool) error {
	var vbr C.opus_int32
	if enable {
		vbr = 1
	} else {
		vbr = 0
	}
	res := C.bridge_encoder_custom_set_vbr(enc.p, C.opus_int32(vbr))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}

func (enc *EncoderCustom) SetApplication(application Application) error {
	res := C.bridge_encoder_custom_set_application(enc.p, C.opus_int32(application))
	if res != C.OPUS_OK {
		return Error(res)
	}
	return nil
}
