// +build go1.7

package main

func init() {
	flagBoolEnvVar(&nosubtests, "nosubtests", "GOTESTS_NOSUBTESTS", "disable generating tests using the Go 1.7 subtests feature")
	flagBoolEnvVar(&parallel, "parallel", "GOTESTS_PARALLEL", "enable generating parallel subtests using the Go 1.7 feature")
	flagBoolEnvVar(&named, "named", "GOTESTS_NAMED", "switch table tests from using slice to map (with test name for the key)")
}
