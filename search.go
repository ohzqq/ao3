package ao3

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
)

func Search(u string) ([]cdb.Book, error) {
	sUrl := ParseUrl(u)
	params := sUrl.Query()
	for _, k := range SearchParams() {
		if params.Has(k) && params.Get(k) == "" {
			params.Del(k)
		}
	}
	sUrl.RawQuery = params.Encode()

	return parseList(sUrl)
}

func SortAndFilter(u string) ([]cdb.Book, error) {
	sUrl := ParseUrl(u)
	params := sUrl.Query()
	for _, k := range SortAndFilterParams() {
		if params.Has(k) && params.Get(k) == "" {
			params.Del(k)
		}
	}
	sUrl.RawQuery = params.Encode()

	return parseList(sUrl)
}

func parseList(u *url.URL) ([]cdb.Book, error) {
	var works []cdb.Book

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		setCookies(u.String()),
	)
	if err != nil {
		return works, err
	}

	total := getTotalPages(ctx, u.String())

	params := u.Query()
	for i := 1; i <= total; i++ {
		page := strconv.Itoa(i)
		params.Set("page", page)
		u.RawQuery = params.Encode()
		w, err := Page(u.String())
		if err != nil {
			return works, err
		}
		works = append(works, w...)
	}

	return works, nil
}

func getTotalPages(ctx context.Context, u string) int {
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		Sleep(5*time.Second),
		chromedp.Navigate(u),
		chromedp.Nodes(
			".pagination > li > a",
			&nodes,
			chromedp.ByQueryAll,
			chromedp.NodeReady,
			chromedp.AtLeast(0),
		),
	)
	if err != nil {
		log.Println(err)
		return 1
	}

	child := getFirstChildValues(nodes)
	if len(child) > 0 {
		return cast.ToInt(child[len(child)-2])
	}

	return 1
}

func SearchParams() []string {
	return []string{
		searchQuery,
		searchTitle,
		searchCreators,
		searchRevisedAt,
		searchComplete,
		searchCrossover,
		searchSingleChapter,
		searchWordCount,
		searchLangId,
		searchFandomNames,
		searchRatingIDs,
		searchCharNames,
		searchRelNames,
		searchTags,
		searchHits,
		searchKudosCount,
		searchCommentsCount,
		searchBookmarksCount,
		searchSortCol,
		searchSortDirection,
	}
}

func SortAndFilterParams() []string {
	return []string{
		filterSortCol,
		filterOtherTags,
		filterExcludedTags,
		filterCrossover,
		filterComplete,
		filterWordsFrom,
		filterWordsTo,
		filterDateFrom,
		filterDateTo,
		filterQuery,
		filterLangID,
		filterTagID,
		includeRatingIDs,
		includeWarningIDs,
		includeCatIDs,
		includeFandomIDs,
		includeCharIDs,
		includeRelIDs,
		includeFreeformIDs,
		excludeRatingIDs,
		excludeWarningIDs,
		excludeCatIDs,
		excludeFandomIDs,
		excludeCharIDs,
		excludeRelIDs,
		excludeFreeformIDs,
	}
}

const (
	searchQuery          = `work_search[query]`
	searchTitle          = `work_search[title]`
	searchCreators       = `work_search[creators]`
	searchRevisedAt      = `work_search[revised_at]`
	searchComplete       = `work_search[complete]`
	searchCrossover      = `work_search[crossover]`
	searchSingleChapter  = `work_search[single_chapter]`
	searchWordCount      = `work_search[word_count]`
	searchLangId         = `work_search[language_id]`
	searchFandomNames    = `work_search[fandom_names]`
	searchRatingIDs      = `work_search[rating_ids]`
	searchCharNames      = `work_search[character_names]`
	searchRelNames       = `work_search[relationship_names]`
	searchTags           = `work_search[freeform_names]`
	searchHits           = `work_search[hits]`
	searchKudosCount     = `work_search[kudos_count]`
	searchCommentsCount  = `work_search[comments_count]`
	searchBookmarksCount = `work_search[bookmarks_count]`
	searchSortCol        = `work_search[sort_column]`
	searchSortDirection  = `work_search[sort_direction]`
)

const (
	filterCommit       = `commit=Sort and Filter`
	filterSortCol      = `work_search[sort_column]`
	filterOtherTags    = `work_search[other_tag_names]`
	filterExcludedTags = `work_search[excluded_tag_names]`
	filterCrossover    = `work_search[crossover]`
	filterComplete     = `work_search[complete]`
	filterWordsFrom    = `work_search[words_from]`
	filterWordsTo      = `work_search[words_to]`
	filterDateFrom     = `work_search[date_from]`
	filterDateTo       = `work_search[date_to]`
	filterQuery        = `work_search[query]`
	filterLangID       = `work_search[language_id]`
	filterTagID        = `tag_id`
	includeRatingIDs   = `include_work_search[rating_ids][]`
	includeWarningIDs  = `include_work_search[archive_warning_ids][]`
	includeCatIDs      = `include_work_search[category_ids][]`
	includeFandomIDs   = `include_work_search[fandom_ids][]`
	includeCharIDs     = `include_work_search[character_ids][]`
	includeRelIDs      = `include_work_search[relationship_ids][]`
	includeFreeformIDs = `include_work_search[freeform_ids][]`
	excludeRatingIDs   = `exclude_work_search[rating_ids][]`
	excludeWarningIDs  = `exclude_work_search[archive_warning_ids][]`
	excludeCatIDs      = `exclude_work_search[category_ids][]`
	excludeFandomIDs   = `exclude_work_search[fandom_ids][]`
	excludeCharIDs     = `exclude_work_search[character_ids][]`
	excludeRelIDs      = `exclude_work_search[relationship_ids][]`
	excludeFreeformIDs = `exclude_work_search[freeform_ids][]`
)
