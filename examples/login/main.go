/******************************************************************************
 * AutoIt Obfuscator WebApi interface usage example.
 *
 * In this example we will verify our activation key status.
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

	result, err := client.Login(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Demo version status - %v\n", result.Demo)
	fmt.Printf("Usage credits left - %d\n", result.CreditsLeft)
	fmt.Printf("Total usage credits - %d\n", result.CreditsTotal)
	fmt.Printf("Max. script size - %d\n", result.StringLimit)
}
