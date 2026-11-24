//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package codegen

type PkgInfo struct {
	PkgName string
	PkgPath string
}

type TypeInfo struct {
	Name     string
	FullName string
	PkgInfo
}
