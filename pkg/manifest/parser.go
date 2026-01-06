package manifest

import (
	"bytes"
	"errors"
	"io"

	"github.com/luaxlou/glow/pkg/api"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) ([]interface{}, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var docs []interface{}

	for {
		var node yaml.Node
		if err := decoder.Decode(&node); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		// Decode into a generic map to check Kind
		var typeMeta api.TypeMeta
		if err := node.Decode(&typeMeta); err != nil {
			return nil, err
		}

		switch typeMeta.Kind {
		case "Host":
			var host api.Host
			if err := node.Decode(&host); err != nil {
				return nil, err
			}
			docs = append(docs, host)
		case "App":
			var app api.App
			if err := node.Decode(&app); err != nil {
				return nil, err
			}
			docs = append(docs, app)
		default:
			// skip or error?
		}
	}
	return docs, nil
}
