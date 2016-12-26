package lib

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
)

// TagItem tag项{"xxx",[1,2,3]}
type TagItem struct {
	Key   string `json:"type"`
	Value []int  `json:"list"`
}

//
// Tag 标签，分类处理
//
type Tag struct {
	// 名称 tag category
	name string
	// 标题 标签 分类
	title string
	// 所有文章列表
	articleList []*Article
	// 标签列表 "aaa":[1,2,3]，后面int是文章列表序号
	tags []TagItem
	// 当前标签
	curTag string
	// 当前标签包含的文章数
	curTagLength int
	// 总页数
	totalPages int
	// 当前页数
	page int
	// 每页条目数
	countPerPage int
	// html条目列表，新一页的时候要清空
	items [][]byte
}

//
// 开始
//
func (t *Tag) init() {
	t.countPerPage = Conf.CountPerPage

	url := Conf.SaveDir + t.name + "/index.html"
	rootpath := GetRootPath(url)
	// 条目模板html
	itemTemplate := []byte(`<a href="{rootpath}$url" style="font-size: $sizepx;">$title</a>`)
	items := make([][]byte, 0)
	for i := 0; i < len(t.tags); i++ {
		key := t.tags[i].Key
		value := t.tags[i].Value
		// fmt.Println("tag", key, value)
		t.handleTag(key, value)
		url1 := Conf.SaveDir + t.name + "/" + key + "/index.html"
		size := rand.Intn(8) + 15
		strs := bytes.Replace(itemTemplate, []byte("$url"), []byte(url1), 1)
		strs = bytes.Replace(strs, []byte("$size"), []byte(strconv.Itoa(size)), 1)
		strs = bytes.Replace(strs, []byte("$title"), []byte(key), 1)
		items = append(items, strs)
	}
	// 生成标签云
	var total string
	if t.name == "tag" {
		total = " Tags"
	} else {
		total = " Categories"
	}
	data := bytes.Replace(Tepl.Tag, []byte("{page-title}"), GetPageTitle([]byte(t.title)), 1)
	data = bytes.Replace(data, []byte("{total}"), []byte(strconv.Itoa(len(items))+total), 1)
	data = bytes.Replace(data, []byte("{tags}"), bytes.Join(items, []byte("\n")), 1)
	data = bytes.Replace(data, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, data)
	fmt.Println("create:", url)
	// html文件数
	TotalFileCount++
}

//
// 处理一个标签
//
func (t *Tag) handleTag(tagName string, indexList []int) {
	t.curTag = tagName
	t.curTagLength = len(indexList)
	t.totalPages = (t.curTagLength-1)/t.countPerPage + 1
	t.page = 1
	t.items = make([][]byte, 0)
	for i := 0; i < t.curTagLength; i++ {
		index := indexList[i]
		t.handleArticle(index)
		// 检测是否完成一页
		count := i + 1
		// fmt.Println(t.curTag, count, t.countPerPage, t.curTagLength)
		if count%t.countPerPage == 0 || count == t.curTagLength {
			t.gen()
		}
	}
}

//
// 处理一篇文章条目
//
func (t *Tag) handleArticle(index int) {
	// 项模板html
	itemTemplate := []byte(`<div class="archive-post">
        <span class="archive-post-time">
        $date
        </span>
        <span class="archive-post-title">
          <a href="{rootpath}$url" class="archive-post-link">
            $title
          </a>
        </span>
      </div>`)
	article := t.articleList[index]
	// 添加项目
	url := string(article.URL) + "?" + t.name + "=" + t.curTag
	str := bytes.Replace(itemTemplate, []byte("$date"), []byte(article.Date.Format("2006-01-02")), 1)
	str = bytes.Replace(str, []byte("$url"), []byte(url), 1)
	str = bytes.Replace(str, []byte("$title"), article.Title, 1)
	t.items = append(t.items, str)
}

//
// 生成新的一页
//
func (t *Tag) gen() {
	tagTitleTemplate := []byte(`<div class="archive-title tag">
          <h2 class="archive-name">$tag</h2>
        </div>`)
	var url string
	if t.page == 1 {
		// 第一页
		url = Conf.SaveDir + t.name + "/" + t.curTag + "/index.html"
	} else {
		url = Conf.SaveDir + t.name + "/" + t.curTag + "/page/" + strconv.Itoa(t.page) + "/index.html"
	}
	rootpath := GetRootPath(url)
	nav := HandleNav(t.totalPages, t.page, t.name+"/"+t.curTag)
	// 加入标题
	items := make([][]byte, 0)
	items = append(items, bytes.Replace(tagTitleTemplate, []byte("$tag"), []byte(t.curTag), 1))
	items = append(items, t.items[0:]...)
	t.items = items
	data := bytes.Replace(Tepl.Archive, []byte("{page-title}"), GetPageTitle([]byte(t.curTag)), 1)
	data = bytes.Replace(data, []byte("{items}"), bytes.Join(t.items, []byte("\n")), 1)
	data = bytes.Replace(data, []byte("{nav}"), nav, 1)
	data = bytes.Replace(data, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, data)
	fmt.Println("create:", url)
	t.items = make([][]byte, 0)
	t.page++
	// html文件数
	TotalFileCount++
}

//
// HandleTags 处理标签
//
func HandleTags(articleList []*Article, tags []TagItem) {
	tag := &Tag{articleList: articleList, tags: tags}
	tag.name = "tag"
	tag.title = "标签"
	tag.init()
}

//
// HandleCategories 处理分类
//
func HandleCategories(articleList []*Article, categories []TagItem) {
	tag := &Tag{articleList: articleList, tags: categories}
	tag.name = "category"
	tag.title = "分类"
	tag.init()
}
