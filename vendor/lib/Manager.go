package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TotalMarkdownCount 总md文件数
var TotalMarkdownCount int

// TotalFileCount 总生成html文件数
var TotalFileCount int

//
// Manager 文章管理
//
type Manager struct {
	// 文章列表
	ArticleList []*Article
	// tag列表
	tags []TagItem
	// category列表
	categories []TagItem
	// 开始时间
	startTime time.Time
}

//
// 读取所有的文章
//
func (m *Manager) readArticles(dir string) {
	// 记录开始时间
	m.startTime = time.Now()
	err1 := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		// 是markdown文件
		if strings.HasSuffix(strings.ToLower(fi.Name()), ".md") {
			if strings.ToLower(fi.Name()) == "about.md" {
				NewAbout(filename)
			} else {
				m.ArticleList = append(m.ArticleList, NewArticle(filename))
			}
			// 文章总数+1
			TotalMarkdownCount++
			TotalFileCount++
			// fmt.Println(filename)
		}
		return nil
	})
	if err1 != nil {
		fmt.Println(err1)
	}
}

//
// 处理列表，排序，生成归档，标签，分类
//
func (m *Manager) handleList() {
	// 排序
	sort.Sort(By{list: m.ArticleList})
	count := len(m.ArticleList)
	// tag和category列表
	m.tags = make([]TagItem, 0)
	m.categories = make([]TagItem, 0)
	for i := 0; i < count; i++ {
		article := m.ArticleList[i]
		// 标签
		for j := 0; j < len(article.Tags); j++ {
			m.addTag(string(article.Tags[j]), i)
		}
		// 分类
		m.addCategory(string(article.Category), i)
	}
	// 开始处理标签，分类
	HandleArchives(m.ArticleList)
	HandleTags(m.ArticleList, m.tags)
	HandleCategories(m.ArticleList, m.categories)
	// 处理首页
	HandleIndex(m.ArticleList)
}

//
// 生成data.js
//
func (m *Manager) handleDataJs() {
	items := make([][]byte, 0)
	// 文章列表 [{title,url}]
	type ArticleItem struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	}
	articles := make([]ArticleItem, 0)
	for i := 0; i < len(m.ArticleList); i++ {
		ar := m.ArticleList[i]
		articles = append(articles, ArticleItem{string(ar.URL), string(ar.Title)})
	}
	str, _ := json.Marshal(articles)
	items = append(items, []byte("var data_articles="+string(str)+";"))
	// tag
	tags, _ := json.Marshal(m.tags)
	items = append(items, []byte("var data_tags="+string(tags)+";"))
	// category
	categories, _ := json.Marshal(m.categories)
	items = append(items, []byte("var data_categories="+string(categories)+";"))

	url := Conf.SaveDir + "data.js"
	data := bytes.Join(items, []byte("\n"))
	WriteFile(url, data)
	fmt.Println("create:", url)
}

//
// 显示生成信息
//
func (m *Manager) showInfo() {
	// 结束时间
	endTime := time.Now()
	// 耗时
	t := endTime.UnixNano() - m.startTime.UnixNano()
	t /= 1000
	tf := float64(t) / 1000000.0
	fmt.Println("")
	fmt.Println("--------build info--------")
	fmt.Println("article count:", TotalMarkdownCount)
	fmt.Println("create html count:", TotalFileCount)
	fmt.Println("time:", strconv.FormatFloat(tf, 'f', 3, 32)+"s")
	fmt.Println("--------------------------")
}

//
// 处理标签，分类
//
func (m *Manager) handleTag(article *Article) {
	for i := 0; i < len(article.Tags); i++ {
		tag := string(article.Tags[i])
		for j := 0; j < len(m.tags); j++ {
			if m.tags[j].Key == tag {
				m.tags[j].Value = append(m.tags[j].Value, i)
			}
		}
	}
}

//
// 添加一个标签
//
func (m *Manager) addTag(tag string, index int) {
	for i := 0; i < len(m.tags); i++ {
		if m.tags[i].Key == tag {
			// 已经有此标签
			m.tags[i].Value = append(m.tags[i].Value, index)
			return
		}
	}
	// 还没有这个标签
	item := TagItem{tag, []int{index}}
	m.tags = append(m.tags, item)
	// fmt.Println("new tag", tag)
}

//
// 添加一个分类
//
func (m *Manager) addCategory(category string, index int) {
	for i := 0; i < len(m.categories); i++ {
		if m.categories[i].Key == category {
			// 已经有此标签
			m.categories[i].Value = append(m.categories[i].Value, index)
			return
		}
	}
	// 还没有这个标签
	item := TagItem{category, []int{index}}
	m.categories = append(m.categories, item)
}

//
// StartManager 开始管理
//
func StartManager() {
	// 先删除所有文件
	os.RemoveAll(Conf.SaveDir)
	m := &Manager{}
	m.readArticles(Conf.PostDir)
	m.handleList()
	m.handleDataJs()
	m.showInfo()
}

// By 排序用
type By struct {
	list []*Article
}

// Len
func (b By) Len() int {
	return len(b.list)
}

// Swap
func (b By) Swap(i, j int) {
	b.list[i], b.list[j] = b.list[j], b.list[i]
}

// Less i是否小于j
func (b By) Less(i, j int) bool {
	date1 := b.list[i].Date
	date2 := b.list[j].Date
	return date1.Unix() > date2.Unix()
}

//
// HandleNav 处理页码 archive tag category中用到
// @param totalPages 总页数
// @param page 当前页码
// @param flag archive tag/xxx category/xxx index
//
func HandleNav(totalPages int, page int, flag string) []byte {
	// 前后页码模板html
	prevTemplate := []byte(`<a class="prev" href="{rootpath}$url">
        <i class="iconfont icon-left"></i>
        <span class="prev-text">Prev</span>
      </a>`)
	nextTemplate := []byte(`
      <a class="next" href="{rootpath}$url">
        <span class="next-text">Next</span>
        <i class="iconfont icon-right"></i>
      </a>`)
	nav := make([][]byte, 0)
	// 第一个和后面的
	var str1, str2 string
	if flag == "index" {
		str1 = "index.html"
		str2 = Conf.SaveDir + "page/$page/index.html"
	} else {
		str1 = Conf.SaveDir + flag + "/index.html"
		str2 = Conf.SaveDir + flag + "/page/$page/index.html"
	}
	if totalPages > 1 {
		nav = append(nav, []byte("<nav class=\"pagination\">"))
		if page > 1 {
			// 有前一页
			var prevURL []byte
			if page == 2 {
				prevURL = []byte(str1)
			} else {
				prevURL = []byte(strings.Replace(str2, "$page", strconv.Itoa(page-1), 1))
			}
			nav = append(nav, bytes.Replace(prevTemplate, []byte("$url"), prevURL, 1))
		}
		if page < totalPages {
			// 有后一页
			nextURL := []byte(strings.Replace(str2, "$page", strconv.Itoa(page+1), 1))
			nav = append(nav, bytes.Replace(nextTemplate, []byte("$url"), nextURL, 1))
		}
		nav = append(nav, []byte("</nav"))
	}
	return bytes.Join(nav, []byte("\n"))
}
