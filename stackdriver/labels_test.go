package stackdriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabels_Clone(t *testing.T) {
	src := labels{{"foo", "bar"}, {"baz", "qux"}}

	res := src.Clone()
	assert.Equal(t, src, res)
}

func TestLabels_MarshalLogObject(t *testing.T) {
	enc := new(ObjectEncoder)
	lbl := labels{{"foo", "bar"}, {"baz", "qux"}}

	enc.On("AddString", "foo", "bar").Once()
	enc.On("AddString", "baz", "qux").Once()
	require.Nil(t, lbl.MarshalLogObject(enc))
	enc.AssertExpectations(t)
}
