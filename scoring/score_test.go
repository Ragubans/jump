package scoring

import (
	"fmt"
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestScoreCalculate(t *testing.T) {
	score1 := Score{2, Now}
	score2 := Score{4, Now}

	assert.True(t, inDelta(1.38, score2.Calculate()-score1.Calculate()))
}

func TestScoreRelevance(t *testing.T) {
	score := Score{2, Now}

	assert.True(t, inDelta(2, score.Relevance()))
}

func TestScoreUpdate(t *testing.T) {
	score := Score{2, Now}
	score.Update(1)

	assert.Equal(t, Now, score.Age)
	assert.Equal(t, 3, score.Weight)
}

func TestScoreString(t *testing.T) {
	score := Score{2, Now}

	assert.NotEqual(t, "", score.String())
}

func TestNewScore(t *testing.T) {
	score := NewScore()
	str := fmt.Sprintf("{1 %s}", score.Age)

	assert.Equal(t, str, score.String())
}

func inDelta(delta, expr float64) bool {
	return delta-0.01 < expr && expr < delta+0.01
}
