package rpurl

import (
	"errors"
	"net/url"
)

func IsRelative(urlIn string) (isRelative bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	return u.Host == "", nil
}

func IsHttpScheme(urlIn string) (isHttpScheme bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	// Assume lack of a URL scheme implies some form of HTTP
	return u.Scheme == "" || u.Scheme == "http" || u.Scheme == "https", nil
}

func AddBaseHost(baseHost string, urlPath string) (urlOut string, err error) {
	b, err := url.Parse(baseHost)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}

	u.Scheme = b.Scheme
	u.Host = b.Host
	u.User = b.User

	result, err := url.QueryUnescape(u.String())
	if err != nil {
		return "", err
	}

	return result, nil
}

func AbsoluteHttpUrl(baseUrl string, href string) (url string, err error) {
	isRelative, err := IsRelative(href)
	if err != nil {
		return "", err
	}

	isHttpScheme, err := IsHttpScheme(href)
	if err != nil {
		return "", err
	}

	if !isHttpScheme {
		return "", errors.New("Invalid URL " + href + ".")
	}

	if isRelative {
		href, err = AddBaseHost(baseUrl, href)
		if err != nil {
			return "", err
		}
	}

	return href, nil
}

func AbsoluteHttpUrls(baseUrl string, hrefs []string) (urls []string, errs []error) {
	for _, href := range hrefs {
		url, err := AbsoluteHttpUrl(baseUrl, href)
		if err != nil {
			errs = append(errs, err)
		} else {
			urls = append(urls, url)
		}
	}

	return urls, errs
}
