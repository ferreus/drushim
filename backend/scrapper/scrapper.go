package scrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	RegionsKey        = "איזורים:"
	RegionsKeyEn      = "Locations:"
	SkilsKey          = "תחומים:"
	SkilsKeyEn        = "Sub Categories:"
	ExperienceKey     = "שנות נסיון:"
	ExperienceKeyEn   = "Required Experience:"
	RequirementsKey   = "דרישות התפקיד:"
	RequirementsKeyEn = "Job Requirements:"
	NextPageKey       = "הבא »"
	FullPageLinkKey   = "fullPage"
)

type JobItem struct {
	Name         string
	DatePosted   string
	Experience   string
	Requirements []string
	Regions      []string
	Skils        []string
	JobCode      string
	Link         string
}

func (job JobItem) String() string {
	jsonTxt, err := json.Marshal(job)
	if err != nil {
		return fmt.Sprintf("(Invalid Json):%s", job.Name)
	}
	return fmt.Sprintf("%s", jsonTxt)
}

func getTextByClass(n *html.Node, className string) string {
	text, ok := scrape.Find(n, scrape.ByClass(className))
	if ok {
		return scrape.Text(text)
	}
	return ""
}

func textFieldMatcher(n *html.Node) bool {
	if n.DataAtom == atom.Div {
		c := scrape.Attr(n, "class")
		return strings.Contains(c, "fieldContainer")
	}
	return false
}

func Fetch(url string) ([]JobItem, string, error) {
	nextPage := ""
	var jobItems []JobItem

	resp, err := http.Get(url)
	if err != nil {
		return jobItems, nextPage, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("fetch: Error reading %s: %v\n", url, err)
		return jobItems, nextPage, err
	}
	resp.Body.Close()
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Div {
			c := scrape.Attr(n, "class")
			return strings.Contains(c, "jobItem")
		}
		return false
	}
	items := scrape.FindAll(root, matcher)
	jobItems = make([]JobItem, len(items))
	for i, article := range items {
		title := getTextByClass(article, "jobName")
		date := getTextByClass(article, "jobDate")
		jobItems[i].Name = title
		jobItems[i].DatePosted = date
		links := scrape.FindAll(article, scrape.ByTag(atom.A))
		for _, l := range links {
			jobCode := scrape.Attr(l, "jobcode")
			if len(jobCode) > 0 {
				jobItems[i].JobCode = jobCode
			}
			if scrape.Attr(l, "class") == FullPageLinkKey {
				jobItems[i].Link = scrape.Attr(l, "href")
			}
		}
		textFields := scrape.FindAll(article, textFieldMatcher)
		for _, tf := range textFields {
			key := getTextByClass(tf, "fieldTitle")
			value := getTextByClass(tf, "fieldText")
			switch key {
			case ExperienceKey, ExperienceKeyEn:
				jobItems[i].Experience = value
			case SkilsKey, SkilsKeyEn:
				jobItems[i].Skils = strings.Split(value, ",")
			case RegionsKey, RegionsKeyEn:
				jobItems[i].Regions = strings.Split(value, ",")
			case RequirementsKey, RequirementsKeyEn:
				jobItems[i].Requirements = strings.Split(value, ",")

			default:
				log.Printf("Unmatched key:[%s], value:[%s]\n", key, value)
			}
		}
	}
	links := scrape.FindAll(root, scrape.ByTag(atom.A))
	for _, l := range links {
		c := scrape.Attr(l, "class")
		if strings.Contains(c, "pager") {
			if scrape.Text(l) == NextPageKey {
				nextPage = scrape.Attr(l, "href")
			}
		}
	}
	return jobItems, nextPage, nil
}
