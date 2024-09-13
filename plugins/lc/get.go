package lc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/tidwall/gjson"
)

var (
	Csrftoken     string
	lcUrl         = "https://leetcode.cn/graphql"
	gettopicparam = `
	query questionOfToday {
      todayRecord {
        date
        userStatus
        question {
          questionId
          frontendQuestionId: questionFrontendId
          difficulty
          title
          titleCn: translatedTitle
          titleSlug
          paidOnly: isPaidOnly
          freqBar
          isFavor
          acRate
          status
          solutionNum
          hasVideoSolution
          topicTags {
            name
            nameTranslated: translatedName
            id
          }
          extra {
            topCompanyTags {
              imgUrl
              slug
              numSubscribed
            }
          }
        }
        lastSubmission {
          id
        }
      }
    }`
	getonetopicparam = `
	query questionData($titleSlug: String!) {
          question(titleSlug: $titleSlug) {
            questionId
            questionFrontendId
            categoryTitle
            boundTopicId
            title
            titleSlug
            content
            translatedTitle
            translatedContent
            isPaidOnly
            difficulty
            likes
            dislikes
            isLiked
            similarQuestions
            contributors {
              username
              profileUrl
              avatarUrl
              __typename
            }
            langToValidPlayground
            topicTags {
              name
              slug
              translatedName
              __typename
            }
            companyTagStats
            codeSnippets {
              lang
              langSlug
              code
              __typename
            }
            stats
            hints
            solution {
              id
              canSeeDetail
              __typename
            }
            status
            sampleTestCase
            metaData
            judgerAvailable
            judgeType
            mysqlSchemas
            enableRunCode
            envInfo
            book {
              id
              bookName
              pressName
              source
              shortDescription
              fullDescription
              bookImgUrl
              pressImgUrl
              productUrl
              __typename
            }
            isSubscribed
            isDailyQuestion
            dailyRecordStatus
            editorType
            ugcQuestionId
            style
            exampleTestcases
            jsonExampleTestcases
            __typename
          }
        }`
	gettopicListparam = `
	query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
  problemsetQuestionList(
    categorySlug: $categorySlug
    limit: $limit
    skip: $skip
    filters: $filters
  ) {
    questions {
      acRate
      difficulty
      titleCn
      titleSlug
      }
}
}`
)

func GetCsrftoken() error {
	req, err := http.Get(lcUrl)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	for _, v := range req.Cookies() {
		if v.Name == "csrftoken" {
			Csrftoken = v.Value
			return nil
		}
	}
	return errors.New("get Csrftoken err: no Csrftoken in cookies")
}

func GetTodayTopic() (g gjson.Result, err error) {
	dataJson := map[string]any{
		"query":     gettopicparam,
		"variables": map[string]any{},
	}
	data, _ := json.Marshal(dataJson)
	data, err = web.Web(web.NewDefaultClient(), lcUrl, http.MethodPost, makeHeard, bytes.NewReader(data))
	g = gjson.ParseBytes(data)
	return
}

func GetOneTopic(title string) (g gjson.Result, err error) {
	dataJson := map[string]any{
		"operationName": "questionData",
		"variables": map[string]any{
			"titleSlug": title,
		},
		"query": getonetopicparam,
	}
	data, _ := json.Marshal(dataJson)
	data, err = web.Web(web.NewDefaultClient(), lcUrl, http.MethodPost, makeHeard, bytes.NewReader(data))
	g = gjson.ParseBytes(data)
	return
}
func GetTopicList(num int, difficult string) (g gjson.Result, err error) {
	filtersmap := map[string]any{}
	if difficult != "" {
		filtersmap["difficulty"] = difficult
	}
	dataJson := map[string]any{
		"operationName": "problemsetQuestionList",
		"variables": map[string]any{
			"categorySlug": "",
			"skip":         num,
			"limit":        50,
			"filters":       filtersmap,
		},
		"query": gettopicListparam,
	}
	data, _ := json.Marshal(dataJson)
	data, err = web.Web(web.NewDefaultClient(), lcUrl, http.MethodPost, makeHeard, bytes.NewReader(data))
	g = gjson.ParseBytes(data)
	return
}

func makeHeard(request *http.Request) {
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	request.Header.Set("accept", "*/*")
	request.Header.Set("User-Agent", web.RandUA())
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Origin", "https://leetcode.cn")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Csrftoken", Csrftoken)
}

func ProcessContent(onetopic gjson.Result) string {
	content := onetopic.Get("data.question.translatedContent").String()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
	var text bytes.Buffer
	doc.Contents().Each(func(i int, s *goquery.Selection) {
		text.WriteString(s.Text())
	})
	sptext := strings.Split(text.String(), "\n")
	text.Reset()
	for _, v := range sptext {
		if t := strings.TrimSpace(v); t != "" {
			text.WriteString(t)
			text.WriteByte('\n')
		}
	}
	return strings.TrimSpace(text.String())
}
