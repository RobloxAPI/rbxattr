# Attribute Binary Format
This document describes the format for serialized attributes in Roblox.

## Version 0
This section describes the first version of the format. Currently, there is only
one version. No bytes are reserved for indicating the version.

This version encodes the `Instance.AttributesSerialize` property. Because no
version indicator is included in the format itself, it is assumed that separate
versions of the format will be encoded in separate properties.

### Definitions
The following types are defined for the purposes of this specification.

Type         | Description
-------------|------------
`uint:N`     | An unsigned integer of length `N` bits.
`int:N`      | A signed integer of length `N` bits.
`float:N`    | An IEEE 754 floating-point number of length `N` bits.
`[N]Type`    | An array of constant length `N`, each entry of type `Type`.
`[Size]Type` | An array of length determined by the value of field `Size`, each entry of type `Type`.

All values are little-endian.

### Entrypoint
Attributes are encoded as a [Dictionary][Dictionary] that maps attribute names
to values.

Several limitations may be applied to the Key field of an [Entry][Entry] in this
Dictionary:
- Must not have a length greater than 100 bytes.
- May only contain alphanumeric bytes and underscores. (`^[0-9A-Za-z_]*$`)
- Must not begin with `RBX`, which is reserved for Roblox.

### Value
[Value]: #user-content-value

- **Size**: 4+N bytes

Field | Offset |   Size | Type                                          | Condition
------|-------:|-------:|-----------------------------------------------|----------
Type  |      0 |      1 | <code>uint:8</code>                           |
Value |      1 |  4+1*N | <code>[String][String]</code>                 | Type == 0x02
Value |      1 |      1 | <code>[Bool][Bool]</code>                     | Type == 0x03
Value |      1 |      4 | <code>[Float][Float]</code>                   | Type == 0x05
Value |      1 |      8 | <code>[Double][Double]</code>                 | Type == 0x06
Value |      1 |      8 | <code>[UDim][UDim]</code>                     | Type == 0x09
Value |      1 |     16 | <code>[UDim2][UDim2]</code>                   | Type == 0x0A
Value |      1 |      4 | <code>[BrickColor][BrickColor]</code>         | Type == 0x0E
Value |      1 |     12 | <code>[Color3][Color3]</code>                 | Type == 0x0F
Value |      1 |      8 | <code>[Vector2][Vector2]</code>               | Type == 0x10
Value |      1 |     12 | <code>[Vector3][Vector3]</code>               | Type == 0x11
Value |      1 | 4+12*N | <code>[NumberSequence][NumberSequence]</code> | Type == 0x17
Value |      1 | 4+20*N | <code>[ColorSequence][ColorSequence]</code>   | Type == 0x19
Value |      1 |      8 | <code>[NumberRange][NumberRange]</code>       | Type == 0x1B
Value |      1 |     16 | <code>[Rect][Rect]</code>                     | Type == 0x1C

The value of the Type field determines how the value of the Value field is
decoded.

### String
[String]: #user-content-string

- **Size**: 4+N bytes
- **Numeric type**: 0x02 (2)
- **Decoded type**: string (string)

| Field | Offset | Size | Type                        |
|-------|-------:|-----:|-----------------------------|
| Size  |      0 |    4 | <code>uint:32</code>        |
| Value |      4 |    N | <code>\[Size\]uint:8</code> |

The Value field determines the bytes of the string.

### Bool
[Bool]: #user-content-bool

- **Size**: 1 byte
- **Numeric type**: 0x03 (3)
- **Decoded type**: bool (boolean)

| Type                |
|---------------------|
| <code>uint:8</code> |

A value of of `0` is decoded to `false`, while any other value is decoded to
`true`.

`false` is encoded to `0`, and `true` is encoded to 1.

### Float
[Float]: #user-content-float

- **Size**: 4 bytes
- **Numeric type**: 0x05 (5)
- **Decoded type**: float (number)

| Type                  |
|-----------------------|
| <code>float:32</code> |

Roblox encodes a number as this type if all of the four most significant bytes
are zero.

### Double
[Double]: #user-content-double

- **Size**: 8 bytes
- **Numeric type**: 0x06 (6)
- **Decoded type**: double (number)

| Type                  |
|-----------------------|
| <code>float:64</code> |

Roblox encodes a number as this type if any of the four most significant bytes
are non-zero.

### Dictionary
[Dictionary]: #user-content-dictionary

- **Size**: 4+N bytes

| Field | Offset | Size | Type                                |
|-------|-------:|-----:|-------------------------------------|
| Size  |      0 |    4 | <code>uint:32</code>                |
| Value |      4 |    N | <code>\[Size\][Entry][Entry]</code> |

When decoding, entries with a duplicate Key are discarded.

### Entry
[Entry]: #user-content-entry

- **Size**: A+B bytes

| Field | Offset | Size | Type                          |
|-------|-------:|-----:|-------------------------------|
| Key   |      0 |    A | <code>[String][String]</code> |
| Value |      A |    B | <code>[Value][Value]</code>   |

### UDim
[UDim]: #user-content-udim

- **Size**: 8 bytes
- **Numeric type**: 0x09 (9)
- **Decoded type**: UDim (userdata)

| Field  | Offset | Size | Type                  |
|--------|-------:|-----:|-----------------------|
| Scale  |      0 |    4 | <code>float:32</code> |
| Offset |      4 |    4 | <code>int:32</code>   |

### UDim2
[UDim2]: #user-content-udim2

- **Size**: 16 bytes
- **Numeric type**: 0x0A (10)
- **Decoded type**: UDim2 (userdata)

| Field | Offset | Size | Type                      |
|-------|-------:|-----:|---------------------------|
| X     |      0 |    8 | <code>[UDim][UDim]</code> |
| Y     |      8 |   16 | <code>[UDim][UDim]</code> |

### BrickColor
[BrickColor]: #user-content-brickcolor

- **Size**: 4 bytes
- **Numeric type**: 0x0E (14)
- **Decoded type**: BrickColor (userdata)

| Type                 |
|----------------------|
| <code>uint:32</code> |

The value corresponds to the Number field of the BrickColor value. Unknown
values are interpreted as the default BrickColor.

### Color3
[Color3]: #user-content-color3

- **Size**: 12 bytes
- **Numeric type**: 0x0F (15)
- **Decoded type**: Color3 (userdata)

| Field | Offset | Size | Type                  |
|-------|-------:|-----:|-----------------------|
| R     |      0 |    4 | <code>float:32</code> |
| G     |      4 |    4 | <code>float:32</code> |
| B     |      8 |    4 | <code>float:32</code> |

### Vector2
[Vector2]: #user-content-vector2

- **Size**: 8 bytes
- **Numeric type**: 0x10 (16)
- **Decoded type**: Vector2 (userdata)

| Field | Offset | Size | Type                  |
|-------|-------:|-----:|-----------------------|
| X     |      0 |    4 | <code>float:32</code> |
| Y     |      4 |    4 | <code>float:32</code> |

### Vector3
[Vector3]: #user-content-vector3

- **Size**: 12 bytes
- **Numeric type**: 0x11 (17)
- **Decoded type**: Vector3 (userdata)

| Field | Offset | Size | Type                  |
|-------|-------:|-----:|-----------------------|
| X     |      0 |    4 | <code>float:32</code> |
| Y     |      4 |    4 | <code>float:32</code> |
| Z     |      8 |    4 | <code>float:32</code> |

### NumberSequence
[NumberSequence]: #user-content-numbersequence

- **Size**: 4+12*N bytes
- **Numeric type**: 0x17 (23)
- **Decoded type**: NumberSequence (userdata)

| Field  | Offset | Size | Type                                        |
|--------|-------:|-----:|---------------------------------------------|
| Size   |      0 |    4 | <code>uint:32</code>                        |
| Values |      4 | 12*N | <code>\[Size\]NumberSequenceKeypoint</code> |

### NumberSequenceKeypoint
[NumberSequenceKeypoint]: #user-content-numbersequencekeypoint

| Field    | Offset | Size | Type                  |
|----------|-------:|-----:|-----------------------|
| Envelope |      0 |    4 | <code>float:32</code> |
| Time     |      4 |    4 | <code>float:32</code> |
| Value    |      8 |    4 | <code>float:32</code> |

### ColorSequence
[ColorSequence]: #user-content-colorsequence

- **Size**: 4+20*N bytes
- **Numeric type**: 0x19 (25)
- **Decoded type**: ColorSequence (userdata)

| Field  | Offset | Size | Type                                       |
|--------|-------:|-----:|--------------------------------------------|
| Size   |      0 |    4 | <code>uint:32</code>                       |
| Values |      4 | 20*N | <code>\[Size\]ColorSequenceKeypoint</code> |

### ColorSequenceKeypoint
[ColorSequenceKeypoint]: #user-content-colorsequencekeypoint

- **Size**: 20 bytes

| Field    | Offset | Size | Type                          |
|----------|-------:|-----:|-------------------------------|
| Envelope |      0 |    4 | <code>float:32</code>         |
| Time     |      4 |    4 | <code>float:32</code>         |
| Value    |      8 |   12 | <code>[Color3][Color3]</code> |

### NumberRange
[NumberRange]: #user-content-numberrange

- **Size**: 8 bytes
- **Numeric type**: 0x1B (27)
- **Decoded type**: NumberRange (userdata)

| Field | Offset | Size | Type                  |
|-------|-------:|-----:|-----------------------|
| Min   |      0 |    4 | <code>float:32</code> |
| Max   |      4 |    4 | <code>float:32</code> |

### Rect
[Rect]: #user-content-rect

- **Size**: 16 bytes
- **Numeric type**: 0x1C (28)
- **Decoded type**: Rect (userdata)

| Field | Offset | Size | Type                            |
|-------|-------:|-----:|---------------------------------|
| Min   |      0 |    8 | <code>[Vector2][Vector2]</code> |
| Max   |      8 |    8 | <code>[Vector2][Vector2]</code> |
