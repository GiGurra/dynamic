package dynamic

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Root struct {
	Kind         string           `json:"kind"`
	Header       T[Header]        `json:"header"`
	ManagedZones []T[ManagedZone] `json:"managedZones"`
}

type Header struct {
	OperationID string  `json:"operationId"`
	Optional    *string `json:"optional,omitempty"`
}

type ManagedZone struct {
	Kind string `json:"kind"`
}

func Test_header2Json(t *testing.T) {
	hdr := NewT(
		Header{OperationID: "operation-yayayaya"},
		map[string]any{"hello": "world"},
	)

	expJson := `{"hello":"world","operationId":"operation-yayayaya"}`

	jsBytes, err := json.Marshal(hdr)
	if err != nil {
		t.Fatalf("error marshalling json: %v", err)
	}

	fmt.Printf("jsBytes: %s\n", jsBytes)
	if string(jsBytes) != expJson {
		t.Fatalf("expected json to be %s but got %s", expJson, jsBytes)
	}
}

func Test_deDuplication(t *testing.T) {
	hdr := NewT(
		Header{
			OperationID: "operation-yayayaya",
		},
		map[string]any{
			"hello":    "world",
			"optional": "hello!",
		},
	)

	jsBytes, err := json.Marshal(hdr)
	if err != nil {
		t.Fatalf("error marshalling json: %v", err)
	}

	fmt.Printf("jsBytes: %s\n", jsBytes)

	var hdrBack T[Header]
	err = json.Unmarshal(jsBytes, &hdrBack)
	if err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	if hdrBack.S.Optional == nil {
		t.Fatalf("expected Optional to be set")
	}

	if *hdrBack.S.Optional != "hello!" {
		t.Fatalf("expected Optional to be hello! but got %s", *hdrBack.S.Optional)
	}

	expHdr2 := NewT(
		Header{
			OperationID: "operation-yayayaya",
			Optional:    ptr("hello!"),
		},
		map[string]any{
			"hello": "world",
		},
	)

	if diff := cmp.Diff(hdrBack, expHdr2); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}

func ptr[T any](t T) *T {
	return &t
}

func TestT_toFromJson(t *testing.T) {

	jsString := `{
		"kind": "dns#managedZonesListResponse",
		"header": {
			"operationId": "operation-yayayaya",
			"extra-shit": 123
		},
		"managedZones": [
			{
				"kind": "dns#managedZone",
				"name": "my-zone",
				"dnsName": "my-zone.com.",	
				"description": "my-zone",
				"yolo": {
					"haha": 1,
					"dada": 2
				},
				"dynamo": [
					[
						{
							"nesty": "world"
						}
					]
				]
			}
		]
	}`

	var mapFromOrigJson map[string]any
	err := json.Unmarshal([]byte(jsString), &mapFromOrigJson)
	if err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	var root Root
	err = json.Unmarshal([]byte(jsString), &root)
	if err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	fmt.Printf("root: %+v\n", root)

	jsBytes, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		t.Fatalf("error marshalling json: %v", err)
	}

	fmt.Printf("jsBytes: %s\n", jsBytes)

	var mapFromNewJson map[string]any
	err = json.Unmarshal(jsBytes, &mapFromNewJson)
	if err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	if diff := cmp.Diff(mapFromOrigJson, mapFromNewJson); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}
