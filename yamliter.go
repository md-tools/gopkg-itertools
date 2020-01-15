package itertools

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// YamlMapPair ...
type YamlMapPair struct {
	Key   *yaml.Node
	Value YamlIter
}

// YamlIter ...
type YamlIter struct {
	Kind yaml.Kind
	Iter
}

func (iter YamlIter) String() string {
	return fmt.Sprintf("YamlIter<%v>", iter.Kind)
}

// NewYamlIter ...
func NewYamlIter(node *yaml.Node) (iter YamlIter) {
	switch node.Kind {
	case yaml.MappingNode:
		return NewYamlMapIter(node.Content)
	case yaml.SequenceNode:
		return NewYamlListIter(node.Content)
	case yaml.ScalarNode:
		return NewYamlScalarIter(node)
	case yaml.DocumentNode:
		return NewYamlDocumentIter(node)
	default:
		return YamlIter{Iter: NewErrIter("Unexpected node type")}
	}
}

// NewYamlDocumentIter ...
func NewYamlDocumentIter(node *yaml.Node) YamlIter {
	docContent := node.Content[0]
	switch docContent.Kind {
	case yaml.MappingNode:
		return NewYamlMapIter(docContent.Content)
	case yaml.SequenceNode:
		return NewYamlListIter(docContent.Content)
	default:
		return YamlIter{Iter: NewErrIter("Yaml document: Unexpected node type")}
	}
}

// NewYamlScalarIter ...
func NewYamlScalarIter(node *yaml.Node) YamlIter {
	emitted := false
	return YamlIter{
		Kind: yaml.ScalarNode,
		Iter: Iter{
			Next: func() (n interface{}, err error) {
				if emitted {
					return n, ErrIterStop
				}
				emitted = true
				return node, err
			},
		},
	}
}

// NewYamlMapIter ...
func NewYamlMapIter(nodes []*yaml.Node) YamlIter {
	index, max := 0, len(nodes)
	return YamlIter{
		Kind: yaml.MappingNode,
		Iter: Iter{
			Next: func() (pair interface{}, err error) {
				if index >= max {
					return pair, ErrIterStop
				}
				pair = YamlMapPair{
					Key:   nodes[index],
					Value: NewYamlIter(nodes[index+1]),
				}
				index += 2
				return pair, err
			},
		},
	}
}

// NewYamlListIter ...
func NewYamlListIter(nodes []*yaml.Node) YamlIter {
	index, max := 0, len(nodes)
	return YamlIter{
		Kind: yaml.SequenceNode,
		Iter: Iter{
			Next: func() (n interface{}, err error) {
				if index >= max {
					return n, ErrIterStop
				}
				node := nodes[index]
				index++
				return node, err
			},
		},
	}
}
