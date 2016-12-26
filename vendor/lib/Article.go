package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"path/filepath"

	"strings"

	"strconv"

	"github.com/russross/blackfriday"
)

//
// Article 对markdown文件进行处理
//
type Article struct {
	// 路径
	path string
	// 输入数据（从文件读取）
	input []byte
	// 生成的html数据
	data []byte
	// 预览的html数据
	Preveiw []byte
	// ---文件信息
	filename []byte
	Title    []byte
	Date     time.Time
	Category []byte
	Tags     [][]byte
	// ---文件信息 end
	// toc，已生成的html代码
	toc []byte
	// 文章url
	URL []byte
}

//
// 读取文件到input
//
func (a *Article) read(path string) {
	a.path = path
	a.input, _ = ioutil.ReadFile(path)
}

//
// 读取文件头并从input中删除文件头，文件头格式：
// ---
// filename: xxx
// title: JavaScript 标题文字信息
// date: 2016-11-30 18:52:21 或 2016/11/30 18:52
// category: js
// tags: JavaScript,前端
// ---
//
func (a *Article) readHeader() {
	reg := regexp.MustCompile("-{3,}([\\s\\S]*)-{3,}")
	reg = regexp.MustCompile("-{3,}\\r?\\n")
	indexes := reg.FindAllIndex(a.input, 2)
	// fmt.Println(indexs)
	if len(indexes) != 2 || indexes[0][0] != 0 {
		a.filename = []byte(strings.Replace(filepath.Base(a.path), filepath.Ext(a.path), "", 1))
		a.Title = a.filename
		return
	}
	// 最后的index
	endIndex := indexes[1][1]
	// 两个---之间的数据
	data := a.input[0:endIndex]
	// 同时把这一段信息从输入数据里去掉
	if len(a.input) > endIndex {
		a.input = a.input[endIndex:]
	} else {
		a.input = make([]byte, 0)
	}
	// fmt.Println(len(a.input))

	// 生成文件名
	reg = regexp.MustCompile("filename:(.+)\\n")
	bytes2 := reg.FindSubmatch(data)
	if len(bytes2) >= 2 {
		a.filename = bytes.Trim(bytes.Trim(bytes2[1], " "), "\r")
	}
	// fmt.Println("filename", string(a.filename))
	if a.filename == nil || string(a.filename) == "" {
		a.filename = []byte(strings.Replace(filepath.Base(a.path), filepath.Ext(a.path), "", 1))
	}

	// 标题
	reg = regexp.MustCompile("title:(.+)\\n")
	bytes2 = reg.FindSubmatch(data)
	if len(bytes2) >= 2 {
		a.Title = bytes.Trim(bytes.Trim(bytes2[1], " "), "\r")
	}

	// 时间 2016-11-30 18:52:21 或 2016-11-30 18:52
	reg = regexp.MustCompile("date:(.+)\\n")
	bytes2 = reg.FindSubmatch(data)
	if len(bytes2) >= 2 {
		t := bytes.Trim(bytes.Trim(bytes2[1], " "), "\r")
		date, err := time.Parse("2006-01-02 15:04:05", string(t))
		if err != nil {
			date, err = time.Parse("2006-01-02 15:04", string(t))
			if err != nil {
				date, err = time.Parse("2006/01/02 15:04:05", string(t))
				if err != nil {
					date, _ = time.Parse("2006/01/02 15:04", string(t))
				}
			}
		}
		a.Date = date
		// fmt.Println("time", string(t), a.date.String())
	} else {
		a.Date = time.Now()
	}

	// 目录
	reg = regexp.MustCompile("category:(.+)\\n")
	bytes2 = reg.FindSubmatch(data)
	if len(bytes2) >= 2 {
		a.Category = bytes.Trim(bytes.Trim(bytes2[1], " "), "\r")
	}

	// 标签
	// <div class="post-tags">
	// 	<a href="/tags/生活/">生活</a>
	// 	<a href="/tags/总结/">总结</a>
	// </div>
	reg = regexp.MustCompile("tags:(.+)\\n")
	bytes2 = reg.FindSubmatch(data)
	if len(bytes2) >= 2 {
		t := bytes.Trim(bytes.Trim(bytes2[1], " "), "\r")
		tags := bytes.Split(t, []byte(","))
		for i := 0; i < len(tags); i++ {
			t := bytes.Trim(tags[i], " ")
			a.Tags = append(a.Tags, t)
		}
	}
}

//
// 加入原创内容
//
func (a *Article) handleOriginal() {
	url := Conf.Site + "{url}"
	origin := bytes.Replace([]byte(Conf.OriginalTemplate), []byte("{url}"), []byte(url), -1)
	a.input = bytes.Replace(a.input, []byte("{{original}}"), origin, -1)
}

//
// 获取标签html
// <div class="post-tags">
// 	<a href="/tags/生活/">生活</a>
// 	<a href="/tags/总结/">总结</a>
// </div>
//
func (a *Article) getTagsHTML() []byte {
	strs := make([][]byte, 0)
	strs = append(strs, []byte("<div class=\"post-tags\">"))
	tmpl := []byte("	<a href=\"{rootpath}" + Conf.SaveDir + "tag/$tag/index.html\">$tag</a>")
	for i := 0; i < len(a.Tags); i++ {
		t := bytes.Trim(a.Tags[i], " ")
		item := bytes.Replace(tmpl, []byte("$tag"), t, -1)
		strs = append(strs, item)
	}
	strs = append(strs, []byte("</div>"))
	return bytes.Join(strs, []byte("\n"))
}

//
// 处理图片或资源路径
// [xxx](../../file/a.txt) => [xxx]({rootpath}/file/a.txt)
//
func (a *Article) handleRes() {
	reg := regexp.MustCompile("\\]\\((../)+")
	a.input = reg.ReplaceAll(a.input, []byte("]({rootpath}"))
}

//
// 从markdown生成html
//
func (a *Article) md2html() {
	a.data = blackfriday.MarkdownCommon(a.input)
}

//
// 处理代码
// 1、给代码加上行号
// <pre> => <pre class="line-numbers">
//
func (a *Article) handleCode() {
	a.data = bytes.Replace(a.data, []byte("<pre>"), []byte("<pre class=\"line-numbers\">"), -1)
	// reg := regexp.MustCompile("<code>(.+)</code>")
	// a.data = reg.ReplaceAll(a.data, []byte("<code class=\"code\">$1</code>"))
	// a.data = bytes.Replace(a.data, []byte("<body>"), []byte("<body class=\"language-markup\">"), 1)
}

//
// 处理里面的toc h1-h3
// <h2>工厂模式</h2> =>
// <h2 id="工厂模式"><a href="#工厂模式" class="headerlink" title="工厂模式"></a>工厂模式</h2>
//
// 页面右侧toc格式：
// <div class="post-toc" id="post-toc">
//     <h2 class="post-toc-title">Contents</h2>
//     <div class="post-toc-content">
//         <ol class="toc">
//             <li class="toc-item toc-level-1">
//                 <a class="toc-link" href="#标题1"><span class="toc-text">标题1</span></a>
//             </li>
//             <li class="toc-item toc-level-2">
//                 <a class="toc-link" href="#标题2"><span class="toc-text">标题2</span></a>
//             </li>
//         </ol>
//     </div>
// </div>
//
func (a *Article) handleToc() {
	// reg1 := regexp.MustCompile("\\s")
	reg := regexp.MustCompile("<h([123])>(.+)</h[123]>")
	repl := []byte("<h$level id=\"$id\"><a href=\"#$id\" class=\"headerlink\" title=\"$id\"></a>$title</h$level>")
	topItem := &TocItem{level: []byte(strconv.Itoa(0))}
	lastItem := topItem
	a.data = reg.ReplaceAllFunc(a.data, func(b []byte) []byte {
		bs := reg.FindSubmatch(b)
		level := bs[1]
		// fmt.Println("level", string(level))
		title := bs[2]
		// 去掉符号
		id := bytes.Replace(bs[2], []byte("\""), []byte(""), -1)
		id = bytes.Replace(id, []byte("'"), []byte(""), -1)
		id = bytes.Replace(id, []byte("<"), []byte(""), -1)
		id = bytes.Replace(id, []byte(">"), []byte(""), -1)
		item := &TocItem{level: level, id: id, title: title}
		a.addTocItem(lastItem, item)
		lastItem = item

		repl2 := bytes.Replace(repl, []byte("$level"), level, -1)
		repl2 = bytes.Replace(repl2, []byte("$title"), title, -1)
		repl2 = bytes.Replace(repl2, []byte("$id"), id, -1)
		return repl2
	})
	// 生成toc html
	if len(topItem.children) > 0 {
		str := []byte(`<div class="post-toc" id="post-toc">
    <h2 class="post-toc-title">文章目录</h2>
    <div class="post-toc-content">
        <ol class="toc">
			{items}
        </ol>
    </div>
</div>`)
		items := a.handleTocItem(topItem, true)
		a.toc = bytes.Replace(str, []byte("{items}"), items, 1)
	}
}

// TocItem toc项
type TocItem struct {
	level    []byte
	id       []byte
	title    []byte
	children []*TocItem
	parent   *TocItem
}

//
// 添加toc项
//
func (a *Article) addTocItem(lastItem *TocItem, item *TocItem) {
	level1, _ := strconv.Atoi(string(item.level))
	level2, _ := strconv.Atoi(string(lastItem.level))
	// fmt.Println("toc", level1, level2)
	if level1 > level2 {
		// 子项
		item.parent = lastItem
		lastItem.children = append(lastItem.children, item)
	} else {
		a.addTocItem(lastItem.parent, item)
	}
}

//
// 添加右侧toc项
//
func (a *Article) handleTocItem(item *TocItem, isRoot bool) []byte {
	itemStr := []byte(`<li class="toc-item toc-level-$level"><a class="toc-link" href="#$id"><span class="toc-text">$title</span></a>`)
	items := make([][]byte, 0)
	if !isRoot {
		str := bytes.Replace(itemStr, []byte("$level"), item.level, -1)
		str = bytes.Replace(str, []byte("$title"), item.title, -1)
		str = bytes.Replace(str, []byte("$id"), item.id, -1)
		items = append(items, str)
	}
	if len(item.children) > 0 {
		items = append(items, []byte("<ol class=\"toc-child\">"))
		for i := 0; i < len(item.children); i++ {
			items = append(items, a.handleTocItem(item.children[i], false))
		}
		items = append(items, []byte("</ol>"))
	}
	if !isRoot {
		items = append(items, []byte("</li>"))
	}
	return bytes.Join(items, []byte("\n"))
}

//
// 分离出预览
//
func (a *Article) genPreview() {
	reg := regexp.MustCompile("([\\s\\S]*)<!--\\s*more\\s*-->")
	bytes1 := reg.FindSubmatch(a.data)
	if len(bytes1) < 2 {
		// 没有，则全部都是预览
		a.Preveiw = a.data
	} else {
		a.Preveiw = bytes1[1]
		// fmt.Println(string(a.Preveiw))
	}
}

//
// 生成文章
//
func (a *Article) genArticle() {
	scriptTemplate := []byte(`<script language="javascript">
            var page_url="$url";
			var rootpath="{rootpath}";
        </script>
        <script type="text/javascript" src="{rootpath}$save-dirdata.js"></script>
        <script type="text/javascript" src="{rootpath}res/js/nav.js"></script>`)
	// 文章url
	url := Conf.SaveDir + a.Date.Format("2006/01/02") + "/" + string(a.filename) + ".html"
	rootpath := GetRootPath(url)
	// 加入nav script
	script := bytes.Replace(scriptTemplate, []byte("$url"), []byte(url), 1)
	script = bytes.Replace(script, []byte("$save-dir"), []byte(Conf.SaveDir), 1)

	article := Tepl.Article
	article = bytes.Replace(article, []byte("{page-title}"), GetPageTitle(a.Title), 1)
	article = bytes.Replace(article, []byte("{title}"), a.Title, -1)
	article = bytes.Replace(article, []byte("{date}"), []byte(a.Date.Format("2006-01-02 15:04")), 1)
	article = bytes.Replace(article, []byte("{toc}"), a.toc, 1)
	article = bytes.Replace(article, []byte("{content}"), a.data, 1)
	article = bytes.Replace(article, []byte("{tags}"), a.getTagsHTML(), 1)
	article = bytes.Replace(article, []byte("{url}"), []byte(url), -1)
	article = bytes.Replace(article, []byte("{script}"), script, 1)
	article = bytes.Replace(article, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, article)
	a.URL = []byte(url)
	// 给 preview 添加 url
	if a.Preveiw != nil {
		a.Preveiw = bytes.Replace(a.Preveiw, []byte("{url}"), []byte(url), -1)
	}
	fmt.Println("create:", a.path, "=>", url)
}

//
// 生成关于
//
func (a *Article) genAbout() {
	// 文章url
	url := Conf.SaveDir + "about/index.html"
	rootpath := GetRootPath(url)
	a.Title = []byte("关于")
	about := Tepl.About
	about = bytes.Replace(about, []byte("{page-title}"), GetPageTitle(a.Title), 1)
	about = bytes.Replace(about, []byte("{toc}"), []byte(""), 1)
	about = bytes.Replace(about, []byte("{content}"), a.data, 1)
	about = bytes.Replace(about, []byte("{script}"), []byte(""), 1)
	about = bytes.Replace(about, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, about)
	fmt.Println("create:", a.path, "=>", url)
}

//
// NewArticle 新建一篇文章
//
func NewArticle(path string) (article *Article) {
	article = &Article{}
	article.read(path)
	article.readHeader()
	article.handleOriginal()
	article.handleRes()
	article.md2html()
	article.handleCode()
	article.handleToc()
	article.genPreview()
	article.genArticle()
	return
}

//
// NewAbout 新建关于页面
//
func NewAbout(path string) (article *Article) {
	article = &Article{}
	article.read(path)
	article.handleRes()
	article.md2html()
	article.handleCode()
	article.handleToc()
	article.genAbout()
	return
}
