package main

import (
	"fmt"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

var testExit = func(i int) {
	if i != 0 {
		panic(fmt.Sprintf("os.Exit(%v) called", i))
	}
}

func TestJenkins(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = ""
	fileFlag = "test/fixtures/jenkins.xml"
	assert.Panics(t, main, "Jenkins should error")
}

func TestJenkinsExclusions(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = "test/fixtures/jenkins-name.yml"
	fileFlag = "test/fixtures/jenkins.xml"
	assert.NotPanics(t, main, "Jenkins should not error")
}

func TestJenkinsExpiredExclusions(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = "test/fixtures/jenkins-name-expired.yml"
	fileFlag = "test/fixtures/jenkins.xml"
	assert.Panics(t, main, "Jenkins should error")
}

func TestBasic(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = ""
	fileFlag = "test/fixtures/basic.xml"
	assert.Panics(t, main, "Basic should error")
}

func TestBasicExceptions(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = "test/fixtures/basic-class.yml"
	fileFlag = "test/fixtures/basic.xml"
	assert.NotPanics(t, main, "Basic should not error")
}

func TestBasicExpirationRequired(t *testing.T) {
	patch := monkey.Patch(os.Exit, testExit)
	defer patch.Unpatch()

	configPath = "test/fixtures/basic-expiration-required.yml"
	fileFlag = "test/fixtures/basic.xml"
	assert.Panics(t, main, "Basic should error because no expiration is set")
}
