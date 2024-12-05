package sets

import (
	"encoding/json"
)

// UnmarshalJSON @implements json.Unmarshaler
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var elements []T
	err := json.Unmarshal(data, &elements)
	if err == nil {
		s.Clear()
		s.Add(elements...)
	}
	return err
}

// MarshalJSON @implements json.Marshaler
func (s Set[T]) MarshalJSON() ([]byte, error) {
	//var builder bytes.Buffer
	//builder.WriteByte('[')
	//count := 0
	//for elem := range s {
	//	count++
	//	b, err := json.Marshal(elem)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if len(b) > 0 {
	//		builder.Write(b)
	//		if count < len(s) {
	//			builder.WriteByte(',')
	//		}
	//	}
	//}
	//
	//builder.WriteByte(']')
	//
	//return builder.Bytes(), nil
	return json.Marshal(s.UnsortedList())
}
