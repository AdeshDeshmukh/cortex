package main

import (
	"fmt"
	"net/http"
	"time"
)

func testFunction() {
	resp, _ := http.Get("https://example.com")
	defer resp.Body.Close()
	
	time.Sleep(5 * time.Second)
	
	panic("test panic")
	
	fmt.Println("debug statement")
	
	var x interface{}
	x = 42
	fmt.Println(x)
}
