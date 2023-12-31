# dynamic

Dual view of types, static and dynamic, in go

Tired of:

* having to juggle between nested `map[string]any` and your own incomplete types?
* keeping your static types up to date with source implementation in an API you are consuming?

You probably don't have to.

Meet `dynamic.T[...]`.

## Usage

replace:

```
var result MyTyp
err := json.Unmarshall(bytes, &result)
```

with:

```
var result dynamic.T[MyTyp]
err := json.Unmarshall(bytes, &result)
```

And you're done!

Now, you have access to:
* all static fields using `.S`
* all remaining (yet?) unmapped fields by `.X`.

Naturally, it serializes back to its original form.

Of course, you can also use `dynamic.T[..]` types nested inside other `dynamic.T[..]` types.

Don't put the same field in both static and dynamic. It makes for a bad day when it comes to consistency.

## Why?

In my case I want to write GCP configuration backup tools that download all the inventory metadata and configuration, and push them to git as my backup place :). 
IAC is nice, but will probably rarely be 100%, so I use this as a complement.

In short, I use this to always get all the data when backing up:
* dns config
* load balancer config
* iap config
* job config
* service config
* etc

## Examples

Copy-pasta from test code

```go
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
	OperationID string `json:"operationId"`
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

```
