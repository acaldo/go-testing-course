package main

import "testing"

func TestAddSuccess(t *testing.T) {
	result := Add(1, 2)

	expect := 3
	if result != expect {
		t.Errorf("got %d , expect %d", result, expect)
	}
}
