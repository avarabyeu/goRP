package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnderstandsYes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, answerYes("yes"))
}

func TestUnderstandsYesUpper(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, answerYes("YES"))
}

func TestEmptyAnswer(t *testing.T) {
	t.Parallel()
	assert.Equal(t, false, answerYes(""))
}

func TestUnderstandsNo(t *testing.T) {
	t.Parallel()
	assert.Equal(t, false, answerYes("no"))
}
