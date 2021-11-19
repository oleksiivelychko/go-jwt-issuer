package service

import "testing"

func TestTokenServiceGenerateTokenPair(t *testing.T) {
	ts := TokenService{}
	ts.InitRedis()

	_, _, _, err := ts.GenerateTokenPair(1)
	if err != nil {
		t.Errorf("unable to generate pair tokens: %s", err.Error())
	}
}
