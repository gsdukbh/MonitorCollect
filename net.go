package main

import "encoding/json"

// 处理telegraf 采集的网络数据

// NetInterfaceFields 表示每个网卡(interface)的基础网络统计
// 对应样例中 name=net, tag.interface=end0 的字段
type NetInterfaceFields struct {
	BytesRecv   int64 `json:"bytes_recv"`   // 接收的总字节数
	BytesSent   int64 `json:"bytes_sent"`   // 发送的总字节数
	DropIn      int64 `json:"drop_in"`      // 接收时丢弃的数据包数
	DropOut     int64 `json:"drop_out"`     // 发送时丢弃的数据包数
	ErrIn       int64 `json:"err_in"`       // 接收时的错误数
	ErrOut      int64 `json:"err_out"`      // 发送时的错误数
	PacketsRecv int64 `json:"packets_recv"` // 接收的数据包总数
	PacketsSent int64 `json:"packets_sent"` // 发送的数据包总数
	Speed       int64 `json:"speed"`        // 网卡速度（Mbps）
}

// NetProtoFields 表示主机层面的网络协议统计
// 对应样例中 name=net, tag.interface=all 的字段
type NetProtoFields struct {
	// ICMP (Internet Control Message Protocol) - 互联网控制消息协议
	IcmpInAddrMaskReps     int64 `json:"icmp_inaddrmaskreps"`     // 接收的地址掩码应答数
	IcmpInAddrMasks        int64 `json:"icmp_inaddrmasks"`        // 接收的地址掩码请求数
	IcmpInCsumErrors       int64 `json:"icmp_incsumerrors"`       // 接收的校验和错误数
	IcmpInDestUnreachs     int64 `json:"icmp_indestunreachs"`     // 接收的目标不可达消息数
	IcmpInEchoReps         int64 `json:"icmp_inechoreps"`         // 接收的回显应答（ping 应答）数
	IcmpInEchos            int64 `json:"icmp_inechos"`            // 接收的回显请求（ping 请求）数
	IcmpInErrors           int64 `json:"icmp_inerrors"`           // 接收的 ICMP 错误总数
	IcmpInMsgs             int64 `json:"icmp_inmsgs"`             // 接收的 ICMP 消息总数
	IcmpInParmProbs        int64 `json:"icmp_inparmprobs"`        // 接收的参数问题消息数
	IcmpInRedirects        int64 `json:"icmp_inredirects"`        // 接收的重定向消息数
	IcmpInSrcQuenchs       int64 `json:"icmp_insrcquenchs"`       // 接收的源抑制消息数
	IcmpInTimeExcds        int64 `json:"icmp_intimeexcds"`        // 接收的超时消息数
	IcmpInTimestampReps    int64 `json:"icmp_intimestampreps"`    // 接收的时间戳应答数
	IcmpInTimestamps       int64 `json:"icmp_intimestamps"`       // 接收的时间戳请求数
	IcmpOutAddrMaskReps    int64 `json:"icmp_outaddrmaskreps"`    // 发送的地址掩码应答数
	IcmpOutAddrMasks       int64 `json:"icmp_outaddrmasks"`       // 发送的地址掩码请求数
	IcmpOutDestUnreachs    int64 `json:"icmp_outdestunreachs"`    // 发送的目标不可达消息数
	IcmpOutEchoReps        int64 `json:"icmp_outechoreps"`        // 发送的回显应答（ping 应答）数
	IcmpOutEchos           int64 `json:"icmp_outechos"`           // 发送的回显请求（ping 请求）数
	IcmpOutErrors          int64 `json:"icmp_outerrors"`          // 发送的 ICMP 错误总数
	IcmpOutMsgs            int64 `json:"icmp_outmsgs"`            // 发送的 ICMP 消息总数
	IcmpOutParmProbs       int64 `json:"icmp_outparmprobs"`       // 发送的参数问题消息数
	IcmpOutRateLimitGlobal int64 `json:"icmp_outratelimitglobal"` // 全局速率限制丢弃的消息数
	IcmpOutRateLimitHost   int64 `json:"icmp_outratelimithost"`   // 主机速率限制丢弃的消息数
	IcmpOutRedirects       int64 `json:"icmp_outredirects"`       // 发送的重定向消息数
	IcmpOutSrcQuenchs      int64 `json:"icmp_outsrcquenchs"`      // 发送的源抑制消息数
	IcmpOutTimeExcds       int64 `json:"icmp_outtimeexcds"`       // 发送的超时消息数
	IcmpOutTimestampReps   int64 `json:"icmp_outtimestampreps"`   // 发送的时间戳应答数
	IcmpOutTimestamps      int64 `json:"icmp_outtimestamps"`      // 发送的时间戳请求数

	// ICMP 消息类型拆分（按 ICMP 类型码分类）
	IcmpMsgInType0  int64 `json:"icmpmsg_intype0"`  // 接收的 Type 0 消息数（Echo Reply - 回显应答）
	IcmpMsgInType3  int64 `json:"icmpmsg_intype3"`  // 接收的 Type 3 消息数（Destination Unreachable - 目标不可达）
	IcmpMsgInType8  int64 `json:"icmpmsg_intype8"`  // 接收的 Type 8 消息数（Echo Request - 回显请求）
	IcmpMsgInType11 int64 `json:"icmpmsg_intype11"` // 接收的 Type 11 消息数（Time Exceeded - 超时）
	IcmpMsgOutType0 int64 `json:"icmpmsg_outtype0"` // 发送的 Type 0 消息数（Echo Reply - 回显应答）
	IcmpMsgOutType3 int64 `json:"icmpmsg_outtype3"` // 发送的 Type 3 消息数（Destination Unreachable - 目标不可达）
	IcmpMsgOutType8 int64 `json:"icmpmsg_outtype8"` // 发送的 Type 8 消息数（Echo Request - 回显请求）

	// IP (Internet Protocol) - 网际协议
	IpDefaultTTL      int64 `json:"ip_defaultttl"`      // 默认 TTL（Time To Live）值
	IpForwarding      int64 `json:"ip_forwarding"`      // IP 转发开关（1=开启，0=关闭）
	IpForwDatagrams   int64 `json:"ip_forwdatagrams"`   // 转发的数据报数量
	IpFragCreates     int64 `json:"ip_fragcreates"`     // 创建的 IP 分片数量
	IpFragFails       int64 `json:"ip_fragfails"`       // 分片失败的数据报数量
	IpFragOKs         int64 `json:"ip_fragoks"`         // 成功分片的数据报数量
	IpInAddrErrors    int64 `json:"ip_inaddrerrors"`    // 接收的地址错误数据报数量
	IpInDelivers      int64 `json:"ip_indelivers"`      // 成功交付到上层协议的数据报数量
	IpInDiscards      int64 `json:"ip_indiscards"`      // 接收时丢弃的数据报数量
	IpInHdrErrors     int64 `json:"ip_inhdrerrors"`     // 接收的头部错误数据报数量
	IpInReceives      int64 `json:"ip_inreceives"`      // 接收的 IP 数据报总数
	IpInUnknownProtos int64 `json:"ip_inunknownprotos"` // 接收的未知协议数据报数量
	IpOutDiscards     int64 `json:"ip_outdiscards"`     // 发送时丢弃的数据报数量
	IpOutNoRoutes     int64 `json:"ip_outnoroutes"`     // 因无路由而丢弃的数据报数量
	IpOutRequests     int64 `json:"ip_outrequests"`     // 本地上层协议请求发送的数据报数量
	IpOutTransmits    int64 `json:"ip_outtransmits"`    // 成功发送的 IP 数据报总数
	IpReasmFails      int64 `json:"ip_reasmfails"`      // IP 重组失败的次数
	IpReasmOKs        int64 `json:"ip_reasmoks"`        // IP 重组成功的次数
	IpReasmReqds      int64 `json:"ip_reasmreqds"`      // 需要重组的 IP 分片数量
	IpReasmTimeout    int64 `json:"ip_reasmtimeout"`    // IP 重组超时的次数

	// TCP (Transmission Control Protocol) - 传输控制协议
	TcpActiveOpens  int64 `json:"tcp_activeopens"`  // 主动打开的 TCP 连接数（客户端发起）
	TcpAttemptFails int64 `json:"tcp_attemptfails"` // 连接尝试失败的次数
	TcpCurrEstab    int64 `json:"tcp_currestab"`    // 当前已建立的 TCP 连接数
	TcpEstabResets  int64 `json:"tcp_estabresets"`  // 已建立连接被重置的次数
	TcpInCsumErrors int64 `json:"tcp_incsumerrors"` // 接收的校验和错误数
	TcpInErrs       int64 `json:"tcp_inerrs"`       // 接收的 TCP 错误总数
	TcpInSegs       int64 `json:"tcp_insegs"`       // 接收的 TCP 段总数
	TcpMaxConn      int64 `json:"tcp_maxconn"`      // TCP 最大连接数（-1 表示动态）
	TcpOutRsts      int64 `json:"tcp_outrsts"`      // 发送的 RST 段数量
	TcpOutSegs      int64 `json:"tcp_outsegs"`      // 发送的 TCP 段总数
	TcpPassiveOpens int64 `json:"tcp_passiveopens"` // 被动打开的 TCP 连接数（服务端接受）
	TcpRetransSegs  int64 `json:"tcp_retranssegs"`  // 重传的 TCP 段数量
	TcpRtoAlgorithm int64 `json:"tcp_rtoalgorithm"` // RTO（重传超时）算法编号
	TcpRtoMax       int64 `json:"tcp_rtomax"`       // 最大 RTO 值（毫秒）
	TcpRtoMin       int64 `json:"tcp_rtomin"`       // 最小 RTO 值（毫秒）

	// UDP (User Datagram Protocol) - 用户数据报协议
	UdpIgnoredMulti int64 `json:"udp_ignoredmulti"` // 忽略的多播数据报数量
	UdpInCsumErrors int64 `json:"udp_incsumerrors"` // 接收的校验和错误数
	UdpInDatagrams  int64 `json:"udp_indatagrams"`  // 接收的 UDP 数据报总数
	UdpInErrors     int64 `json:"udp_inerrors"`     // 接收的 UDP 错误总数
	UdpMemErrors    int64 `json:"udp_memerrors"`    // UDP 内存分配错误数
	UdpNoPorts      int64 `json:"udp_noports"`      // 发送到无监听端口的数据报数量
	UdpOutDatagrams int64 `json:"udp_outdatagrams"` // 发送的 UDP 数据报总数
	UdpRcvbufErrors int64 `json:"udp_rcvbuferrors"` // 接收缓冲区溢出错误数
	UdpSndbufErrors int64 `json:"udp_sndbuferrors"` // 发送缓冲区溢出错误数

	// UDPLite (Lightweight User Datagram Protocol) - 轻量级用户数据报协议
	UdpliteIgnoredMulti int64 `json:"udplite_ignoredmulti"` // 忽略的 UDPLite 多播数据报数量
	UdpliteInCsumErrors int64 `json:"udplite_incsumerrors"` // 接收的 UDPLite 校验和错误数
	UdpliteInDatagrams  int64 `json:"udplite_indatagrams"`  // 接收的 UDPLite 数据报总数
	UdpliteInErrors     int64 `json:"udplite_inerrors"`     // 接收的 UDPLite 错误总数
	UdpliteMemErrors    int64 `json:"udplite_memerrors"`    // UDPLite 内存分配错误数
	UdpliteNoPorts      int64 `json:"udplite_noports"`      // 发送到无监听端口的 UDPLite 数据报数量
	UdpliteOutDatagrams int64 `json:"udplite_outdatagrams"` // 发送的 UDPLite 数据报总数
	UdpliteRcvbufErrors int64 `json:"udplite_rcvbuferrors"` // UDPLite 接收缓冲区溢出错误数
	UdpliteSndbufErrors int64 `json:"udplite_sndbuferrors"` // UDPLite 发送缓冲区溢出错误数
}

// FromFieldsMap 填充 NetInterfaceFields
func (n *NetInterfaceFields) FromFieldsMap(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, n)
}

// FromFieldsMap 填充 NetProtoFields
func (n *NetProtoFields) FromFieldsMap(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, n)
}
