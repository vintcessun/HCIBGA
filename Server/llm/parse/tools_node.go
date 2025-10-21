package parse

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/bingsearch"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// newToolsNode component initialization function of node 'Search' in graph 'HCIBGAGetInformation'
func newToolsNode(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	// TODO Modify component configuration here.
	config := &compose.ToolsNodeConfig{}
	toolIns11, err := newTool(ctx)
	if err != nil {
		return nil, err
	}
	toolIns12, err := newTool1(ctx)
	if err != nil {
		return nil, err
	}
	toolIns13, err := newTool2(ctx)
	if err != nil {
		return nil, err
	}
	config.Tools = []tool.BaseTool{toolIns11, toolIns12, toolIns13}
	tsn, err = compose.NewToolNode(ctx, config)
	if err != nil {
		return nil, err
	}
	return tsn, nil
}

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &duckduckgo.Config{}
	bt, err = duckduckgo.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func newTool1(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &googlesearch.Config{}
	bt, err = googlesearch.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func newTool2(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &bingsearch.Config{}
	bt, err = bingsearch.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

// newToolsNode1 component initialization function of node 'InformationLibrary' in graph 'HCIBGAGetInformation'
func newToolsNode1(ctx context.Context) (tsn *compose.ToolsNode, err error) {
	// TODO Modify component configuration here.
	config := &compose.ToolsNodeConfig{}
	toolIns11, err := newTool3(ctx)
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

type Tool3Impl struct {
	config *Tool3Config
}

type Tool3Config struct {
}

func newTool3(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &Tool3Config{}
	bt = &Tool3Impl{config: config}
	return bt, nil
}

func (impl *Tool3Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}

func (impl *Tool3Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}
