package annotation

import (
	"encoding/binary"
	"testing"

	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/require"
)

func s(key, value string) *zipkincore.BinaryAnnotation {
	return &zipkincore.BinaryAnnotation{
		Key:            key,
		Value:          []byte(value),
		AnnotationType: zipkincore.AnnotationType_STRING,
	}
}

func i(key string, val uint64) *zipkincore.BinaryAnnotation {
	value := [8]byte{}
	binary.BigEndian.PutUint64(value[:], val)
	return &zipkincore.BinaryAnnotation{
		Key:            key,
		Value:          value[:],
		AnnotationType: zipkincore.AnnotationType_I64,
	}
}

func as(as ...*zipkincore.BinaryAnnotation) []*zipkincore.BinaryAnnotation {
	return as
}

func TestTypes(t *testing.T) {
	for _, tc := range []struct {
		matcher    string
		annotation []*zipkincore.BinaryAnnotation
		match      bool
	}{
		{`{foo="bar"}`, as(s("foo", "bar")), true},
		{`{foo="bar"}`, as(s("foo", "baz")), false},
		{`{foo="bar"}`, as(), false},

		{`{foo!="bar"}`, as(s("foo", "bar")), false},
		{`{foo!="bar"}`, as(s("foo", "baz")), true},
		{`{foo!="bar"}`, as(), false},

		{`{foo=100}`, as(i("foo", 100)), true},
		{`{foo=100}`, as(i("foo", 200)), false},
		{`{foo=100}`, as(), false},

		{`{foo=~"^bar"}`, as(s("foo", "bar/foo")), true},
		{`{foo=~"^bar"}`, as(s("foo", "foo/bar")), false},
		{`{foo=~"^bar"}`, as(), false},

		{`{foo!~"^bar"}`, as(s("foo", "bar/foo")), false},
		{`{foo!~"^bar"}`, as(s("foo", "/foo/bar")), true},
		{`{foo!~"^bar"}`, as(), false},

		{`{http.url=~"^/admin"}`, as(s("http.url", "/admin/foo/bar")), true},
		{`{http.url=~"^/admin"}`, as(s("http.url", "/foo/bar")), false},
		{`{http.status_code!="200"}`, as(s("http.status_code", "501")), true},
		{`{http.status_code!="200"}`, as(s("http.status_code", "200")), false},

		{`{http.url =~ "^/admin", http.status_code != 200}`,
			as(s("http.url", "/admin/foo/bar"), i("http.status_code", 501)), true},
		{`{http.url=~"^/admin", http.status_code!="200"}`,
			as(s("http.url", "/admin/foo/bar"), i("http.status_code", 200)), false},
	} {
		t.Run(tc.matcher, func(t *testing.T) {
			matcher, err := Parse(tc.matcher)
			require.Nil(t, err)
			require.Equal(t, tc.match, matcher.Match(tc.annotation))
		})
	}
}
