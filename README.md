# AutoIt Obfuscator — Go Web API SDK

Go client for the [AutoIt Obfuscator](https://www.pelock.com/products/autoit-obfuscator) Web API. Protect AutoIt `.au3` source against analysis and decompilation.

API documentation: https://www.pelock.com/products/autoit-obfuscator/api

Endpoint: `https://www.pelock.com/api/autoit-obfuscator/v1`

## Installation

```bash
go get github.com/PELock/AutoIt-Obfuscator-Go
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	autoitobfuscator "github.com/PELock/AutoIt-Obfuscator-Go"
)

func main() {
	client := autoitobfuscator.New("YOUR-WEB-API-KEY")
	client.RenameVariables = true
	client.CryptStrings = true

	result, err := client.ObfuscateScriptSource(context.Background(), `ConsoleWrite("Hello World")`)
	if err != nil {
		panic(err)
	}
	if result.Error == autoitobfuscator.ErrorSuccess {
		fmt.Println(result.Output)
	}
}
```

Obfuscation strategy flags default to `false`. Optional zlib compression of the `source` field is enabled with `EnableCompression`.

See `examples/` for login and simple obfuscation samples.

## License

Apache-2.0. Copyright Bartosz Wójcik / PELock.
