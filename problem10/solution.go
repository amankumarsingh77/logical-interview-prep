package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	log.Println(ParseQuery("key1=valueA%26valueB&key2=valueC"))
}

func ParseQuery(query string) (map[string]interface{}, error) {
	parsedQuery := make(map[string]interface{})

	segments := strings.Split(query, "&")
	for _, seg := range segments {
		if len(seg) == 0 {
			continue
		}

		parts := strings.SplitN(seg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed segment: %s", seg)
		}

		key, errKey := url.QueryUnescape(parts[0])
		val, errVal := url.QueryUnescape(parts[1])
		if errKey != nil || errVal != nil {
			return nil, fmt.Errorf("malformed segment: %s", seg)
		}
		if err := formatKeyVal(parsedQuery, key, val); err != nil {
			return nil, err
		}
	}

	return parsedQuery, nil
}

func formatKeyVal(root map[string]interface{}, key, val string) error {
	parts := strings.Split(key, ".")
	var current interface{} = root

	for i, part := range parts {
		isLast := i == len(parts)-1

		openBracket := strings.Index(part, "[")
		if openBracket != -1 && strings.HasSuffix(part, "]") {
			baseKey := part[:openBracket]
			idxStr := part[openBracket+1 : len(part)-1]

			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return fmt.Errorf("malformed index in part '%s': %w", part, err)
			}

			parentMap, ok := current.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected a map to set array key '%s', but got something else", baseKey)
			}

			if _, ok = parentMap[baseKey]; !ok {
				parentMap[baseKey] = []interface{}{}
			}
			slice, ok := parentMap[baseKey].([]interface{})
			if !ok {
				return fmt.Errorf("key '%s' exists but is not an array/slice", baseKey)
			}

			for len(slice) <= idx {
				slice = append(slice, nil)
			}

			if isLast {
				slice[idx] = val
				parentMap[baseKey] = slice
				return nil
			}

			if slice[idx] == nil {
				slice[idx] = make(map[string]interface{})
			}
			parentMap[baseKey] = slice
			current = slice[idx]
		} else {
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return fmt.Errorf("cannot set key '%s' on a non-map element", part)
			}

			if isLast {
				currentMap[part] = val
				return nil
			}

			if _, ok := currentMap[part]; !ok {
				currentMap[part] = make(map[string]interface{})
			}
			current = currentMap[part]
		}
	}

	return nil
}
