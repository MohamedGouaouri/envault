package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var VAULT_ADDR_DEFAULT = ""
var CLI_VERSION = "0.0.1"

func GetEnvFromWorkspaceFile() {

}

func HandleError(err error, messages ...string) {
	PrintErrorAndExit(1, err, messages...)
}

func PrintErrorAndExit(exitCode int, err error, messages ...string) {
	printError(err)

	if len(messages) > 0 {
		for _, message := range messages {
			fmt.Println(message)
		}
	}

	os.Exit(exitCode)
}

func PrintWarning(message string) {
	color.New(color.FgYellow).Fprintf(os.Stderr, "Warning: %v \n", message)
}

func PrintSuccessMessage(message string) {
	color.New(color.FgGreen).Println(message)
}

func PrintErrorMessageAndExit(messages ...string) {
	if len(messages) > 0 {
		for _, message := range messages {
			fmt.Fprintln(os.Stderr, message)
		}
	}

	os.Exit(1)
}

func printError(e error) {
	color.New(color.FgRed).Fprintf(os.Stderr, "error: %v\n", e)
}

func MakeUpperCase(value string) string {
	return strings.ToUpper(value)
}
