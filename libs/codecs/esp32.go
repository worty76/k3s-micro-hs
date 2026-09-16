package codecs

type Esp32Codec struct{}

func NewEsp32Codec() *Esp32Codec {
	return &Esp32Codec{}
}

func (c *Esp32Codec) Decode(payload []byte) (string, error) {
	// Implement the decoding logic specific to ESP32 here
	// For example, you might want to parse the payload and return a string representation
	return string(payload), nil
}
