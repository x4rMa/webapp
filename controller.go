// 26.04.15 11:41
// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

package app

import (
    "net/http"
    "github.com/dblokhin/webapp/context"
)

type Controller interface {

    GET(app *ContextApplication)
    POST(app *ContextApplication)
    PUT(app *ContextApplication)
    DELETE(app *ContextApplication)
    PATCH(app *ContextApplication)
    OPTIONS(app *ContextApplication)
    HEAD(app *ContextApplication)
    TRACE(app *ContextApplication)
    CONNECT(app *ContextApplication)
}

// obs инициализирует контекст для заданного клиента и вызывает контроллер
func obs(handler Controller) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, req *http.Request) {


        ctx := context.New(w, req)
        app := GetApplication()
        doc := app.Doc.Clone("")
        doc["Ctx"] = ctx
        doc["User"] = ctx.User()

        contextApp := &ContextApplication{ctx, doc, app.Config, app.DB}

        switch ctx.Input.Method() {
            case "GET":     handler.GET(contextApp);
            case "POST":    handler.POST(contextApp);
            case "PUT":     handler.PUT(contextApp);
            case "DELETE":  handler.DELETE(contextApp);
            case "PATCH":   handler.PATCH(contextApp);
            case "OPTIONS": handler.OPTIONS(contextApp);
            case "HEAD":    handler.HEAD(contextApp);
            case "TRACE":   handler.TRACE(contextApp);
            case "CONNECT": handler.CONNECT(contextApp);

            default: http.Error(ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
        }
    }
}

// HTTPController объект для встраивания в контроллеры, содержащие стандартные методы для контроллера
// Задача контроллеров переписать необходимые методы.
type HTTPController struct {}

func (h HTTPController) GET(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) POST(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) PUT(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) DELETE(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) PATCH(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) OPTIONS(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) HEAD(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) TRACE(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}

func (h HTTPController) CONNECT(app *ContextApplication) {
    http.Error(app.Ctx.Response(), "Method not allowed", http.StatusMethodNotAllowed)
}
