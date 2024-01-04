package mermaid

import (
	"os/exec"
)

func GenerateSVG(mmdfile string, output string) error {

	// commandStr := fmt.Sprintf("mmdc -i %s -o result/result.svg -s 2", mmdfile)

	cmd := exec.Command("mmdc", "-i", mmdfile, "-o", output, "-c", "config/mermaid-config.json", "-f")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
