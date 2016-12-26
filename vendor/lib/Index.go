package lib

import (
	"bytes"
	"fmt"
	"strconv"
)

//
// Index 首页
//
type Index struct {
	articleList []*Article
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
// 开始处理
//
func (ind *Index) init() {
	count := len(ind.articleList)
	// 每页数量
	ind.countPerPage = Conf.IndexCountPerPage
	// 页数
	ind.totalPages = (count-1)/ind.countPerPage + 1
	ind.page = 1
	for i := 0; i < count; i++ {
		ind.handleArticle(i)
		// 检测是否完成一页
		count := i + 1
		if count%ind.countPerPage == 0 || count == len(ind.articleList) {
			ind.gen()
		}
	}
}

//
// 处理一篇文章，添加一条
//
func (ind *Index) handleArticle(index int) {
	// 项模板html
	itemTemplate := []byte(`<article class="post">
    <header class="post-header">
      <h1 class="post-title">
          <a class="post-link" href="{rootpath}$url">$title</a>
      </h1>
      <div class="post-meta">
        <span class="post-time">
          $time
        </span>
      </div>
    </header>
    <div class="post-content">
        $content
        <div class="read-more">
            <a href="{rootpath}$url" class="read-more-link">Read more..</a>
        </div>
    </div>
  </article>`)
	article := ind.articleList[index]
	str := bytes.Replace(itemTemplate, []byte("$url"), article.URL, -1)
	str = bytes.Replace(str, []byte("$title"), article.Title, 1)
	str = bytes.Replace(str, []byte("$time"), []byte(article.Date.Format("2006-01-02 15:04")), 1)
	str = bytes.Replace(str, []byte("$content"), article.Preveiw, -1)
	ind.items = append(ind.items, str)
}

//
// 生成新的一页
//
func (ind *Index) gen() {
	var url string
	if ind.page == 1 {
		// 第一页
		url = "index.html"
	} else {
		url = Conf.SaveDir + "page/" + strconv.Itoa(ind.page) + "/index.html"
	}
	rootpath := GetRootPath(url)
	nav := HandleNav(ind.totalPages, ind.page, "index")
	data := bytes.Replace(Tepl.Index, []byte("{page-title}"), []byte(Conf.Title), 1)
	data = bytes.Replace(data, []byte("{items}"), bytes.Join(ind.items, []byte("\n")), 1)
	data = bytes.Replace(data, []byte("{nav}"), nav, 1)
	data = bytes.Replace(data, []byte("{rootpath}"), rootpath, -1)
	WriteFile(url, data)
	fmt.Println("create:", url)
	ind.items = make([][]byte, 0)
	ind.page++
	// html文件数
	TotalFileCount++
}

//
// HandleIndex 开始处理主页
//
func HandleIndex(list []*Article) {
	ind := &Index{articleList: list}
	ind.init()
}
