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
	UseIP     int    `json:"useip"`     //1使用ip定位(默认),2使用自定义
	Longitude string `json:"longitude"` //经度 当 enable等于1或者2 且UseIP等于0时使用
	Latitude  string `json:"latitude"`  //纬度
	Accuracy  string `json:"accuracy"`  //精度（米）"
}

// EnvRequest represents the request parameters for EnvCreate
type EnvInfo struct {
	EnvId           uint64     `json:"envId" form:"envId"`
	System          string     `json:"system" form:"system"`               //
	UaVersion       string     `json:"uaVersion" form:"uaVersion"`         //
	PublicIp        string     `json:"publicIp" form:"publicIp"`           //
	IpChannel       string     `json:"ipChannel" form:"ipChannel"`         //
	Kernel          string     `json:"kernel" form:"kernel"`               //
	KernelVersion   string     `json:"kernelVersion" form:"kernelVersion"` //
	CustomerId      string     `json:"customerId" form:"customerId"`       //
	EnvName         string     `json:"envName" form:"envName"`             //
	Remark          string     `json:"remark" form:"remark"`               //
	Serial          string     `json:"serial" form:"serial"`               //
	UserAgent       string     `json:"userAgent"`                          //UserAgent 不写根据系统和浏览器版本自动生成
	Language        []string   `json:"language"`                           //浏览器的语言 不传会根据代理IP地址自动生成 详细看支持的语言详细列表(如果是使用动态IP自动生成为中文)
	Zone            string     `json:"zone"`                               //时区 不传会根据代理IP地址自动生成 详细查看时区支持的列表(如果是使用动态IP自动生成为北京时间)
	Geographic      Geographic `json:"geographic"`                         //地理位置 （默认使用IP定位 动态代理IP不支持次选项为禁止）
	Dpi             string     `json:"dpi"`                                //平面分辨率 空自动生成
	FontList        []string   `json:"fontList"`                           //字体列表 不传系统自动生成
	WebRTC          int        `json:"webRTC"`                             //WebRTC 3隐私 2替换 1真实 4禁用
	WebRTCIP        string     `json:"webRTCIP"`                           //Chrome即时通信组件，支持：proxy 替换 ，使用代理IP覆盖真实IP，代理场景使用 local 真实 ，网站会获取真实IP disabled 禁用(默认)，网站会拿不到IP
	Canvas          int        `json:"canvas"`                             //浏览器canvas指纹开关 1隐身 2倾向随机 3倾向一致性
	WebGl           int        `json:"webGl"`                              //浏览器webgl元数据指纹开关 1隐身 2真实（默认）
	AudioContext    int        `json:"audioContext"`                       //AudioContext 1隐身 2真实
	SpeechVoices    int        `json:"speechVoices"`                       //SpeechVoices指纹，1：每个浏览器使用当前电脑默认的SpeechVoices,真实 2：添加相应的噪音，同一电脑上为每个浏览器生成不同的SpeechVoices（默认）
	MediaDevice     int        `json:"mediaDevice"`                        //媒体设备开关，1：关闭（每个浏览器使用当前电脑默认的媒体设备id，真实）2：启用（使用相匹配的值代替您真实的媒体设备ID，噪声）（默认）
	Cpu             int        `json:"cpu"`                                //CPU核心数量 不传会自动生成
	Mem             float64    `json:"mem"`                                //内存参数 不传会自动生成
	DeviceName      string     `json:"deviceName"`                         //计算机名 不传会自动生成
	Mac             string     `json:"mac"`                                //MAC地址 不传会自动生成
	Hardware        int        `json:"hardware"`                           //硬件加速 1开启（默认） 2关闭
	Bluetooth       int        `json:"bluetooth"`                          //蓝牙 1开启 2关闭(默认)
	DoNotTrack      int        `json:"doNotTrack"`                         //请勿跟踪”浏览器设置   1不启用 2启用 3默认（默认）
	EnableScanPort  int        `json:"enableScanPort"`                     //端口扫描防护 1开启(默认) 2关闭
	ScanPort        []int      `json:"scanPort"`                           //白名单 0~65535 关闭状态不写 当EnableScanPort是1时这里为空会自动生成本地端口
	EnableCookie    int        `json:"enableCookie"`                       //Cookie 1按环境，2按用户
	Enableopen      int        `json:"enableopen"`                         //多开设置 1开启，2关闭
	Enablenotice    int        `json:"enablenotice"`                       //网页通知 1开启，2关闭
	Enablepic       int        `json:"enablepic"`                          //禁止加载图片 1开启，2关闭
	PicSize         string     `json:"picsize"`                            //图片大小
	Enablesound     int        `json:"enablesound"`                        //禁止播放声音 1开启，2关闭",
	Enablevideo     int        `json:"enablevideo"`                        //禁止加载视频 1开启，2关闭"
	IgnoreCookieErr int        `json:"ignoreCookieErr"`                    //忽略Cookie格式错误, 1是 2否
	Proxy           string     `json:"proxy" form:"proxy"`                 //代理配置，格式为：socks5://user:pwd@ipaddr:6666
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
	EnvId int64 `json:"envId" form:"envId"` //
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
