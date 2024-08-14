package pkg

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

var Primitives = map[rune]string{
	'Z': "boolean",
	'B': "byte",
	'C': "char",
	'S': "short",
	'I': "int",
	'J': "long",
	'F': "float",
	'D': "double",
}
var PrimitiveKeys = maps.Keys(Primitives)

func TypeString(str string) string {
	length := len(str)
	wip := string(str)
	for rune(wip[0]) == '[' {
		wip = wip[1:]
	}
	depth := length - len(wip)

	result := string(wip)
	if len(wip) == 1 {
		r := rune(wip[0])
		if slices.Contains(PrimitiveKeys, r) {
			result = Primitives[r]
		}
	}

	if len(wip) > 0 && rune(wip[0]) == 'L' && rune(wip[len(wip)-1]) == ';' {
		classRaw := wip[1 : len(wip)-1]
		result = strings.ReplaceAll(classRaw, "/", ".")
	}

	return fmt.Sprintf("%s%s", result, strings.Repeat("[]", depth))
}
