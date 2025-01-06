package dag

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gatsu420/sql_dag/utils"
)

func ParseQuery(filename string) (map[string][]string, error) {
	file, err := utils.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cteNameTemp := ""
	dagMap := map[string][]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		/*
			CTE parsing depends on regex, while Scan() runs on single line. Since CTE may span over
			several lines, we use cteNameTemp as placeholder for name of CTE whenever cteName and
			sourceName don't sit on the same line.

			Query illustration:
			// cteName and sourceName don't sit on the same line
			with pp as (
				select * from `project.dataset.table`
			)

			// They do on this line
			, ss as (select * from `project.dataset.table`)

		*/
		cteName := ""
		cteRegexp, _ := regexp.Compile(`\b(\w+)\s+as\s+\(`)
		ctes := cteRegexp.FindAllStringSubmatch(line, -1)
		for _, cte := range ctes {
			cteName = fmt.Sprint(cte[1])
		}

		sourceName := ""
		sourceRegexp, _ := regexp.Compile("from\\s+`([^`]+)`")
		sources := sourceRegexp.FindAllStringSubmatch(line, -1)
		for _, source := range sources {
			sourceName = fmt.Sprint(source[1])
		}

		// Save CTE name into temporary variable. Also skip the line to avoid having "" as sourceName.
		if sourceName == "" {
			cteNameTemp = cteName
			continue
		}

		// Return it back when sourceName is found
		if sourceName != "" && cteName == "" {
			cteName = cteNameTemp
			cteNameTemp = ""
		}

		dagMap[sourceName] = append(dagMap[sourceName], cteName)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return dagMap, nil
}

func GenerateDAG(dag map[string][]string, dotFilename string, pngFilename string) error {
	file, err := os.Create(dotFilename)
	if err != nil {
		return fmt.Errorf("error creating dot file: %w", err)
	}
	defer file.Close()

	file.WriteString("digraph DAG { \n")
	for sourceName, cteNames := range dag {
		for _, cteName := range cteNames {
			file.WriteString(fmt.Sprintf("	\"%v\" -> \"%v\"; \n", sourceName, cteName))
		}
	}
	file.WriteString("} \n")

	cmd := exec.Command("dot", "-Tpng", dotFilename, "-o", pngFilename)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing dot command: %w", err)
	}

	return nil
}
