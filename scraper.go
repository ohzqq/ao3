package ao3

import (
	"context"
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

	work := GetWork(ctx, u)
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
		work := GetWork(ctx, link)
		works = append(works, work)
	}

	return works, nil
}

func GetWork(ctx context.Context, u string) cdb.Book {
	var work cdb.Book

	err := chromedp.Run(ctx,
		Sleep(5*time.Second),
		chromedp.Navigate(u),
		getTitle(&work.Title),
		getComments(&work.Comments),
	)
	if err != nil {
		log.Fatal(err)
	}

	getTags(ctx, &work)
	getContributors(ctx, &work)
	getPubdate(ctx, &work)
	getSeries(ctx, &work)
	getFormats(ctx, &work)

	return work
}

func GetLinkList(ctx context.Context, u string) []string {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		Sleep(5*time.Second),
		chromedp.Navigate(u),
		chromedp.Nodes(
			selListLink,
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
		),
	)
	if err != nil {
		log.Println(err)
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
