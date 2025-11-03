package third

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func newToolsNode(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	config := &compose.ToolsNodeConfig{}
	toolIns11, err := newTool(ctx)
	if err != nil {
		return nil, err
	}
	config.Tools = []tool.BaseTool{toolIns11}
	tsn, err = compose.NewToolNode(ctx, config)
	if err != nil {
		return nil, err
	}
	return tsn, nil
}

type ToolImpl struct {
	config *ToolConfig
}
type ToolConfig struct {
}

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	config := &ToolConfig{}
	bt = &ToolImpl{config: config}
	return bt, nil
}
func (impl *ToolImpl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}
func (impl *ToolImpl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}
func newToolsNode1(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	config := &compose.ToolsNodeConfig{}
	toolIns11, err := newTool1(ctx)
	if err != nil {
		return nil, err
	}
	config.Tools = []tool.BaseTool{toolIns11}
	tsn, err = compose.NewToolNode(ctx, config)
	if err != nil {
		return nil, err
	}
	return tsn, nil
}

type Tool1Impl struct {
	config *Tool1Config
}
type Tool1Config struct {
}

func newTool1(ctx context.Context) (bt tool.BaseTool, err error) {
	config := &Tool1Config{}
	bt = &Tool1Impl{config: config}
	return bt, nil
}
func (impl *Tool1Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}
func (impl *Tool1Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}
