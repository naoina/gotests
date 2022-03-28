package main

import (
	"flag"
	"os"
	"strconv"
)

type boolValue bool

func (b *boolValue) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *boolValue) Set(v string) error {
	val, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	*b = boolValue(val)
	return nil
}

func (b *boolValue) IsBoolFlag() bool {
	return true
}

func flagBoolEnv(name, env, usage string) *bool {
	var b bool
	flagBoolEnvVar(&b, name, env, usage)
	return &b
}

func flagBoolEnvVar(p *bool, name, env, usage string) {
	s := os.Getenv(env)
	val := (*boolValue)(p)
	// Set fails, but env has any value, then set to true.
	if err := val.Set(s); err != nil {
		*val = s != ""
	}
	flag.Var(val, name, usage)
}
