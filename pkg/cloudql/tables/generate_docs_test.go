package opengovernance

import (
	"context"
	"fmt"
	"html"
	"os"
	"strings"
	"testing"
)

func TestGenerateDocs(t *testing.T) {
	plg := Plugin(context.Background())

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pathPrefix := "./docs/table_def"
	if strings.HasSuffix(currentDir, "opengovernance") {
		pathPrefix = "." + pathPrefix
	}

	for _, table := range plg.TableMap {
		doc := `# Columns  

<table>
	<tr><td>Column Name</td><td>Description</td></tr>
`
		for _, column := range table.Columns {
			desc := column.Description
			desc = html.EscapeString(desc)
			doc += fmt.Sprintf(`	<tr><td>%s</td><td>%s</td></tr>
`, column.Name, desc)
		}

		doc += "</table>"

		err := os.WriteFile(pathPrefix+"/"+table.Name+".md", []byte(doc), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
