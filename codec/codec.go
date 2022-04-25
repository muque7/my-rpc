package codec

import (
	"bytes"
	"encoding/json"
)

// Data presents the data transported between server and client.
type Data struct {
	ServiceName string        // service name
	MethodName  string        // method name
	Args        []interface{} // request's or response's body except error
	Err         string        // remote server error
}

func Encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	// decoder.UseNumber()
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	// log.Println(data)
	return data, nil
}
