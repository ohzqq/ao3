package ao3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	cookiemonster "github.com/MercuryEngineering/CookieMonster"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	userAgent string = `user-agent=churkeybot/1.0 (+https://archiveofourown.org/users/churkey/profile)`
	delay     int64  = 5000000000
	ao3Host   string = `archiveofourown.org`
)

func Scrape(u string) ([]cdb.Book, error) {
	var works []cdb.Book

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		setCookies(u),
	)
	if err != nil {
		return works, err
	}

	work, err := GetWork(ctx, u)
	if err != nil {
		return works, err
	}
	works = append(works, work)

	return works, nil
}

func Page(u string) ([]cdb.Book, error) {
	var works []cdb.Book

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		setCookies(u),
	)
	if err != nil {
		return works, err
	}

	links := GetLinkList(ctx, u)
	for _, link := range links {
		work, err := GetWork(ctx, link)
		if err != nil {
			return works, err
		}
		works = append(works, work)
	}

	return works, nil
}

func GetWork(ctx context.Context, u string) (cdb.Book, error) {
	viper.Set("url", u)

	var (
		work    cdb.Book
		pubdate string
		formats []*cdp.Node
		tags    []*cdp.Node
		fandom  []*cdp.Node
		ships   []*cdp.Node
		con     []*cdp.Node
		rel     []*cdp.Node
	)

	actions := []chromedp.Action{
		Sleep(5 * time.Second),
		chromedp.Navigate(u),
		GetInnerHTML(Comments, &work.Comments),
		GetString(Title, &work.Title),
		GetString(Pubdate, &pubdate),
		GetAllNodes(Downloads, &formats),
		GetNodes(Tags, &tags),
		GetNodes(Ships, &ships),
		GetNodes(Fandom, &fandom),
		GetNodes(Author, &con),
	}

	if IsPodfic() {
		actions = append(actions, GetAllNodes(RelatedWorks, &rel))
	}

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		return work, err
	}

	work.Pubdate = parsePubdate(pubdate)
	work.Formats = parseFormats(formats)
	work.Tags = parseTags(tags, ships, fandom)

	var auth []string
	if len(con) > 0 {
		auth = append(auth, getFirstChildValues(con)...)
	}
	if IsPodfic() {
		work.Narrators = auth
	} else {
		work.Authors = auth
	}
	if len(rel) > 0 {
		work.Authors = parseRelated(rel)
	}

	getSeries(ctx, &work)

	return work, nil
}

func GetLinkList(ctx context.Context, u string) []string {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		Sleep(1*time.Second),
		chromedp.Navigate(u),
		chromedp.Nodes(
			ListLink,
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		fmt.Errorf("%w %w\n", scrapeErr("link list"), err)
		return []string{}
	}

	var links []string
	for _, node := range nodes {
		t := ParseUrl(node.AttributeValue("href"))
		links = append(links, t.String())
	}
	return links
}

func ParseUrl(u string) *url.URL {
	pu, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	vals := pu.Query()

	vals.Set("view_full_work", "true")
	vals.Set("view_adult", "true")
	pu.RawQuery = vals.Encode()

	if pu.Scheme == "" {
		pu.Scheme = "https"
	}

	if pu.Host == "" {
		pu.Host = ao3Host
	}

	return pu
}

func setCookies(u string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		cparams := make([]*network.CookieParam, len(Cookies()))
		for i, c := range Cookies() {
			t := cdp.TimeSinceEpoch(c.Expires)
			cparams[i] = &network.CookieParam{
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Path:     c.Path,
				Secure:   c.Secure,
				HTTPOnly: c.HttpOnly,
				SameSite: network.CookieSameSite(cast.ToString(c.SameSite)),
				Expires:  &t,
				URL:      u,
			}
		}
		// add cookies to chrome
		err := network.SetCookies(cparams).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func Cookies() []*http.Cookie {
	cfg, _ := os.UserConfigDir()
	file := path.Join(cfg, "ur", "cookies.txt")

	cookies, err := cookiemonster.ParseFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return cookies
}
