package opus

import "testing"

func TestDecoderCustomNew(t *testing.T) {
	mode, err := NewOpusMode(48000, 64)

	if err != nil || mode == nil {
		t.Errorf("Error creating mode: %v", err)
	}

	dec, err2 := NewDecoderCustom(2, mode)
	if err2 != nil {
		t.Errorf("Error creatinh mode: %v", err)
	}

	dec.Close()
}

func DecoderCustomUnitialized(t *testing.T) {
	var dec DecoderCustom

	_, err := dec.Decode(nil, nil)
	if err != errDecCustomUninitialized {
		t.Errorf("Expected \"unitialized decoder custom\" error: %v", err)
	}

	_, err = dec.DecodeFloat32(nil, nil)
	if err != errDecCustomUninitialized {
		t.Errorf("Expected \"unitialized decoder custom\" error: %v", err)
	}
}

// func TestDecoderCustom_GetLastPacketDuration(t *testing.T) {
// 	const (
// 		G4            = 391.995
// 		SAMPLE_RATE   = 48000
// 		FRAME_SIZE_MS = 60
// 		FRAME_SIZE    = SAMPLE_RATE * FRAME_SIZE_MS / 1000

// 		SYSTEM_FRAME_SIZE_SAMPLES        = 64
// 		DOUBLE_SYSTEM_FRAME_SIZE_SAMPLES = SYSTEM_FRAME_SIZE_SAMPLES * 2
// 	)

// 	pcm := make([]int16, FRAME_SIZE)
// 	mode, err := NewOpusMode(SAMPLE_RATE, SYSTEM_FRAME_SIZE_SAMPLES)
// 	enc, err := NewEncoderCustom(1, mode)
// 	if err != nil || enc == nil {
// 		t.Fatalf("Error creating new Encoder custom: %v", err)
// 	}
// 	addSine(pcm, SAMPLE_RATE, G4)

// 	enc.SetVBR(false)
// 	enc.SetApplication(AppRestrictedLowdelay)
// 	enc.SetComplexity(1)

// 	data := make([]byte, 10000)
// 	n, err := enc.Encode(pcm, data)
// 	if err != nil {
// 		t.Fatalf("Couldn't encode data: %v", err)
// 	}
// 	data = data[:n]

// 	dec, err := NewDecoderCustom(1, mode)
// 	if err != nil || dec == nil {
// 		t.Fatalf("Error creating new decoder: %v", err)
// 	}
// 	n, err = dec.Decode(data, pcm)
// 	if err != nil {
// 		t.Fatalf("Couldn't decode data: %v", err)
// 	}
// 	samples, err := dec.LastPacketDuration()
// 	if err != nil {
// 		t.Fatalf("Couldn't get last packet duration: %v", err)
// 	}
// 	if samples != n {
// 		t.Fatalf("Wrong duration length. Expected %d. Got %d", n, samples)
// 	}
// }
