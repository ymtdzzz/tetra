package adapter

import (
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
)

func ConvertResultToCSV(result []table.Row) string {
	if len(result) == 0 {
		return ""
	}

	var b strings.Builder
	cols := []string{}

	i := 0
	for k := range result[0].Data {
		data := result[0].Data
		cols = append(cols, k)
		b.WriteString(k)
		if i < len(data)-1 {
			b.WriteString(",")
		}
		i++
	}
	b.WriteString("\n")

	for _, r := range result {
		for i, k := range cols {
			b.WriteString(fmt.Sprintf("%v", r.Data[k]))
			if i < len(cols)-1 {
				b.WriteString(",")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
