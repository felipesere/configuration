package tags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type Tag struct {
	Key   string
	Value string
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:\"%s\"", t.Key, t.Value)
}

func ParseTags(field reflect.StructField) []Tag {
	var tags []Tag

	tagsStr := fmt.Sprintf("%s", field.Tag)

	for _, tagStr := range strings.Split(tagsStr, " ") {
		tagStr = strings.TrimSpace(tagStr)
		if len(tagStr) == 0 {
			continue
		}
		if tag, err := parseTagString(tagStr); err == nil {
			tags = append(tags, tag)
		}
	}

	return tags
}

func parseTagString(tag string) (Tag, error) {
	kv := strings.Split(tag, ":")
	if len(kv) != 2 {
		return Tag{}, errors.New("bad tag format")
	}

	return Tag{
		Key:   kv[0],
		Value: strings.Trim(kv[1], "\""),
	}, nil
}

func structTag(tags []Tag) reflect.StructTag {
	var tagsStr string
	for _, tag := range tags {
		tagsStr += ` ` + tag.String()
	}
	return reflect.StructTag(strings.TrimPrefix(tagsStr, " "))
}

func ChangedTagKeys(i interface{}, oldKey, newKey string) interface{} {
	oldValue := reflect.ValueOf(i)
	newType := changedTagKeys(i, oldKey, newKey)
	return reflect.NewAt(newType, unsafe.Pointer(oldValue.Pointer())).Interface()
}

func changedTagKeys(i interface{}, oldKey, newKey string) reflect.Type {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	var fields []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		vField := v.Field(i)

		if tField.Type.Kind() == reflect.Struct {
			tField.Type = changedTagKeys(vField.Addr().Interface(), oldKey, newKey)
			fields = append(fields, tField)
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			tField.Type = changedTagKeys(vField.Addr().Interface(), oldKey, newKey)
			fields = append(fields, tField)
			continue
		}

		tags := ParseTags(tField)
		tField.Tag = structTag(replaceKey(tags, oldKey, newKey))
		fields = append(fields, tField)
	}

	return reflect.StructOf(fields)
}

func replaceKey(tags []Tag, oldKey, newKey string) []Tag {
	for i := range tags {
		if tags[i].Key == oldKey {
			tags[i].Key = newKey
		}
	}
	return tags
}
