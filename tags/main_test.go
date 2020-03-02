package tags

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestName(t *testing.T) {
	obj1 := struct {
		Name string
	}{}

	//t.Logf("obj type: %#v", ChangedTagKeys(&obj1, "key", "json"))

	//bb := []byte(`{"name": "name_json"}`)
	bb := []byte(`name: name_yaml`)

	err := yaml.Unmarshal(bb, &obj1)
	t.Log(err)
	t.Log(obj1)
}
