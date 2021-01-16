[![Go Reference](https://pkg.go.dev/badge/github.com/robloxapi/rbxattr.svg)](https://pkg.go.dev/github.com/robloxapi/rbxattr)

# rbxattr
The rbxattr package implements the serialized format of Roblox's instance
attributes.

## Specification
The [spec.md](spec.md) file is a specification describing the structure of the
attributes binary format.

## Usage
Example using ReadFrom:
```go
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/robloxapi/rbxattr"
)

func main() {
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
```

Example using WriteTo:
```go
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/robloxapi/rbxattr"
)

func main() {
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
```
