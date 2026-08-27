package internal

import (
	"reflect"
	"testing"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var (
	yamlUnsorted = []byte(`c:
  z: test1
  h: test2
a: test3
d:
  - z
  - h
  - a
y:
  - bname: b
    avalue: bv
  - yname: a
    xvalue: av
`)
	yamlSorted = []byte(`a: test3
c:
  h: test2
  z: test1
d:
  - a
  - h
  - z
y:
  - avalue: bv
    bname: b
  - xvalue: av
    yname: a
`)
)

func TestSortYamlNodes(t *testing.T) {
	var got, want yaml.Node

	if err := yaml.Unmarshal(yamlUnsorted, &got); err != nil {
		t.Errorf("Error unmarshaling the unsorted source YAML: %s", err)
	}
	if err := yaml.Unmarshal(yamlSorted, &want); err != nil {
		t.Errorf("Error unmarshaling the sorted source YAML: %s", err)
	}
	SortYamlNodes(&got, Config{SortList: true})
	btGot, _ := yaml.Marshal(got.Content[0])
	btWant, _ := yaml.Marshal(want.Content[0])
	ok := reflect.DeepEqual(btGot, btWant)
	if !ok {
		t.Errorf("got %s, want %s", btGot, btWant)
	}
}

func TestSortYamlNodesPreservesDuplicateKeys(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte("b: third\na: first\na: second\n"), &node); err != nil {
		t.Fatal(err)
	}

	SortYamlNodes(&node, Config{})

	root := node.Content[0]
	got := make([]string, 0, len(root.Content))
	for _, content := range root.Content {
		got = append(got, content.Value)
	}
	want := []string{"a", "first", "a", "second", "b", "third"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got mapping content %q, want %q", got, want)
	}
}

func TestSortYamlNodesReversesSequences(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte("- a\n- z\n- m\n"), &node); err != nil {
		t.Fatal(err)
	}

	SortYamlNodes(&node, Config{SortList: true, Reverse: true})

	root := node.Content[0]
	got := make([]string, 0, len(root.Content))
	for _, content := range root.Content {
		got = append(got, content.Value)
	}
	want := []string{"z", "m", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got sequence %q, want %q", got, want)
	}
}

func TestSortYamlNodesSortsMappingsInReverseOrder(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{name: "default order", cfg: Config{Reverse: true}},
		{name: "empty custom order", cfg: Config{Reverse: true, CustomSort: []string{}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var node yaml.Node
			if err := yaml.Unmarshal([]byte("c: 3\na: 1\nd: 4\ny: 25\n"), &node); err != nil {
				t.Fatal(err)
			}

			SortYamlNodes(&node, test.cfg)

			root := node.Content[0]
			got := []string{root.Content[0].Value, root.Content[2].Value, root.Content[4].Value, root.Content[6].Value}
			want := []string{"y", "d", "c", "a"}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got keys %q, want %q", got, want)
			}
		})
	}
}

func TestSortYamlNodesUsesCustomOrder(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte("c: 3\na: 1\nb: 2\n"), &node); err != nil {
		t.Fatal(err)
	}

	SortYamlNodes(&node, Config{CustomSort: []string{"b"}})

	root := node.Content[0]
	got := []string{root.Content[0].Value, root.Content[2].Value, root.Content[4].Value}
	want := []string{"b", "a", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got keys %q, want %q", got, want)
	}
}

func TestAddEmptyLinesBeforeTopLevelKeysPreservesComments(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte("first: 1\n# second comment\nsecond: 2\n"), &node); err != nil {
		t.Fatal(err)
	}

	AddEmptyLinesBeforeTopLevelKeys(&node)

	secondKey := node.Content[0].Content[2]
	if secondKey.HeadComment != "\n# second comment" {
		t.Fatalf("got head comment %q, want a leading empty line", secondKey.HeadComment)
	}
}
