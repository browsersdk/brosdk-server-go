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
	Enable    int    `json:"enable"`    //启用1(默认) 询问2 禁止3
	User      int    `json:"user"`      //1使用ip定位(默认),2使用自定义
	Longitude string `json:"longitude"` //经度(当enable等于2且UseIP等于0时使用) -180 - 180
	Latitude  string `json:"latitude"`  //纬度(当enable等于2且UseIP等于0时使用) -90 - 90
	Accuracy  string `json:"accuracy"`  //精度（米）(当enable等于2且UseIP等于0时使用) 10 - 5000
}

type Font struct {
	Enable int      `json:"enable"` //1隐身包含 2使用用户设置
	List   []string `json:"list"`   //enable=2开启情况下如果为空系统自动生成
}

// EnvRequest represents the request parameters for EnvCreate
type EnvInfo struct {
	EnvId         string `json:"envId,omitempty" form:"envId"`                 //envid为空时自动生成
	CustomerId    string `json:"customerId,omitempty" form:"customerId"`       //三方用户id
	EnvName       string `json:"envName,omitempty" form:"envName"`             //环境名称
	Serial        string `json:"serial,omitempty" form:"serial"`               //环境序号
	Proxy         string `json:"proxy,omitempty" form:"proxy"`                 //代理配置，格式为：socks5://user:pwd@ipaddr:6666
	BridgeProxy   string `json:"bridgeProxy,omitempty" form:"bridgeProxy"`     //桥代理配置，格式为：socks5://user:pwd@ipaddr:6666
	IpChannel     string `json:"ipChannel,omitempty" form:"ipChannel"`         //IP监测渠道  海外代理：ip2location，国内代理：ipdata
	Region        string `json:"region,omitempty" form:"region"`               //国家代号，当无法获取代理配置时，传此参数生成对应区域ip，否则获取客户端ip
	Kernel        string `json:"kernel,omitempty" form:"kernel"`               //内核,以Finger中的kernel为准
	KernelVersion string `json:"kernelVersion,omitempty" form:"kernelVersion"` //内核版本，以Finger中的kernelVersion为准
	Finger        Finger `json:"finger,omitempty" form:"finger"`
}

type CreateEnv struct {
	CustomerId    string `json:"customerId" form:"customerId"`       //三方用户id
	EnvName       string `json:"envName" form:"envName"`             //环境名称
	Serial        string `json:"serial" form:"serial"`               //环境序号
	Proxy         string `json:"proxy" form:"proxy"`                 //代理配置，格式为：socks5://user:pwd@ipaddr:6666
	BridgeProxy   string `json:"bridgeProxy" form:"bridgeProxy"`     //桥代理配置，格式为：socks5://user:pwd@ipaddr:6666
	IpChannel     string `json:"ipChannel" form:"ipChannel"`         //IP监测渠道  海外代理：ip2location，国内代理：ipdata
	Region        string `json:"region" form:"region"`               //国家代号，当无法获取代理配置时，传此参数生成对应区域ip，否则获取客户端ip
	Kernel        string `json:"kernel" form:"kernel"`               //内核,以Finger中的kernel为准
	KernelVersion string `json:"kernelVersion" form:"kernelVersion"` //内核版本，以Finger中的kernelVersion为准
	Finger        Finger `json:"finger" form:"finger"`
}

type UpdateEnv struct {
	EnvId         string  `json:"envId" form:"envId"`                 //envid为空时自动生成
	CustomerId    *string `json:"customerId" form:"customerId"`       //三方用户id
	EnvName       *string `json:"envName" form:"envName"`             //环境名称
	Serial        *string `json:"serial" form:"serial"`               //环境序号
	Proxy         *string `json:"proxy" form:"proxy"`                 //代理配置，格式为：socks5://user:pwd@ipaddr:6666
	BridgeProxy   *string `json:"bridgeProxy" form:"bridgeProxy"`     //桥代理配置，格式为：socks5://user:pwd@ipaddr:6666
	IpChannel     *string `json:"ipChannel" form:"ipChannel"`         //IP监测渠道  海外代理：ip2location，国内代理：ipdata
	Region        *string `json:"region" form:"region"`               //国家代号，当无法获取代理配置时，传此参数生成对应区域ip，否则获取客户端ip
	Kernel        *string `json:"kernel" form:"kernel"`               //内核,以Finger中的kernel为准
	KernelVersion *string `json:"kernelVersion" form:"kernelVersion"` //内核版本，以Finger中的kernelVersion为准
	Finger        *Finger `json:"finger" form:"finger"`
}

type UpdateEnvMeta struct {
	EnvId       string  `json:"envId" form:"envId"`             //envid为空时自动生成
	CustomerId  *string `json:"customerId" form:"customerId"`   //三方用户id
	EnvName     *string `json:"envName" form:"envName"`         //环境名称
	Serial      *string `json:"serial" form:"serial"`           //环境序号
	Proxy       *string `json:"proxy" form:"proxy"`             //代理配置，格式为：socks5://user:pwd@ipaddr:6666
	BridgeProxy *string `json:"bridgeProxy" form:"bridgeProxy"` //桥代理配置，格式为：socks5://user:pwd@ipaddr:6666
}

type Finger struct {
	System        string `json:"system,omitempty"`        //系统
	Kernel        string `json:"kernel,omitempty"`        //内核
	KernelVersion string `json:"kernelVersion,omitempty"` //内核版本
}

// EnvResponse represents the response for EnvCreate
type EnvResponse struct {
	Code  int     `json:"code"`
	Data  EnvInfo `json:"data"`
	Msg   string  `json:"msg"`
	ReqId string  `json:"reqId"`
}

// EnvReq represents the request parameters for destroy operation
type EnvDelReq struct {
	EnvId string `json:"envId" form:"envId"` //
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
	EnvIds     []string `json:"envIds" form:"envIds"`         //主键集合
	CustomerId string   `json:"customerId" form:"customerId"` //客户ID
}

// PageResp represents paginated response structure
type PageResp struct {
	Code  int    `json:"code"`
	Data  Page   `json:"data"`
	Msg   string `json:"msg"`
	ReqId string `json:"reqId"`
	Total int64  `json:"total"`
}

type Page struct {
	List        []EnvInfo `json:"list"`        //数据列表
	Total       int64     `json:"total"`       //总条数
	PageSize    int       `json:"pageSize"`    //分页大小
	CurrentPage int       `json:"currentPage"` //当前第几页
}

type KernelVersionInfo struct {
	Id      int    `json:"id"`
	Version string `json:"version"`
}

type SYSDEVICE struct {
	Name            string `json:"name"`            //系统平台 Windows Android MAC iPhone Linux
	System          string `json:"system"`          //具体的系统版本号
	Browser         string `json:"browser"`         //浏览器上面对应的版本
	PlatformVersion string `json:"platformVersion"` //浏览器的platformVersion
}

// 操作系统和浏览器核关系
type SYSTEMKERNEL struct {
	Windows []SYSDEVICE `json:"Windows"`
	MacOS   []SYSDEVICE `json:"MacOS"`
	Linux   []SYSDEVICE `json:"Linux"`
}

type CONAB struct {
	Country  string `json:"country"`  //国家
	Province string `json:"province"` //省
	AB       string `json:"ab"`       //简写
	ECountry string `json:"ecountry"` //国家英文名
	Must     int    `json:"must"`     //是不是这个国家的必须语言 1是 0不是
	Name     string `json:"name"`
	Code     string `json:"code"`
}

type GetUiFingerList struct {
	ChromeKernelversion   []KernelVersionInfo `json:"chromeKernelVersion"`   //支持的浏览器内核大版本
	FirefoxKernelversion  []KernelVersionInfo `json:"firefoxKernelversion"`  //支持的浏览器火狐内核大版本
	System                SYSTEMKERNEL        `json:"system"`                //操作系统版本
	ChromeUAversion       []string            `json:"chromeUAversion"`       //浏览器UA版本
	FirefoxUAversion      []string            `json:"firefoxUAversion"`      //火狐浏览器UA版本
	Language              []CONAB             `json:"language"`              //语言
	Zone                  []string            `json:"zone"`                  //时区
	Dpi                   any                 `json:"dpi"`                   //屏幕分辨率
	Webgl                 any                 `json:"webgl"`                 //webgl
	Cpu                   any                 `json:"cpu"`                   //CPU参数
	Mem                   any                 `json:"mem"`                   //内存参数
	Region                any                 `json:"region"`                //Region
	PlatformKernelversion any                 `json:"platformKernelversion"` //平台可用的内核
}

type GetUiFingerListResponse struct {
	Code  int             `json:"code"`
	Data  GetUiFingerList `json:"data"`
	Msg   string          `json:"msg"`
	ReqId string          `json:"reqId"`
}
