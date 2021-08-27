// Encoding: UTF-8

package main

import (
	"encoding/xml"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/araddon/dateparse"
)

type TestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	Text       string   `xml:",chardata"`
	Disabled   string   `xml:"disabled,attr"`
	Errors     string   `xml:"errors,attr"`
	Failures   string   `xml:"failures,attr"`
	Tests      string   `xml:"tests,attr"`
	Time       string   `xml:"time,attr"`
	TestSuites []struct {
		Text      string `xml:",chardata"`
		Disabled  string `xml:"disabled,attr"`
		Errors    string `xml:"errors,attr"`
		Failures  string `xml:"failures,attr"`
		Name      string `xml:"name,attr"`
		Package   string `xml:"package,attr"`
		Skipped   string `xml:"skipped,attr"`
		Tests     string `xml:"tests,attr"`
		Time      string `xml:"time,attr"`
		TestCases []struct {
			Text      string `xml:",chardata"`
			Classname string `xml:"classname,attr"`
			File      string `xml:"file,attr"`
			Name      string `xml:"name,attr"`
			Failure   *struct {
				Text    string `xml:",chardata"`
				Message string `xml:"message,attr"`
				Type    string `xml:"type,attr"`
			} `xml:"failure,omitempty"`
		} `xml:"testcase"`
	} `xml:"testsuite"`
}

type Config struct {
	ExceptionList     []*Exception `yaml:"exceptions"`
	RequireExpiration bool         `yaml:"require_expiration"`
}

type Exception struct {
	// Suite Scope
	Suite   string `json:",omitempty" yaml:"suite,omitempty"`
	Package string `json:",omitempty" yaml:"package,omitempty"`
	// Test Scope
	Name       string            `json:",omitempty" yaml:"name,omitempty"`
	Classname  string            `json:",omitempty" yaml:"classname,omitempty"`
	File       string            `json:",omitempty" yaml:"file,omitempty"`
	Properties map[string]string `json:",omitempty" yaml:"properties,omitempty"`
	// Global
	Expires     string `json:",omitempty" yaml:"expires,omitempty"`
	Description string `json:",omitempty" yaml:"description,omitempty"`
	expired     bool
}

type ExceptionMatch struct {
	Exception Exception
	Match     interface{}
	Reason    string
}

type Result struct {
	Exceptions []ExceptionMatch
	Errors     []interface{}
}

func (c *Config) Exceptions() (e []Exception) {
	for _, exception := range c.ExceptionList {
		if c.RequireExpiration && exception.Expires == "" {
			log.Debugln("Skipping exception because require_expiration: true and no expiration date is set:", exception)
			continue
		}
		if !exception.Expired() {
			e = append(e, *exception)
		}
	}
	return
}

func (e *Exception) Expired() (expired bool) {
	if e.expired {
		return true
	}
	if e.Expires == "" {
		return
	}
	d, err := dateparse.ParseAny(e.Expires)
	if err != nil {
		log.Fatal(err)
	}
	if d.Before(time.Now()) {
		log.Warnln("Exception Expired:", e)
		e.expired = true
		return true
	}
	return
}

func (e *Exception) PropertiesMatch(props map[string]string) bool {
	if e.Properties != nil {
		for k, v := range e.Properties {
			if val, ok := props[k]; ok {
				if v != val {
					log.Debugln("Value mismatch:", k, v, val)
					return false
				}
			} else {
				log.Debugln("Key Missing:", k)
				return false
			}
		}
	}
	return true
}
