package ao3

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/cdb"
	"github.com/spf13/viper"
)

const (
	testWork      = `https://archiveofourown.org/works/49186696`
	testPage      = `https://archiveofourown.org/series/1331351`
	testPodfic    = `https://archiveofourown.org/works/49186696`
	testSearch    = `https://archiveofourown.org/works/search?work_search%5Bquery%5D=&work_search%5Btitle%5D=&work_search%5Bcreators%5D=&work_search%5Brevised_at%5D=&work_search%5Bcomplete%5D=T&work_search%5Bcrossover%5D=&work_search%5Bsingle_chapter%5D=0&work_search%5Bword_count%5D=&work_search%5Blanguage_id%5D=en&work_search%5Bfandom_names%5D=Teen+Wolf+%28TV%29&work_search%5Brating_ids%5D=&work_search%5Bcharacter_names%5D=&work_search%5Brelationship_names%5D=Derek+Hale%2FStiles+Stilinski%2CDerek+Hale%2FPeter+Hale&work_search%5Bfreeform_names%5D=&work_search%5Bhits%5D=&work_search%5Bkudos_count%5D=&work_search%5Bcomments_count%5D=&work_search%5Bbookmarks_count%5D=&work_search%5Bsort_column%5D=_score&work_search%5Bsort_direction%5D=desc&commit=Search`
	testSearchAll = `https://archiveofourown.org/works/search?work_search%5Bquery%5D=&work_search%5Btitle%5D=&work_search%5Bcreators%5D=&work_search%5Brevised_at%5D=&work_search%5Bcomplete%5D=&work_search%5Bcrossover%5D=F&work_search%5Bsingle_chapter%5D=0&work_search%5Bword_count%5D=%3E1&work_search%5Blanguage_id%5D=en&work_search%5Bfandom_names%5D=Teen+Wolf+%28TV%29&work_search%5Brating_ids%5D=13&work_search%5Barchive_warning_ids%5D%5B%5D=14&work_search%5Barchive_warning_ids%5D%5B%5D=16&work_search%5Bcategory_ids%5D%5B%5D=21&work_search%5Bcategory_ids%5D%5B%5D=23&work_search%5Bcharacter_names%5D=Danny+M%C4%81healani&work_search%5Brelationship_names%5D=Derek+Hale%2FStiles+Stilinski%2CDerek+Hale%2FPeter+Hale&work_search%5Bfreeform_names%5D=Fluff&work_search%5Bhits%5D=%3E2&work_search%5Bkudos_count%5D=%3E3&work_search%5Bcomments_count%5D=%3C4000000&work_search%5Bbookmarks_count%5D=%3E5&work_search%5Bsort_column%5D=_score&work_search%5Bsort_direction%5D=desc&commit=Search`
)

func TestSearch(t *testing.T) {
	t.SkipNow()
	s, err := Search(testSearch)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", s)

	ns, err := Search(testSearchAll)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", ns)
}

func TestPage(t *testing.T) {
	books, err := Page(testPage)
	if err != nil {
		t.Error(err)
	}

	if len(books) < 1 {
		t.Error("no books returned, expected at least 1")
	}

	for k, v := range books[0].StringMapString() {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestPodfic(t *testing.T) {
	viper.Set("is-podfic", true)
	books, err := Scrape(testPodfic)
	if err != nil {
		t.Error(err)
	}

	if len(books) < 1 {
		t.Error("no books returned, expected at least 1")
	}

	for k, v := range books[0].StringMapString() {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestWork(t *testing.T) {
	books, err := Scrape(testWork)
	if err != nil {
		t.Error(err)
	}

	if len(books) < 1 {
		t.Error("no books returned, expected at least 1")
	}

	for k, v := range books[0].StringMapString() {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestWorkCmd(t *testing.T) {
	books, err := Scrape(testWork)
	if err != nil {
		t.Error(err)
	}

	for _, book := range books {
		if book.Title == "" {
			t.Error("no title")
		}
		err := book.Save(casing.Snake(book.Title)+".yaml", true)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestReadMeta(t *testing.T) {
	files, err := filepath.Glob("testdata/*")
	if err != nil {
		t.Error(err)
	}
	for _, file := range files {
		dec, err := cdb.ReadMetadataFile(file)
		if err != nil {
			t.Error(err)
		}
		var book cdb.Book
		err = dec.Decode(&book)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%#v\n", book)
	}
}
