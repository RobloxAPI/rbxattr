package rbxattr_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/robloxapi/rbxattr"
)

func TestModelRoundtrip(t *testing.T) {
	var data = `DgAAAAQAAABNTU1NGwAAQMEAAAhCBAAAAEtLS0sXBQAAAAAAAAAAAAAAAAAAADQz0z402Vs+ZmYWP2hm5j29JAU/AADQPgAAAAANB1g/MzMTPzQzMz4AAIA/AAAgPwQAAABKSkpKEaRwRUG4HmNCG54RQQQAAABISEhID8HAQD2JiAg+4eBgPgQAAABOTk5OHAAAQMEAAGDCAAAIQgAAnEIEAAAATExMTBkFAAAAAAAAAAAAAACRkBA+jo0NP7GwMD4AAAAA+GGqPZGQED+lpCQ+//7+PgAAAADEa7s+gYAAPuno6D3JyEg+AAAAAGLBSj8AAIA/zs1NP93cXD4AAAAAAACAP7m4uD6FhAQ/goEBPwQAAABHR0dHDvMDAAAEAAAARkZGRgqPwvU9IgAAAClcDz9OAAAABAAAAEVFRUUJj8L1PSIAAAAEAAAARERERAbvzauJZ0UjAQQAAABDQ0NDBgAAAAAAAAAABAAAAElJSUkQpHBFQbgeY0IEAAAAQkJCQgMBBAAAAEFBQUECBgAAAGZvb2Jhcg==`
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	base64.StdEncoding.DecodedLen(len(data))
	var attrs rbxattr.Model
	n, err := attrs.ReadFrom(r)
	if err != nil {
		t.Fatal(err)
	}
	if n != 409 {
		t.Fatalf("expected 409 bytes read, got %d", n)
	}

	var w bytes.Buffer
	bw := base64.NewEncoder(base64.StdEncoding, &w)
	n, err = attrs.WriteTo(bw)
	bw.Close()
	if err != nil {
		t.Fatal(err)
	}
	if n != 409 {
		t.Fatalf("expected 409 bytes written, got %d", n)
	}
	if w.String() != data {
		t.Fatalf("encoded bytes do not match decoded bytes\n\t%s\n\t%s", data, w.String())
	}
}

func ExampleModel_ReadFrom() {
	var data = `AgAAAAQAAABTaXplCgAAAD9kAAAAAAAAP2QAAAAIAAAAUG9zaXRpb24KAACAPs7///8AAIA+zv///w==`
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	var model rbxattr.Model
	n, _ := model.ReadFrom(r)
	fmt.Printf("Read %d bytes\n", n)

	var dict = make(map[string]rbxattr.Value, len(model.Value))
	for _, entry := range model.Value {
		if _, ok := dict[entry.Key]; !ok {
			dict[entry.Key] = entry.Value
		}
	}

	fmt.Println("Size.X.Scale:", dict["Size"].(*rbxattr.ValueUDim2).X.Scale)
	fmt.Println("Size.X.Offset:", dict["Size"].(*rbxattr.ValueUDim2).X.Offset)
	fmt.Println("Size.Y.Scale:", dict["Size"].(*rbxattr.ValueUDim2).Y.Scale)
	fmt.Println("Size.Y.Offset:", dict["Size"].(*rbxattr.ValueUDim2).Y.Offset)
	fmt.Println("Position.X.Scale:", dict["Position"].(*rbxattr.ValueUDim2).X.Scale)
	fmt.Println("Position.X.Offset:", dict["Position"].(*rbxattr.ValueUDim2).X.Offset)
	fmt.Println("Position.Y.Scale:", dict["Position"].(*rbxattr.ValueUDim2).Y.Scale)
	fmt.Println("Position.Y.Offset:", dict["Position"].(*rbxattr.ValueUDim2).Y.Offset)
	// Output:
	// Read 58 bytes
	// Size.X.Scale: 0.5
	// Size.X.Offset: 100
	// Size.Y.Scale: 0.5
	// Size.Y.Offset: 100
	// Position.X.Scale: 0.25
	// Position.X.Offset: -50
	// Position.Y.Scale: 0.25
	// Position.Y.Offset: -50
}

func ExampleModel_WriteTo() {
	model := rbxattr.Model{
		Value: rbxattr.ValueDictionary{
			{Key: "Size", Value: &rbxattr.ValueUDim2{
				X: rbxattr.ValueUDim{Scale: 0.5, Offset: 100},
				Y: rbxattr.ValueUDim{Scale: 0.5, Offset: 100},
			}},
			{Key: "Position", Value: &rbxattr.ValueUDim2{
				X: rbxattr.ValueUDim{Scale: 0.25, Offset: -50},
				Y: rbxattr.ValueUDim{Scale: 0.25, Offset: -50},
			}},
		},
	}

	var w bytes.Buffer
	bw := base64.NewEncoder(base64.StdEncoding, &w)
	n, _ := model.WriteTo(bw)
	fmt.Printf("Wrote %d bytes\n", n)
	bw.Close()
	fmt.Println(w.String())

	// Output:
	// Wrote 58 bytes
	// AgAAAAQAAABTaXplCgAAAD9kAAAAAAAAP2QAAAAIAAAAUG9zaXRpb24KAACAPs7///8AAIA+zv///w==
}
