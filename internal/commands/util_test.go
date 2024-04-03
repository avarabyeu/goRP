package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnderstandsYes(t *testing.T) {
	t.Parallel()
	assert.True(t, answerYes("yes"))
}

func TestUnderstandsYesUpper(t *testing.T) {
	t.Parallel()
	assert.True(t, answerYes("YES"))
}

func TestEmptyAnswer(t *testing.T) {
	t.Parallel()
	assert.False(t, answerYes(""))
}

func TestUnderstandsNo(t *testing.T) {
	t.Parallel()
	assert.False(t, answerYes("no"))
}
