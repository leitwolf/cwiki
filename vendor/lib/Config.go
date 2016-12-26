package lib

//
// Conf Config实例
//
var Conf *Config

// MenuItem 菜单项
type MenuItem struct {
	Key, Value string
}

//
// Config 配置文件
//
type Config struct {
	// 是否是测试
	Debug bool
	// 网站域名
	Site string
	// 作者
	Author string
	// 网站名称
	Title string
	// 头部菜单 {archive:归档}
	Menu []MenuItem
	// github
	Github string
	// email
	Email string
	// 文章所在目录
	PostDir string
	// 生成文件存放目录，后面有/
	SaveDir string
	// 首页每页条目数
	IndexCountPerPage int
	// 每页条目数（归档，标签，目录中用到）
	CountPerPage int
	// 原创信息模板
	OriginalTemplate string
}

//
// CreateConfig 建立配置信息
//
func CreateConfig() {
	Conf = &Config{}
	Conf.Debug = false
	Conf.Site = "http://lonewolf.me/"
	Conf.Author = "Lonewolf"
	Conf.Title = "lonewolf的博客"
	Conf.Menu = make([]MenuItem, 0)
	Conf.Menu = append(Conf.Menu, MenuItem{"archive", "归档"})
	Conf.Menu = append(Conf.Menu, MenuItem{"category", "分类"})
	Conf.Menu = append(Conf.Menu, MenuItem{"tag", "标签"})
	Conf.Menu = append(Conf.Menu, MenuItem{"about", "关于"})
	Conf.Github = "leitwolf"
	Conf.Email = "wolf.pan@qq.com"
	Conf.PostDir = "_post"
	Conf.SaveDir = "content/"
	Conf.IndexCountPerPage = 10
	Conf.CountPerPage = 15
	Conf.OriginalTemplate = `原链接地址：[{url}]({url})  
原创博客，转载请注明。

---`
}
