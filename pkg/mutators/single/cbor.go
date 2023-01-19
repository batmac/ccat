package mutators

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/fxamacker/cbor/v2"
)

func init() {
	singleRegister("cbor2json", cborDecode, withDescription("decode cbor"), withCategory("convert"))
	singleRegister("json2cbor", cborEncode, withDescription("encode to cbor"), withCategory("convert"))
}

func cborEncode(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	var fromjson interface{}

	fj := json.NewDecoder(r)

	if err := fj.Decode(&fromjson); err != nil {
		return 0, err
	}

	cborEncoder := cbor.NewEncoder(w)

	if err := cborEncoder.Encode(fromjson); err != nil {
		return 0, err
	}

	return -1, nil
}

func cborDecode(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	var fromcbor interface{}

	opts := cbor.DecOptions{
		DefaultMapType: reflect.TypeOf(map[string]interface{}(nil)),
	}
	em, err := opts.DecMode()
	if err != nil {
		return 0, err
	}
	fc := em.NewDecoder(r)

	if err := fc.Decode(&fromcbor); err != nil {
		return 0, err
	}

	// log.Printf("fromcbor: %#v", fromcbor)
	jsonEncoder := json.NewEncoder(w)

	if err := jsonEncoder.Encode(fromcbor); err != nil {
		return 0, err
	}

	return -1, nil
}
