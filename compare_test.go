package main

import "testing"

func BenchmarkHandle(b *testing.B) {
	p := &P{}
	ba := BusA{func(i interface{}) error {
		return nil
	}}
	for n := 0; n < b.N; n++ {
		ba.handle(p)
	}
}

func BenchmarkHandleInt(b *testing.B) {
	p := &P{}
	ba := BusB{}
	for n := 0; n < b.N; n++ {
		ba.handle(p)
	}
}
