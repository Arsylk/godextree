// cmd.main
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Arsylk/godextree/pkg"
)

func main() {
	args := os.Args[1:]

	tree, _ := pkg.Disassembly(args...)
	for _, file := range tree.Files {
		fmt.Printf("file %s\n", pkg.Color(file.GetKey(), pkg.Cyan))
		for _, class := range file.Classes {
			fmt.Printf("  id:c0x%08x \x1b[36m%s\x1b[0m\n", class.Class.Id, pkg.TypeString(class.GetKey()))
			for _, method := range class.Methods {
				PrintMethod(method)
				for _, str := range method.Strings {
					PrintMethodString(str)
				}
			}
			for _, field := range class.Fields {
				fmt.Printf("    id:f0x%08x %s %s\n", field.Field.Id, pkg.Color(pkg.TypeString(field.Field.Type.String()), pkg.White), pkg.Color(field.GetKey(), pkg.Green))
			}
		}
	}
}

func PrintMethod(method pkg.DexMethod) {
	name := pkg.Color(method.Method.Name.String(), pkg.Blue)
	args := make([]string, len(method.Method.Params))
	for i, param := range method.Method.Params {
		args[i] = pkg.TypeString(param.String())
	}
	fmt.Printf("    %s(%s)\n", name, strings.Join(args, ", "))
}

func PrintMethodString(str string) {
	stripped := strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\r", "")
	quoted := fmt.Sprintf("\"%s\"", stripped)
	fmt.Printf("      %s\n", pkg.Color(quoted, pkg.Yellow))
}
