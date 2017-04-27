package main

import (
  "testing"
	"github.com/stretchr/testify/assert"
)

func TestParseOneLine(t *testing.T) {

  // Given:
  line := "some.valid.key:test translation"

  // When:
  key, value, err := parseLine(line, ":");

  // Then:
  assert.Equal(t, "some.valid.key", key, "key should be as first part of incoming string")
  assert.Equal(t, "test translation", value, "value should be as the second part of incoming string")

  // and:
  assert.Nil(t, err)
}

func TestParseOneLineAnotherDelimenter(t *testing.T) {

  // Given:
  line := "some.valid.key=test translation"

  // When:
  key, value, err := parseLine(line, "=");

  // Then:
  assert.Equal(t, "some.valid.key", key, "key should be as first part of incoming string")
  assert.Equal(t, "test translation", value, "value should be as the second part of incoming string")

  // and:
  assert.Nil(t, err)
}

func TestParseOneLineTrimgWhiteSpaces(t *testing.T) {

  // Given:
  line := "some.valid.key:   test translation with whitespaces     "

  // When:
  key, value, err := parseLine(line, ":");

  // Then:
  assert.Equal(t, "some.valid.key", key, "key should be as first part of incoming string")
  assert.Equal(t, "test translation with whitespaces", value, "value should be as the second part of incoming string")

  // and:
  assert.Nil(t, err)
}

func TestParseOneLineWrongDelimeter(t *testing.T) {

  // Given:
  line := "some.valid.key:used wrong delimeter"

  // When:
  _, _, err := parseLine(line, "="); // <-- we use wrong delimeter

  // Then:
  assert.NotNil(t, err)
}
