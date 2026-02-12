package brosdk

// Response represents the standard API response structure
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	ReqId string      `json:"reqId"`
}

// GetUserSigRequest represents the request parameters for GetUserSig
type GetUserSigRequest struct {
	CustomerId string `json:"customerId"`
	Duration   int    `json:"duration"`
}

// UserSigData represents the data structure in GetUserSig response
type UserSigData struct {
	ExpireTime int64  `json:"expireTime"`
	UserSig    string `json:"userSig"`
}

// GetUserSigResponse represents the response for GetUserSig
type GetUserSigResponse struct {
	Code  int         `json:"code"`
	Data  UserSigData `json:"data"`
	Msg   string      `json:"msg"`
	ReqId string      `json:"reqId"`
}

// Geographic represents the geographic information structure
type Geographic struct {
	Accuracy  string `json:"accuracy"`
	Enable    int    `json:"enable"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Useip     int    `json:"useip"`
}

// EnvCreateRequest represents the request parameters for EnvCreate
type EnvCreateRequest struct {
	AudioContext    int        `json:"AudioContext"`
	Bluetooth       int        `json:"Bluetooth"`
	Canvas          int        `json:"Canvas"`
	Cpu             int        `json:"Cpu"`
	DPI             string     `json:"DPI"`
	DeviceName      string     `json:"DeviceName"`
	DoNotTrack      int        `json:"DoNotTrack"`
	EnableScanPort  int        `json:"EnableScanPort"`
	Enablesound     int        `json:"Enablesound"`
	Enablevideo     int        `json:"Enablevideo"`
	FontList        []string   `json:"FontList"`
	Hardware        int        `json:"Hardware"`
	Language        []string   `json:"Language"`
	Mac             string     `json:"Mac"`
	MediaDevice     int        `json:"MediaDevice"`
	Mem             int        `json:"Mem"`
	ScanPort        []int      `json:"ScanPort"`
	SpeechVoices    int        `json:"SpeechVoices"`
	UserAgent       string     `json:"UserAgent"`
	WebGl           int        `json:"WebGl"`
	WebRTC          int        `json:"WebRTC"`
	WebRTCIP        string     `json:"WebRTCIP"`
	Zone            string     `json:"Zone"`
	CustomerId      string     `json:"customerId"`
	EnableCookie    int        `json:"enableCookie"`
	Enablenotice    int        `json:"enablenotice"`
	Enableopen      int        `json:"enableopen"`
	Enablepic       int        `json:"enablepic"`
	EnvId           int        `json:"envId"`
	EnvName         string     `json:"envName"`
	Geographic      Geographic `json:"geographic"`
	IgnoreCookieErr int        `json:"ignoreCookieErr"`
	IpChannel       string     `json:"ipChannel"`
	Kernel          string     `json:"kernel"`
	KernelVersion   string     `json:"kernelVersion"`
	Picsize         string     `json:"picsize"`
	Proxy           string     `json:"proxy"`
	PublicIp        string     `json:"publicIp"`
	Remark          string     `json:"remark"`
	Serial          string     `json:"serial"`
	System          string     `json:"system"`
	UaVersion       string     `json:"uaVersion"`
}

// EnvCreateResponse represents the response for EnvCreate
type EnvCreateResponse struct {
	Code  int              `json:"code"`
	Data  EnvCreateRequest `json:"data"`
	Msg   string           `json:"msg"`
	ReqId string           `json:"reqId"`
}

// EnvReq represents the request parameters for destroy operation
type EnvReq struct {
	EnvId uint64 `json:"envId" form:"envId"` //
}

// ReqPage represents pagination request parameters
type ReqPage struct {
	Page     int `json:"page" form:"page" query:"-"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" query:"-"` // 每页大小
}

// GetEnvPageReq represents the request parameters for paginated environment listing
type GetEnvPageReq struct {
	ReqPage `query:"-"`

	SortOrder  string   `json:"-" query:"type:order;column:id"`
	EnvIds     []uint64 `json:"envIds" form:"envIds"`         //主键集合
	CustomerId string   `json:"customerId" form:"customerId"` //客户ID
}

// EnvInfo represents environment information structure
type EnvInfo struct {
	EnvId      uint64 `json:"envId"`
	CustomerId string `json:"customerId"`
	EnvName    string `json:"envName"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// PageResp represents paginated response structure
type PageResp struct {
	Code  int       `json:"code"`
	Data  []EnvInfo `json:"data"`
	Msg   string    `json:"msg"`
	ReqId string    `json:"reqId"`
	Total int64     `json:"total"`
}
