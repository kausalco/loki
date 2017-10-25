package annotation

import (
	"encoding/binary"
	"regexp"

	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

type Matcher interface {
	Match([]*zipkincore.BinaryAnnotation) bool
}

type Matchers []Matcher

func (ms Matchers) Match(bas []*zipkincore.BinaryAnnotation) bool {
	for _, m := range ms {
		if !m.Match(bas) {
			return false
		}
	}
	return true
}

type noopMatcher struct{}

func (noopMatcher) Match(_ []*zipkincore.BinaryAnnotation) bool {
	return true
}

var NoopMatcher Matcher = noopMatcher{}

type eqStr struct {
	key, value string
}

func EqStr(key, value string) Matcher {
	return eqStr{key, value}
}

func (m eqStr) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			string(a.GetValue()) == m.value {
			return true
		}
	}
	return false
}

type neStr struct {
	key, value string
}

func NeStr(key, value string) Matcher {
	return neStr{key, value}
}

func (m neStr) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			string(a.GetValue()) != m.value {
			return true
		}
	}
	return false
}

func intVal(a *zipkincore.BinaryAnnotation) int64 {
	switch a.GetAnnotationType() {
	case zipkincore.AnnotationType_I16:
		return int64(binary.BigEndian.Uint16(a.Value))
	case zipkincore.AnnotationType_I32:
		return int64(binary.BigEndian.Uint32(a.Value))
	case zipkincore.AnnotationType_I64:
		return int64(binary.BigEndian.Uint64(a.Value))
	}
	return 0
}

func isInt(a *zipkincore.BinaryAnnotation) bool {
	switch a.GetAnnotationType() {
	case zipkincore.AnnotationType_I16, zipkincore.AnnotationType_I32, zipkincore.AnnotationType_I64:
		return true
	}
	return false
}

type eqInt struct {
	key   string
	value int64
}

func EqInt(key string, value int64) Matcher {
	return eqInt{key, value}
}

func (m eqInt) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key && isInt(a) && intVal(a) == m.value {
			return true
		}
	}
	return false
}

type neInt struct {
	key   string
	value int64
}

func NeInt(key string, value int64) Matcher {
	return neInt{key, value}
}

func (m neInt) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key && isInt(a) && intVal(a) != m.value {
			return true
		}
	}
	return false
}

type re struct {
	key string
	re  *regexp.Regexp
}

func Re(key, expr string) Matcher {
	r, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return re{key, r}
}

func (m re) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			m.re.Match(a.GetValue()) {
			return true
		}
	}
	return false
}

type nre struct {
	key string
	re  *regexp.Regexp
}

func Nre(key, expr string) Matcher {
	r, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return nre{key, r}
}

func (m nre) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			!m.re.Match(a.GetValue()) {
			return true
		}
	}
	return false
}
