package parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/pemistahl/lingua-go"
	"io"
	"net/http"
	"strings"
)

func GetLangsInSitemap(r io.Reader) ([]string, error) {
	return nil, errors.New("not implemented")
}

func GetLangsInTags(r io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)

	doc.Find("link[hreflang]").Each(func(i int, selection *goquery.Selection) {
		lang, _ := selection.Attr("hreflang")
		res = append(res, lang)
	})

	lang, exists := doc.Find("html").Attr("lang")
	if exists {
		res = append(res, lang)
	}

	return res, nil
}

func GetContentLang(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	languages := []lingua.Language{
		lingua.Russian,
		lingua.English,
		lingua.French,
		lingua.German,
		lingua.Spanish,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	doc.Find("code,script,style").Remove()
	txt := doc.Find("body").Text()
	confidenceValues := detector.ComputeLanguageConfidenceValues(txt)

	for _, elem := range confidenceValues {
		if elem.Value() == 1 {
			return elem.Language().IsoCode639_1().String(), nil
		}
	}
	return "", nil
}

func getLangFromHTTPHeaders(resp *http.Response) []string {
	contentLanguage := strings.Split(
		strings.ReplaceAll(
			resp.Header.Get("Content-Language"),
			" ",
			"",
		),
		",",
	)
	//todo add link header https://developers.google.com/search/docs/advanced/crawling/localized-versions?hl=ru#html
	return contentLanguage
}
