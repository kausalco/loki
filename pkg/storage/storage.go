package storage

import (
	"regexp"

	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

type SpanStore interface {
	Append(*zipkincore.Span) error
	ReadStore
}

type ReadStore interface {
	Services() ([]string, error)
	SpanNames(serviceName string) ([]string, error)
	Trace(id int64) (Trace, error)
	Traces(query Query) ([]Trace, error)
}

type Query struct {
	ServiceName     string
	SpanName        string
	MinDurationUS   int64
	MaxDurationUS   int64
	EndMS           int64
	StartMS         int64
	Limit           int
	AnnotationQuery KVQuery
}

type KVQuery interface {
	Match([]*zipkincore.BinaryAnnotation) bool
}

type FnMatcher func([]*zipkincore.BinaryAnnotation) bool

func (f FnMatcher) Match(as []*zipkincore.BinaryAnnotation) bool {
	return f(as)
}

func NoopQuery() KVQuery {
	return FnMatcher(func(as []*zipkincore.BinaryAnnotation) bool {
		return true
	})
}

func StrEqQuery(key, value string) KVQuery {
	return FnMatcher(func(as []*zipkincore.BinaryAnnotation) bool {
		for _, a := range as {
			if a.GetKey() == key &&
				a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
				string(a.GetValue()) == value {
				return true
			}
		}
		return false
	})
}

func StrReQuery(key, expr string) (KVQuery, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return FnMatcher(func(as []*zipkincore.BinaryAnnotation) bool {
		for _, a := range as {
			if a.GetKey() == key &&
				a.GetAnnotationType() == zipkincore.AnnotationType_STRING &&
				re.Match(a.GetValue()) {
				return true
			}
		}
		return false
	}), nil
}

type AndQuery struct {
	left, right KVQuery
}

func (q AndQuery) Match(as []*zipkincore.BinaryAnnotation) bool {
	return q.left.Match(as) && q.right.Match(as)
}

type OrQuery struct {
	left, right KVQuery
}

func (q OrQuery) Match(as []*zipkincore.BinaryAnnotation) bool {
	return q.left.Match(as) || q.right.Match(as)
}
