package binstruct

// Unmarshaler provides an interface for types to unmarshal themselves to binary.
type Unmarshaler interface {
	UnmarshalBinary([]byte) error
}

func Unmarshal(data []byte, v interface{}) error {
	return nil
}
