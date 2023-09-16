package ao3

import (
	"errors"
	"log"
	"net/url"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

type Query struct {
	URI     string
	ID      string
	IsAudio bool
	Query   url.Values
	path    []string
	*url.URL
	*Host
}

type Host struct {
	sub    string
	domain string
	tld    string
}

type Path struct {
	paths []string
}

func New(uri ...string) *Query {
	if len(uri) > 0 {
		return ParseURL(uri[0])
	}
	return &Query{
		URL: &url.URL{
			Scheme: "https",
		},
		Query: make(url.Values),
		Host:  &Host{},
	}
}

func ParseURL(uri string) *Query {
	q := &Query{}
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	q.URL = u
	q.Query = u.Query()
	q.AppendPath(strings.Split(q.URL.Path, "/")...)
	q.SetHost(q.URL.Host)
	return q
}

func (q *Query) SetHost(host string) *Query {
	h, err := ParseHost(host)
	if err != nil {
		log.Printf("%v\n", err)
	}
	q.Host = h
	return q
}

func (q *Query) AppendPath(paths ...string) *Query {
	q.path = append(q.path, paths...)
	return q
}

func (q *Query) PrependPath(path string) *Query {
	q.path = slices.Insert(q.path, 0, path)
	return q
}

func (q *Query) SetParam(name, val string) *Query {
	q.Query.Set(name, val)
	return q
}

func (q *Query) AddParam(name, val string) *Query {
	q.Query.Add(name, val)
	return q
}

func (q *Query) String() string {
	q.URL.Path = path.Join(q.path...)
	q.URL.Host = q.Host.String()
	q.RawQuery = q.Query.Encode()
	return q.URL.String()
}

func ParseHost(hs string) (*Host, error) {
	host := &Host{}
	h := strings.Split(hs, ".")
	switch len(h) {
	case 3:
		host.SubDomain(h[0])
		host.Domain(h[1])
		host.TLD(h[2])
		return host, nil
	case 2:
		host.Domain(h[0])
		host.TLD(h[1])
		return host, nil
	default:
		return host, errors.New("invalid host")
	}
}

func (h *Host) SubDomain(s string) *Host {
	h.sub = s
	return h
}

func (h *Host) Domain(s string) *Host {
	h.domain = s
	return h
}

func (h *Host) TLD(s string) *Host {
	h.tld = s
	return h
}

func (h Host) String() string {
	var host []string
	if h.sub != "" {
		host = append(host, h.sub)
	}
	host = append(host, h.domain, h.tld)
	s := strings.TrimPrefix(strings.Join(host, "."), ".")
	if s == "." {
		return ""
	}
	return s
}
