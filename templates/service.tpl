package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
import "strconv"
import "net/http"
{{range .ServiceTypes}}
// {{.TypeName}} service object
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}
// New{{.TypeName}} creates an instance of {{.TypeName}}
func New{{.TypeName}}() *{{.TypeName}} {
	return &{{.TypeName}}{ Biz:{{.Biz.PackageName}}.New{{.Biz.TypeName}}() }
}
{{if .IsSimplePK}}// Insert serves the endpoint POST [basehost]/{{.ResourceName}}
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
// Update serves the endpoint PUT [basehost]/{{.ResourceName}}/:id
func (s *{{.TypeName}}) Update(c *gin.Context) {
	id := c.Params.ByName("id")
	var v {{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	if c.BindJSON(&v) == nil {
		value,err := s.Biz.UpdateByID(id, &v)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// Delete serves the endpoint DELETE [basehost]/{{.ResourceName}}/:id
func (s *{{.TypeName}}) Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	if id != "" {
		value,err := s.Biz.DeleteByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// Find serves the endpoint GET [basehost]/{{.ResourceName}}/:id
func (s *{{.TypeName}}) Find(c *gin.Context) {
	id := c.Params.ByName("id")
	if id != "" {
		value,err := s.Biz.FindByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, value)
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid input format.")
	}
}
// List serves the endpoint GET [basehost]/{{.ResourceName}}?take=[value]&skip=[value]&count=[value]
func (s *{{.TypeName}}) List(c *gin.Context) {
    take,_ := strconv.Atoi(c.DefaultQuery("take", "10"))
    skip,_ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	count := c.DefaultQuery("count", "none")
	if count == "none" {
		var list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	    list, err := s.Biz.List(take, skip)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, list)
		}
	} else {
		var list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
	    list, err := s.Biz.List(take, skip)
		cnt, err := s.Biz.Count()
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"list": list, "count":cnt})
		}
	}
}
{{else}}{{if .Biz.IsReadOnly}}{{else}}// Insert serves the endpoint POST [basehost]/{{.ResourceName}}
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
// Update serves the endpoint PUT [basehost]/{{.ResourceName}}
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
// Delete serves the endpoint POST [basehost]/{{.ResourceName}}/delete
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
// Find serves the endpoint POST [basehost]/{{.ResourceName}}/find
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
// List serves the endpoint GET [basehost]/{{.ResourceName}}?take=[value]&skip=[value]
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
// Count serves the endpoint GET [basehost]/{{.ResourceName}}/count
func (s *{{.TypeName}}) Count(c *gin.Context) {
	value,err := s.Biz.Count()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, value)
	}
}
{{end}}{{end}}