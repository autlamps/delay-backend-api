package output

// Response is the root response for every api call
type Response struct {
	Success bool        `json:"success"`
	Errors  interface{} `json:"errors"`
	Result  interface{} `json:"result"`
	Meta    Meta        `json:"meta"`
}

// Errors is our error struct for if something goes wrong
type Errors struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Meta contains our version number, by ad
type Meta struct {
	Version string `json:"version"`
	By      string `json:"by"`
}

// New500Response returns a response object with the info for a 500 response
func New500Response() Response {
	return Response{
		Success: false,
		Errors: Errors{
			Code: 500,
			Msg:  "Internal Server Error",
		},
		Result: nil,
		Meta:   GetMeta(),
	}
}

var JSON500Response = `{"success":false,"errors":{"code":500,"msg":"Internal Server Error"},"result":null,"meta":{"version":"Alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`
var JSON401Response = `{"success":false,"errors":{"code":403,"msg":"Invalid authentication token"},"result":null,"meta":{"version":"alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`
var JSON403Response = `{"success":false,"errors":{"code":401,"msg":"Authenticated header not included"},"result":null,"meta":{"version":"alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`
var JSON409Response = `{"success":false,"errors":{"code":409,"msg":"Email already registered"},"result":null,"meta":{"version":"alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`

// GetMeta returns the meta info for our response
func GetMeta() Meta {
	return Meta{
		By:      "Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead",
		Version: "Alpha",
	}
}
