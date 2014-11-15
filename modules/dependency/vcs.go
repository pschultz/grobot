package dependency

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type metaImport struct {
	Prefix   string
	VCS      string
	RepoRoot string
}

// repoRootForImportDynamic finds a *repoRoot for a custom domain that's not
// statically known by repoRootForImportPathStatic.
//
// This handles "vanity import paths" like "name.tld/pkg/foo".
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/vcs.go
func repoRootForImportDynamic(importPath string) (string, error) {
	slash := strings.Index(importPath, "/")
	if slash < 0 {
		return "", errors.New("import path doesn't contain a slash")
	}
	host := importPath[:slash]
	if !strings.Contains(host, ".") {
		return "", errors.New("import path doesn't contain a hostname")
	}
	urlStr, body, err := httpsOrHTTP(importPath)
	if err != nil {
		return "", fmt.Errorf("http/https fetch: %v", err)
	}
	defer body.Close()
	imports, err := parseMetaGoImports(body)
	if err != nil {
		return "", fmt.Errorf("parsing %s: %v", importPath, err)
	}
	metaImport, err := matchGoImport(imports, importPath)
	if err != nil {
		if err != errNoMatch {
			return "", fmt.Errorf("parse %s: %v", urlStr, err)
		}
		return "", fmt.Errorf("parse %s: no go-import meta tags", urlStr)
	}

	log.Debug("Found meta tag %v at %s", metaImport, urlStr)

	// If the import was "uni.edu/bob/project", which said the
	// prefix was "uni.edu" and the RepoRoot was "evilroot.com",
	// make sure we don't trust Bob and check out evilroot.com to
	// "uni.edu" yet (possibly overwriting/preempting another
	// non-evil student).  Instead, first verify the root and see
	// if it matches Bob's claim.
	if metaImport.Prefix != importPath {
		log.Debug("Verifying non-authoritative meta tag", importPath)
		urlStr0 := urlStr
		urlStr, body, err = httpsOrHTTP(metaImport.Prefix)
		if err != nil {
			return "", fmt.Errorf("fetch %s: %v", urlStr, err)
		}
		imports, err := parseMetaGoImports(body)
		if err != nil {
			return "", fmt.Errorf("parsing %s: %v", importPath, err)
		}
		if len(imports) == 0 {
			return "", fmt.Errorf("fetch %s: no go-import meta tag", urlStr)
		}
		metaImport2, err := matchGoImport(imports, importPath)
		if err != nil || metaImport != metaImport2 {
			return "", fmt.Errorf("%s and %s disagree about go-import for %s", urlStr0, urlStr, metaImport.Prefix)
		}
	}

	if !strings.Contains(metaImport.RepoRoot, "://") {
		return "", fmt.Errorf("%s: invalid repo root %q; no scheme", urlStr, metaImport.RepoRoot)
	}

	return metaImport.RepoRoot, nil
}

// httpsOrHTTP returns the body of either the importPath's
// https resource or, if unavailable, the http resource.
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/http.go
func httpsOrHTTP(importPath string) (urlStr string, body io.ReadCloser, err error) {
	fetch := func(scheme string) (urlStr string, res *http.Response, err error) {
		u, err := url.Parse(scheme + "://" + importPath)
		if err != nil {
			return "", nil, err
		}
		u.RawQuery = "go-get=1"
		urlStr = u.String()
		log.Debug("Fetching %s", urlStr)
		res, err = grobot.GetHTTP(urlStr)
		return
	}
	closeBody := func(res *http.Response) {
		if res != nil {
			res.Body.Close()
		}
	}
	urlStr, res, err := fetch("https")
	if err != nil || res.StatusCode != 200 {
		if err != nil {
			log.Debug("https fetch failed.")
		} else {
			log.Debug("Ignoring https fetch with status code %d", res.StatusCode)
		}
		closeBody(res)
		urlStr, res, err = fetch("http")
	}
	if err != nil {
		closeBody(res)
		return "", nil, err
	}
	// Note: accepting a non-200 OK here, so people can serve a
	// meta import in their http 404 page.
	log.Debug("Parsing meta tags from %s (status code %d)", urlStr, res.StatusCode)
	return urlStr, res.Body, nil
}

// charsetReader returns a reader for the given charset. Currently
// it only supports UTF-8 and ASCII. Otherwise, it returns a meaningful
// error which is printed by go get, so the user can find why the package
// wasn't downloaded if the encoding is not supported. Note that, in
// order to reduce potential errors, ASCII is treated as UTF-8 (i.e. characters
// greater than 0x7f are not rejected).
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/discovery.go
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "ascii":
		return input, nil
	default:
		return nil, fmt.Errorf("can't decode XML document using charset %q", charset)
	}
}

// parseMetaGoImports returns meta imports from the HTML in r.
// Parsing ends at the end of the <head> section or the beginning of the <body>.
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/discovery.go
func parseMetaGoImports(r io.Reader) (imports []metaImport, err error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var t xml.Token
	for {
		t, err = d.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			return
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			return
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}
		if attrValue(e.Attr, "name") != "go-import" {
			continue
		}
		if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 3 {
			imports = append(imports, metaImport{
				Prefix:   f[0],
				VCS:      f[1],
				RepoRoot: f[2],
			})
		}
	}
}

// attrValue returns the attribute value for the case-insensitive key
// `name', or the empty string if nothing is found.
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/discovery.go
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

// errNoMatch is returned from matchGoImport when there's no applicable match.
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/vcs.go
var errNoMatch = errors.New("no import match")

// matchGoImport returns the metaImport from imports matching importPath.
// An error is returned if there are multiple matches.
// errNoMatch is returned if none match.
//
// DISCLAIMER: This code is basically copied from /usr/lib/golang/src/cmd/go/vcs.go
func matchGoImport(imports []metaImport, importPath string) (_ metaImport, err error) {
	match := -1
	for i, im := range imports {
		if !strings.HasPrefix(importPath, im.Prefix) {
			continue
		}
		if match != -1 {
			err = fmt.Errorf("multiple meta tags match import path %q", importPath)
			return
		}
		match = i
	}
	if match == -1 {
		err = errNoMatch
		return
	}
	return imports[match], nil
}
