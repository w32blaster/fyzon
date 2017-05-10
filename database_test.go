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

func TestContainsDefaultCode(t *testing.T)  {

  // When:
  result := isContainingDefaultLanguage("gb,ru,gr", "ru")

  // Then:
  assert.True(t, result, "Line of codes contains RU")
}

func TestContainsDefaultCodeDoesntContain(t *testing.T)  {

  // When:
  result := isContainingDefaultLanguage("gb,ru,gr", "ch")

  // Then:
  assert.False(t, result, "There is no CH language")
}

func TestContainsDefaultCodeOneLang(t *testing.T)  {

  // When:
  result := isContainingDefaultLanguage("ru", "ru")

  // Then:
  assert.True(t, result, "There is only one language RU")
}

func TestContainsDefaultCodeEmpty(t *testing.T)  {

  // When:
  result := isContainingDefaultLanguage("", "ru")

  // Then:
  assert.False(t, result, "Epty string")
}
