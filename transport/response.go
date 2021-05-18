package transport

const (
	OK                      = "OK"
	InvalidTimestamp        = "InvalidTimestamp"
	Forbidden               = "Forbidden"
	HandlerNotFound         = "HandlerNotFound"
	TokenNotFound           = "TokenNotFound"
	ServiceNotOpened        = "ServiceNotOpened"
	ServiceNotAuthorized    = "ServiceNotAuthorized"
	ServerError             = "ServerError"
	HandlerClosed           = "HandlerClosed"
	Timeout                 = "Timeout"
	Uninitialized           = "Uninitialized"
	EncodeError             = "EncodeError"
	DecodeError             = "DecodeError"
	FileNotFound            = "FileNotFound"
	DownloadError           = "DownloadError"
	DeployError             = "DeployError"
	ServiceSwitchError      = "ServiceSwitchError"
	Upgrading               = "Upgrading"
	ParameterEmpty          = "ParameterEmpty"
	ParameterTypeError      = "ParameterTypeError"
	FaultInjectCmdError     = "FaultInjectCmdError"
	FaultInjectExecuteError = "FaultInjectExecuteError"
	FaultInjectNotSupport   = "FaultInjectNotSupport"
	JavaAgentCmdError       = "JavaAgentCmdError"
)

type CodeType struct {
	Code int32
	Msg  string
}

var Code = map[string]CodeType{
	OK:                      {200, "success"},
	InvalidTimestamp:        {401, "invalid timestamp"},
	Forbidden:               {403, "forbidden"},
	HandlerNotFound:         {404, "request handler not found"},
	TokenNotFound:           {405, "access token not found"},
	ServiceNotOpened:        {410, "ahas service not opened"},
	ServiceNotAuthorized:    {411, "ahas service not authorized"},
	ServerError:             {500, "server error"},
	HandlerClosed:           {501, "handler closed"},
	Timeout:                 {510, "timeout"},
	Uninitialized:           {511, "uninitialized"},
	EncodeError:             {512, "encode error"},
	DecodeError:             {513, "decode error"},
	FileNotFound:            {514, "file not found"},
	DownloadError:           {515, "download file error"},
	DeployError:             {516, "deploy file error"},
	ServiceSwitchError:      {517, "service switch error"},
	Upgrading:               {518, "service is upgrading"},
	ParameterEmpty:          {600, "parameter is empty"},
	ParameterTypeError:      {601, "parameter type error"},
	FaultInjectCmdError:     {701, "cannot handle the faultInject cmd"},
	FaultInjectExecuteError: {702, "execute faultInject error"},
	FaultInjectNotSupport:   {703, "the inject type not support"},
	JavaAgentCmdError:       {704, "cannot handle the javaagent cmd"},
}

type Response struct {
	Code    int32       `json:"code"`
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Result  interface{} `json:"result"`
}

//Return default code message
func Return(codeType CodeType) *Response {
	if codeType == Code[OK] {
		return &Response{Code: codeType.Code, Success: true, Result: codeType.Msg}
	}
	return &Response{Code: codeType.Code, Success: false, Error: codeType.Msg}
}

//ReturnFail with error message
func ReturnFail(codeType CodeType, err string) *Response {
	return &Response{Code: codeType.Code, Success: false, Error: err}
}

//ReturnSuccess with result
func ReturnSuccess(result interface{}) *Response {
	return &Response{Code: Code[OK].Code, Success: true, Result: result}
}
