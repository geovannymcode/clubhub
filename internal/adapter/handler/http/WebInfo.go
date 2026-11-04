package http

import (
	"bytes"
	"context"
	"net/url"

	"github.com/Geovanny0401/clubhub/internal/security"
	"github.com/PuerkitoBio/goquery"
)

func GetTitleAndLogo(ctx context.Context, address string) (string, string, error) {
	addressComplete := "https://" + address + "/"
	body, err := security.Fetch(ctx, addressComplete)
	if err != nil {
		return "", "", err
	}
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}

	pageTitle := document.Find("title").Contents().First().Text()
	var pageLogo string

	document.Find("link").Each(func(index int, element *goquery.Selection) {
		value, exists := element.Attr("type")
		if exists && validateType(value) {
			valueHref, _ := element.Attr("href")
			pageLogo = valueHref
		}
	})

	document.Find("meta").Each(func(index int, element *goquery.Selection) {
		value, exists := element.Attr("itemprop")
		if exists && "image" == value {
			valueContent, _ := element.Attr("content")
			pageLogo = valueContent
		}
	})

	base, _ := url.Parse(addressComplete)
	logo, parseErr := url.Parse(pageLogo)
	pageLogo = ""
	if parseErr == nil && logo != nil {
		logo = base.ResolveReference(logo)
		if logo.Scheme == "https" && logo.User == nil {
			pageLogo = logo.String()
		}
	}

	return pageTitle, pageLogo, err
}

func validateType(value string) bool {
	types := []string{"image/x-icon", "image/icon", "image/vnd.microsoft.icon", "image/svg+xml", "image/png", "image/jpg"}

	for _, element := range types {
		if element == value {
			return true
		}
	}

	return false
}
