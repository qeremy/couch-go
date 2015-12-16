package util

import (
    _fmt "fmt"
    _str "strings"
    _strc "strconv"
    _json "encoding/json"
    _rex "regexp"
)

func Shutup() {}

func Type(args ...interface{}) string {
    return _str.Trim(TypeReal(args[0]), " *<>{}[]")
}
func TypeReal(args ...interface{}) string {
    return _fmt.Sprintf("%T", args[0])
}

func String(input interface{}) string {
    switch input.(type) {
        case int,
             bool,
             string:
            return _fmt.Sprintf("%v", input)
        default:
            var inputType = _fmt.Sprintf("%T", input)
            if StringSearch(inputType, "u?int(\\d+)?|float(32|64)") {
                return _fmt.Sprintf("%v", input)
            }
            panic("Unsupported input type '"+ inputType +"' given!");
    }
}
func Int(input interface{}) int {
    result, err := _strc.Atoi(input.(string))
    if err != nil {
        return int(result)
    }
    return 0
}
func Number(input interface{}, inputType string) interface{} {
    number, err := _strc.Atoi(input.(string))
    if err != nil {
        return nil
    }
    switch inputType {
        // signed
        case    "int": return int(number)
        case   "int8": return int8(number)
        case  "int16": return int16(number)
        case  "int32": return int32(number)
        case  "int64": return int64(number)
        // unsigned
        case   "uint": return uint(number)
        case  "uint8": return uint8(number)
        case "uint16": return uint16(number)
        case "uint32": return uint32(number)
        case "uint64": return uint64(number)
    }
    return 0
}


func IsEmpty(input interface{}) bool {
    if input == nil || input == "" || input == 0 {
        return true
    }
    return false
}
func IsEmptySet(input interface{}, inputDefault interface{}/*, inputType string*/) interface{} {
    if IsEmpty(input) {
        input = inputDefault
        // switch inputType {
        //     case "string":
        //         input = String(inputDefault)
        //     default:
        //         panic("Unimplemeted type '"+ inputType +"' given!")
        // }
    }
    return input
}

func Dump(args ...interface{}) {
    _fmt.Println(args...)
}
func Dumps(args ...interface{}) {
    var format string
    for _, arg := range args {
        _ = arg // silence..
        format += "%+v "
    }
    _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}
func Dumpf(format string, args ...interface{}) {
    if format == "" {
        for _, arg := range args {
            _ = arg // silence..
            format += "%+v "
        }
    }
    _fmt.Printf("%s\n", _fmt.Sprintf(format, args...))
}

func Quote(input string) string {
    return _strc.Quote(input)
}
func QuoteEncode(input string) string {
    return _str.Replace(input, "\"", "%22", -1)
}

func Param(param map[string]interface{}) map[string]interface{} {
    if param == nil {
        param = make(map[string]interface{})
    }
    return param
}
func ParamList(argv ...interface{}) map[string]interface{} {
    var argc = len(argv)
    if argc % 2 != 0 {
        panic("MakeParamList() accepts equal param length (key1, value2, key2, value2)!")
    }
    var paramList = make(map[string]interface{});
    // tricky?
    for i := 1; i < argc; i += 2 {
        if key, ok := argv[i-1].(string); ok {
            paramList[key] = argv[i]
            continue
        }
        panic("Each param key must be string!");
    }
    return paramList
}

// parsers
func ParseUrl(url string) map[string]string {
    if url == "" {
        panic("No URL given!")
    }
    var result = make(map[string]string)
    var pattern = "(?:(?P<Scheme>https?)://(?P<Host>[^:/]+))?" +
                  "(?:\\:(?P<Port>\\d+))?(?P<Path>/[^?#]*)?"   +
                  "(?:\\?(?P<Query>[^#]+))?"                   +
                  "(?:\\??#(?P<Fragment>.*))?"
    re, _ := _rex.Compile(pattern)
    if re == nil {
        return result
    }
    var match = re.FindStringSubmatch(url)
    for i, name := range re.SubexpNames() {
        if i != 0 {
            result[name] = match[i]
        }
    }
    return result
}

func ParseQuery(query string) map[string]string {
    var ret = make(map[string]string)
    var tmp = _str.Split(query, "&")
    for _, tmp := range tmp {
        var tmp = _str.Split(tmp, "=")
        ret[tmp[0]] = tmp[1]
    }
    return ret
}

func ParseHeaders(headers string) map[string]string {
    var result = make(map[string]string)
    if tmps := _str.Split(headers, "\r\n"); tmps != nil {
        for _, tmp := range tmps {
            var tmp = _str.SplitN(tmp, ":", 2)
            // request | response check?
            if len(tmp) == 1 {
                // status line >> HTTP/1.0 200 OK
                result["0"] = tmp[0]
                continue
            }
            var key, value =
                _str.TrimSpace(tmp[0]),
                _str.TrimSpace(tmp[1]);
            result[key] = value
        }
    }
    return result
}

func ParseBody(in string, out interface{}) (interface{}, error) {
    // simply prevent useless unmarshal error
    if in == "" {
        in = `null`
    }
    err := _json.Unmarshal([]byte(in), &out)
    if err != nil {
        return nil, _fmt.Errorf("JSON error: %s!", err)
    }
    return out, nil
}
func UnparseBody(in interface{}) (string, error) {
    out, err := _json.Marshal(in)
    if err != nil {
        return "", _fmt.Errorf("JSON error: %s!", err)
    }
    return string(out), nil
}

// diggers
func Dig(key string, object interface{}) interface{} {
    if object == nil {
        return nil
    }
    var keys = _str.Split(key, ".")
    key = _shift(&keys)
    if len(keys) == 0 {
        // add more if needs
        switch object.(type) {
            case map[string]int:
                return object.(map[string]int)[key]
            case map[string]string:
                return object.(map[string]string)[key]
            case map[string]interface{}:
                return object.(map[string]interface{})[key]
            case []string:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return object.([]string)[key]
                }
            case []interface{}:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return object.([]interface{})[key]
                }
            case []map[string]interface{}:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return object.([]map[string]interface{})[key]
                }
            default:
                // @tmp @debug
                panic("Unimplemented type: "+ TypeReal(object))
        }
    } else {
        // @overwrite
        var keys = _str.Join(keys, ".")
        // add more if needs
        switch object.(type) {
            case map[string]int:
                return Dig(keys, object.(map[string]int)[key])
            case map[string]string:
                return Dig(keys, object.(map[string]string)[key])
            case map[string]interface{}:
                return Dig(keys, object.(map[string]interface{})[key])
            case []string:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return object.([]string)[key]
                }
            case []interface{}:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return object.([]interface{})[key]
                }
            case []map[string]interface{}:
                key, err := _strc.Atoi(key)
                if err == nil {
                    return Dig(keys, object.([]map[string]interface{})[key])
                }
            default:
                // @tmp @debug
                panic("Unimplemented type: "+ TypeReal(object))
        }
    }

    return nil
}
func DigInt(key string, object interface{}) int {
    if value := Dig(key, object); value != nil {
        switch value := value.(type) {
            case int:
                return value
            case float32:
                return int(value)
            case float64:
                return int(value)
        }
    }
    return 0
}
func DigFloat(key string, object interface{}) float64 {
    if value := Dig(key, object); value != nil {
        switch value := value.(type) {
            case float64:
                return value
            case float32:
                return float64(value)
            case int:
                return float64(value)
        }
    }
    return 0.0
}
func DigString(key string, object interface{}) string {
    if value := Dig(key, object); value != nil {
        return value.(string)
    }
    return ""
}
func DigBool(key string, object interface{}) bool {
    if value := Dig(key, object); value != nil {
        return true
    }
    return false
}
func DigMap(key string, object interface{}) map[string]interface{} {
    return Dig(key, object).(map[string]interface{})
}
func DigMapList(key string, object interface{}) []map[string]interface{} {
    return Dig(key, object).([]map[string]interface{})
}

func _shift(slice *[]string) string {
    var value = (*slice)[0]
    *slice = (*slice)[1 : len(*slice)]
    return value
}

func Join(sep string, args ...interface{}) string {
    var result []string
    for _, arg := range args {
        switch arg.(type) {
            case nil:
                // pass
            case string:
                result = append(result, arg.(string))
            default:
                panic("Only string args accepted!")
        }
    }
    return _str.Join(result, sep)
}

func StringSearch(input, search string) bool {
    re, _ := _rex.Compile(search)
    if re == nil {
        return false
    }
    return "" != re.FindString(input)
}

// shortcut maps
func Map() map[string]interface{} {
    return make(map[string]interface{})
}
func MapInt() map[int]string {
    return make(map[int]string)
}
func MapString() map[string]string {
    return make(map[string]string)
}
func MapMapInt() map[int]map[string]interface{} {
    return make(map[int]map[string]interface{})
}
func MapMapString() map[string]map[string]interface{} {
    return make(map[string]map[string]interface{})
}

func _length(length interface{}) int {
    switch length.(type) {
        case int:
            return length.(int)
        case []int:
            return len(length.([]int))
        case []string:
            return len(length.([]string))
        case []interface{}:
            return len(length.([]interface{}))
        // case:
            // @todo add more cases if needs
    }
    return -1
}

func MapList(length interface{}) []map[string]interface{} {
    len := _length(length)
    if len != -1 {
        return make([]map[string]interface{}, len)
    }
    return []map[string]interface{}{}
}
func MapListInt(length interface{}) []map[int]string {
    len := _length(length)
    if len != -1 {
        return make([]map[int]string, len)
    }
    return []map[int]string{}
}
func MapListString(length interface{}) []map[string]string {
    len := _length(length)
    if len != -1 {
        return make([]map[string]string, len)
    }
    return []map[string]string{}
}
func MapSliceInt(length interface{}) []int {
    len := _length(length)
    if len != -1 {
        return make([]int, len)
    }
    return []int{}
}
func MapSliceString(length interface{}) []string {
    len := _length(length)
    if len != -1 {
        return make([]string, len)
    }
    return []string{}
}
