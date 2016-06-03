package producers

import "github.com/MikeK123/dingo/model"

func ProduceDaoPackage(config *model.Configuration, schema *model.DatabaseSchema, mpkg *model.ModelPackage) (pkg *model.DaoPackage) {
	pkg = &model.DaoPackage{PackageName: "dao", BasePackage: config.BasePackage}
	pkg.AppendImport(mpkg.BasePackage + "/" + mpkg.PackageName)
	pkg.AppendImport("database/sql")
	i := 0
	for _, table := range schema.Tables {
		dao := &model.DaoType{TypeName: mpkg.ModelTypes[i].TypeName + "Dao", PackageName: "dao"}
		dao.Model = mpkg.ModelTypes[i]
		dao.Entity = table
		pkg.DaoTypes = append(pkg.DaoTypes, dao)
		if CheckAutoIncrementPK(table) {
			dao.HasAutoIncrementPK = true
		}
		i++
	}
	i = 0
	for _, view := range schema.Views {
		dao := &model.DaoType{TypeName: mpkg.ViewModelTypes[i].TypeName + "Dao", PackageName: "dao"}
		dao.Model = mpkg.ViewModelTypes[i]
		dao.View = view
		pkg.ViewDaoTypes = append(pkg.ViewDaoTypes, dao)
		i++
	}
	return pkg
}

func CheckAutoIncrementPK(table *model.Table) bool {
	if len(table.PrimaryKeys) == 1 {
		if table.PrimaryKeys[0].IsAutoIncrement == true {
			return true
		}
	}
	return false
}

func FindTable(schema *model.DatabaseSchema, tn string) (*model.Table, int) {
	for i, t := range schema.Tables {
		if t.TableName == tn {
			return t, i
		}
	}
	return nil, 0
}

func ProduceDaoFromQueryPackage(config *model.Configuration, schema *model.DatabaseSchema, mpkg *model.ModelPackage) (pkg *model.DaoMixedPackage) {
	pkg = &model.DaoMixedPackage{PackageName: "dao", BasePackage: config.BasePackage}
	pkg.AppendImport(mpkg.BasePackage + "/" + mpkg.PackageName)
	pkg.AppendImport("database/sql")
	for _, mdt := range config.MixedDaoTables {
		dao := &model.DaoMixedType{
			PackageName: "dao",
			Model:       make([]*model.ModelType, 0, 3),
			Entity:      make([]*model.Table, 0, 3),
			View:        make([]*model.View, 0, 2),
			Where:       mdt.Where,
		}
		for _, tn := range mdt.Tables {
			t, i := FindTable(schema, tn)
			m := mpkg.ModelTypes[i]
			dao.TypeName += m.TypeName
			dao.Model = append(dao.Model, m)
			dao.Entity = append(dao.Entity, t)
			if CheckAutoIncrementPK(t) {
				dao.HasAutoIncrementPK = true
			}
		}
		dao.TypeName += "Dao"
		pkg.DaoMixedTypes = append(pkg.DaoMixedTypes, dao)
	}
	// i = 0
	// for _, view := range schema.Views {
	// 	dao := &model.DaoType{TypeName: mpkg.ViewModelTypes[i].TypeName + "Dao", PackageName: "dao"}
	// 	dao.Model = mpkg.ViewModelTypes[i]
	// 	dao.View = view
	// 	pkg.ViewDaoTypes = append(pkg.ViewDaoTypes, dao)
	// 	i++
	// }
	return pkg
}
