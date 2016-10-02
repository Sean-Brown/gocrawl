package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

/* The configuration for crawling. Note the field names need to start with capital
letters in order for the JSON parser to not ignore the field (lower-case members are ignored) */
type Config struct {
	StartUrl         string
	UrlParsingRules  URLParsingRules
	DataParsingRules []DataParsingRule
	DataStore DataStorage
}

/* Rules for parsing urls */
type URLParsingRules struct {
	SameDomain bool
	MaxDepth   int
}

/* rules for parsing data from the DOM */
type DataParsingRule struct {
	/* url pattern to match the rule to */
	UrlMatch string
	/* the data goquery selector string (e.g. "div.content > div#main p.text") */
	DataSelector string
}

/* Default URL Parsing rules */
func NewURLParsingRules() URLParsingRules {
	return URLParsingRules{SameDomain: true, MaxDepth: 1}
}
func CreateURLParsingRules(sameDomain bool, maxDepth int) URLParsingRules {
	return URLParsingRules{SameDomain: sameDomain, MaxDepth: maxDepth}
}
func NewDataParsingRules() DataParsingRule {
	return DataParsingRule{UrlMatch:".*", DataSelector:".*"}
}

/* Initialize URL Parsing rules defined by the user */
func InitURLParsingRules(sameDomain bool, maxDepth int) URLParsingRules {
	return URLParsingRules{SameDomain: sameDomain, MaxDepth: maxDepth}
}

func ReadConfig(path string) Config {
	/* check if the file exists */
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Invalid path to the configuration file: ", path)
	}
	return readConfig(file)
}

func readConfig(contents []byte) Config {
	var config Config
	err := json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatal("Unable to unmarshal the configuration file into an object: ", string(contents))
	}
	return config
}
