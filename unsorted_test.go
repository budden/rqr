package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trimToTheNumberOfRunes(t *testing.T) {
	assert.Equal(t, "юл", trimToTheNumberOfRunes("юла", 2))
	assert.Equal(t, "юла", trimToTheNumberOfRunes("юла", 3))
	assert.Equal(t, "щ...", trimToTheNumberOfRunes("щурёнок", 4))
	assert.Equal(t, "щур...", trimToTheNumberOfRunes("щурёнок", 6))
	assert.Equal(t, "щурёнок Й", trimToTheNumberOfRunes("щурёнок Й", 9))
	assert.Equal(t, "щурёнок Й", trimToTheNumberOfRunes("щурёнок Й", 10))
}
