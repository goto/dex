package dlq

import "strings"

func buildDlqPrefixDirectory(template string, firehoseName string) string {
	return strings.Replace(template, "{{ .name }}", firehoseName, 1)
}
