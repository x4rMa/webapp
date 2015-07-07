package app

// 26.04.15 11:47
// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
    "bytes"
    html "html/template"
    "log"
    "github.com/dblokhin/typo"
)


var (
    // Кэш template'ов для единоразового чтения с диска. TODO: задействовать.
    tpls map[string]*html.Template
)

// AbstractPage - структура задающая общие поля для шаблонов всех веб-страниц
type AbstractPage map[string]interface{}

// Основные рекомендуемые поля:
// Host              string
// MetaTitle         string
// MetaDescription   string
// MetaAuthor        string
// MetaCopyright     string
// MetaKeywords      string
// MetaOgImage       string
// MetaOgTitle       string
// MetaOgDescription string
// ContactEmail      string
// ContactPhone      string

// Clone возвращает новый экземляр AbstractPage c наследованными полями и значениями
func (page AbstractPage) Clone(tplName string) AbstractPage {
    doc := make(AbstractPage)
    for k, v := range page {
        doc[k] = v
    }

    doc["__tpl"] = tplName
    return doc
}

// Compile return page formatted with template from tpls/%d.tpl
func (page AbstractPage) Compile() string {
    var data bytes.Buffer

    for k, v := range page {
        switch val := v.(type) {
            case AbstractPage: {
                page[k] = html.HTML(val.Compile())
            }
            case func()string: {
                page[k] = val()
            }
        }
    }

    // Директива загрузки модулей динамичная (ctx записан в doc["Ctx"])
    getTpl(page["__tpl"].(string)).Execute(&data, page)

    return data.String()
}

// Tpls return template Name (load from cache/fs)
func getTpl(Name string) *html.Template {

    defer func() {
        if err := recover(); err != nil {
            log.Println(err)
        }
    }()

    res, ok := tpls[Name]

    if !ok {
        res = loadTemplate(Name)
    }

    return res
}

// loadTemplate load template from tpls/%s.tpl
func loadTemplate(Name string) *html.Template {
    funcMap := html.FuncMap{
        "html": func(val string) html.HTML {
            return html.HTML(val)
        },
        "typo": func(val string) string {
            return typo.Typo(val)
        },
        // TODO: в разработке
        /*"mod": func(args ...interface{}) interface{} {
            if len(args) == 0 {
                return ""
            }

            name := args[0].(string)
            ctx := new(context.Context)

            if len(args) > 1 {
                ctx = args[1].(*context.Context)
            }

            modules := reflect.ValueOf(modules.Get())
            mod := modules.MethodByName(name)

            if (mod == reflect.Value{}) {
                return ""
            }

            inputs := make([]reflect.Value, 0)
            inputs = append(inputs, reflect.ValueOf(ctx))

            ret := mod.Call(inputs)
            return ret[0].Interface()
        },*/
    }

    return html.Must(html.New("*").Funcs(funcMap).Delims("{{%", "%}}").ParseFiles("tpls/" + Name + ".tpl"))
}

func init() {
    // init Init cache of templates
    tpls = make(map[string]*html.Template)
}