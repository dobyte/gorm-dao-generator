package template

const ExternalTemplate = `
package ${VarDaoPackageName}

import (
	"${VarDaoPackagePath}/internal"
	"gorm.io/gorm"
)

type (
	${VarDaoPrefixName}Columns = internal.${VarDaoPrefixName}Columns
	${VarDaoPrefixName}OrderBy = internal.${VarDaoPrefixName}OrderBy
	${VarDaoPrefixName}FilterFunc = internal.${VarDaoPrefixName}FilterFunc
	${VarDaoPrefixName}UpdateFunc = internal.${VarDaoPrefixName}UpdateFunc
	${VarDaoPrefixName}ColumnFunc = internal.${VarDaoPrefixName}ColumnFunc
	${VarDaoPrefixName}OrderFunc = internal.${VarDaoPrefixName}OrderFunc
)

type ${VarDaoClassName} struct {
	*internal.${VarDaoClassName}
}

func New${VarDaoClassName}(db *gorm.DB) *${VarDaoClassName} {
	return &${VarDaoClassName}{${VarDaoClassName}: internal.New${VarDaoClassName}(db)}
}
`

const InternalTemplate = `
// --------------------------------------------------------------------------------------------
// The following code is automatically generated by the mongo-dao-generator tool.
// Please do not modify this code manually to avoid being overwritten in the next generation. 
// For more tool details, please click the link to view https://github.com/dobyte/gorm-dao-generator
// --------------------------------------------------------------------------------------------

package internal

import (
	${VarPackages}
)

type ${VarDaoPrefixName}OrderBy struct {
	Column string	
	Order  string
}

type ${VarDaoPrefixName}FilterFunc func(cols *${VarDaoPrefixName}Columns) interface{}
type ${VarDaoPrefixName}UpdateFunc func(cols *${VarDaoPrefixName}Columns) interface{}
type ${VarDaoPrefixName}ColumnFunc func(cols *${VarDaoPrefixName}Columns) []string
type ${VarDaoPrefixName}OrderFunc func(cols *${VarDaoPrefixName}Columns) []${VarDaoPrefixName}OrderBy

type ${VarDaoClassName} struct {
	model     *${VarModelPackageName}.${VarModelClassName}
	Columns   *${VarDaoPrefixName}Columns
	Database  *gorm.DB
	TableName string
}

type ${VarDaoPrefixName}Columns struct {
	${VarModelColumnsDefine}
}

var ${VarDaoVariableName}Columns = &${VarDaoPrefixName}Columns{
	${VarModelColumnsInstance}
}

func New${VarDaoClassName}(db *gorm.DB) *${VarDaoClassName} {
	dao := &${VarDaoClassName}{}
	dao.model = &${VarModelPackageName}.${VarModelClassName}{}
	dao.Columns = ${VarDaoVariableName}Columns
	dao.TableName = "${VarTableName}"
	dao.Database = db

	return dao
}

// New create a new instance and return
func (dao *${VarDaoClassName}) New(tx *gorm.DB) *${VarDaoClassName} {
	d := &${VarDaoClassName}{}
	d.model = dao.model
	d.Columns = dao.Columns
	d.TableName = dao.TableName
	d.Database = tx

	return d
}

// Table create a new table db instance
func (dao *GameRoom) Table() *gorm.DB {
	return dao.Database.Model(dao.model).Table(dao.TableName)
}

// Insert executes an insert command to insert multiple documents into the collection.
func (dao *${VarDaoClassName}) Insert(ctx context.Context, models ...*${VarModelPackageName}.${VarModelClassName}) (int64, error) {
	if len(models) == 0 {
		return 0, errors.New("models is empty")
	}

	var rst *gorm.DB

	if len(models) == 1 {
		rst = dao.Table().WithContext(ctx).Create(models[0])
	} else {
		rst = dao.Table().WithContext(ctx).Create(models)
	}

	return rst.RowsAffected, rst.Error
}

// Delete executes a delete command to delete at most one document from the collection.
func (dao *${VarDaoClassName}) Delete(ctx context.Context, filterFunc ...${VarDaoPrefixName}FilterFunc) (int64, error) {
	db := dao.Table().WithContext(ctx)

	if len(filterFunc) > 0 && filterFunc[0] != nil {
		db = db.Where(filterFunc[0](dao.Columns))
	}

	rst := db.Delete(&${VarModelPackageName}.${VarModelClassName}{})

	return rst.RowsAffected, rst.Error
}

// Update executes an update command to update documents in the collection.
func (dao *${VarDaoClassName}) Update(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, updateFunc ${VarDaoPrefixName}UpdateFunc, columnFunc ...${VarDaoPrefixName}ColumnFunc) (int64, error) {
	db := dao.Table().WithContext(ctx)

	if filterFunc != nil {
		db = db.Where(filterFunc(dao.Columns))
	}

	if len(columnFunc) > 0 && columnFunc[0] != nil {
		db = db.Select(columnFunc[0](dao.Columns))
	}

	if updateFunc != nil {
		rst := db.Updates(updateFunc(dao.Columns))

		return rst.RowsAffected, rst.Error
	}

	return 0, nil
}

// Count returns the number of documents in the collection.
func (dao *${VarDaoClassName}) Count(ctx context.Context, filterFunc ...${VarDaoPrefixName}FilterFunc) (count int64, err error) {
    db := dao.Table().WithContext(ctx)

	if len(filterFunc) > 0 && filterFunc[0] != nil {
		db = db.Where(filterFunc[0](dao.Columns))
	}

	err = db.Count(&count).Error

	return
}

// Sum returns the sum of the given field.
func (dao *${VarDaoClassName}) Sum(ctx context.Context, columnFunc ${VarDaoPrefixName}ColumnFunc, filterFunc ...${VarDaoPrefixName}FilterFunc) (sums []float64, err error) {
	columns := columnFunc(dao.Columns)
	if len(columns) == 0 {
		return
	}

	fields := make([]string, len(columns))
	for i, column := range columns {
		fields[i] = fmt.Sprintf("COALESCE(SUM(%s), 0) as ${SymbolBacktick}sum_%d${SymbolBacktick}", column, i)
	}

	db := dao.Table().WithContext(ctx).Select(strings.Join(fields, ","))

	if len(filterFunc) > 0 && filterFunc[0] != nil {
		db = db.Where(filterFunc[0](dao.Columns))
	}

	rst := make(map[string]interface{}, len(columns))

	if err = db.Scan(&rst).Error; err != nil {
		return
	}

	for i := range columns {
		val, _ := rst[fmt.Sprintf("sum_%d", i)]
		sum, _ := strconv.ParseFloat(val.(string), 64)
		sums = append(sums, sum)
	}

	return
}

// Avg returns the avg of the given field.
func (dao *${VarDaoClassName}) Avg(ctx context.Context, columnFunc ${VarDaoPrefixName}ColumnFunc, filterFunc ...${VarDaoPrefixName}FilterFunc) (avgs []float64, err error) {
	columns := columnFunc(dao.Columns)
	if len(columns) == 0 {
		return
	}

	fields := make([]string, len(columns))
	for i, column := range columns {
		fields[i] = fmt.Sprintf("COALESCE(AVG(%s), 0) as ${SymbolBacktick}avg_%d${SymbolBacktick}", column, i)
	}

	db := dao.Table().WithContext(ctx).Select(strings.Join(fields, ","))

	if len(filterFunc) > 0 && filterFunc[0] != nil {
		db = db.Where(filterFunc[0](dao.Columns))
	}

	rst := make(map[string]interface{}, len(columns))

	if err = db.Scan(&rst).Error; err != nil {
		return
	}

	for i := range columns {
		val, _ := rst[fmt.Sprintf("avg_%d", i)]
		avg, _ := strconv.ParseFloat(val.(string), 64)
		avgs = append(avgs, avg)
	}

	return
}

// First executes a first command and returns a model for one record in the table.
func (dao *${VarDaoClassName}) First(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, columnFunc ...${VarDaoPrefixName}ColumnFunc) (*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		model = &${VarModelPackageName}.${VarModelClassName}{}
		db    = dao.Table().WithContext(ctx)
	)

	if filterFunc != nil {
		db = db.Where(filterFunc(dao.Columns))
	}

	if len(columnFunc) > 0 && columnFunc[0] != nil {
		columns := columnFunc[0](dao.Columns)

		if len(columns) > 0 {
			db = db.Select(columns)
		}
	}

	if rst := db.First(model); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, rst.Error
	}

	return model, nil
}

// Last executes a last command and returns a model for one record in the table.
func (dao *${VarDaoClassName}) Last(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, columnFunc ...${VarDaoPrefixName}ColumnFunc) (*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		model = &${VarModelPackageName}.${VarModelClassName}{}
		db    = dao.Table().WithContext(ctx)
	)

	if filterFunc != nil {
		db = db.Where(filterFunc(dao.Columns))
	}

	if len(columnFunc) > 0 && columnFunc[0] != nil {
		columns := columnFunc[0](dao.Columns)

		if len(columns) > 0 {
			db = db.Select(columns)
		}
	}

	if rst := db.Last(model); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, rst.Error
	}

	return model, nil
}

// FindOne executes a take command and returns a model for one record in the table.
func (dao *${VarDaoClassName}) FindOne(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, columnFunc ...${VarDaoPrefixName}ColumnFunc) (*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		model = &${VarModelPackageName}.${VarModelClassName}{}
		db    = dao.Table().WithContext(ctx)
	)

	if filterFunc != nil {
		db = db.Where(filterFunc(dao.Columns))
	}

	if len(columnFunc) > 0 && columnFunc[0] != nil {
		columns := columnFunc[0](dao.Columns)

		if len(columns) > 0 {
			db = db.Select(columns)
		}
	}

	if rst := db.Take(model); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, rst.Error
	}

	return model, nil
}

// FindMany executes a find command and returns many models the matching documents in the collection.
func (dao *${VarDaoClassName}) FindMany(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, columnFunc ${VarDaoPrefixName}ColumnFunc, orderFunc ${VarDaoPrefixName}OrderFunc, limitAndOffset ...int) ([]*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		models = make([]*${VarModelPackageName}.${VarModelClassName}, 0)
		db     = dao.Table().WithContext(ctx)
	)

	if filterFunc != nil {
		db = db.Where(filterFunc(dao.Columns))
	}

	if columnFunc != nil {
		columns := columnFunc(dao.Columns)

		if len(columns) > 0 {
			db = db.Select(columns)
		}
	}

	if orderFunc != nil {
		orders := orderFunc(dao.Columns)

		for _, order := range orders {
			db = db.Order(fmt.Sprintf("%s %s", order.Column, order.Order))
		}
	}

	if len(limitAndOffset) > 0 {
		db = db.Limit(limitAndOffset[0])
	}

	if len(limitAndOffset) > 1 {
		db = db.Offset(limitAndOffset[1])
	}

	rst := db.Scan(&models)

	if rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, rst.Error
	}

	return models, nil
}
`
