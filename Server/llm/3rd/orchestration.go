package third

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildHCIBGA3rdSource(ctx context.Context) (r compose.Runnable[any, any], err error) {
	const (
		ChooseTemplate = "ChooseTemplate"
		ChooseModel    = "ChooseModel"
		QicaiTool      = "QicaiTool"
		JwTool         = "JwTool"
	)
	g := compose.NewGraph[any, any]()
	chooseTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChooseTemplate, chooseTemplateKeyOfChatTemplate)
	chooseModelKeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChooseModel, chooseModelKeyOfChatModel)
	qicaiToolKeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(QicaiTool, qicaiToolKeyOfToolsNode)
	jwToolKeyOfToolsNode, err := newToolsNode1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(JwTool, jwToolKeyOfToolsNode)
	_ = g.AddEdge(compose.START, ChooseTemplate)
	_ = g.AddEdge(ChooseModel, compose.END)
	_ = g.AddEdge(ChooseTemplate, ChooseModel)
	_ = g.AddEdge(ChooseModel, QicaiTool)
	_ = g.AddEdge(ChooseModel, JwTool)
	r, err = g.Compile(ctx, compose.WithGraphName("HCIBGA3rdSource"))
	if err != nil {
		return nil, err
	}
	return r, err
}
