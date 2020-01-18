package cqrs

import (
	"github.com/fnxnet/go-cqrs/cqrs/dummy"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCommand struct{}

func TestExtractName(t *testing.T) {
	expected := "cqrs.testCommand"
	assert.Equal(t, expected, extractName(testCommand{}))
	assert.Equal(t, expected, extractName(&testCommand{}))
	assert.Equal(t, expected, extractName(*(&testCommand{})))
}

func TestExtractNameDifferentPackage(t *testing.T) {
	expected := "dummy.DummyCommand"
	assert.Equal(t, expected, extractName(dummy.DummyCommand{}))
	assert.Equal(t, expected, extractName(&dummy.DummyCommand{}))
	assert.Equal(t, expected, extractName(*(&dummy.DummyCommand{})))
}
