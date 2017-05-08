package main

import (
  "testing"
	"github.com/stretchr/testify/assert"
)

func TestAsArrayString(t *testing.T) {

  // Given:
  langs := []ProjectLanguage{ProjectLanguage{CountryCode: "ru", IsDefault: false}, ProjectLanguage{CountryCode: "en", IsDefault: false}}

  // When:
  arr := *asStringArray(&langs)

  // Then:
  assert.Equal(t, "ru", arr[0])
  assert.Equal(t, "en", arr[1])
}
