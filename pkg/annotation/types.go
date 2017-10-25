package annotation

import (
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

type eq struct {
	key, value string
}

func Eq(key, value string) Matcher {
	return eq{key, value}
}

func (m eq) Match(as []*zipkincore.BinaryAnnotation) bool {
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

type ne struct {
	key, value string
}

func Ne(key, value string) Matcher {
	return ne{key, value}
}

func (m ne) Match(as []*zipkincore.BinaryAnnotation) bool {
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

func (m nre) Match(as []*zipkincore.BinaryAnnotation) bool {
	for _, a := range as {
		if a.GetKey() == m.key &&
			a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
			m.re.Match(a.GetValue()) {
			return false
		}
	}
	return true
}
