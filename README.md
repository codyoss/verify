# verify

Package verify uses struct field tags to verify data.

[![GoDoc](https://godoc.org/github.com/codyoss/verify?status.svg)](https://godoc.org/github.com/codyoss/verify)
[![Build Status](https://cloud.drone.io/api/badges/codyoss/verify/status.svg)](https://cloud.drone.io/codyoss/verify)
[![codecov](https://codecov.io/gh/codyoss/verify/branch/master/graph/badge.svg)](https://codecov.io/gh/codyoss/verify)
[![Go Report Card](https://goreportcard.com/badge/github.com/codyoss/verify)](https://goreportcard.com/report/github.com/codyoss/verify)

## Tags supported

- `minSize` -- specifies the minimum allowable length of a field. This can only be used on the following types: string,
slice, array, or map.

- `maxSize` -- specifies the maximum allowable length of a field. This can only be used on the following types: string,
slice, array, or map.

- `min` -- specifies the minimum allowable value of a field. This should only be used on types that can be parsed into
an int64 or float64.

- `max` -- specifies the maximum allowable value of a field. This should only be used on types that can be parsed into
an int64 or float64.

- `required` -- specifies the field may not be set to the zero value for the given type. This may be used on any types
except arrays and structs.

## Example usage

Here is an example of the usage of each tag:

```golang
type Foo struct {
    A []string  `verify:"minSize=5"`
    B string    `verify:"maxSize=10"`
    C int8      `verify:"min=3"`
    D float32   `verify:"max=1.2"`
    E int64     `verify:"min=3,max=7"`
    F *bool     `verify:"required"`
}
```

## Limitations

1. verify only supports working with flat structures at the moment; it will not work with inner/embedded structs. Also,
2. Because this package makes use of reflection the tags may only be used on exported fields.

## Blog Post

I wrote a blog post about how to make your own struct field tags. [Check it out here!](https://medium.com/@codyoss/creating-your-own-struct-field-tags-in-go-c6c86727eff)