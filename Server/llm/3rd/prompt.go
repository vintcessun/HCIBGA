package third

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatTemplateImpl struct {
	config *ChatTemplateConfig
}
type ChatTemplateConfig struct {
}

func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{}
	ctp = &ChatTemplateImpl{config: config}
	return ctp, nil
}
func (impl *ChatTemplateImpl) Format(ctx context.Context, vs map[string]any, opts ...prompt.Option) ([]*schema.Message, error) {
	panic("implement me")
}
