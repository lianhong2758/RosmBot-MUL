package lc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/tidwall/gjson"
)

var (
	Csrftoken     string
	lcUrl         = "https://leetcode.cn"
	graphqlUrl    = "/graphql"
	gettopicparam = `query questionOfToday {
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
	getonetopicparam = `query questionData($titleSlug: String!) {
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
)

func GetCsrftoken() error {
	req, err := http.Get(lcUrl + graphqlUrl)
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
	data, err = web.Web(web.NewDefaultClient(), lcUrl+graphqlUrl, http.MethodPost, makeHeard, bytes.NewReader(data))
	g = gjson.ParseBytes(data)
	return
}

func GetOneTopic(tg gjson.Result) (g gjson.Result, err error) {
	title := tg.Get("data.todayRecord.0.question.titleSlug").String()
	dataJson := map[string]any{
		"operationName": "questionData",
		"variables": map[string]any{
			"titleSlug": title,
		},
		"query": getonetopicparam,
	}
	data, _ := json.Marshal(dataJson)
	data, err = web.Web(web.NewDefaultClient(), lcUrl+graphqlUrl, http.MethodPost, makeHeard, bytes.NewReader(data))
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
