package parse

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildHCIBGAGetInformation(ctx context.Context) (r compose.Runnable[*[]schema.Message, MaterialInformation], err error) {
	const (
		InformationTemplate = "InformationTemplate"
		GenerateInformation = "GenerateInformation"
		Search              = "Search"
		Lambda4             = "Lambda4"
		InformationLibrary  = "InformationLibrary"
	)
	g := compose.NewGraph[*[]schema.Message, MaterialInformation]()
	informationTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(InformationTemplate, informationTemplateKeyOfChatTemplate)
	generateInformationKeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(GenerateInformation, generateInformationKeyOfChatModel)
	searchKeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(Search, searchKeyOfToolsNode)
	_ = g.AddLambdaNode(Lambda4, compose.InvokableLambda(newLambda))
	informationLibraryKeyOfToolsNode, err := newToolsNode1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(InformationLibrary, informationLibraryKeyOfToolsNode)
	_ = g.AddEdge(compose.START, InformationTemplate)
	_ = g.AddEdge(Lambda4, compose.END)
	_ = g.AddEdge(InformationTemplate, GenerateInformation)
	_ = g.AddEdge(GenerateInformation, Search)
	_ = g.AddEdge(GenerateInformation, InformationLibrary)
	_ = g.AddEdge(GenerateInformation, Lambda4)
	r, err = g.Compile(ctx, compose.WithGraphName("HCIBGAGetInformation"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
