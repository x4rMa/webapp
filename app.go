package webapp

import (
    "log"
    "github.com/dblokhin/webapp/sql"
    "github.com/dblokhin/webapp/context"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
)

type singleRoute map[string]Controller
type MapRoutes []singleRoute

type Application struct {
    Doc AbstractPage
    Config Config
    DB SQL

    routes MapRoutes
}

type ContextApplication struct {
    Ctx *context.Context
    Doc AbstractPage
    Config Config
    DB SQL
}

// Routes устанавливает обработчики запросов в соответсвии с URL'ами
func (app *Application) Routes(r MapRoutes) {
    app.routes = r
}

func (app *Application) Run() {
    r := mux.NewRouter()
    r.StrictSlash(true)

    for _, val := range app.routes {
        for url, ctrl := range val {
            r.HandleFunc(url, obs(ctrl))
        }
    }

    http.Handle("/", r)
    listen := fmt.Sprintf("%s:%d", app.Config.Net.Listen_host, app.Config.Net.Listen_port)

    log.Println("Server is started on", listen)
    if err := http.ListenAndServe(listen, nil); err != nil {
        log.Println(err)
    }
}

var app *Application

// GetApplication возвращает экземпляр Application
func GetApplication() *Application {
    if app == nil {
        app = new(Application)

        // Init
        app.Config = loadConfig("config.ini")
        log.Println("Application config is loaded")

        app.Doc = make(AbstractPage)
        //app.routes = make(MapRoutes)

        // Настройка значений глобальных полей документа
        app.Doc["Host"] = app.Config.Site.Host
        app.Doc["MetaTitle"] = app.Config.Site.Title
        app.Doc["MetaDescription"] = app.Config.Site.Description
        app.Doc["MetaAuthor"] = app.Config.Site.Author
        app.Doc["MetaCopyright"] = app.Config.Site.Copyright
        app.Doc["MetaKeywords"] = app.Config.Site.Keywords
        app.Doc["ContactEmail"] = app.Config.Site.Email
        app.Doc["ContactPhone"] = app.Config.Site.Phone
        app.Doc["UploadPath"] = app.Config.Site.UploadPath

        if !app.Config.Db.Disable {
            db, err := qmysql.New(app.Config.Db.Driver, app.Config.Db.Datasource)
            if err != nil {
                panic("Database error: " + err.Error())
            }

            app.DB = db
            log.Println("Database is connected")
        }
    }

    return app
}