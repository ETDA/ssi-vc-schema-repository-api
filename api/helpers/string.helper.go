package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func FileExtensionMatch(filename string, expectedExtension string) error {
	split := strings.Split(filename, ".")
	if len(split) != 2 {
		return errors.New("required filename in form filename.ext")
	}

	if split[1] != expectedExtension {
		return fmt.Errorf("file extension %s should match expected extension %s", split[1], expectedExtension)
	}

	return nil
}

func UniqueStringArray(v []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range v {
		if occured[v[e]] != true {
			occured[v[e]] = true
			result = append(result, v[e])
		}
	}
	return result
}

func AlphaNumericize(originalString string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(originalString, "")
}
