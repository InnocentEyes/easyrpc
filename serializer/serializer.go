package serializer

type Serializer interface {
	Marshal(message interface{}) ([]byte, error)
	UnMarshal(data []byte, message interface{}) error
}
