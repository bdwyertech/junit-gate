// Encoding: UTF-8
//
// JUnit Gate
//
// Copyright © 2022 Brian Dwyer - Intelligent Digital Services
//

package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/TylerBrock/colorjson"
)

func prettyJson(i interface{}) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err)
	}
	var obj interface{}
	err = json.Unmarshal(jsonBytes, &obj)
	if err != nil {
		log.Fatal(err)
	}
	f := colorjson.NewFormatter()
	f.Indent = 2
	out, err := f.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
