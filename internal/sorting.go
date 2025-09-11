package internal

import (
	"slices"
	"sort"
	"strings"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

// sortYamlNodes recursively sorts YAML nodes
func SortYamlNodes(node *yaml.Node, cfg Config) {
	var rootNode *yaml.Node
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		sortMapNodes(node, cfg)
		for _, content := range node.Content {
			SortYamlNodes(content, cfg)
			node = content
		}
	case yaml.SequenceNode:
		sortSequenceNodes(node, cfg)
		for _, content := range node.Content {
			SortYamlNodes(content, cfg)
			node = content
		}
	case yaml.DocumentNode:
		rootNode = extractRootNode(*node)
		SortYamlNodes(rootNode, cfg)
	}
}

// sortMapNodes sorts the keys of a YAML mapping node
func sortMapNodes(node *yaml.Node, cfg Config) {
	if node.Kind != yaml.MappingNode || len(node.Content) < 2 {
		return
	}

	keys := make([]string, 0)
	pairs := make(map[string][]*yaml.Node)
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value
		keys = append(keys, key)
		pairs[key] = []*yaml.Node{keyNode, valueNode}
	}
	if cfg.CustomSort != nil {
		customSort(&keys, cfg)
	} else {
		if cfg.Reverse {
			slices.Reverse(keys)
		} else {
			slices.Sort(keys)
		}
	}
	newContent := make([]*yaml.Node, 0)
	for _, key := range keys {
		newContent = append(newContent, pairs[key][0], pairs[key][1])
	}

	node.Content = newContent
}

// sortSequenceNodes sorts the elements of a YAML sequence node
func sortSequenceNodes(node *yaml.Node, cfg Config) {
	if node.Kind != yaml.SequenceNode || len(node.Content) < 1 || !cfg.SortList {
		return
	}

	// Create a slice of sequence items with their string representation for sorting
	type sortableItem struct {
		original *yaml.Node
		sortKey  string
	}

	items := make([]sortableItem, len(node.Content))
	for i, item := range node.Content {
		// Convert the node to a string for comparison
		var sb strings.Builder
		encoder := yaml.NewEncoder(&sb)
		encoder.Encode(item)
		items[i] = sortableItem{
			original: item,
			sortKey:  sb.String(),
		}
	}

	// Sort the items based on their string representation
	sort.Slice(items, func(i, j int) bool {
		return items[i].sortKey < items[j].sortKey
	})

	// Rebuild the content in sorted order
	newContent := make([]*yaml.Node, len(items))
	for i, item := range items {
		newContent[i] = item.original
	}

	node.Content = newContent
}

func customSort(keys *[]string, cfg Config) {
	customSort := cfg.CustomSort
	found := false
	customSorted := []string{}
	for _, keyword := range customSort {
		if j := slices.Index(*keys, keyword); j != -1 {
			customSorted = append(customSorted, keyword)
			*keys = slices.Delete(*keys, j, j+1)
			found = true
		}
	}
	if cfg.Reverse {
		slices.Reverse(*keys)
	} else {
		slices.Sort(*keys)
	}
	if found {
		*keys = slices.Concat(customSorted, *keys)
	}
}

// addEmptyLinesBeforeTopLevelKeys adds empty lines before top-level keys
func AddEmptyLinesBeforeTopLevelKeys(node *yaml.Node) {
	rootNode := extractRootNode(*node)
	if rootNode.Kind != yaml.MappingNode {
		return
	}

	for i := 2; i < len(rootNode.Content); i += 2 {
		if i >= len(rootNode.Content) {
			break
		}
		keyNode := rootNode.Content[i]

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

func extractRootNode(node yaml.Node) (rootNode *yaml.Node) {
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		rootNode = node.Content[0]
	} else {
		rootNode = &node
	}
	return
}
