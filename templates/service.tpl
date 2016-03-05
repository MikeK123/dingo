package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
import "strconv"
import "net/http"
{{range .ServiceTypes}}
// Service object for {{.TypeName}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}
//Create a {{.TypeName}}
func New{{.TypeName}}() *{{.TypeName}} {
	return &{{.TypeName}}{ Biz:{{.Biz.PackageName}}.New{{.Biz.TypeName}}() }
}
{{if .Biz.IsReadOnly}}{{else}}// Endpoint POST [basehost]/{{.ResourceName}}
func (s *{{.TypeName}}) Insert(c *gin.Context) {
	var v {{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	if c.BindJSON(&v) == nil {
		value,err := s.Biz.Insert(&v)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// Endpoint PUT [basehost]/{{.ResourceName}}
func (s *{{.TypeName}}) Update(c *gin.Context) {
	var v {{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	if c.BindJSON(&v) == nil {
		value,err := s.Biz.Update(&v)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// Endpoint POST [basehost]/{{.ResourceName}}/delete
func (s *{{.TypeName}}) Delete(c *gin.Context) {
	var v {{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	if c.BindJSON(&v) == nil {
		value,err := s.Biz.Delete(&v)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// Endpoint POST [basehost]/{{.ResourceName}}/find
func (s *{{.TypeName}}) Find(c *gin.Context) {
	var v {{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	if c.BindJSON(&v) == nil {
		value,err := s.Biz.Find(&v)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
{{end}}
// Endpoint GET [basehost]/{{.ResourceName}}?take=[value]&skip=[value]
func (s *{{.TypeName}}) List(c *gin.Context) {
    take,_ := strconv.Atoi(c.DefaultQuery("take", "10"))
    skip,_ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	var list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
    list, err := s.Biz.List(take, skip)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}
// Endpoint GET [basehost]/{{.ResourceName}}/count
func (s *{{.TypeName}}) Count(c *gin.Context) {
	value,err := s.Biz.Count()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, value)
	}
}
{{end}}