package pixiv

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

var (
	//hibiapi地址
	pixivProxyUrl = "https://hibiapi.dengfenglai.icu/api/pixiv/"
)

type IRS struct {
	I struct {
		ID        int64    `json:"id"`
		Title     string   `json:"title"`
		Type      string   `json:"type"`
		ImageUrls ImageUrl `json:"image_urls"`
		Caption   string   `json:"caption"`
		Restrict  int      `json:"restrict"`
		User      struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			Account          string `json:"account"`
			ProfileImageUrls struct {
				Medium string `json:"medium"`
			} `json:"profile_image_urls"`
			IsFollowed bool `json:"is_followed"`
		} `json:"user"`
		Tags []struct {
			Name           string `json:"name"`
			TranslatedName string `json:"translated_name"`
		} `json:"tags"`
		Tools          []interface{} `json:"tools"`
		CreateDate     time.Time     `json:"create_date"`
		PageCount      int           `json:"page_count"`
		Width          int           `json:"width"`
		Height         int           `json:"height"`
		SanityLevel    int           `json:"sanity_level"`
		XRestrict      int           `json:"x_restrict"`
		MetaSinglePage struct {
			OriginalImageURL string `json:"original_image_url"`
		} `json:"meta_single_page"`
		MetaPages []struct {
			ImageUrls ImageUrl `json:"image_urls"`
		} `json:"meta_pages"`
		TotalView            int  `json:"total_view"`
		TotalBookmarks       int  `json:"total_bookmarks"`
		IsBookmarked         bool `json:"is_bookmarked"`
		Visible              bool `json:"visible"`
		IsMuted              bool `json:"is_muted"`
		TotalComments        int  `json:"total_comments"`
		IllustAiType         int  `json:"illust_ai_type"`
		IllustBookStyle      int  `json:"illust_book_style"`
		CommentAccessControl int  `json:"comment_access_control"`
	} `json:"illust"`
}
type ImageUrl struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
	Original     string `json:"original"`
}

type Illust struct {
	Pid      int64
	UID      int
	UserName string
	R18      bool
	Title    string
	URLs     []string
}

// Works 获取插画信息
func Works(id int64) (i *Illust, err error) {
	data, err := web.GetData(pixivProxyUrl+"illust?id="+strconv.FormatInt(id, 10), web.RandUA())
	if err != nil {
		return nil, err
	}
	ir := new(IRS)
	if err = json.Unmarshal(data, ir); err != nil {
		return nil, err
	}
	urls := []string{}
	if ir.I.MetaSinglePage.OriginalImageURL != "" {
		urls = append(urls, ir.I.MetaSinglePage.OriginalImageURL)
	}
	for _, v := range ir.I.MetaPages {
		urls = append(urls, v.ImageUrls.Original)
	}
	return &Illust{
		Pid:      ir.I.ID,
		UID:      ir.I.User.ID,
		UserName: ir.I.User.Name,
		R18:      ir.I.XRestrict > 0,
		Title:    ir.I.Title,
		URLs:     urls,
	}, nil

}

func DownLoadWorks(url, path string) error {
	url = strings.ReplaceAll(url, "i.pximg.net", "i.pixiv.re")
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, http.MethodGet, "https://www.pixiv.net/", web.RandUA(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0666)
}
