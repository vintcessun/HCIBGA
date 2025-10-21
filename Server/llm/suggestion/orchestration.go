package suggestion

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildHCIBGASuggestion(ctx context.Context) (r compose.Runnable[PersonInformation, Plan], err error) {
	const (
		PersonTemplate     = "PersonTemplate"
		SuggestionModel    = "SuggestionModel"
		Search             = "Search"
		GeneratePlan       = "GeneratePlan"
		InformationLibrary = "InformationLibrary"
	)
	g := compose.NewGraph[PersonInformation, Plan]()
	personTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(PersonTemplate, personTemplateKeyOfChatTemplate)
	suggestionModelKeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(SuggestionModel, suggestionModelKeyOfChatModel)
	searchKeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(Search, searchKeyOfToolsNode)
	_ = g.AddLambdaNode(GeneratePlan, compose.CollectableLambda(newLambda))
	informationLibraryKeyOfToolsNode, err := newToolsNode1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(InformationLibrary, informationLibraryKeyOfToolsNode)
	_ = g.AddEdge(compose.START, PersonTemplate)
	_ = g.AddEdge(GeneratePlan, compose.END)
	_ = g.AddEdge(PersonTemplate, SuggestionModel)
	_ = g.AddEdge(SuggestionModel, Search)
	_ = g.AddEdge(SuggestionModel, GeneratePlan)
	_ = g.AddEdge(SuggestionModel, InformationLibrary)
	r, err = g.Compile(ctx, compose.WithGraphName("HCIBGASuggestion"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
