package ao3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
)

const (
	Title        = `h2.title`
	Author       = `h3.byline a`
	Series       = `dd.series .position`
	Comments     = `.preface .summary .userstuff`
	Ships        = `dd.relationship a`
	Tags         = `dd.freeform a`
	Fandom       = `dd.fandom a`
	Pubdate      = `dd.published`
	ListLink     = `li.work h4.heading a:first-of-type`
	RelatedWorks = `ul.associations li a`
	Downloads    = `li.download ul li a`
)

func GetString(sel string, val *string) chromedp.Action {
	return chromedp.Action(chromedp.TextContent(
		sel,
		val,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func GetNodes(sel string, nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		sel,
		nodes,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func GetAllNodes(sel string, nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		sel,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
}

func GetInnerHTML(sel string, val *string) chromedp.Action {
	return chromedp.Action(chromedp.InnerHTML(
		sel,
		val,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func getSeries(ctx context.Context, book *cdb.Book) {
	var s string
	err := chromedp.Run(ctx,
		chromedp.TextContent(
			Series,
			&s,
			chromedp.ByQuery,
			chromedp.NodeReady,
			chromedp.AtLeast(0),
		),
	)
	if err != nil {
		fmt.Errorf("%w %w\n", scrapeErr("series"), err)
		return
	}

	title, pos, err := parseSeriesText(s)
	if err != nil {
		fmt.Errorf("series parsing err: %w\n", err)
	}

	book.Series = title
	book.SeriesIndex = pos
}

var NoTitleErr = errors.New("no title")

func parseSeriesText(s string) (string, float64, error) {
	var (
		title string
		pos   float64
	)
	if s == "" {
		return title, pos, NoTitleErr
	}
	seriesRegexp := regexp.MustCompile(`Part (?P<pos>\d+) of (?P<name>.*)`)
	matches := seriesRegexp.FindStringSubmatch(s)

	if len(matches) < 1 {
		return title, pos, fmt.Errorf("no matches for %s\n", s)
	}

	pos = cast.ToFloat64(matches[seriesRegexp.SubexpIndex("pos")])
	title = matches[seriesRegexp.SubexpIndex("name")]
	return title, pos, nil
}

func parsePubdate(pubdate string) time.Time {
	t, err := time.Parse(time.DateOnly, pubdate)
	if err != nil {
		t = time.Now()
	}
	return t
}

func parseFormats(nodes []*cdp.Node) []string {
	formats := make([]string, len(nodes))
	for i, node := range nodes {
		t := node.AttributeValue("href")
		formats[i] = ParseUrl(t).String()
	}
	return formats
}

func parseRelated(nodes []*cdp.Node) []string {
	var rels []string
	for _, node := range nodes {
		if rel := node.AttributeValue("rel"); rel == "author" {
			rels = append(rels, node.Children[0].NodeValue)
		}
	}
	return rels
}

func parseTags(nodes ...[]*cdp.Node) []string {
	var tags []string
	for _, n := range nodes {
		tags = append(tags, getFirstChildValues(n)...)
	}
	return tags
}

func getFirstChildValues(nodes []*cdp.Node) []string {
	var vals []string
	for _, node := range nodes {
		var t string
		if len(node.Children) > 0 {
			t = node.Children[0].NodeValue
		}
		vals = append(vals, t)
	}
	return vals
}

func DownloadWork(u, name string) {
	response, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal(response.StatusCode)
	}

	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func scrapeErr(name string) error {
	return fmt.Errorf("error scraping %s from %s\n", name, CurrentURL())
}
