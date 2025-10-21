package parse

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newChatTemplate component initialization function of node 'InformationTemplate' in graph 'HCIBGAGetInformation'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	// TODO Modify component configuration here.
	config := &ChatTemplateConfig{}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}
