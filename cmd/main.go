package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


type Templates struct{
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

type Contact struct{
   Name string
   Email string
}

type Count struct{
    Count int 
}

func newContact(name, email string) Contact{
  return Contact{
    Name: name,
    Email: email,
  }
}


type Contacts = []Contact

func (d *Data) hasEmail(email string) bool{
    for _, contact := range d.Contacts{
        if(contact.Email == email){
            return true
        }
    }
    return false
}

type Data struct{
    Contacts Contacts
}

func newData() Data{
    return Data{
        Contacts: []Contact{
            newContact("Salman", "chotiheight@driver.com"),
            newContact("Ashwariya", "robot@shivaji.com"),
        },
    }
}


type FormData struct{
    Values map[string]string
    Errors map[string]string
}


func newFormData() FormData{
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

func main(){
    e := echo.New()
    e.Use(middleware.Logger())

    data := newData()
    e.Renderer = newTemplate()

    e.GET("/", func (c echo.Context) error{
        return c.Render(200, "index", data)
    })
    
    e.POST("/contacts", func (c echo.Context) error{
        name := c.FormValue("name")
        email := c.FormValue("email")

        if data.hasEmail(email){
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists"
            
            return c.Render(400, "form", formData)
        }

        data.Contacts = append(data.Contacts, newContact(name, email))
        return c.Render(200, "display", data)
    })

    e.Logger.Fatal(e.Start(":42069"))
}
