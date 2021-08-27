// Encoding: UTF-8
//
// JUnit Gate
//
// Copyright © 2021 Brian Dwyer - Intelligent Digital Services
//

package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joshdk/go-junit"
	"gopkg.in/yaml.v3"
)

var configPath string
var fileFlag string
var debugFlag bool

func init() {
	flag.StringVar(&configPath, "c", "", "Path to junit-gate config file")
	flag.StringVar(&fileFlag, "f", "", "Path to the Junit XML file")
	flag.BoolVar(&debugFlag, "debug", false, "Enable verbose log output")
}

func main() {
	// Parse Flags
	flag.Parse()

	if debugFlag {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	if versionFlag {
		showVersion()
		os.Exit(0)
	}

	if fileFlag == "" {
		if len(os.Args) >= 2 {
			fileFlag = os.Args[1]
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if configPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		configPath = filepath.Join(pwd, ".junit-gate.yml")
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	suites, err := junit.IngestFile(fileFlag)
	if err != nil {
		log.Fatalf("failed to ingest JUnit xml %v", err)
	}

	var result Result

	var evalSuite func(suite junit.Suite)
	evalSuite = func(suite junit.Suite) {
		// Recurse through child suites
		for _, s := range suite.Suites {
			evalSuite(s)
		}
		// Check if suite is excluded
		for _, exception := range config.Exceptions() {
			if exception.Classname == "" && exception.Name == "" {
				if exception.Suite == "" && exception.Package == "" {
					log.Fatalln("Invalid Exception:", exception)
				}
				if exception.Suite != "" && suite.Name == exception.Suite {
					log.Debugln("Suite excluded:", suite.Name)
					result.Exceptions = append(result.Exceptions, ExceptionMatch{exception, suite, "Suite Match"})
					return
				} else if exception.Package != "" && suite.Package == exception.Package {
					log.Debugln("Package excluded:", suite.Package)
					result.Exceptions = append(result.Exceptions, ExceptionMatch{exception, suite, "Suite Match"})
					return
				}
			}
		}

		for _, test := range suite.Tests {
			if test.Error != nil {
				var excluded bool
				for _, exception := range config.Exceptions() {
					// Check if the exception is scoped to a specific Suite or Package
					if (exception.Suite != "" && exception.Suite != suite.Name) || (exception.Package != "" && exception.Package != suite.Package) {
						continue
					}
					// Compare Properties
					if !exception.PropertiesMatch(test.Properties) && !exception.PropertiesMatch(suite.Properties) {
						continue
					}

					if test.Name == exception.Name {
						log.Debugln("Test excluded by Name!", test.Name)
						excluded = true
						result.Exceptions = append(result.Exceptions, ExceptionMatch{exception, test, "Name Match"})
						break
					} else if exception.Name == "" && strings.HasPrefix(test.Classname, exception.Classname) {
						log.Debugln("Test excluded by Classname!", test.Classname)
						excluded = true
						result.Exceptions = append(result.Exceptions, ExceptionMatch{exception, test, "Classname Match"})
						break
					} else if exception.Suite != "" && exception.Properties != nil && exception.Suite == suite.Name && exception.PropertiesMatch(test.Properties) {
						log.Debugln("Test excluded by Suite & Properties!", suite.Name)
						excluded = true
						result.Exceptions = append(result.Exceptions, ExceptionMatch{exception, test, "Suite & Properties Match"})
					}
				}

				if !excluded {
					result.Errors = append(result.Errors, test)
				}
			}
		}
	}

	for _, suite := range suites {
		suite.Aggregate()

		if suite.Totals.Failed == 0 && suite.Totals.Error == 0 {
			continue
		}

		evalSuite(suite)
	}

	if len(result.Exceptions) > 0 {
		b, err := json.MarshalIndent(result.Exceptions, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Infoln("Exceptions:", string(b))
	}

	if len(result.Errors) > 0 {
		b, err := json.MarshalIndent(result.Errors, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalln("Failures:", string(b))
	}
}
