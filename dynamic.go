package dynamic

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type T[S any] struct {
	S S              `json:",squash"`
	X map[string]any `json:",remain,squash"`
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

	var resultMap map[string]any
	err := json.Unmarshal(data, &resultMap)
	if err != nil {
		return err
	}

	err = mapstructure.Decode(resultMap, &t.S)
	if err != nil {
		return err
	}

	t.X = resultMap

	return nil
}
