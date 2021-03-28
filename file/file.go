// +build !js

package file

import (
	"io/ioutil"
)

func ReadAll(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}
