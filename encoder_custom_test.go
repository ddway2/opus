package opus

import "testing"

func TestEncodeNew(t *testing.T) {
	mode, err := NewOpusMode(48000, 128)
	if err != nil || mode == nil {
		t.Errorf("Error creating new opus Mode: %v", err)
	}

	enc, err := NewEncoderCustom(1, mode)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
}
