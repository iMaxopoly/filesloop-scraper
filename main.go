package main

import (
	"os"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type reportStructure struct {
	licensor  string
	siteLink  string
	pageTitle string

	searchServerClass string
	searchServerLink  string

	cyberlockerLink     string
	cyberlockerFiletype string
	cyberlockerFilesize string
}

const VERSION = "0.0.0.1"

var (
	// Command-line setup

	confVerbose = kingpin.Flag(
		"verbose",
		"Toggles verbosity, default is true").
		Default("true").Short('v').Bool()

	confWorkers = kingpin.Flag(
		"workers",
		"Number of workers making requests simultaneously and getting the links").
		Default("2").Short('w').Uint()

	confSearchEngineAlluc             = kingpin.Flag("alluc", "Toggles alluc to be enabled.").Default("false").Bool()
	confSearchEngineSharedir          = kingpin.Flag("sharedir", "Toggles sharedir to be enabled.").Default("false").Bool()
	confSearchEngineSearch4shared     = kingpin.Flag("search4shared", "Toggles search4shared to be enabled.").Default("false").Bool()
	confSearchEngineFilesbug          = kingpin.Flag("filesbug", "Toggles filesbug to be enabled.").Default("false").Bool()
	confSearchEngineRapid4me          = kingpin.Flag("rapid4me", "Toggles rapid4me to be enabled.").Default("false").Bool()
	confSearchEngineGeneralsearch     = kingpin.Flag("generalsearch", "Toggles generalsearch to be enabled.").Default("false").Bool()
	confSearchEngineRapidsearchengine = kingpin.Flag("rapidsearchengine", "Toggles rapidsearchengine to be enabled.").Default("false").Bool()
	confSearchEngineFilespr           = kingpin.Flag("filespr", "Toggles filespr to be enabled.").Default("false").Bool()
	confSearchEngineGeneralfiles      = kingpin.Flag("generalfiles", "Toggles generalfiles to be enabled.").Default("false").Bool()
	confSearchEngineIrfree            = kingpin.Flag("irfree", "Toggles irfree to be enabled.").Default("false").Bool()
	confSearchEngineSoftarchive       = kingpin.Flag("softarchive", "Toggles softarchive to be enabled.").Default("false").Bool()
	confSearchEngineSceper            = kingpin.Flag("sceper", "Toggles sceper to be enabled.").Default("false").Bool()
	confSearchEngine2ddl              = kingpin.Flag("2ddl", "Toggles 2ddl to be enabled.").Default("false").Bool()

	confSearchEngineAll = kingpin.Flag("all", "Toggles all search engines to be enabled.").Default("true").Bool()

	confRequestWaitTimeout = kingpin.Flag(
		"timeout",
		"Time out for each request after which request is abandoned; Defaults to 30").
		Default("30").Short('t').Uint()
)

func main() {
	//Commandline setup
	kingpin.Version(`
	filesloop.com scraper
 *  Contact:
 *  Manish Prakash Singh
 *  contact@kryptodev.com
 *  Skype: kryptodev
	` +
		"\n©filesloop.com scraper v" + VERSION + " - removeyourmedia.com, All Rights Reserved.")

	kingpin.Parse()

	var fl filesloop
	fl.workers = *confWorkers
	fl.timeout = *confRequestWaitTimeout

	if *confSearchEngineAll != true {
		fl.seAlluc = *confSearchEngineAlluc
		fl.seSharedir = *confSearchEngineSharedir
		fl.seSearch4shared = *confSearchEngineSearch4shared
		fl.seFilesbug = *confSearchEngineFilesbug
		fl.seRapid4me = *confSearchEngineRapid4me
		fl.seGeneralsearch = *confSearchEngineGeneralsearch
		fl.seRapidsearchengine = *confSearchEngineRapidsearchengine
		fl.seFilespr = *confSearchEngineFilespr
		fl.seGeneralfiles = *confSearchEngineGeneralfiles
		fl.seIrfree = *confSearchEngineIrfree
		fl.seSoftarchive = *confSearchEngineSoftarchive
		fl.seSceper = *confSearchEngineSceper
		fl.se2ddl = *confSearchEngine2ddl
	} else {
		fl.seAlluc = false
		fl.seSharedir = true
		fl.seSearch4shared = false
		fl.seFilesbug = false
		fl.seRapid4me = false
		fl.seGeneralsearch = false
		fl.seRapidsearchengine = false
		fl.seFilespr = false
		fl.seGeneralfiles = false
		fl.seIrfree = false
		fl.seSoftarchive = false
		fl.seSceper = false
		fl.se2ddl = false
	}
	fl.storeClients()

	for {
		fl.timestamp = time.Now().Format("02_01_06-15.04")
		fl.processAll()

		//infoLog("Post-processing text files, removing duplicates and reordering.")
		//fl.postProcess()

		fl.generateExcelReport(fl.timestamp)

		err := os.Remove(DEBUGFILEPATH)
		handleErrorAndPanic(err)

		infoLog("Task finished. Restarting in 8 hours.")

		time.Sleep(8 * time.Hour)
	}

	//searchServerLink, _, err := fl.getSearchServerLink("http://www.filesloop.com/file/-animerg-naruto-shippuuden-463-720p-10bit-jrr-shippuden-mkv/alluc/l/AnimeRG-Naruto-Shippuuden-463-720p-10bit-JRR-Shippuden-mkv/8a0hmbt5")
	//handleErrorAndPanic(err)
	//
	//cyberlockLinks, err := fl.getCyberlockerLinks(SEALLUC, searchServerLink)
	//handleErrorAndPanic(err)
	//
	//infoLog(searchServerLink, cyberlockLinks)
}
