package binstruct

type Marshaler interface {
	MarshalBinary() ([]byte, error)
}

func Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}
