package kstock

import "io/ioutil"

func readTestResponseFile(filename string) string {
	b, _ := ioutil.ReadFile("./testhtml/" + filename)
	return string([]byte(b))
}
