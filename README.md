# go-llvm-sample

[Go言語で利用するLLVM入門](https://postd.cc/an-introduction-to-llvm-in-go/)をもとにGo言語でLLVMに触ってみる。

## LLVMのインストール、build

`import ("llvm.org/llvm/bindings/go/llvm")`していたので、`go get`した。エラーが出たのでbuildまではしてくれない。
`llvm.org/llvm/bindings/go/README.txt`には`go get`するときは`-d`フラグを付けることを推奨されるているので、そちらのほうがいいのかも（試していない）。

buildは、`build.sh`でできそうなので以下で進めた。`build.sh`を実行すると50分弱くらいかかった。

```sh
❯ go get -u llvm.org/llvm/bindings/go/llvm
❯ cd $GOPATH/llvm.org/llvm/bindings/go
❯ ./build.sh
❯ go install
```

## サンプルコード

記事中のサンプルコードをそのまま動かしてみる。

```go
package main

import (
	"fmt"

	"llvm.org/llvm/bindings/go/llvm"
)

func main() {
	// setup our builder and module
	builder := llvm.NewBuilder()
	mod := llvm.NewModule("my_module")

	// create our function prologue
	main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	builder.SetInsertPoint(block, block.FirstInstruction())

	// int a = 32
	a := builder.CreateAlloca(llvm.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(llvm.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)

	// return a + b
	bVal := builder.CreateLoad(b, "b_val")
	aVal := builder.CreateLoad(a, "a_val")
	result := builder.CreateAdd(aVal, bVal, "ab_val")
	builder.CreateRet(result)

	// verify it's all good
	if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
		fmt.Println(ok.Error())
	}
	mod.Dump()

	// create our exe engine
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		fmt.Println(err.Error())
	}

	// run the function!
	funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})
	fmt.Printf("%d\n", funcResult.Int(false))
}
```

```sh
❯ go run main.go
; ModuleID = 'my_module'
source_filename = "my_module"

define i32 @main() {
entry:
  %a = alloca i32
  store i32 32, i32* %a
  %b = alloca i32
  store i32 16, i32* %b
  %b_val = load i32, i32* %b
  %a_val = load i32, i32* %a
  %ab_val = add i32 %a_val, %b_val
  ret i32 %ab_val
}
48
```

(最初にbuildしたら3分以上掛かった。。。)


## llvmのインストール

```sh
❯ brew install llvm
```

## References
* [Go言語で利用するLLVM入門](https://postd.cc/an-introduction-to-llvm-in-go/)
* [LLVMのGo bindingでhello world](http://yukirinmk2.hatenablog.com/entry/2016/03/26/004135)