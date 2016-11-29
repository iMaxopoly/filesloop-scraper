package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	pool "gopkg.in/go-playground/pool.v3"
)

type filesloop struct {
	workers   uint
	timeout   uint
	timestamp string

	seAlluc             bool
	seSharedir          bool
	seSearch4shared     bool
	seFilesbug          bool
	seRapid4me          bool
	seGeneralsearch     bool
	seRapidsearchengine bool
	seFilespr           bool
	seGeneralfiles      bool
	seIrfree            bool
	seSoftarchive       bool
	seSceper            bool
	se2ddl              bool

	myclients []myclient
}

type myclient struct {
	fileName string
	data     []string
}

const (
	DEBUGFOLDER   = "./_reports_www.filesloop.com"
	DEBUGFILEPATH = DEBUGFOLDER + "/debug"

	SEALLUC             = "alluc"
	SESHAREDIR          = "sharedir"
	SESEARCH4SHARED     = "search4shared"
	SEFILESBUG          = "filesbug"
	SERAPID4ME          = "rapid4me"
	SEGENERALSEARCH     = "generalsearch"
	SERAPIDSEARCHENGINE = "rapidsearchengine"
	SEFILESPR           = "filespr"
	SEGENERALFILES      = "generalfiles"
	SEIRFREE            = "irfree"
	SESOFTARCHIVE       = "softarchive"
	SESCEPER            = "sceper"
	SE2DDL              = "2ddl"

	APIURL = "https://www.filesloop.com/fileslist/get/"
)

type flRequestResults struct {
	Search   string `json:"search"`
	Server   string `json:"server"`
	Order    string `json:"order"`
	Filesize string `json:"filesize"`
	Ext      string `json:"ext"`
	Page     int    `json:"page"`
	Driver   string `json:"driver"`
}

type responseResults struct {
	Config struct {
		Filesize int    `json:"filesize"`
		HostsID  string `json:"hosts_id"`
		Order    string `json:"order"`
		Page     int    `json:"page"`
		Proxy    string `json:"proxy"`
		Qry      string `json:"qry"`
		Server   int    `json:"server"`
	} `json:"config"`
	Files []struct {
		Fileinfo    string      `json:"fileinfo"`
		Ext         interface{} `json:"ext"`
		Filesize    string      `json:"filesize"`
		Info        string      `json:"info"`
		Invalid     bool        `json:"invalid"`
		Leech       string      `json:"leech"`
		Name        string      `json:"name"`
		Origurl     string      `json:"origurl"`
		Seed        string      `json:"seed"`
		Server      string      `json:"server"`
		Serverclass string      `json:"serverclass"`
		Source      string      `json:"source"`
		URL         string      `json:"url"`
	} `json:"files"`
	Perpage int `json:"perpage"`
	Total   int `json:"total"`
}

func readFileIntoList(fn string) []string {
	var res []string

	file, err := os.Open(fn)
	handleErrorAndPanic(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, strings.TrimSpace(scanner.Text()))
	}

	err = scanner.Err()
	handleErrorAndPanic(err)

	err = file.Close()
	handleErrorAndPanic(err)

	return res
}

func (fl *filesloop) storeClients() {
	infoLog("Storing Clients into Memory")
	defer infoLog("Finished storing Clients into Memory")
	// read brand names
	myclientsDir, err := ioutil.ReadDir("./myclients")
	handleErrorAndPanic(err)

	for _, f := range myclientsDir {
		if !f.IsDir() {
			sClient := myclient{}
			sClient.fileName = strings.TrimSuffix(f.Name(), ".txt")
			sClient.data = readFileIntoList("./myclients/" + f.Name())
			sort.Strings(sClient.data)
			fl.myclients = append(fl.myclients, sClient)
		}
	}
}

func (fl *filesloop) postProcess() {
	infoLog("Starting post processing")
	defer infoLog("Finished post processing")

	removeDuplicatesUnordered := func(elements []string) []string {
		encountered := map[string]bool{}

		// Create a map of all unique elements.
		for v := range elements {
			encountered[elements[v]] = true
		}

		// Place all keys from the map into a slice.
		result := []string{}
		for key := range encountered {
			result = append(result, key)
		}
		return result
	}

	data := readFileIntoList(DEBUGFILEPATH)

	data = removeDuplicatesUnordered(data)

	err := os.Remove(DEBUGFILEPATH)
	handleErrorAndPanic(err)

	writeToFile(DEBUGFILEPATH, strings.Join(data, "\n"))
}

func (fl *filesloop) processAll() {
	threadPool1 := pool.NewLimited(fl.workers)
	defer threadPool1.Close()

	threadPoolBatch1 := threadPool1.Batch()

	for _, client := range fl.myclients {
		for _, data := range client.data {
			if fl.seAlluc {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEALLUC, client.fileName, data))
			}
			if fl.seSharedir {
				threadPoolBatch1.Queue(fl.useSearchAPI(SESHAREDIR, client.fileName, data))
			}
			if fl.seSearch4shared {
				threadPoolBatch1.Queue(fl.useSearchAPI(SESEARCH4SHARED, client.fileName, data))
			}
			if fl.seFilesbug {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEFILESBUG, client.fileName, data))
			}
			if fl.seRapid4me {
				threadPoolBatch1.Queue(fl.useSearchAPI(SERAPID4ME, client.fileName, data))
			}
			if fl.seGeneralsearch {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEGENERALSEARCH, client.fileName, data))
			}
			if fl.seRapidsearchengine {
				threadPoolBatch1.Queue(fl.useSearchAPI(SERAPIDSEARCHENGINE, client.fileName, data))
			}
			if fl.seFilespr {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEFILESPR, client.fileName, data))
			}
			if fl.seGeneralfiles {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEGENERALFILES, client.fileName, data))
			}
			if fl.seIrfree {
				threadPoolBatch1.Queue(fl.useSearchAPI(SEIRFREE, client.fileName, data))
			}
			if fl.seSoftarchive {
				threadPoolBatch1.Queue(fl.useSearchAPI(SESOFTARCHIVE, client.fileName, data))
			}
			if fl.seSceper {
				threadPoolBatch1.Queue(fl.useSearchAPI(SESCEPER, client.fileName, data))
			}
			if fl.se2ddl {
				threadPoolBatch1.Queue(fl.useSearchAPI(SE2DDL, client.fileName, data))
			}
		}
	}
	threadPoolBatch1.QueueComplete()

	for work := range threadPoolBatch1.Results() {
		if err := work.Error(); err != nil {
			errorLog(err)
			continue
		}
		res := work.Value().(string)
		infoLog(res, "finished")
	}

	flLinks := readFileIntoList(DEBUGFILEPATH)
	err := os.Remove(DEBUGFILEPATH)
	if err != nil {
		panic(err)
	}

	threadPool2 := pool.NewLimited(fl.workers)
	defer threadPool2.Close()

	threadPoolBatch2 := threadPool2.Batch()

	for _, link := range flLinks {
		threadPoolBatch2.Queue(func(link string) pool.WorkFunc {
			return func(wu pool.WorkUnit) (interface{}, error) {
				if wu.IsCancelled() {
					// return values not used
					return nil, nil
				}

				metaSplit := strings.Split(link, "\t")
				if len(metaSplit) != 5 {
					panic(fmt.Errorf("%s - %v", "metasplit length not 5", metaSplit))
				}

				var rs []reportStructure

				searchServerLink, pageTitle, err := fl.getSearchServerLink(metaSplit[2])
				if err != nil {
					errorLog(err)
					return nil, err
				}

				cyberLockerLinks, err := fl.getCyberlockerLinks(metaSplit[1], searchServerLink)
				if err != nil {
					errorLog(err)
					return nil, err
				}

				for _, cLink := range cyberLockerLinks {
					rs = append(rs, reportStructure{
						licensor:            metaSplit[0],
						searchServerClass:   metaSplit[1],
						siteLink:            metaSplit[2],
						cyberlockerFiletype: metaSplit[3],
						cyberlockerFilesize: metaSplit[4],
						searchServerLink:    searchServerLink,
						pageTitle:           pageTitle,
						cyberlockerLink:     cLink,
					})
				}
				return rs, nil
			}
		}(link))
	}
	threadPoolBatch2.QueueComplete()

	for work := range threadPoolBatch2.Results() {
		if err := work.Error(); err != nil {
			errorLog(err)
			continue
		}
		res := work.Value().([]reportStructure)
		if len(res) > 0 {
			for _, r := range res {
				if r.licensor == "" || r.siteLink == "" || r.pageTitle == "" || r.searchServerClass == "" ||
					r.searchServerLink == "" || r.cyberlockerLink == "" || r.cyberlockerFiletype == "" || r.cyberlockerFilesize == "" {
					panic(fmt.Errorf("%s - %v", "report length not 8, empty in some scopes", r))
				}
				writeToFile(
					DEBUGFILEPATH,
					fmt.Sprintf(
						"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s", r.licensor, r.siteLink, r.pageTitle, r.searchServerClass,
						r.searchServerLink, r.cyberlockerLink, r.cyberlockerFiletype, r.cyberlockerFilesize),
				)
			}
		}
	}
}

func (fl *filesloop) useSearchAPI(driver, licensor, searchTerm string) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}

		flr := flRequestResults{
			Search:   searchTerm,
			Server:   "all",
			Order:    "relevance",
			Filesize: "all",
			Ext:      "all",
			Page:     1,
			Driver:   driver,
		}

		var lastResp responseResults

		for {
			reqBody, err := json.Marshal(flr)
			handleErrorAndPanic(err)

			body, _, err := requestPost(APIURL+driver, reqBody, fl.timeout)
			if err != nil {
				if err.Error() != "timeout" {
					errorLog(err, "Licensor: [", licensor, "] SearchEngine: [", driver, "] SearchTerm: [", searchTerm, "]")
					fmt.Print("\n\n")
					return nil, err
				}
				debugLog("Timeout for Licensor: [", licensor, "] SearchEngine: [", driver, "] SearchTerm: [", searchTerm, "]")
				fmt.Print("\n\n")
				continue
			}

			var resp responseResults

			debugLog("body", string(body))
			err = json.Unmarshal(body, &resp)
			if err != nil {
				return nil, err
			}

			infoLog(licensor, "Total:", resp.Total, "Per Page:", resp.Perpage)
			infoLog(licensor, "Files on this page:", resp.Files)

			if len(resp.Files) <= 0 || resp.Total <= 0 || reflect.DeepEqual(lastResp.Files, resp.Files) {
				debugLog("Nothing found for Licensor: [", licensor, "] SearchEngine: [", driver, "] SearchTerm: [", searchTerm, "]")
				fmt.Print("\n\n")
				break
			}

			lastResp = resp

			for _, f := range resp.Files {
				// licensor [TAB] driver [TAB] siteLink [TAB] cyberlockerFiletype [TAB] cyberlockerFilesize
				ext := ""
				fsize := ""
				if e, ok := f.Ext.(string); ok {
					if e == "" {
						ext = "empty"
					} else {
						ext = e
					}
					ext = "empty"
				} else {
					ext = "empty"
				}
				if f.Filesize == "" {
					fsize = "empty"
				} else {
					fsize = f.Filesize
				}
				infoLog(licensor+"\t"+driver+"\t", flr.Page, "\t"+"http://www.filesloop.com"+f.URL+"\t"+ext+"\t"+fsize)
				writeToFile(DEBUGFILEPATH, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", licensor, driver, "http://www.filesloop.com"+f.URL, ext, fsize))
			}

			fmt.Print("\n\n")

			flr.Page += 1
		}
		return licensor + " " + searchTerm + " " + driver, nil
	}
}

func (fl *filesloop) getSearchServerLink(link string) (searchServerLink, pageTitle string, err error) {
	body, _, err := requestGet(link, fl.timeout)
	if err != nil {
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", "", err
	}

	pageTitle = doc.Find("title").First().Text()

	iframes := doc.Find("iframe")

	for node := range iframes.Nodes {
		iframeAttr, exists := iframes.Eq(node).Attr("sandbox")
		if !exists || iframeAttr != "allow-same-origin allow-scripts allow-popups allow-forms" {
			continue
		}

		link, exists := iframes.Eq(node).Attr("src")
		if !exists || link == "" {
			continue
		}
		searchServerLink = link
	}

	return searchServerLink, pageTitle, nil
}

func (fl *filesloop) getCyberlockerLinks(searchEngine, link string) (cyberlockerLinks []string, err error) {
	body, _, err := requestGet(link, fl.timeout)
	if err != nil {
		return []string{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	handleErrorAndPanic(err)

	switch searchEngine {
	case SEALLUC:
		cLink, err := fl.parseAlluc(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink)
	case SESHAREDIR:
		cLink, err := fl.parseShareDir(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink...)
	case SESEARCH4SHARED:
		//cyberlockerLinks, err = fl.parseSearch4Shared(doc)
	case SEFILESBUG:
		cLink, err := fl.parseFilesBug(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink...)
	case SERAPID4ME:
		cLink, err := fl.parseRapid4Me(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink)
	case SEGENERALSEARCH:
		//cyberlockerLinks, err = fl.parseGeneralSearch(doc)
	case SERAPIDSEARCHENGINE:
		cLink, err := fl.parseRapidSearchEngine(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink)
	case SEFILESPR:
		//		cyberlockerLinks, err = fl.parseFilesPr(doc)
	case SEGENERALFILES:
		//		cyberlockerLinks, err = fl.parseGeneralFiles(doc)
	case SEIRFREE:
		//		cyberlockerLinks, err = fl.parseIrFree(doc)
	case SESOFTARCHIVE:
		//		cyberlockerLinks, err = fl.parseSoftArchive(doc)
	case SESCEPER:
		cLink, err := fl.parseSceper(doc)
		if err != nil {
			return []string{}, err
		}
		cyberlockerLinks = append(cyberlockerLinks, cLink...)
	case SE2DDL:
		//		cyberlockerLinks, err = fl.parse2ddl(doc)
	default:
		return []string{}, errors.New("Unable to match requested cyberlocker type!")
	}

	return cyberlockerLinks, nil
}

func (fl *filesloop) parseAlluc(doc *goquery.Document) (parsedLink string, err error) {
	html, err := doc.Html()
	if err != nil {
		return "", nil
	}

	encryptString, encryptKey, err := func() (string, string, error) {
		const decAPattern = "decrypt('"
		const decBPattern = "' )"

		decA := strings.Index(html, decAPattern)
		if decA == -1 {
			return "", "", fmt.Errorf("%s - %s", "Index not found for decA", html)
		}

		decA += len(decAPattern)

		decB := strings.Index(html[decA:], decBPattern)
		if decB == -1 {
			return "", "", fmt.Errorf("%s - %s", "Index not found for decB", html)
		}

		tIndex := strings.Index(html[decA:decA+decB], "', '")
		if tIndex == -1 {
			return "", "", fmt.Errorf("%s - %s", "Index not found for tIndex", html)
		}

		return strings.TrimPrefix(html[decA:decA+decB], decAPattern)[:tIndex],
			string(html[decA : decA+decB][len(html[decA:decA+decB])-1]),
			nil
	}()
	if err != nil {
		return "", err
	}

	parsedLink, err = func() (string, error) {
		vm := goja.New()
		value, err := vm.RunString(`
		function base64_decode(r){var e,t,o,n,a,i,d,h,c="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",f=0,u=0,l="",s=[];if(!r)return r;r+="";do n=c.indexOf(r.charAt(f++)),a=c.indexOf(r.charAt(f++)),i=c.indexOf(r.charAt(f++)),d=c.indexOf(r.charAt(f++)),h=n<<18|a<<12|i<<6|d,e=h>>16&255,t=h>>8&255,o=255&h,64==i?s[u++]=String.fromCharCode(e):64==d?s[u++]=String.fromCharCode(e,t):s[u++]=String.fromCharCode(e,t,o);while(f<r.length);return l=s.join(""),l.replace(/\0+$/,"")}function ord(r){var e=r+"",t=e.charCodeAt(0);if(t>=55296&&56319>=t){var o=t;if(1===e.length)return t;var n=e.charCodeAt(1);return 1024*(o-55296)+(n-56320)+65536}return t>=56320&&57343>=t?t:t}function decrypt(r,e){var t="";r=base64_decode(r);var o=0;for(o=0;o<r.length;o++){var n=r.substr(o,1),a=e.substr(o%e.length-1,1);n=Math.floor(ord(n)-ord(a)),n=String.fromCharCode(n),t+=n}return t};
		decrypt('` + encryptString + `', '` + encryptKey + `');
		`)
		if err != nil {
			return "", fmt.Errorf("VM error in Goja - %v", err)
		}

		valueStr := value.ToString().String()

		if strings.HasPrefix(valueStr, "<iframe") {
			iA := strings.Index(valueStr, "src=\"")
			if iA == -1 {
				return "", errors.New("Unable to find src in iframe")
			}

			iB := strings.Index(valueStr[iA+len("src=\""):], "\"")
			if iB == -1 {
				return "", errors.New("Unable to find closure for src in iframe")
			}

			valueStr = valueStr[iA+len("src=\"") : iA+len("src=\"")+iB]
		}
		return valueStr, nil
	}()
	if err != nil {
		return "", err
	}

	return parsedLink, nil
}

func (fl *filesloop) parseShareDir(doc *goquery.Document) (parsedLink []string, err error) {
	pLink := strings.TrimSpace(doc.Find("#dirlinks").First().Text())
	if pLink == "" {
		err = errors.New("Parsed link is empty for ShareDir")
		return []string{}, err
	}
	parsedLink = strings.Split(pLink, "\n")
	return parsedLink, nil
}

func (fl *filesloop) parseSearch4Shared(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseFilesBug(doc *goquery.Document) (parsedLinks []string, err error) {
	parsedLinks = strings.Split(strings.TrimSpace(doc.Find("textarea").First().Text()), "\n")
	if len(parsedLinks) <= 0 {
		err = errors.New("Parsed link is empty for FilesBug")
		return []string{}, err
	}
	return parsedLinks, nil
}

func (fl *filesloop) parseRapid4Me(doc *goquery.Document) (parsedLink string, err error) {
	parsedLink, exists := doc.Find("input").Filter(".push-05left").Attr("value")
	if !exists || parsedLink == "" {
		err = errors.New("Parsed link is empty for Rapid4Me")
		return "", err
	}
	return parsedLink, nil
}

func (fl *filesloop) parseGeneralSearch(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseRapidSearchEngine(doc *goquery.Document) (parsedLink string, err error) {
	parsedLink = strings.TrimSpace(doc.Find("#dirlinksdiv").First().Text())
	if parsedLink == "" {
		parsedLink = strings.TrimSpace(doc.Find("#dirlinks").First().Text())
	}
	if parsedLink == "" {
		err = errors.New("Parsed link is empty for RapidSearchEngine")
		return "", err
	}
	return parsedLink, nil
}

func (fl *filesloop) parseFilesPr(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseGeneralFiles(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseIrFree(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseSoftArchive(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}

func (fl *filesloop) parseSceper(doc *goquery.Document) (parsedLink []string, err error) {
	spans := doc.Find("span")
	var nodes []int
	for node := range spans.Nodes {
		attri, exists := spans.Eq(node).Attr("style")
		if !exists || attri != "color: #99cc00" {
			continue
		}
		infoLog(attri, node)
		nodes = append(nodes, node)
	}

	if len(nodes) < 2 {
		err = errors.New("Parsed link is empty for Sceper")
		return []string{}, err
	}

	atags := doc.Find("span").Eq(nodes[1]).Find("a")
	for node := range atags.Nodes {
		link, exists := atags.Eq(node).Attr("href")
		if !exists || link == "" {
			continue
		}
		parsedLink = append(parsedLink, link)
	}

	return parsedLink, nil
}

func (fl *filesloop) parse2ddl(doc *goquery.Document) (parsedLink string, err error) {
	return parsedLink, nil
}
