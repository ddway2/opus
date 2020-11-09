package opus

import "testing"

func TestDecoderCustomNew(t *testing.T) {
	mode, err := NewOpusMode(48000, 1)

	if err != nil {
		t.Errorf("Error creating mode: %v", err)
	}

	dec, err2 := NewDecoderCustom(2, mode)
	if err2 != nil {
		t.Errorf("Error creatinh mode: %v", err)
	}
}
