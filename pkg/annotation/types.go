package annotation

import (
	"regexp"

	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

type Matcher interface {
	match([]*zipkincore.BinaryAnnotation) bool
}

type Matchers []Matcher

func (ms Matchers) match(bas []*zipkincore.BinaryAnnotation) bool {
	for _, m := range ms {
		if !m.match(bas) {
			return false
		}
	}
	return true
}

type eq struct {
	key, value string
}

func Eq(key, value string) Matcher {
	return eq{key, value}
}

func (m eq) match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			string(a.GetValue()) == m.value {
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

func (m re) match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			m.re.Match(a.GetValue()) {
			return true
		}
	}
	return false
}

type ne struct {
	key, value string
}

func Ne(key, value string) Matcher {
	return ne{key, value}
}

func (m ne) match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			string(a.GetValue()) == m.value {
			return false
		}
	}
	return true
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

func (m nre) match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			m.re.Match(a.GetValue()) {
			return false
		}
	}
	return true
}
