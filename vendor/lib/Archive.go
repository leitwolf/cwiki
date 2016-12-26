package lib

import (
	"bytes"
	"fmt"
	"strconv"
)

//
// Archive 归档列表处理
//
type Archive struct {
	articleList []*Article
	// 总页数
	totalPages int
	// 当前页数
	page int
	// 每页条目数
	countPerPage int
	// 当前所属年份，不同年份要加个标题，新一页要重新设定
	year int
	// html条目列表，新一页的时候要清空
	items [][]byte
}

//
// 开始处理
//
func (a *Archive) init() {
	count := len(a.articleList)
	// 每页数量
	a.countPerPage = Conf.CountPerPage
	// 页数
	a.totalPages = (count-1)/a.countPerPage + 1
	a.page = 1
	a.year = -1
	for i := 0; i < count; i++ {
		a.handleArticle(i)
		// 检测是否完成一页
		count := i + 1
		if count%a.countPerPage == 0 || count == len(a.articleList) {
			a.gen()
		}
	}
}

//
// 处理一篇文章，添加一条
//
func (a *Archive) handleArticle(index int) {
	// 总共多少条模板html
	allTemplate := []byte(`<div class="archive-title">
          <h2 class="archive-name"> All </h2>
          <span class="archive-post-counter">
            $count Posts In Total
          </span>
        </div>`)
	// 年份模板html
	yearTemplate := []byte(`<div class="collection-title">
            <h2 class="archive-year">$year</h2>
          </div>`)
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
	article := a.articleList[index]
	if index == 0 {
		// 总的第一条，要加个all
		str := bytes.Replace(allTemplate, []byte("$count"), []byte(strconv.Itoa(len(a.articleList))), 1)
		a.items = append(a.items, str)
	}
	y := article.Date.Year()
	if y != a.year {
		// 年份不同了，加标题
		a.year = y
		str := bytes.Replace(yearTemplate, []byte("$year"), []byte(strconv.Itoa(y)), 1)
		a.items = append(a.items, str)
	}
	// 添加项目
	str := bytes.Replace(itemTemplate, []byte("$date"), []byte(article.Date.Format("01-02")), 1)
	str = bytes.Replace(str, []byte("$url"), article.URL, 1)
	str = bytes.Replace(str, []byte("$title"), article.Title, 1)
	a.items = append(a.items, str)
}

//
// 生成新的一页
//
func (a *Archive) gen() {
	var url string
	if a.page == 1 {
		// 第一页
		url = Conf.SaveDir + "archive/index.html"
	} else {
		url = Conf.SaveDir + "archive/page/" + strconv.Itoa(a.page) + "/index.html"
	}
	rootpath := GetRootPath(url)
	nav := HandleNav(a.totalPages, a.page, "archive")
	data := bytes.Replace(Tepl.Archive, []byte("{page-title}"), GetPageTitle([]byte("归档")), 1)
	data = bytes.Replace(data, []byte("{items}"), bytes.Join(a.items, []byte("\n")), 1)
	data = bytes.Replace(data, []byte("{nav}"), nav, 1)
	data = bytes.Replace(data, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, data)
	fmt.Println("create:", url)
	a.items = make([][]byte, 0)
	a.page++
	a.year = -1
	// html文件数
	TotalFileCount++
}

//
// HandleArchives 开始处理归档
//
func HandleArchives(list []*Article) {
	archive := &Archive{articleList: list}
	archive.init()
}
