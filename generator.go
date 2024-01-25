package main

import (
	"github.com/dobyte/gorm-dao-generator/template"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	symbolBacktick = "`"
)

const (
	symbolBacktickKey = "SymbolBacktick"
)

const (
	varPackagesKey             = "VarPackages"
	varModelClassNameKey       = "VarModelClassName"
	varModelPackageNameKey     = "VarModelPackageName"
	varModelPackagePathKey     = "VarModelPackagePath"
	varModelVariableNameKey    = "VarModelVariableName"
	varModelColumnsDefineKey   = "VarModelColumnsDefine"
	varModelColumnsInstanceKey = "VarModelColumnsInstance"
	varDaoClassNameKey         = "VarDaoClassName"
	varDaoVariableNameKey      = "VarDaoVariableName"
	varDaoPackageNameKey       = "VarDaoPackageName"
	varDaoPackagePathKey       = "VarDaoPackagePath"
	varDaoPrefixNameKey        = "VarDaoPrefixName"
	varTableNameKey            = "VarTableName"
)

type options struct {
	modelDir      string
	modelPkgPath  string
	modelPkgAlias string
	modelNames    []string
	tableNames    []string
	daoDir        string
	daoPkgPath    string
	subPkgEnable  bool
	subPkgStyle   style
	fileNameStyle style
}

type generator struct {
	opts       *options
	modelNames map[string]string
}

func newGenerator(opts *options) *generator {
	modelNames := make(map[string]string, len(opts.modelNames))
	for _, names := range opts.modelNames {
		elems := strings.Split(names, ":")

		if len(elems) == 0 {
			continue
		}

		if !isExportable(elems[0]) {
			continue
		}

		if len(elems) > 1 {
			modelNames[elems[0]] = elems[1]
		} else {
			modelNames[elems[0]] = ""
		}
	}

	if len(modelNames) == 0 {
		log.Fatalf("error: %d model type names found", len(modelNames))
	}

	return &generator{
		opts:       opts,
		modelNames: modelNames,
	}
}

func (g *generator) makeDao() {
	models := g.parseModels()

	for _, m := range models {
		g.makeModelInternalDao(m)

		g.makeModelExternalDao(m)
	}
}

// generate an internal dao file based on model
func (g *generator) makeModelInternalDao(m *model) {
	replaces := make(map[string]string)
	replaces[varModelClassNameKey] = m.modelClassName
	replaces[varModelPackageNameKey] = m.modelPkgName
	replaces[varModelPackagePathKey] = m.modelPkgPath
	replaces[varModelVariableNameKey] = m.modelVariableName
	replaces[varDaoPrefixNameKey] = m.daoPrefixName
	replaces[varDaoClassNameKey] = m.daoClassName
	replaces[varDaoVariableNameKey] = m.daoVariableName
	replaces[varModelColumnsDefineKey] = m.modelColumnsDefined()
	replaces[varModelColumnsInstanceKey] = m.modelColumnsInstance()
	replaces[varPackagesKey] = m.packages()
	replaces[symbolBacktickKey] = symbolBacktick

	if tableName := g.modelNames[m.modelName]; tableName != "" {
		replaces[varTableNameKey] = tableName
	} else {
		replaces[varTableNameKey] = m.tableName
	}

	file := m.daoOutputDir + "/internal/" + m.daoOutputFile

	err := doWrite(file, template.InternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// generate an external dao file based on model
func (g *generator) makeModelExternalDao(m *model) {
	file := m.daoOutputDir + "/" + m.daoOutputFile

	_, err := os.Stat(file)
	if err != nil {
		switch {
		case os.IsNotExist(err):
		// ignore
		case os.IsExist(err):
			return
		default:
			log.Fatal(err)
		}
	} else {
		return
	}

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = m.daoClassName
	replaces[varDaoPrefixNameKey] = m.daoPrefixName
	replaces[varDaoPackageNameKey] = m.daoPkgName
	replaces[varDaoPackagePathKey] = m.daoPkgPath

	err = doWrite(file, template.ExternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// parse multiple models from the go file
func (g *generator) parseModels() []*model {
	var (
		pkg          = g.loadPackage()
		models       = make([]*model, 0, len(pkg.Syntax))
		daoPkgPath   = g.opts.daoPkgPath
		modelPkgPath = g.opts.modelPkgPath
		modelPkgName = g.opts.modelPkgAlias
	)

	if g.opts.daoPkgPath == "" && pkg.Module != nil {
		outPath, err := filepath.Abs(g.opts.daoDir)
		if err != nil {
			log.Fatal(err)
		}
		daoPkgPath = pkg.Module.Path + outPath[len(pkg.Module.Dir):]
	}

	daoPkgPath = strings.ReplaceAll(daoPkgPath, `\`, `/`)

	for _, file := range pkg.Syntax {
		if g.opts.modelPkgPath == "" && pkg.Module != nil && pkg.Fset != nil {
			filePath := filepath.Dir(pkg.Fset.Position(file.Package).Filename)
			modelPkgPath = pkg.Module.Path + filePath[len(pkg.Module.Dir):]
		}

		modelPkgPath = strings.ReplaceAll(modelPkgPath, `\`, `/`)
		modelPkgName = file.Name.Name

		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				return true
			}

			for _, s := range decl.Specs {
				spec, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}

				_, ok = g.modelNames[spec.Name.Name]
				if !ok {
					continue
				}

				st, ok := spec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				model := newModel(g.opts)
				model.setModelName(spec.Name.Name)
				model.setModelPkg(modelPkgName, modelPkgPath)
				model.setDaoPkgPath(daoPkgPath)

				for _, item := range st.Fields.List {
					name := item.Names[0].Name

					if !isExportable(name) {
						continue
					}

					field := &field{name: name, column: toUnderscoreCase(name)}

					if item.Tag != nil && len(item.Tag.Value) > 2 {
						runes := []rune(item.Tag.Value)
						if runes[0] != '`' || runes[len(runes)-1] != '`' {
							continue
						}

						tag := reflect.StructTag(runes[1 : len(runes)-1])

						if v := tag.Get("gorm"); v != "" {
							for _, item := range strings.Split(v, ";") {
								vv := strings.Split(item, ":")
								if len(vv) >= 2 && vv[0] == "column" {
									field.column = vv[1]
								}
							}
						}
					}

					if item.Doc != nil {
						field.documents = make([]string, 0, len(item.Doc.List))
						for _, doc := range item.Doc.List {
							field.documents = append(field.documents, doc.Text)
						}
					}

					if item.Comment != nil {
						field.comment = item.Comment.List[0].Text
					}

					model.addFields(field)
				}

				models = append(models, model)
			}

			return true
		})
	}

	return models
}

func (g *generator) loadPackage() *packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, g.opts.modelDir)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	return pkgs[0]
}
