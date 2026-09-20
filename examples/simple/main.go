/******************************************************************************
 * AutoIt Obfuscator WebApi interface usage example.
 *
 * In this example we will obfuscate sample source with default options.
 *
 * Version        : v1.5.0
 * Language       : Go
 * Author         : Bartosz Wójcik
 * Web page       : https://www.pelock.com
 *
 *****************************************************************************/

package main

import (
	"context"
	"fmt"
	"os"

	autoitobfuscator "github.com/PELock/AutoIt-Obfuscator-Go"
)

func main() {
	client := autoitobfuscator.New("ABCD-ABCD-ABCD-ABCD")

	result, err := client.ObfuscateScriptSource(context.Background(), `ConsoleWrite("Hello World")`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if result.Error != autoitobfuscator.ErrorSuccess {
		fmt.Fprintf(os.Stderr, "An error occurred, error code: %d\n", result.Error)
		os.Exit(1)
	}
	fmt.Println(result.Output)
}
