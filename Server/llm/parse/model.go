package parse

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ChatModelImpl struct {
	config *ChatModelConfig
}

type ChatModelConfig struct {
}

// newChatModel component initialization function of node 'GenerateInformation' in graph 'HCIBGAGetInformation'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ChatModelConfig{}
	cm = &ChatModelImpl{config: config}
	return cm, nil
}

func (impl *ChatModelImpl) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	panic("implement me")
}

func (impl *ChatModelImpl) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	panic("implement me")
}

func (impl *ChatModelImpl) BindTools(tools []*schema.ToolInfo) error {
	panic("implement me")
}
