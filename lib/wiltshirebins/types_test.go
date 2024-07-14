package wiltshirebins

import (
	"testing"
)

func TestErrorFill(t *testing.T) {
	var c Calendar

	c.errorFill("user_burped")

	for _, collection := range c {
		if collection.Errors["user_burped"] != 1 {
			t.Errorf("user didn't burp when they should have")
		}
	}
}
