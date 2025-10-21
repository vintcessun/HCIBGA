package lib

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildHCIBGAInformationLibrary(ctx context.Context) (r compose.Runnable[[]*schema.Message, []*schema.Message], err error) {
	const (
		LocalInformation            = "LocalInformation"
		GetLocalInformationTemplate = "GetLocalInformationTemplate"
	)
	g := compose.NewGraph[[]*schema.Message, []*schema.Message]()
	localInformationKeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(LocalInformation, localInformationKeyOfToolsNode)
	getLocalInformationTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(GetLocalInformationTemplate, getLocalInformationTemplateKeyOfChatTemplate)
	_ = g.AddEdge(compose.START, GetLocalInformationTemplate)
	_ = g.AddEdge(GetLocalInformationTemplate, compose.END)
	_ = g.AddEdge(GetLocalInformationTemplate, LocalInformation)
	r, err = g.Compile(ctx, compose.WithGraphName("HCIBGAInformationLibrary"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
