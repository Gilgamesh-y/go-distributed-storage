package conf

const (
	// StoreLocal : 节点本地
	StoreLocal = iota
	// StoreOSS : 阿里OSS
	StoreOSS
	// StoreAll : 所有类型的存储都存一份数据
	StoreAll
)