package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//
// Tepl Template单例
//
var Tepl *Template

//
// Template 模板加载处理类
//
type Template struct {
	// 头部文件
	header []byte
	footer []byte
	// 主页
	Index []byte
	// 文章
	Article []byte
	// 关于
	About []byte
	// 归档
	Archive []byte
	// 标签或分类
	Tag []byte
}

//
// 获取当前运行程序的路径
//
func (t *Template) getAppPath() string {
	if Conf.Debug {
		return ""
	}
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file) + "/"
	return path
}

//
// 加载头部并处理
// 替换掉里面的一些公共变量
// 给关键字加高亮 <body language="markup">
//
func (t *Template) readHeader() {
	header, err := ioutil.ReadFile(t.getAppPath() + "template/header.html")
	if err == nil {
		header = bytes.Replace(header, []byte("{author}"), []byte(Conf.Author), -1)
		// 菜单
		menuTemplate := []byte(`<li class="menu-item">
                        <a class="menu-item-link" href="{rootpath}$1">$2</a>
                    </li>`)
		mobileMenuTemplate := []byte(`<a class="mobile-menu-item" href="{rootpath}$1">$2</a>`)
		menus := make([][]byte, 0)
		mobileMenus := make([][]byte, 0)
		for i := 0; i < len(Conf.Menu); i++ {
			item := Conf.Menu[i]
			url1 := []byte(Conf.SaveDir + item.Key + "/index.html")
			item1 := bytes.Replace(menuTemplate, []byte("$1"), url1, 1)
			item1 = bytes.Replace(item1, []byte("$2"), []byte(item.Value), 1)
			menus = append(menus, item1)
			item2 := bytes.Replace(mobileMenuTemplate, []byte("$1"), url1, 1)
			item2 = bytes.Replace(item2, []byte("$2"), []byte(item.Value), 1)
			mobileMenus = append(mobileMenus, item2)
		}
		header = bytes.Replace(header, []byte("{menu}"), bytes.Join(menus, []byte("\n")), -1)
		header = bytes.Replace(header, []byte("{mobile-menu}"), bytes.Join(mobileMenus, []byte("\n")), -1)
		// 加上关键字高亮
		header = bytes.Replace(header, []byte("<body>"), []byte("<body class=\"language-markup\">"), 1)
	}
	t.header = header
}

//
// 加载脚部并处理
// 替换掉里面的一些公共变量
//
func (t *Template) readFooter() {
	footer, err := ioutil.ReadFile(t.getAppPath() + "template/footer.html")
	if err == nil {
		footer = bytes.Replace(footer, []byte("{author}"), []byte(Conf.Author), -1)
		footer = bytes.Replace(footer, []byte("{mailto}"), []byte(Conf.Email), -1)
		footer = bytes.Replace(footer, []byte("{github}"), []byte(Conf.Github), -1)
	}
	t.footer = footer
}

//
// 加载主页
//
func (t *Template) readIndex() {
	index, err := ioutil.ReadFile(t.getAppPath() + "template/index.html")
	if err == nil {
		index = bytes.Replace(index, []byte("{header}"), t.header, 1)
		index = bytes.Replace(index, []byte("{footer}"), t.footer, 1)
		t.Index = index
	}
}

//
// 加载文章
//
func (t *Template) readArticle() {
	postHeader, _ := ioutil.ReadFile(t.getAppPath() + "template/post_header.html")
	postFooter, _ := ioutil.ReadFile(t.getAppPath() + "template/post_footer.html")
	article, err := ioutil.ReadFile(t.getAppPath() + "template/article.html")
	if err != nil {
		fmt.Println("load template/article.html", err)
	}
	article = bytes.Replace(article, []byte("{header}"), t.header, 1)
	article = bytes.Replace(article, []byte("{footer}"), t.footer, 1)
	article = bytes.Replace(article, []byte("{post-header}"), postHeader, 1)
	article = bytes.Replace(article, []byte("{post-footer}"), postFooter, 1)
	t.Article = article
}

//
// 加载关于
//
func (t *Template) readAbout() {
	article, err := ioutil.ReadFile(t.getAppPath() + "template/article.html")
	if err == nil {
		article = bytes.Replace(article, []byte("{header}"), t.header, 1)
		article = bytes.Replace(article, []byte("{footer}"), t.footer, 1)
		article = bytes.Replace(article, []byte("{post-header}"), []byte(""), 1)
		article = bytes.Replace(article, []byte("{post-footer}"), []byte(""), 1)
		t.About = article
	}
}

//
// 加载归档
//
func (t *Template) readArchive() {
	archive, err := ioutil.ReadFile(t.getAppPath() + "template/archive.html")
	if err == nil {
		archive = bytes.Replace(archive, []byte("{header}"), t.header, 1)
		archive = bytes.Replace(archive, []byte("{footer}"), t.footer, 1)
		t.Archive = archive
	}
}

//
// 加载标签
//
func (t *Template) readTag() {
	tag, err := ioutil.ReadFile(t.getAppPath() + "template/tag.html")
	if err == nil {
		tag = bytes.Replace(tag, []byte("{header}"), t.header, 1)
		tag = bytes.Replace(tag, []byte("{footer}"), t.footer, 1)
		t.Tag = tag
	}
}

//
// CreateTemplate 建立 Template 单例
//
func CreateTemplate() {
	Tepl = &Template{}
	Tepl.readHeader()
	Tepl.readFooter()
	Tepl.readIndex()
	Tepl.readArticle()
	Tepl.readAbout()
	Tepl.readArchive()
	Tepl.readTag()
}
