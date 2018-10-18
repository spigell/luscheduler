// taken from https://github.com/vadv/zalua/blob/master/src/zalua/dsl/filepath.go. Thanks, vadv.
package dsl

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslState) dslFilepathBasename(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Base(path)))
	return 1
}

func (d *dslState) dslFilepathDir(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Dir(path)))
	return 1
}

func (d *dslState) dslFilepathExt(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Ext(path)))
	return 1
}

func (d *dslState) dslFilepathGlob(L *lua.LState) int {
	pattern := L.CheckString(1)
	files, err := filepath.Glob(pattern)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	result := L.CreateTable(len(files), 0)
	for _, file := range files {
		result.Append(lua.LString(file))
	}
	L.Push(result)
	return 1
}
