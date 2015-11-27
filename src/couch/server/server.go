package server

/**
 * Links
 * - http://blog.golang.org/json-and-go
 * - http://golang.org/pkg/encoding/json/#example_Unmarshal
 */

import _client   "./../client"
// import _response "./../http/response"

import u "./../util"
// @tmp
var _dump, _dumps, _dumpf = u.Dump, u.Dumps, u.Dumpf

type Server struct {
    Client *_client.Client
}

func Shutup() {}

func New(client *_client.Client) *Server {
    return &Server{
        Client: client,
    }
}

func (this *Server) Ping() bool {
    return (200 == this.Client.Head("/", nil, nil).GetStatusCode())
}

func (this *Server) Info() map[string]interface{} {
    type Data struct {
        CouchDB string
        Uuid    string
        Version string
        Vendor  map[string]string
    }
    data, err := this.Client.Get("/", nil, nil).GetData(&Data{})
    if err != nil {
        return nil
    }
    var _return = make(map[string]interface{});
    _return["couchdb"] = data.(*Data).CouchDB
    _return["uuid"]    = data.(*Data).Uuid
    _return["version"] = data.(*Data).Version
    _return["vendor"]  = map[string]string{
           "name": data.(*Data).Vendor["name"],
        "version": data.(*Data).Vendor["version"],
    }
    return _return
}

func (this *Server) Version() string {
    return u.Dig("version", this.Info()).(string)
}

func (this *Server) GetActiveTasks() map[int]map[string]interface{} {
    type Data struct {
        ChangesDone  uint   `json:"changes_done"`
        Database     string
        Pid          string
        Progress     uint
        TotalChanges uint   `json:"total_changes"`
        Type         string
        StartedOn    uint32 `json:"started_on"`
        UpdatedOn    uint32 `json:"updated_on"`
    }
    data, err := this.Client.Get("/_active_tasks", nil, nil).GetData(&[]Data{})
    if err != nil {
        panic(err)
    }
    var _return = make(map[int]map[string]interface{});
    for i, data := range *data.(*[]Data) {
        _return[i] = map[string]interface{}{
             "changes_done": data.ChangesDone,
                 "database": data.Database,
                      "pid": data.Pid,
                 "progress": data.Progress,
            "total_changes": data.TotalChanges,
                     "type": data.Type,
               "started_on": data.StartedOn,
               "updated_on": data.UpdatedOn,
        }
    }
    return _return
}
