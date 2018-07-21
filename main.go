package main

import (
	"fmt"

	"llvm.org/llvm/bindings/go/llvm"
)

func main() {
	// module := llvm.NewModule("main")
	// foo := llvm.Int16Type()
	// bar := llvm.Int32Type()
	//
	// fupa := llvm.IntType(32)
	//
	// ages := llvm.ArrayType(llvm.Int32Type(), 16)

	// foo := llvm.ConstInt(llvm.Int32Type(), 666, false) // false means negative sign
	// bar := llvm.ConstFloat(llvm.FloatType(), 32.5)

	// a := llvm.ConstInt(llvm.Int32Type(), 12, false)
	// b := llvm.ConstInt(llvm.Int32Type(), 24, false)
	// c := llvm.ConstAdd(a, b)

	// llvm.AddBasicBlock(context, "entry")

	// builder := llvm.NewBuilder() // create a function "main" and a block "entry"
	// foo := builder.CreateAlloca(llvm.Int32Type(), "foo")
	// builder.CreateStore(foo, llvm.ConstInt(llvm.Int32Type(), 12, false))

	builder := llvm.NewBuilder() // create a function "main" and a block "entry"
	mod := llvm.NewModule("my_module")
	main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	// mainFunc := mod.NamedFunction("main")

	builder.SetInsertPoint(block, block.FirstInstruction())

	a := builder.CreateAlloca(llvm.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)
	b := builder.CreateAlloca(llvm.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)

	aVal := builder.CreateLoad(a, "a_val")
	bVal := builder.CreateLoad(b, "b_val")
	result := builder.CreateAdd(aVal, bVal, "ab_value")
	builder.CreateRet(result)

	if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
		fmt.Println(ok.Error())
	}
	mod.Dump()

	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		fmt.Println(err.Error())
	}

	funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})
	fmt.Printf("%d\n", funcResult.Int(false))
}
