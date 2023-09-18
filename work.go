package ao3

import (
	"context"
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

func getPubdate(ctx context.Context, book *cdb.Book) {
	var pubdate string
	err := chromedp.Run(ctx,
		chromedp.Text(
			selPubdate,
			&pubdate,
			chromedp.ByQuery,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}

	t, err := time.Parse(time.DateOnly, pubdate)
	if err != nil {
		t = time.Now()
	}

	book.Pubdate = t
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
		log.Println(err)
		return
	}

	seriesRegexp := regexp.MustCompile(`Part (?P<pos>\d+) of (?P<name>.*)`)
	matches := seriesRegexp.FindStringSubmatch(s)
	book.SeriesIndex = cast.ToFloat64(matches[seriesRegexp.SubexpIndex("pos")])
	book.Series = matches[seriesRegexp.SubexpIndex("name")]
}

func getFormats(ctx context.Context, book *cdb.Book) {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(
			selFormats,
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}
	for _, node := range nodes {
		t := node.AttributeValue("href")
		book.Formats = append(book.Formats, ParseUrl(t).String())
	}
}

func getRelated(ctx context.Context, book *cdb.Book) {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(
			selRelated,
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}
	for _, node := range nodes {
		//fmt.Printf("%+V\n", node)
		if rel := node.AttributeValue("rel"); rel == "author" {
			book.Authors = append(book.Authors, node.Children[0].NodeValue)
		}
	}
}

func getContributors(ctx context.Context, book *cdb.Book, isPodfic bool) {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(
			selAuthor,
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}
	if isPodfic {
		book.Narrators = append(book.Narrators, getFirstChildValues(nodes)...)
		getRelated(ctx, book)
		return
	}
	book.Authors = append(book.Authors, getFirstChildValues(nodes)...)
}

func getTags(ctx context.Context, book *cdb.Book) {
	for _, sel := range []string{selRel, selTags, selFandom} {
		var nodes []*cdp.Node
		err := chromedp.Run(ctx,
			chromedp.Nodes(
				sel,
				&nodes,
				chromedp.ByQueryAll,
				chromedp.NodeReady,
			),
		)
		if err != nil {
			log.Println(err)
			return
		}
		book.Tags = append(book.Tags, getFirstChildValues(nodes)...)
	}
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
