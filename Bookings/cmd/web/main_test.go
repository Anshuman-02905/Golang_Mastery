package main

import "testing"

func TestRun(t *testing.T) {

	_, err := run()

	if err != nil {
		t.Errorf("Error in running the Run function")
	}
}
