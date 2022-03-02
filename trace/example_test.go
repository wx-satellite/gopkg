package trace_test

import "github.com/wx-satellite/gopkg/trace"

func f1() {
	defer trace.Apply()()
	f2()
}

func f2() {
	defer trace.Apply()()
	f3()
}

func f3() {
	defer trace.Apply()()
}


func ExampleTrace() {
	f1()
	// Output:
	// g[00001]:    ->github.com/wx-satellite/gopkg/trace_test.f1
	// g[00001]:        ->github.com/wx-satellite/gopkg/trace_test.f2
	// g[00001]:            ->github.com/wx-satellite/gopkg/trace_test.f3
	// g[00001]:            <-github.com/wx-satellite/gopkg/trace_test.f3
	// g[00001]:        <-github.com/wx-satellite/gopkg/trace_test.f2
	// g[00001]:    <-github.com/wx-satellite/gopkg/trace_test.f1
}