package itertools

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestYamlIteration(t *testing.T) {
	rawyaml := `---
FirstNode:
    k1: v1
    k2: v2
SecondNode:
    k3: v3
    k4: v4
`
	expected := map[string]interface{}{
		"FirstNode": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
		"SecondNode": map[string]string{
			"k3": "v3",
			"k4": "v4",
		},
	}
	parsed := map[string]interface{}{}
	node := yaml.Node{}
	err := yaml.NewDecoder(strings.NewReader(rawyaml)).Decode(&node)
	if err != nil {
		t.Error(err.Error())
	}
	NewYamlIter(&node).
		Each(func(pair YamlMapPair) {
			subMap := map[string]string{}
			parsed[pair.Key.Value] = subMap
			pair.Value.Each(func(subpair YamlMapPair) {
				subpair.Value.Each(func(n *yaml.Node) {
					subMap[subpair.Key.Value] = n.Value
				})
			})
		})
	if !reflect.DeepEqual(parsed, expected) {
		t.Error("fail to reconstact yaml structure")
	}
}
