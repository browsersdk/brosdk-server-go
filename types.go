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
	EnvId       string `json:"envId" form:"envId"`             //envid 更新时必填
	CustomerId  string `json:"customerId" form:"customerId"`   //三方用户id
	EnvName     string `json:"envName" form:"envName"`         //环境名称
	Remark      string `json:"remark" form:"remark"`           //备注
	Serial      string `json:"serial" form:"serial"`           //环境序号
	Cookie      string `json:"cookie" form:"cookie"`           //cookie
	Proxy       string `json:"proxy" form:"proxy"`             //代理配置，格式为：socks5://user:pwd@ipaddr:6666
	BridgeProxy string `json:"bridgeProxy" form:"bridgeProxy"` //桥代理配置，格式为：socks5://user:pwd@ipaddr:6666
	IpChannel   string `json:"ipChannel" form:"ipChannel"`     //IP监测渠道  海外代理：ip2location，国内代理：ipdata
	Region      string `json:"region" form:"region"`           //国家代号，当无法获取代理配置时，传此参数生成对应区域ip，否则获取客户端ip
	Finger      Finger `json:"finger" form:"finger"`
}

type NextSystem struct { //次优优先级操作系统
	Android string `json:"Android"`
	MacOS   string `json:"MacOS"`
	IOS     string `json:"IOS"`
	Linux   string `json:"Linux"`
}

type Finger struct {
	System           string     `json:"system"`           //系统
	Nextsystem       NextSystem `json:"nextsystem"`       //（可选不使用设置为空）指纹生成的时候回根据System 或者 Nextsystem去生成，如果给了多项会从这些当中随机选择生成
	Kernel           string     `json:"kernel"`           //内核
	KernelVersion    string     `json:"kernelVersion"`    //内核版本
	UaVersion        string     `json:"uaVersion"`        //浏览器大版本号 不传自动生成 详细查看具体支持的版本号大版本
	Ua               string     `json:"ua"`               //不写根据系统和浏览器版本自动生成
	Language         []string   `json:"language"`         //浏览器的语言 不传会根据代理IP地址自动生成 详细看支持的语言详细列表(如果是使用动态IP自动生成为中文)
	Uilanguage       []string   `json:"uilanguage"`       //界面语言 空会根据代理IP地址自动生成
	Zone             string     `json:"zone"`             //时区 不传会根据代理IP地址自动生成 详细查看时区支持的列表(如果是使用动态IP自动生成为北京时间)
	Geographic       Geographic `json:"geographic"`       //地理位置 （默认使用IP定位 动态代理IP不支持次选项为禁止）
	Dpi              string     `json:"dpi"`              //平面分辨率 空自动生成
	Widowssize       string     `json:"widowssize"`       //浏览器窗口大小  （目前可传可不传已经不用了使用dpi计算得出）
	Font             Font       `json:"font"`             //字体列表
	FontFinger       int        `json:"fontFinger"`       //字体指纹 1隐身 2真实
	ClientRects      int        `json:"clientRects"`      //ClientRects 1隐身 2真实 目前同fontfinger值
	WebRTC           int        `json:"webRTC"`           //WebRTC 0:禁用,网站会拿不到IP 1:真实,网站会获取真实IP 2:替换,使用代理IP覆盖真实IP
	WebRTCIP         string     `json:"webRTCIP"`         //Chrome即时通信组件，支持：proxy 替换 ，使用代理IP覆盖真实IP，代理场景使用 local 真实 ，网站会获取真实IP disabled 禁用(默认)，网站会拿不到IP
	Canvas           int        `json:"canvas"`           //浏览器canvas指纹开关 0:倾向一致性 1:关闭 2:倾向随机性
	WebGl            int        `json:"webGl"`            //浏览器webgl元数据指纹开关 1隐身 2真实（默认）
	WebGlInfo        int        `json:"webGlInfo"`        //浏览器WebGlInfo 1:真实 2:自定义
	WebGLVendor      string     `json:"webGlVendor"`      //浏览器WebGL厂商，Windows系统可选值为Google Inc. (NVIDIA)、Google Inc. (AMD)、Google Inc. (Intel)，MacOS系统可选值为Google Inc. (ATI Technologies Inc.)、Google Inc. (NVIDIA)、Google Inc. (Apple)，Android系统可选值为Qualcomm，IOS系统可选值为Apple Inc 自定义时传值，为空会自动生成
	WebGLRenderer    string     `json:"webGlRenderer"`    //浏览器WebGL渲染 自定义时传值，为空会自动生成，该字段不为空时webGLVendor必传
	AudioContext     int        `json:"audioContext"`     //AudioContext 1隐身 2真实
	SpeechVoices     int        `json:"speechVoices"`     //SpeechVoices指纹，1：每个浏览器使用当前电脑默认的SpeechVoices,真实 2：添加相应的噪音，同一电脑上为每个浏览器生成不同的SpeechVoices（默认）
	MediaDevice      int        `json:"mediaDevice"`      //媒体设备开关，1：关闭（每个浏览器使用当前电脑默认的媒体设备id，真实）2：启用（使用相匹配的值代替您真实的媒体设备ID，噪声）（默认）
	Cpu              int        `json:"cpu"`              //CPU核心数量 不传会自动生成
	Mem              float64    `json:"mem"`              //内存参数 不传会自动生成
	DeviceName       string     `json:"deviceName"`       //计算机名 不传会自动生成
	Mac              string     `json:"mac"`              //MAC地址 不传会自动生成
	Hardware         int        `json:"hardware"`         //硬件加速 1开启（默认） 2关闭
	Bluetooth        int        `json:"bluetooth"`        //蓝牙 1开启 2关闭(默认)
	DoNotTrack       int        `json:"doNotTrack"`       //请勿跟踪浏览器设置   1不启用 2启用 3默认
	EnableScanPort   int        `json:"enableScanPort"`   //端口扫描防护 1开启(默认) 2关闭
	ScanPort         string     `json:"scanPort"`         //白名单 0~65535 关闭状态不写 当EnableScanPort是1时这里为空会自动生成本地端口
	EnableOpen       int        `json:"enableOpen"`       //多开设置 1开启，2关闭
	EnableOpenNumber int        `json:"enableOpenNumber"` //多开人数设置
	EnableNotice     int        `json:"enableNotice"`     //网页通知 1开启，2关闭
	EnablePic        int        `json:"enablePic"`        //禁止加载图片 1开启，2关闭
	PicSize          string     `json:"picSize"`          //图片大小
	EnableGc         int        `json:"enableGc"`         //是否开启垃圾回收 1开启 2关闭
	GcTime           int        `json:"gcTime"`           //垃圾回收时间,当enableGc为1时必填 1-5
	EnableSound      int        `json:"enableSound"`      //禁止播放声音 1开启，2关闭
	EnableVideo      int        `json:"enableVideo"`      //禁止加载视频 1开启，2关闭
	Battery          int        `json:"battery"`          //电池 1真实 2噪声
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
