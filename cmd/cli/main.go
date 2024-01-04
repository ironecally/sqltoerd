package main

import (
	"flag"
	"log"

	"github.com/ironecally/sqltoerd/internal/mermaid"
	tables "github.com/ironecally/sqltoerd/internal/tables"
	writer "github.com/ironecally/sqltoerd/internal/writer"
)

func main() {
	var filename = flag.String("i", "", ".sql file to parse")
	var outputFile = flag.String("o", "result/result.pdf", "output filename and extension (supporting .png, .svg and .pdf)")

	flag.Parse()
	if filename == nil {
		log.Fatal("No filename (-i) specified")
	}

	err := processSQLFile(*filename, *outputFile)
	if err != nil {
		log.Fatal(err)
	}

}

func processSQLFile(filename string, outputFile string) error {
	mmdRes, err := tables.ReadAndConvertDDL(filename)
	if err != nil {
		log.Println("fail to read source file")
		return err
	}

	mmdFileName := filename + ".mmd"
	err = writer.WriteToFile(mmdRes, mmdFileName)
	if err != nil {
		log.Println("fail to write mmd file")
		return err
	}

	err = mermaid.GenerateSVG(mmdFileName, outputFile)
	if err != nil {
		log.Println("fail to generate", outputFile)
		return err
	}

	return nil
}
