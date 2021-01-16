package rbxattr

import (
	"fmt"
	"io"
)

// Type identifies an attribute type within an encoding.
type Type byte

// Many types are not implemented because they are not officially supported by
// Roblox, and could change in the future. However, Roblox could instead expose
// more of these types, so they are documented here.

const (
	_                  Type = 0x00 // Null
	_                  Type = 0x01 // Empty
	TypeString         Type = 0x02
	TypeBool           Type = 0x03
	_                  Type = 0x04 // Int
	TypeFloat          Type = 0x05
	TypeDouble         Type = 0x06
	_                  Type = 0x07 // Array
	_                  Type = 0x08 // Dictionary
	TypeUDim           Type = 0x09
	TypeUDim2          Type = 0x0A
	_                  Type = 0x0B // Ray
	_                  Type = 0x0C // Faces
	_                  Type = 0x0D // Axes
	TypeBrickColor     Type = 0x0E
	TypeColor3         Type = 0x0F
	TypeVector2        Type = 0x10
	TypeVector3        Type = 0x11
	_                  Type = 0x12 // Vector2int16
	_                  Type = 0x13 // Vector3int16
	_                  Type = 0x14 // CFrame
	_                  Type = 0x15 // EnumItem
	_                  Type = 0x16 // Unknown
	TypeNumberSequence Type = 0x17
	_                  Type = 0x18 // NumberSequenceKeypoint
	TypeColorSequence  Type = 0x19
	_                  Type = 0x1A // ColorSequenceKeypoint
	TypeNumberRange    Type = 0x1B
	TypeRect           Type = 0x1C
	_                  Type = 0x1D // PhysicalProperties
	_                  Type = 0x1E // Unknown
	_                  Type = 0x1F // Region3
	_                  Type = 0x20 // Region3int16
)

// Value is an attribute value that can be decoded from and encoded to bytes,
// with an identifying type.
type Value interface {
	Type() Type
	ReadFrom(r io.Reader) (n int64, err error)
	WriteTo(w io.Writer) (n int64, err error)
}

// NewValue returns a new Value of the given Type, or nil if the Type does not
// correspond to a known Value.
func NewValue(typ Type) Value {
	switch typ {
	case TypeString:
		return new(ValueString)
	case TypeBool:
		return new(ValueBool)
	case TypeFloat:
		return new(ValueFloat)
	case TypeDouble:
		return new(ValueDouble)
	case TypeUDim:
		return new(ValueUDim)
	case TypeUDim2:
		return new(ValueUDim2)
	case TypeBrickColor:
		return new(ValueBrickColor)
	case TypeColor3:
		return new(ValueColor3)
	case TypeVector2:
		return new(ValueVector2)
	case TypeVector3:
		return new(ValueVector3)
	case TypeNumberSequence:
		return new(ValueNumberSequence)
	case TypeColorSequence:
		return new(ValueColorSequence)
	case TypeNumberRange:
		return new(ValueNumberRange)
	case TypeRect:
		return new(ValueRect)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// type ValueNull struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueEmpty struct{}

////////////////////////////////////////////////////////////////////////////////

type ValueString string

func (ValueString) Type() Type {
	return TypeString
}

func (v *ValueString) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a string
	if br.String(&a) {
		return br.N(), fmt.Errorf("String: %w", br.Err())
	}
	*v = ValueString(a)
	return br.End()
}

func (v ValueString) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.String(string(v)) {
		return bw.N(), fmt.Errorf("String: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueBool bool

func (ValueBool) Type() Type {
	return TypeBool
}

func (v *ValueBool) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a byte
	if br.Number(&a) {
		return br.N(), fmt.Errorf("Bool: %w", br.Err())
	}
	*v = a != 0
	return br.End()
}

func (v ValueBool) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if v {
		bw.Number(byte(1))
	} else {
		bw.Number(byte(0))
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

// type ValueInt struct{}

////////////////////////////////////////////////////////////////////////////////

type ValueFloat float32

func (ValueFloat) Type() Type {
	return TypeFloat
}

func (v *ValueFloat) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a float32
	if br.Number(&a) {
		return br.N(), fmt.Errorf("Float: %w", br.Err())
	}
	*v = ValueFloat(a)
	return br.End()
}

func (v ValueFloat) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(float32(v)) {
		return bw.N(), fmt.Errorf("Float: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueDouble float64

func (ValueDouble) Type() Type {
	return TypeDouble
}

func (v *ValueDouble) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a float64
	if br.Number(&a) {
		return br.N(), fmt.Errorf("Double: %w", br.Err())
	}
	*v = ValueDouble(a)
	return br.End()
}

func (v ValueDouble) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(float64(v)) {
		return bw.N(), fmt.Errorf("Double: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

// type ValueArray struct{}

////////////////////////////////////////////////////////////////////////////////

type Entry struct {
	Key   string
	Value Value
}

type ValueDictionary []Entry

func (v *ValueDictionary) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var length uint32
	if br.Number(&length) {
		return br.N(), fmt.Errorf("Dictionary length: %w", br.Err())
	}
	d := make(ValueDictionary, length)
	for i := range d {
		var key string
		if br.String(&key) {
			return br.N(), fmt.Errorf("Dictionary[%d](%q) key: %w", i, key, br.Err())
		}
		var typ byte
		if br.Number(&typ) {
			return br.N(), fmt.Errorf("Dictionary[%d](%q) type: %w", i, key, br.Err())
		}
		value := NewValue(Type(typ))
		if value == nil {
			return br.N(), fmt.Errorf("Dictionary[%d](%q) value: unknown data type 0x%02X", i, key, typ)
		}
		if br.Add(value.ReadFrom(r)) {
			return br.N(), fmt.Errorf("Dictionary[%d](%q) value: %w", i, key, br.Err())
		}
		d[i] = Entry{Key: key, Value: value}
	}
	*v = d
	return br.End()
}

func (v ValueDictionary) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(uint32(len(v))) {
		return bw.N(), fmt.Errorf("Dictionary length: %w", bw.Err())
	}
	for i, entry := range v {
		if bw.String(entry.Key) {
			return bw.N(), fmt.Errorf("Dictionary[%d](%q) key: %w", i, entry.Key, bw.Err())
		}
		if bw.Number(byte(entry.Value.Type())) {
			return bw.N(), fmt.Errorf("Dictionary[%d](%q) type: %w", i, entry.Key, bw.Err())
		}
		if bw.Add(entry.Value.WriteTo(w)) {
			return bw.N(), fmt.Errorf("Dictionary[%d](%q) value: %w", i, entry.Key, bw.Err())
		}
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueUDim struct {
	Scale  float32
	Offset int32
}

func (ValueUDim) Type() Type {
	return TypeUDim
}

func (v *ValueUDim) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueUDim
	if br.Number(&a.Scale) {
		return br.N(), fmt.Errorf("UDim.Scale: %w", br.Err())
	}
	if br.Number(&a.Offset) {
		return br.N(), fmt.Errorf("UDim.Offset: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueUDim) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.Scale) {
		return bw.N(), fmt.Errorf("UDim.Scale: %w", bw.Err())
	}
	if bw.Number(v.Offset) {
		return bw.N(), fmt.Errorf("UDim.Offset: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueUDim2 struct {
	X ValueUDim
	Y ValueUDim
}

func (ValueUDim2) Type() Type {
	return TypeUDim2
}

func (v *ValueUDim2) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueUDim2
	if br.Add((&a.X).ReadFrom(r)) {
		return br.N(), fmt.Errorf("UDim2.X: %w", br.Err())
	}
	if br.Add((&a.Y).ReadFrom(r)) {
		return br.N(), fmt.Errorf("UDim2.Y: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueUDim2) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Add(v.X.WriteTo(w)) {
		return bw.N(), fmt.Errorf("UDim2.X: %w", bw.Err())
	}
	if bw.Add(v.Y.WriteTo(w)) {
		return bw.N(), fmt.Errorf("UDim2.Y: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

// type ValueRay struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueFaces struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueAxes struct{}

////////////////////////////////////////////////////////////////////////////////

type ValueBrickColor uint32

func (ValueBrickColor) Type() Type {
	return TypeBrickColor
}

func (v *ValueBrickColor) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a uint32
	if br.Number(&a) {
		return br.N(), fmt.Errorf("BrickColor: %w", br.Err())
	}
	*v = ValueBrickColor(a)
	return br.End()
}

func (v ValueBrickColor) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(uint32(v)) {
		return bw.N(), fmt.Errorf("BrickColor: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueColor3 struct {
	R float32
	G float32
	B float32
}

func (ValueColor3) Type() Type {
	return TypeColor3
}

func (v *ValueColor3) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueColor3
	if br.Number(&a.R) {
		return br.N(), fmt.Errorf("Color3.R: %w", br.Err())
	}
	if br.Number(&a.G) {
		return br.N(), fmt.Errorf("Color3.G: %w", br.Err())
	}
	if br.Number(&a.B) {
		return br.N(), fmt.Errorf("Color3.B: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueColor3) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.R) {
		return bw.N(), fmt.Errorf("Color3.R: %w", bw.Err())
	}
	if bw.Number(v.G) {
		return bw.N(), fmt.Errorf("Color3.G: %w", bw.Err())
	}
	if bw.Number(v.B) {
		return bw.N(), fmt.Errorf("Color3.B: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueVector2 struct {
	X float32
	Y float32
}

func (ValueVector2) Type() Type {
	return TypeVector2
}

func (v *ValueVector2) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueVector2
	if br.Number(&a.X) {
		return br.N(), fmt.Errorf("Vector2.X: %w", br.Err())
	}
	if br.Number(&a.Y) {
		return br.N(), fmt.Errorf("Vector2.Y: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueVector2) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.X) {
		return bw.N(), fmt.Errorf("Vector2.X: %w", bw.Err())
	}
	if bw.Number(v.Y) {
		return bw.N(), fmt.Errorf("Vector2.Y: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueVector3 struct {
	X float32
	Y float32
	Z float32
}

func (ValueVector3) Type() Type {
	return TypeVector3
}

func (v *ValueVector3) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueVector3
	if br.Number(&a.X) {
		return br.N(), fmt.Errorf("Vector3.X: %w", br.Err())
	}
	if br.Number(&a.Y) {
		return br.N(), fmt.Errorf("Vector3.Y: %w", br.Err())
	}
	if br.Number(&a.Z) {
		return br.N(), fmt.Errorf("Vector3.Z: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueVector3) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.X) {
		return bw.N(), fmt.Errorf("Vector3.X: %w", bw.Err())
	}
	if bw.Number(v.Y) {
		return bw.N(), fmt.Errorf("Vector3.Y: %w", bw.Err())
	}
	if bw.Number(v.Z) {
		return bw.N(), fmt.Errorf("Vector3.Z: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

// type ValueVector2int16 struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueVector3int16 struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueCFrame struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueEnumItem struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueUnknown struct{}

////////////////////////////////////////////////////////////////////////////////

type ValueNumberSequence []ValueNumberSequenceKeypoint

func (ValueNumberSequence) Type() Type {
	return TypeNumberSequence
}

func (v *ValueNumberSequence) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var length uint32
	if br.Number(&length) {
		return br.N(), fmt.Errorf("NumberSequence length: %w", br.Err())
	}
	s := make(ValueNumberSequence, length)
	for i := range s {
		var k ValueNumberSequenceKeypoint
		if br.Add(k.ReadFrom(r)) {
			return br.N(), fmt.Errorf("NumberSequence[%d]: %w", i, br.Err())
		}
		s[i] = k
	}
	*v = s
	return br.End()
}

func (v ValueNumberSequence) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(uint32(len(v))) {
		return bw.N(), fmt.Errorf("NumberSequence: length %w", bw.Err())
	}
	for i, k := range v {
		if bw.Add(k.WriteTo(w)) {
			return bw.N(), fmt.Errorf("NumberSequence[%d]: %w", i, bw.Err())
		}
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueNumberSequenceKeypoint struct {
	Envelope float32
	Time     float32
	Value    float32
}

func (v *ValueNumberSequenceKeypoint) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueNumberSequenceKeypoint
	if br.Number(&a.Envelope) {
		return br.N(), fmt.Errorf("NumberSequenceKeypoint.Envelope: %w", br.Err())
	}
	if br.Number(&a.Time) {
		return br.N(), fmt.Errorf("NumberSequenceKeypoint.Time: %w", br.Err())
	}
	if br.Number(&a.Value) {
		return br.N(), fmt.Errorf("NumberSequenceKeypoint.Value: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueNumberSequenceKeypoint) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.Envelope) {
		return bw.N(), fmt.Errorf("NumberSequenceKeypoint.Envelope: %w", bw.Err())
	}
	if bw.Number(v.Time) {
		return bw.N(), fmt.Errorf("NumberSequenceKeypoint.Time: %w", bw.Err())
	}
	if bw.Number(v.Value) {
		return bw.N(), fmt.Errorf("NumberSequenceKeypoint.Value: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueColorSequence []ValueColorSequenceKeypoint

func (ValueColorSequence) Type() Type {
	return TypeColorSequence
}

func (v *ValueColorSequence) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var length uint32
	if br.Number(&length) {
		return br.N(), fmt.Errorf("ColorSequence length: %w", br.Err())
	}
	s := make(ValueColorSequence, length)
	for i := range s {
		var k ValueColorSequenceKeypoint
		if br.Add(k.ReadFrom(r)) {
			return br.N(), fmt.Errorf("ColorSequence[%d]: %w", i, br.Err())
		}
		s[i] = k
	}
	*v = s
	return br.End()
}

func (v ValueColorSequence) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(uint32(len(v))) {
		return bw.N(), fmt.Errorf("ColorSequence length: %w", bw.Err())
	}
	for i, k := range v {
		if bw.Add(k.WriteTo(w)) {
			return bw.N(), fmt.Errorf("ColorSequence[%d]: %w", i, bw.Err())
		}
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueColorSequenceKeypoint struct {
	Envelope float32
	Time     float32
	Value    ValueColor3
}

func (v *ValueColorSequenceKeypoint) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueColorSequenceKeypoint
	if br.Number(&a.Envelope) {
		return br.N(), fmt.Errorf("ColorSequenceKeypoint.Envelope: %w", br.Err())
	}
	if br.Number(&a.Time) {
		return br.N(), fmt.Errorf("ColorSequenceKeypoint.Time: %w", br.Err())
	}
	if br.Add((&a.Value).ReadFrom(r)) {
		return br.N(), fmt.Errorf("ColorSequenceKeypoint.Value: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueColorSequenceKeypoint) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.Envelope) {
		return bw.N(), fmt.Errorf("ColorSequenceKeypoint.Envelope: %w", bw.Err())
	}
	if bw.Number(v.Time) {
		return bw.N(), fmt.Errorf("ColorSequenceKeypoint.Time: %w", bw.Err())
	}
	if bw.Add(v.Value.WriteTo(w)) {
		return bw.N(), fmt.Errorf("ColorSequenceKeypoint.Value: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueNumberRange struct {
	Min float32
	Max float32
}

func (ValueNumberRange) Type() Type {
	return TypeNumberRange
}

func (v *ValueNumberRange) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueNumberRange
	if br.Number(&a.Min) {
		return br.N(), fmt.Errorf("NumberRange.Min: %w", br.Err())
	}
	if br.Number(&a.Max) {
		return br.N(), fmt.Errorf("NumberRange.Max: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueNumberRange) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Number(v.Min) {
		return bw.N(), fmt.Errorf("NumberRange.Min: %w", bw.Err())
	}
	if bw.Number(v.Max) {
		return bw.N(), fmt.Errorf("NumberRange.Max: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

type ValueRect struct {
	Min ValueVector2
	Max ValueVector2
}

func (ValueRect) Type() Type {
	return TypeRect
}

func (v *ValueRect) ReadFrom(r io.Reader) (n int64, err error) {
	br := newBinaryReader(r)
	var a ValueRect
	if br.Add((&a.Min).ReadFrom(r)) {
		return br.N(), fmt.Errorf("Rect.Min: %w", br.Err())
	}
	if br.Add((&a.Max).ReadFrom(r)) {
		return br.N(), fmt.Errorf("Rect.Max: %w", br.Err())
	}
	*v = a
	return br.End()
}

func (v ValueRect) WriteTo(w io.Writer) (n int64, err error) {
	bw := newBinaryWriter(w)
	if bw.Add(v.Min.WriteTo(w)) {
		return bw.N(), fmt.Errorf("Rect.Min: %w", bw.Err())
	}
	if bw.Add(v.Max.WriteTo(w)) {
		return bw.N(), fmt.Errorf("Rect.Max: %w", bw.Err())
	}
	return bw.End()
}

////////////////////////////////////////////////////////////////////////////////

// type ValuePhysicalProperties struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueUnknown struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueRegion3 struct{}

////////////////////////////////////////////////////////////////////////////////

// type ValueRegion3int16 struct{}
