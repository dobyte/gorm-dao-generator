package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const (
	defaultModelPkgAlias     = "modelpkg"
	defaultModelVariableName = "model"
)

const (
	pkg1 = "fmt"
	pkg2 = "context"
	pkg3 = "gorm.io/gorm"
	pkg4 = "errors"
	pkg5 = "strconv"
	pkg6 = "strings"
)

type field struct {
	name      string
	column    string
	comment   string
	documents []string
}

type model struct {
	opts               *options
	fields             []*field
	imports            map[string]string
	tableName          string
	modelName          string
	modelClassName     string
	modelVariableName  string
	modelPkgPath       string
	modelPkgName       string
	daoClassName       string
	daoVariableName    string
	daoPkgPath         string
	daoPkgName         string
	daoOutputDir       string
	daoOutputFile      string
	daoPrefixName      string
	fieldNameMaxLen    int
	fieldComplexMaxLen int
}

func newModel(opts *options) *model {
	m := &model{
		opts:    opts,
		fields:  make([]*field, 0),
		imports: make(map[string]string, 7),
	}

	m.addImport(pkg1)
	m.addImport(pkg2)
	m.addImport(pkg3)
	m.addImport(pkg4)
	m.addImport(pkg5)
	m.addImport(pkg6)

	return m
}

func (m *model) setModelName(name string) {
	m.modelName = name
	m.modelClassName = toPascalCase(m.modelName)
	m.modelVariableName = toCamelCase(m.modelName)
	m.daoClassName = toPascalCase(m.modelName)
	m.daoVariableName = toCamelCase(m.modelName)
	m.daoOutputFile = fmt.Sprintf("%s.go", toFileName(m.modelName, m.opts.fileNameStyle))
	m.tableName = toUnderscoreCase(m.modelName)

	dir := strings.TrimSuffix(m.opts.daoDir, "/")

	if m.opts.subPkgEnable {
		m.daoOutputDir = dir + "/" + toPackagePath(m.modelName, m.opts.subPkgStyle)
	} else {
		m.daoOutputDir = dir
		m.daoPrefixName = toPascalCase(m.modelName)
	}
}

func (m *model) setModelPkg(name, path string) {
	m.modelPkgPath = path

	if m.opts.modelPkgAlias != "" {
		m.modelPkgName = m.opts.modelPkgAlias
		m.addImport(m.modelPkgPath, m.modelPkgName)
	} else {
		m.modelPkgName = name
		m.addImport(m.modelPkgPath)
	}

	if m.modelPkgName == defaultModelVariableName {
		m.modelPkgName = defaultModelPkgAlias
		m.addImport(m.modelPkgPath, m.modelPkgName)
	}
}

func (m *model) setDaoPkgPath(path string) {
	if m.opts.subPkgEnable {
		m.daoPkgPath = path + "/" + toPackagePath(m.modelName, m.opts.subPkgStyle)
	} else {
		m.daoPkgPath = path
	}

	m.daoPkgName = toPackageName(filepath.Base(m.daoPkgPath))
}

func (m *model) addImport(pkg string, alias ...string) {
	if len(alias) > 0 {
		m.imports[pkg] = alias[0]
	} else {
		m.imports[pkg] = ""
	}
}

func (m *model) addFields(fields ...*field) {
	for _, f := range fields {
		if l := len(f.name); l > m.fieldNameMaxLen {
			m.fieldNameMaxLen = l
		}

		if l := len(f.name) + len(f.column) + 5; l > m.fieldComplexMaxLen {
			m.fieldComplexMaxLen = l
		}
	}

	m.fields = append(m.fields, fields...)
}

func (m *model) modelColumnsDefined() (str string) {
	for i, f := range m.fields {
		str += fmt.Sprintf("\t%s%s%s %s", f.name, strings.Repeat(" ", m.fieldNameMaxLen-len(f.name)+1), "string", f.comment)
		if i != len(m.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimPrefix(str, "\t")
	return
}

func (m *model) modelColumnsInstance() (str string) {
	for i, f := range m.fields {
		s := fmt.Sprintf("%s:%s\"%s\",", f.name, strings.Repeat(" ", m.fieldNameMaxLen-len(f.name)+1), f.column)
		s += strings.Repeat(" ", m.fieldComplexMaxLen-len(s)+1) + f.comment
		str += "\t" + s
		if i != len(m.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimLeft(str, "\t")
	return
}

func (m *model) packages() (str string) {
	packages := make([]string, 0, len(m.imports))
	for pkg := range m.imports {
		packages = append(packages, pkg)
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i] < packages[j]
	})

	for _, pkg := range packages {
		if alias := m.imports[pkg]; alias != "" {
			str += fmt.Sprintf("\t%s \"%s\"\n", alias, pkg)
		} else {
			str += fmt.Sprintf("\t\"%s\"\n", pkg)
		}
	}

	str = strings.TrimPrefix(str, "\t")
	str = strings.TrimSuffix(str, "\n")
	return
}
