package spans

import gen "github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger/gen-go/jaeger"

var tagNames = []string{"level", "stage", "node"}

func NewTags(tokens []string) []*gen.Tag {
	var tags []*gen.Tag
	for i, name := range tagNames {
		tag := &gen.Tag{
			Key:   name,
			VType: gen.TagType_STRING,
			VStr:  &tokens[i+6],
		}
		tags = append(tags, tag)
	}
	return tags
}
