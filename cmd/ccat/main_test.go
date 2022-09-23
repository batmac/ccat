package main

import "testing"

func Test_main(_ *testing.T) {
	// just be sure these don't panic
	buildLine()
	Usage()
	FullUsage()
	main()
}
