package output

// Response is the root response for every api call
type Response struct {
	Success bool
	Errors  Errors
	Result  interface{}
	Meta    Meta
}

// Errors is our error struct for if something goes wrong
type Errors struct {
	Code int
	Msg  string
}

// Meta contains our version number, by ad
type Meta struct {
	Version string
	By      string
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

var JSON500Response = `{"Success":false,"Errors":{"Code":500,"Msg":"Internal Server Error"},"Result":null,"Meta":{"Version":"Early Alpha","By":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`
var JSON401Response = `{"success":false,"errors":{"code":403,"msg":"Invalid authentication token"},"result":null,"meta":{"version":"early alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`
var JSON403Response = `{"success":false,"errors":{"code":401,"msg":"Authenticated header not included"},"result":null,"meta":{"version":"early alpha","by":"Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead"}}`

// GetMeta returns the meta info for our response
func GetMeta() Meta {
	return Meta{
		By:      "Izaac Crooke, Dhayrin Colbert, Dominic Porter, Hayden Woodhead",
		Version: "Early Alpha",
	}
}
