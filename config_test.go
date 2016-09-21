package gocrawl

import (
	"testing"
	"strings"
	"os"
	"os/exec"
)

/* A valid JSON configuration as a byte array*/
var validJSONConfig = []byte(
	`{
	"startUrl": "www.abc.com",
	"urlParsingRules": {
		"sameDomain": true,
		"depth": 10
	},
	"dataParsingRules": [
		{
			"urlMatch": "www.abc.com/pageA",
			"dataSelector": "div.content > p"
		},
		{
			"urlMatch": "www.abc.com/pageB",
			"dataSelector": "div.main p span"
		}
	]
	}`)

/* The valid JSON config but in lower-case */
var lowerJSONConfig = []byte(
	`{
	"starturl": "www.abc.com",
	"urlparsingrules": {
		"sameDomain": true,
		"depth": 10
	},
	"dataparsingrules": [
		{
			"urlmatch": "www.abc.com/pageA",
			"dataselector": "div.content > p"
		},
		{
			"urlmatch": "www.abc.com/pageB",
			"dataselector": "div.main p span"
		}
	]
	}`)

/* An invalid JSON configuration */
var invalidJSONConfig = []byte(
	`{
	"startUrl": "www.abc.com",
	"urlRules": {
		"sameDomain": true,
		"depth": 10
	},
	"dataRules": [
		{
			"urlMatch": "www.abc.com/pageA",
			"dataSelector": "div.content > p"
		},
		{
			"urlMatch": "www.abc.com/pageB",
			"dataSelector": "div.main p span"
		}
	]
	}`)

func TestReadsValidConfigIntoStruct(t *testing.T) {
	config := readConfig(validJSONConfig)
	if !strings.EqualFold(config.StartUrl, "www.abc.com") {
		t.Fatal("Failed to read the valid byte array into a Config struct")
	}
}

func TestReadsValidLowercaseConfigIntoStruct(t *testing.T) {
	config := readConfig(lowerJSONConfig)
	if !strings.EqualFold(config.StartUrl, "www.abc.com") {
		t.Fatal("Failed to read the valid byte array into a Config struct")
	}
}

/* This tests that reading an invalid configuration file will stop the program. To do this, call the
exception throwing function from a sub-process and check the return code of that process. This idea
is from http://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go/26226382#26226382 */
func TestReadingInvalidConfigErrorsOutOfTheProgram(t *testing.T) {
	if os.Getenv("SUBPROC_TestReadingInvalidConfigErrorsOutOfTheProgram") == "1" {
		readConfig(invalidJSONConfig)
		t.Fatal("The code should not reach this point, the program should've exited already")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReadingInvalidConfigErrorsOutOfTheProgram")
	cmd.Env = append(os.Environ(), "SUBPROC_TestReadingInvalidConfigErrorsOutOfTheProgram=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
