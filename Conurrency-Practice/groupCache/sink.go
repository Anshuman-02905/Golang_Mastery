package main

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

type Sink interface {
	SetString(s string) error
	SetBytes(v []byte) error
	SetProto(m proto.Message) error
	view() (ByteView, error)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func setSinkView(s Sink, v ByteView) error {
	type viewSetter interface {
		setView(v ByteView) error
	}

	// Check if s supports the special case interface
	if vs, ok := s.(viewSetter); ok {
		return vs.setView(v)
	}

	//  Fall back to general behavior
	if v.b != nil {
		return s.SetBytes(v.b)
	}
	return s.SetString(v.s)
}

// StringSink returns a Sink that populates the provided string pointer.

func StringSink(sp *string) Sink {
	return &stringSink{sp: sp}
}

type stringSink struct {
	sp *string
	v  ByteView
}

func (s *stringSink) SetString(v string) error {
	s.v.b = nil
	s.v.s = v
	*s.sp = v
	return nil

}

func (s *stringSink) SetBytes(v []byte) error {
	return s.SetString(string(v))

}
func (s *stringSink) SetProto(m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	s.v.b = b
	*s.sp = string(b)
	return nil

}

func (s *stringSink) view() (ByteView, error) {
	return s.v, nil

}

// ByteViewSink returns a Sink that populates a ByteView.
func ByteViewSink(dst *ByteView) Sink {
	if dst == nil {
		panic("nil dst")
	}
	return &byteViewSink{dst: dst}
}

type byteViewSink struct {
	dst *ByteView

	// if this code ever ends up tracking that at least one set*
	// method was called, don't make it an error to call set
	// methods multiple times. Lorry's payload.go does that, and
	// it makes sense. The comment at the top of this file about
	// "exactly one of the Set methods" is overly strict. We
	// really care about at least once (in a handler), but if
	// multiple handlers fail (or multiple functions in a program
	// using a Sink), it's okay to re-use the same one.
}

func (s *byteViewSink) setView(v ByteView) error {
	*s.dst = v
	return nil
}
func (s *byteViewSink) view() (ByteView, error) {
	return *s.dst, nil

}

func (s *byteViewSink) SetString(v string) error {

	*s.dst = ByteView{s: v}
	return nil

}
func (s *byteViewSink) SetBytes(b []byte) error {
	*s.dst = ByteView{b: cloneBytes(b)}
	return nil

}

func (s *byteViewSink) SetProto(m proto.Message) error {
	b, err := proto.Marshal(m)

	if err != nil {
		return err
	}
	*s.dst = ByteView{
		b: b,
	}
	return nil

}

// ProtoSink returns a sink that unmarshals binary proto values into m.

func ProtoSink(m proto.Message) Sink {
	return &protoSink{
		dst: m,
	}
}

type protoSink struct {
	dst proto.Message //authorative value both are binary data
	typ stringSink

	v ByteView //encoded both are binary data
}

func (s *protoSink) view() (ByteView, error) {
	return s.v, nil
}

//WHY PROTOSINK HAS EMPTY string  but in other setString we set the string to ByteView

func (s *protoSink) SetString(v string) error {

	b := []byte(v)
	err := proto.Unmarshal(b, s.dst)
	if err != nil {
		return err
	}
	s.v.b = b
	s.v.s = ""
	return nil

}

func (s *protoSink) SetBytes(b []byte) error {
	err := proto.Unmarshal(b, s.dst)
	if err != nil {
		return err
	}
	s.v.b = cloneBytes(b)
	s.v.s = ""
	return nil
}

func (s *protoSink) SetProto(m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(b, s.dst)
	if err != nil {
		return err
	}

	s.v.b = b
	s.v.s = ""
	return nil

}

// ALLOCATING BYTESLICE SINK
func AllocatingByteSlice(dst *[]byte) Sink {
	return &allocBytesSink{dst: dst}
}

type allocBytesSink struct {
	dst *[]byte
	v   ByteView
}

func (s *allocBytesSink) view() (ByteView, error) {

	return s.v, nil
}

// setview is impelented less why
func (s *allocBytesSink) setView(v ByteView) error {
	if v.b != nil {
		*s.dst = cloneBytes(v.b)

	} else {
		*s.dst = []byte(v.s)
	}
	s.v = v
	return nil
}

// Why marshal the data when setting in proto is the data coming as proto.message is a struct understood by golang
func (s *allocBytesSink) SetProto(m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return s.setBytesOwned(b)
}
func (s *allocBytesSink) SetBytes(b []byte) error {
	return s.setBytesOwned(b)
}

func (s *allocBytesSink) setBytesOwned(b []byte) error {
	if s.dst == nil {
		return errors.New("nil AllocatingByteSliceSink *[]byte dst")
	}
	*s.dst = cloneBytes(b)
	s.v.b = b
	s.v.s = ""
	return nil
}

func (s *allocBytesSink) SetString(v string) error {
	if s.dst == nil {
		return errors.New("nil AllocatingByteSliceSink *[]byte dst")
	}
	*s.dst = []byte(v)
	s.v.b = nil
	s.v.s = v
	return nil
}

func TruncatingByteSliceSink(dst *[]byte) Sink {
	return &truncBytesSink{dst: dst}
}

type truncBytesSink struct {
	dst *[]byte
	v   ByteView
}

func (s *truncBytesSink) view() (ByteView, error) {
	return s.v, nil
}
func (s *truncBytesSink) SetProto(m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return s.setBytesOwned(b)
}

func (s *truncBytesSink) SetBytes(b []byte) error {
	return s.setBytesOwned(cloneBytes(b))
}
func (s *truncBytesSink) setBytesOwned(b []byte) error {
	if s.dst == nil {
		return errors.New("nil TruncatingBytesSliceSink *[]byte dst")
	}
	n := copy(*s.dst, b)
	if n < len(*s.dst) {
		*s.dst = (*s.dst)[:n]
	}
	s.v.b = b
	s.v.s = ""
	return nil

}

func (s *truncBytesSink) SetString(v string) error {
	if s.dst == nil {
		return errors.New("nil TruncatingByteSliceSink *[]byte dst")
	}
	n := copy(*s.dst, v)
	if n < len(*s.dst) {
		*s.dst = (*s.dst)[:n]
	}
	s.v.b = nil
	s.v.s = v
	return nil
}
