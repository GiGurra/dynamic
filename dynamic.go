package dynamic

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type T[S any] struct {
	S S
	X map[string]any
}

func NewT[V any](static V, extra map[string]any) T[V] {
	return T[V]{
		S: static,
		X: func() map[string]any {
			if extra == nil {
				return make(map[string]any)
			} else {
				return extra
			}
		}(),
	}
}

// see https://pkg.go.dev/github.com/mitchellh/mapstructure#section-readme
// Json encoder will encode the static fields and the extra fields.

func (t T[S]) MarshalJSON() ([]byte, error) {

	resultMap := make(map[string]any)

	for k, v := range t.X {
		resultMap[k] = v
	}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &resultMap,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(t.S)
	if err != nil {
		return nil, err
	}

	return json.Marshal(resultMap)
}

func (t *T[S]) UnmarshalJSON(data []byte) error {

	err := json.Unmarshal(data, &t.X)
	if err != nil {
		return err
	}

	meta := mapstructure.Metadata{}
	config := &mapstructure.DecoderConfig{
		Metadata: &meta,
		Result:   &t.S,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(t.X)
	if err != nil {
		return err
	}

	// lastly remove duplicates
	allowedExtraKeys := map[string]bool{}
	for _, key := range meta.Unused {
		allowedExtraKeys[key] = true
	}
	for key := range t.X {
		if !allowedExtraKeys[key] {
			delete(t.X, key)
		}
	}

	return nil
}
