package main

import (
	"fmt"
	"reflect"
)

var testURL = []string{"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/anth.html"}

func main() {
	fmt.Println(reflect.TypeOf(Fetch(testURL)))
	fmt.Println(reflect.TypeOf([]*HTTPResponse{}))
}
