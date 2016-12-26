package lib

//
// Init 初始化各个操作单元
//
func Init() {
	// 配置文件
	CreateConfig()
	// 建立各个模板
	CreateTemplate()
	// 开始
	StartManager()
}
