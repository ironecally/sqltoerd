package tables

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Table struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name      string
	IsPrimary bool
	DataType  string
	Meta      string
}

func ReadAndConvertDDL(filename string) (string, error) {
	tableList, err := readTableListFromFile(filename)
	if err != nil {
		return "", err
	}

	indexList, err := readIndexListFromFile(filename)
	if err != nil {
		return "", err
	}

	tables := extractTableToStruct(tableList, indexList)
	mmdList := convertTablesToMMDFormat(tables)

	mmdRes := encodeToMermaidERD(mmdList)

	return mmdRes, nil
}

func readTableListFromFile(fileDir string) (string, error) {
	res := ""

	f, err := os.Open(fileDir)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)
	isRecording := false
	parenthesisCheck := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CREATE TABLE") {
			isRecording = true
		}

		if isRecording {

			if strings.Contains(line, "(") {
				parenthesisCheck++
			}
			if strings.Contains(line, ")") {
				parenthesisCheck--
			}
			res = fmt.Sprintf("%s%s\n", res, line)

			if parenthesisCheck == 0 {
				isRecording = false
			}
		}
	}

	return res, nil
}

func readIndexListFromFile(fileDir string) (map[string]map[string]bool, error) {
	isPK := make(map[string]map[string]bool)

	f, err := os.Open(fileDir)
	if err != nil {
		return isPK, err
	}

	scanner := bufio.NewScanner(f)
	isRecording := false
	tableName := ""
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "ALTER TABLE") {
			re := regexp.MustCompile(`\.([_A-Za-z]+)`)
			table := re.FindAllString(line, -1)
			if len(table) > 0 {
				isRecording = true
				tableName = strings.Replace(table[0], ".", "", 1)
			}
		}

		if isRecording && strings.Contains(line, "PRIMARY KEY") {

			re := regexp.MustCompile(`PRIMARY KEY \(([_A-Za-z]+)\)`)
			fieldName := re.FindAllString(line, -1)
			fieldName[0] = strings.Replace(fieldName[0], "PRIMARY KEY (", "", 1)
			fieldName[0] = strings.Replace(fieldName[0], ")", "", 1)
			fieldPK := make(map[string]bool)
			fieldPK[fieldName[0]] = true
			isPK[tableName] = fieldPK

			tableName = ""
			isRecording = false
		}
	}
	return isPK, nil
}

func extractTableToStruct(src string, isPK map[string]map[string]bool) []Table {
	tableList := []Table{}
	rows := strings.Split(src, "\n")
	parenthesisCheck := 0

	table := Table{}
	fields := []Field{}
	for _, row := range rows {

		if strings.Contains(row, "(") {
			parenthesisCheck++
		}
		if strings.Contains(row, ")") {
			parenthesisCheck--
		}

		if strings.Contains(row, "CREATE TABLE") {
			row = strings.Replace(row, "CREATE TABLE public.", "", -1)
			row = strings.Replace(row, " (", "", -1)
			table.Name = row

		} else if strings.HasPrefix(row, "    ") {
			//this is for the fieldnames
			//sanitize the data, removing tabs in the front and commas in the back
			row = strings.Replace(row, "    ", "", -1)
			row = strings.Replace(row, ",", "", -1)
			row = strings.Replace(row, `"`, "", -1)
			field := Field{}
			splittedRow := strings.Split(row, " ")

			field.Name = splittedRow[0]
			field.DataType = splittedRow[1]
			field.Meta = strings.Join(splittedRow[2:], " ")

			field.IsPrimary = isPK[table.Name][field.Name]

			fields = append(fields, field)

		}

		if parenthesisCheck == 0 {
			table.Fields = fields

			tableList = append(tableList, table)
			table = Table{}
			fields = []Field{}
		}
	}

	return tableList
}

func convertIsPKtoMMD(isPK bool) string {
	if isPK {
		return "PK"
	}
	return ""
}

func convertTablesToMMDFormat(tables []Table) []string {
	mermaidList := []string{}
	tableID := 0
	for _, table := range tables {
		alias := fmt.Sprintf("n%d", tableID)
		if table.Name == "" {
			continue
		}
		mmd := fmt.Sprintf("\t%s[%s] {\n", alias, table.Name)
		tableID++
		for _, field := range table.Fields {
			mmd += fmt.Sprintf("\t\t%s %s %s %s\n", field.DataType, field.Name, convertIsPKtoMMD(field.IsPrimary), `"`+field.Meta+`"`)
		}
		mmd += "\t}\n"

		mermaidList = append(mermaidList, mmd)
	}
	return mermaidList
}

func encodeToMermaidERD(mmdList []string) string {
	res := fmt.Sprintf("erDiagram\n%s", strings.Join(mmdList, "\n\n"))

	return res
}
