package trace_test

import "github.com/wx-satellite/gopkg/trace"

func f1() {
	defer trace.Do()()
	f2()
}

func f2() {
	defer trace.Do()()
	f3()
}

func f3() {
	defer trace.Do()()
}


func ExampleDo() {
	f1()
	// Output:
	// g[00001]:    ->github.com/wx-satellite/gopkg/trace_test.f1
	// g[00001]:        ->github.com/wx-satellite/gopkg/trace_test.f2
	// g[00001]:            ->github.com/wx-satellite/gopkg/trace_test.f3
	// g[00001]:            <-github.com/wx-satellite/gopkg/trace_test.f3
	// g[00001]:        <-github.com/wx-satellite/gopkg/trace_test.f2
	// g[00001]:    <-github.com/wx-satellite/gopkg/trace_test.f1
}