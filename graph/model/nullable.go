package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type NullableString struct {
	Value *string
	Set   bool
}

func (n *NullableString) UnmarshalGQL(v any) error {
	n.Set = true
	if v == nil {
		n.Value = nil
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("NullableString must be a string, got %T", v)
	}
	n.Value = &s
	return nil
}

func (n NullableString) MarshalGQL(w io.Writer) {
	if n.Value == nil {
		io.WriteString(w, "null")
		return
	}
	io.WriteString(w, strconv.Quote(*n.Value))
}

type NullableInt struct {
	Value *int
	Set   bool
}

func (n *NullableInt) UnmarshalGQL(v any) error {
	n.Set = true
	if v == nil {
		n.Value = nil
		return nil
	}
	switch val := v.(type) {
	case int:
		n.Value = &val
	case int64:
		i := int(val)
		n.Value = &i
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return fmt.Errorf("NullableInt: %w", err)
		}
		iv := int(i)
		n.Value = &iv
	default:
		return fmt.Errorf("NullableInt must be an int, got %T", v)
	}
	return nil
}

func (n NullableInt) MarshalGQL(w io.Writer) {
	if n.Value == nil {
		io.WriteString(w, "null")
		return
	}
	io.WriteString(w, strconv.Itoa(*n.Value))
}

type NullableFloat struct {
	Value *float64
	Set   bool
}

func (n *NullableFloat) UnmarshalGQL(v any) error {
	n.Set = true
	if v == nil {
		n.Value = nil
		return nil
	}
	switch val := v.(type) {
	case float64:
		n.Value = &val
	case int:
		f := float64(val)
		n.Value = &f
	case int64:
		f := float64(val)
		n.Value = &f
	case json.Number:
		f, err := val.Float64()
		if err != nil {
			return fmt.Errorf("NullableFloat: %w", err)
		}
		n.Value = &f
	default:
		return fmt.Errorf("NullableFloat must be a number, got %T", v)
	}
	return nil
}

func (n NullableFloat) MarshalGQL(w io.Writer) {
	if n.Value == nil {
		io.WriteString(w, "null")
		return
	}
	io.WriteString(w, strconv.FormatFloat(*n.Value, 'f', -1, 64))
}

type NullableBool struct {
	Value *bool
	Set   bool
}

func (n *NullableBool) UnmarshalGQL(v any) error {
	n.Set = true
	if v == nil {
		n.Value = nil
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		return fmt.Errorf("NullableBool must be a bool, got %T", v)
	}
	n.Value = &b
	return nil
}

func (n NullableBool) MarshalGQL(w io.Writer) {
	if n.Value == nil {
		io.WriteString(w, "null")
		return
	}
	io.WriteString(w, strconv.FormatBool(*n.Value))
}
