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
