package pkg

import "io/ioutil"

func WriteToFile(file, content string) error {
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
