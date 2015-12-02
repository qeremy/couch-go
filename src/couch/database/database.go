package database

import _client "./../client"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Database struct {
    Client *_client.Client
    Name   string
}

func Shutup() {
    u.Shutup()
}

func New(client *_client.Client, name string) *Database {
    return &Database{
        Client: client,
          Name: name,
    }
}

func (this *Database) Ping() bool {
    return (200 == this.Client.Head(this.Name, nil, nil).GetStatusCode())
}

func (this *Database) Info() (map[string]interface{}, error) {
    type Data map[string]interface{}
    data, err := this.Client.Get(this.Name, nil, nil).GetBodyData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for key, value := range *data.(*Data) {
        _return[key] = value
    }
    return _return, nil
}

func (this *Database) Create() bool {
    return (201 == this.Client.Put(this.Name, nil, nil, nil).GetStatusCode())
}

func (this *Database) Remove() bool {
    return (200 == this.Client.Delete(this.Name, nil, nil).GetStatusCode())
}

func (this *Database) Replicate(target string, targetCreate bool) (map[string]interface{}, error) {
    var body = map[string]interface{}{
        "source": this.Name,
        "target": target,
        "create_target": targetCreate,
    }
    type Data map[string]interface{}
    data, err := this.Client.Post("/_replicate", nil, body, nil).GetBodyData(&Data{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for key, value := range *data.(*Data) {
        if key == "history" {
            _return[key] = make(map[int]map[string]interface{})
            for i, history := range value.([]interface{}) {
                _return[key] = make([]map[string]interface{}, len(value.([]interface{})))
                for kkey, vvalue := range history.(map[string]interface{}) {
                    if _return[key].([]map[string]interface{})[i] == nil {
                        _return[key].([]map[string]interface{})[i] = make(map[string]interface{})
                    }
                    _return[key].([]map[string]interface{})[i][kkey] = vvalue
                }
            }
            continue
        }
        _return[key] = value
    }
    return _return, nil
}

/**
 * Document stuff. @tmp?
 */
type Document struct {
    Id        string
    Key       string
    Value     map[string]string
    Doc       map[string]interface{}
}
type Documents struct {
    Offset    uint
    TotalRows uint `json:"total_rows"`
    UpdateSeq uint `json:"update_seq"`
    Rows      []Document
}

func (this *Database) GetDocument(key string) (map[string]interface{}, error) {
    data, err := this.Client.Get(this.Name +"/_all_docs", map[string]interface{}{
        "include_docs": true,
        "key"         : u.Quote(u.QuoteEscape(key)),
    }, nil).GetBodyData(&Documents{})
    if err != nil {
        return nil, err
    }
    var _return = make(map[string]interface{})
    for _, doc := range data.(*Documents).Rows {
        _return["id"]    = doc.Id
        _return["key"]   = doc.Key
        _return["value"] = map[string]string{"rev": doc.Value["rev"]}
        _return["doc"]   = map[string]interface{}{}
        for key, value := range doc.Doc {
            _return["doc"].(map[string]interface{})[key] = value
        }
    }
    return _return, nil
}

func (this *Database) GetDocumentAll(query map[string]interface{}, keys []string) (map[string]interface{}, error) {
    query = u.MakeParam(query)
    if query["include_docs"] == nil {
        query["include_docs"] = true
    }
    // reusable lambda
    var _return = func(data interface{}, err error) (map[string]interface{}, error) {
        if err != nil {
            return nil, err
        }
        var _return = make(map[string]interface{})
        _return["offset"]     = data.(*Documents).Offset
        _return["total_rows"] = data.(*Documents).TotalRows
        _return["rows"]       = make([]map[string]interface{}, len(data.(*Documents).Rows))
        for i, row := range data.(*Documents).Rows {
            _return["rows"].([]map[string]interface{})[i] = map[string]interface{}{
                   "id": row.Id,
                  "key": row.Key,
                "value": map[string]string{"rev": row.Value["rev"]},
                  "doc": row.Doc,
            }
        }
        return _return, nil
    }
    if keys == nil {
        return _return(
            this.Client.Get(this.Name +"/_all_docs", query, nil).GetBodyData(&Documents{}))
    } else {
        return _return(
            this.Client.Post(this.Name +"/_all_docs", query, map[string]interface{}{
                "keys": keys}, nil).GetBodyData(&Documents{}))
    }
}
