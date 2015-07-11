// 26.04.15 11:40
// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

package webapp

type SQL interface {
    Query(sql string, args ...interface{}) ([]map[string]string, error)
    Result(sql string, args ...interface{}) (string, error)
    Exec(sql string, args ...interface{}) (interface{}, error)
    ExecId(sql string, args ...interface{}) (int64, error)
    Start()
    Rollback()
    Commit()
}