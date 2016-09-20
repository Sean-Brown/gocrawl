package gocrawl

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

/* The configuration for crawling */
type Config struct {
	startUrl string
	urlParsingRules []URLParsingRules
	dataParsingRules []DataParsingRule
}

/* Rules for parsing urls */
type URLParsingRules struct {
	sameDomain bool
}

/* rules for parsing data from the DOM */
type DataParsingRule struct {
	/* url pattern to match the rule to */
	urlMatch string
	/* the data goquery selector string (e.g. "div.content > div#main p.text") */
	dataSelector string
}

/* Default URL Parsing rules */
func NewURLParsingRules() URLParsingRules {
	return URLParsingRules{sameDomain: true}
}
/* Initialize URL Parsing rules defined by the user */
func InitURLParsingRules(sameDomain bool) URLParsingRules {
	return URLParsingRules{sameDomain: sameDomain}
}

func ReadConfig(path string) Config {
	var config Config
	/* check if the file exists */
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Invalid path to the configuration file: ", path)
	}
	json.Unmarshal(file, &config)
	return config
}