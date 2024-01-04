package writer

import "os"

func WriteToFile(mmdRes string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(mmdRes)
	if err != nil {
		return err
	}

	return nil
}
