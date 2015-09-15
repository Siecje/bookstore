package main

import (
	"reflect"
	"testing"
)

func TestFetch(t *testing.T) {
	url := []string{"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/anth.html"}
	if reflect.TypeOf(Fetch(url)) != reflect.TypeOf([]*HTTPResponse{}) {
		t.Errorf("Expected true, got false")
	}
}
