package fake

import "testing"

func TestSha256(t *testing.T) {
	hash := Sha256()
	t.Log(string(hash))
}
