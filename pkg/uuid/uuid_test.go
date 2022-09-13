package uuid

import "testing"

func TestUUID_HexString(t *testing.T) {
	uuid, err := NewUUID("0c49d01f-f51a-c1ae-f54b-731e053fad4f")
	if err != nil {
		t.Error(err)
	}

	hex := uuid.HexString()
	expect := "4fad3f051e734bf5aec11af51fd0490c"
	if hex != expect {
		t.Fatalf("Expect %s, got %s", expect, hex)
	}
}
