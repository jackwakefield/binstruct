package binstruct

import "github.com/pkg/errors"
import "math/bits"

type StringType = string

const (
	// StringFixed is a fixed length string, the length is specified
	// by the options Len or LenField.
	StringFixed StringType = "fixed"
	// StringNullTerminated is a string with no predetermined length.
	// Reading the string ends when a null-byte is reached.
	// Writing the string is suffixed with a null-byte.
	StringNullTerminated StringType = "null"
	// StringInt8 is a string with the length determined by a
	// 1-byte prefix.
	StringInt8 StringType = "int8"
	// StringInt16 is a string with the length determined by a
	// 2-byte prefix.
	StringInt16 StringType = "int16"
	// StringInt32 is a string with the length determined by a
	// 4-byte prefix.
	StringInt32 StringType = "int32"
	// StringInt64 is a string with the length determined by a
	// 8-byte prefix.
	StringInt64 StringType = "int64"
)

// FieldOptions define the options used when reading and writing
// the struct field.
type FieldOptions struct {
	// Skip is a relative position in the stream where the value
	// will be read from or written to.
	Skip int64
	// Offset is an absolute position in the stream where the
	// value will be read from or written to.
	Offset int64
	// OffsetField is the name of a sibling field which will be
	// used as the offset.
	OffsetField string
	// Len is the fixed size of the slice or string being read or
	// written. In the case of a string this is used when StringType
	// is fixed. If the actual length of the string is lower than the
	// fixed length, the remaining bytes will be padded with a null
	// character, which can be overriden by setting StringPad.
	Len int64
	// LenField is the name of a sibling field which will be used as
	// the slice or string length.
	LenField string
	// StringType is the type of string to be read or written.
	StringType StringType
	// StringPad is the byte used for the remainder of fixed-length strings.
	StringPad byte
	// Align determines whether the struct data is aligned.
	Align bool
	// AlignBytes is the number of bytes to align the struct data to.
	AlignBytes int64
	// Mask is applied to integer values when reading (XOR) and writing (OR).
	Mask uint64
}

var defaultFieldOptions = &FieldOptions{
	Skip:        0,
	Offset:      0,
	OffsetField: "",
	Len:         0,
	LenField:    "",
	StringType:  StringFixed,
	StringPad:   0,
	Align:       false,
	AlignBytes:  8,
	Mask:        0,
}

// SetDefaultOptions sets the default options for fields, these are overriden
// by the tags defined alongside the struct field.
func SetDefaultOptions(options *FieldOptions) {
	defaultFieldOptions = options
}

// parseTagFieldOptions creates field options from the given tag.
func parseTagFieldOptions(t tag) (*FieldOptions, error) {
	// make a shallow-copy of the default options
	options := &FieldOptions{}
	*options = *defaultFieldOptions
	if t != nil {
		var err error
		if t.Contains("skip") {
			if options.Skip, err = t.Int64("skip"); err != nil {
				return nil, errors.Wrap(err, "failed to parse skip value")
			}
		}
		if t.Contains("offset") {
			if options.Offset, err = t.Int64("offset"); err != nil {
				return nil, errors.Wrap(err, "failed to parse offset value")
			}
		}
		if t.Contains("offsetfield") {
			if options.OffsetField, err = t.String("offsetfield"); err != nil {
				return nil, errors.Wrap(err, "failed to parse offsetfield value")
			}
		}
		if t.Contains("len") {
			if options.Len, err = t.Int64("len"); err != nil {
				return nil, errors.Wrap(err, "failed to parse len value")
			}
		}
		if t.Contains("lenfield") {
			if options.LenField, err = t.String("lenfield"); err != nil {
				return nil, errors.Wrap(err, "failed to parse lenfield value")
			}
		}
		if t.Contains("stringtype") {
			if options.StringType, err = t.String("stringtype"); err != nil {
				return nil, errors.Wrap(err, "failed to parse stringtype value")
			}
		}
		if t.Contains("stringpad") {
			if options.StringPad, err = t.Byte("stringpad"); err != nil {
				return nil, errors.Wrap(err, "failed to parse stringpad value")
			}
		}
		if t.Contains("align") {
			if options.Align, err = t.Bool("align"); err != nil {
				return nil, errors.Wrap(err, "failed to parse align value")
			}
		}
		if t.Contains("alignbytes") {
			if options.AlignBytes, err = t.Int64("alignbytes"); err != nil {
				return nil, errors.Wrap(err, "failed to parse alignbytes value")
			}
		}
		if t.Contains("mask") {
			if mask, err := t.Int64("mask"); err != nil {
				options.Mask = uint64(mask)
				return nil, errors.Wrap(err, "failed to parse mask value")
			}
		}
	}
	return options, nil
}

// MaskBits returns the minimum number of bits required to represent the mask option.
func (o *FieldOptions) MaskBits() int {
	return bits.Len64(o.Mask)
}
