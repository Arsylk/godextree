package pkg

import (
	"fmt"

	"github.com/csnewman/dextk"
	"golang.org/x/exp/mmap"
)

type ITreeNode interface {
	GetParent() *ITreeNode
	GetKey() string
}

func (t TreeNode) GetPath() string {
	path := t.GetKey()
	parent := t.GetParent()
	for ; parent != nil; parent = (*parent).GetParent() {
		path = fmt.Sprintf("%s/%s", (*parent).GetKey(), path)
	}
	return path
}

type TreeNode struct {
	Parent *ITreeNode
	Key    string
}

func (t TreeNode) GetParent() *ITreeNode {
	return t.Parent
}

func (t TreeNode) GetKey() string {
	return t.Key
}

type DexTree struct {
	ITreeNode
	Files []DexFile
}

type DexFile struct {
	ITreeNode
	Name    string
	Classes []DexClass
}

type DexClass struct {
	ITreeNode
	Class   dextk.ClassNode
	Fields  []DexField
	Methods []DexMethod
	Strings []string
}

type DexField struct {
	ITreeNode
	Field  dextk.FieldNode
	Static bool
}

type DexMethod struct {
	ITreeNode
	Method  dextk.MethodNode
	Direct  bool
	Code    dextk.CodeNode
	Strings []string
}

func Disassembly(paths ...string) (*DexTree, error) {
	tree := DexTree{
		ITreeNode: TreeNode{
			Parent: nil,
			Key:    "<root>",
		},
		Files: []DexFile{},
	}

	// process each file
	for _, path := range paths {
		f, err := mmap.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		dex, err := dextk.Read(f)
		if err != nil {
			return nil, err
		}

		dexFile := &DexFile{
			ITreeNode: TreeNode{
				Parent: &tree.ITreeNode,
				Key:    path,
			},
			Name:    path,
			Classes: []DexClass{},
		}
		(&tree).Files = append(tree.Files, *dexFile)
		dexFile = &tree.Files[len(tree.Files)-1]

		ci := dex.ClassIter()
		for ci.HasNext() {
			node, err := ci.Next()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			dexClass := decodeClass(dex, dexFile, node)
			dexFile.Classes = append(dexFile.Classes, dexClass)
		}
	}
	return &tree, nil
}

func decodeClass(dex *dextk.Reader, file *DexFile, class dextk.ClassNode) DexClass {
	dexClass := &DexClass{
		ITreeNode: TreeNode{
			Parent: &file.ITreeNode,
			Key:    class.Name.String(),
		},
		Class:   class,
		Fields:  []DexField{},
		Methods: []DexMethod{},
		Strings: []string{},
	}

	for _, method := range class.DirectMethods {
		if code, _ := decodeMethod(dex, method); code != nil {
			node := TreeNode{
				Parent: &dexClass.ITreeNode,
				Key:    method.Name.String(),
			}
			dexMth := DexMethod{
				ITreeNode: node,
				Method:    method,
				Direct:    true,
				Code:      *code,
				Strings:   decodeStrings(code.Ops),
			}
			dexClass.Methods = append(dexClass.Methods, dexMth)
		}
	}

	for _, method := range class.VirtualMethods {
		if code, _ := decodeMethod(dex, method); code != nil {
			node := TreeNode{
				Parent: &dexClass.ITreeNode,
				Key:    method.Name.String(),
			}
			dexMth := DexMethod{
				ITreeNode: node,
				Method:    method,
				Direct:    false,
				Code:      *code,
				Strings:   decodeStrings(code.Ops),
			}
			dexClass.Methods = append(dexClass.Methods, dexMth)
		}
	}

	for _, field := range class.StaticFields {
		node := TreeNode{
			Parent: &dexClass.ITreeNode,
			Key:    field.Name.String(),
		}
		dexField := DexField{node, field, true}
		dexClass.Fields = append(dexClass.Fields, dexField)

	}

	for _, field := range class.InstanceFields {
		node := TreeNode{
			Parent: &dexClass.ITreeNode,
			Key:    field.Name.String(),
		}
		dexField := DexField{node, field, false}
		dexClass.Fields = append(dexClass.Fields, dexField)

	}

	return *dexClass
}

func decodeMethod(dex *dextk.Reader, method dextk.MethodNode) (*dextk.CodeNode, error) {
	if method.CodeOff == 0 {
		return nil, nil
	}

	code, err := dex.ReadCodeAndParse(method.CodeOff)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &code, nil
}

func decodeStrings(ops []dextk.OpNode) []string {
	strs := []string{}
	for _, op := range ops {
		str := ""
		if constOp, ok := op.(dextk.ConstOpNode); ok {
			switch constOp.RawOp().(type) {
			case dextk.OpConstString:
				str = fmt.Sprintf("%s", constOp.Value)
			case dextk.OpConstStringJumbo:
				str = fmt.Sprintf("%s", constOp.Value)
			}
		}
		if len(str) > 0 {
			strs = append(strs, str)
		}
	}
	return strs
}
