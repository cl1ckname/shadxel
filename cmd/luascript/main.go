package main

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	// Load the Lua script
	if err := L.DoFile("factorial.lua"); err != nil {
		log.Fatalf("Failed to load Lua script: %s", err)
	}

	// Push the Lua function we want to call
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("factorial"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(7)); err != nil {
		log.Fatalf("Error calling Lua function: %s", err)
	}

	// Retrieve the result
	ret := L.Get(-1)
	L.Pop(1) // remove it from the stack

	if result, ok := ret.(lua.LNumber); ok {
		fmt.Printf("Factorial from Lua: %d\n", int(result))
	} else {
		log.Fatalf("Unexpected return type: %T", ret)
	}
}
