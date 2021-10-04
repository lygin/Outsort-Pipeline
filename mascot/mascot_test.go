package mascot_test

import (
	"go-demo-1/mascot"
	"testing"
)

func TestMascot(t *testing.T) {
	if mascot.Hello() != "go gopher" {
		t.Fatal("Wrong Def :(")
	}
}
