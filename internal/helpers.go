package internal

import (
	"strings"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

// addEmptyLinesBeforeTopLevelKeys adds empty lines before top-level keys
func AddEmptyLinesBeforeTopLevelKeys(node *yaml.Node) {
	if node.Kind != yaml.MappingNode {
		return
	}

	for i := 2; i < len(node.Content); i += 2 {
		if i >= len(node.Content) {
			break
		}
		keyNode := node.Content[i]

		// Only add empty lines for top-level keys
		if keyNode.HeadComment == "" {
			keyNode.HeadComment = "\n"
		} else {
			// If there's already a comment, prepend a newline
			if !strings.HasPrefix(keyNode.HeadComment, "\n") {
				keyNode.HeadComment = "\n" + keyNode.HeadComment
			}
		}
	}
}
