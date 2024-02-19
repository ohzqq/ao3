package ao3

import (
	"context"
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
	selTitle    = `h2.title`
	selAuthor   = `h3.byline a`
	selSeries   = `dd.series .position`
	selComments = `.preface .summary .userstuff`
	selRel      = `dd.relationship a`
	selTags     = `dd.freeform a`
	selFandom   = `dd.fandom a`
	selPubdate  = `dd.published`
	selListLink = `li.work h4.heading a:first-of-type`
	selRelated  = `ul.associations li a`
	selFormats  = `li.download ul li a`
)

func getTitle(title *string) chromedp.Action {
	return chromedp.Action(chromedp.Text(
		selTitle,
		title,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func getComments(comments *string) chromedp.Action {
	return chromedp.Action(chromedp.InnerHTML(
		selComments,
		comments,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func getPubdate(pubdate *string) chromedp.Action {
	return chromedp.Action(chromedp.Text(
		selPubdate,
		pubdate,
		chromedp.ByQuery,
		chromedp.NodeReady,
	))
}

func parsePubdate(pubdate string) time.Time {
	t, err := time.Parse(time.DateOnly, pubdate)
	if err != nil {
		t = time.Now()
	}
	return t
}

func getSeries(ctx context.Context, book *cdb.Book) {
	var s string
	err := chromedp.Run(ctx,
		chromedp.Text(
			selSeries,
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

	seriesRegexp := regexp.MustCompile(`Part (?P<pos>\d+) of (?P<name>.*)`)
	matches := seriesRegexp.FindStringSubmatch(s)
	book.SeriesIndex = cast.ToFloat64(matches[seriesRegexp.SubexpIndex("pos")])
	book.Series = matches[seriesRegexp.SubexpIndex("name")]
}

func getFormats(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selFormats,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
}

func parseFormats(nodes []*cdp.Node) []string {
	formats := make([]string, len(nodes))
	for i, node := range nodes {
		t := node.AttributeValue("href")
		formats[i] = ParseUrl(t).String()
	}
	return formats
}

func getRelated(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selRelated,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
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

func getContributors(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selAuthor,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
}

func getShips(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selRel,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
}

func getFandom(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selFandom,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
}

func getTags(nodes *[]*cdp.Node) chromedp.Action {
	return chromedp.Action(chromedp.Nodes(
		selTags,
		nodes,
		chromedp.ByQueryAll,
		chromedp.NodeReady,
	))
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
