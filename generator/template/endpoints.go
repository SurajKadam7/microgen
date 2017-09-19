package template

import (
	"github.com/devimteam/microgen/generator/write_method"
	"github.com/devimteam/microgen/util"
	"github.com/vetcher/godecl/types"
	. "github.com/vetcher/jennifer/jen"
)

type endpointsTemplate struct {
	Info *GenerationInfo
}

func NewEndpointsTemplate(info *GenerationInfo) Template {
	infoCopy := info.Duplicate()
	infoCopy.Force = true
	return &endpointsTemplate{
		Info: infoCopy,
	}
}

func endpointStructName(str string) string {
	return str + "Endpoint"
}

// Renders endpoints file.
//
//		// This file was automatically generated by "microgen" utility.
//		// Please, do not edit.
//		package stringsvc
//
//		import (
//			context "context"
//			endpoint "github.com/go-kit/kit/endpoint"
//		)
//
//		type Endpoints struct {
//			CountEndpoint endpoint.Endpoint
//		}
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int, positions []int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//		}
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count, positions := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{
//					Count:     count,
//					Positions: positions,
//				}, nil
//			}
//		}
//
func (t *endpointsTemplate) Render(i *GenerationInfo) *Statement {
	f := Statement{}

	f.Type().Id("Endpoints").StructFunc(func(g *Group) {
		for _, signature := range i.Iface.Methods {
			g.Id(endpointStructName(signature.Name)).Qual(PackagePathGoKitEndpoint, "Endpoint")
		}
	}).Line()

	for _, signature := range i.Iface.Methods {
		f.Add(serviceEndpointMethod(signature)).Line().Line()
	}
	f.Line()
	for _, signature := range i.Iface.Methods {
		f.Add(createEndpoint(signature, i)).Line().Line()
	}

	return &f
}

func (endpointsTemplate) DefaultPath() string {
	return "./endpoints.go"
}

func (t *endpointsTemplate) ChooseMethod() (write_method.Method, error) {
	return write_method.NewFileMethod(t.Info.AbsOutPath, t.DefaultPath()), nil
}

// Render full endpoints method.
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int, positions []int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//		}
//
func serviceEndpointMethod(signature *types.Function) *Statement {
	return methodDefinition("Endpoints", signature).
		BlockFunc(serviceEndpointMethodBody(signature))
}

// Render interface method body.
//
//		req := CountRequest{
//			Symbol: symbol,
//			Text:   text,
//		}
//		resp, err := e.CountEndpoint(ctx, &req)
//		if err != nil {
//			return
//		}
//		return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//
func serviceEndpointMethodBody(signature *types.Function) func(g *Group) {
	return func(g *Group) {
		g.Id("_req").Op(":=").Id(requestStructName(signature)).Values(dictByVariables(removeContextIfFirst(signature.Args)))
		g.List(Id("_resp"), Err()).Op(":=").Id(util.FirstLowerChar("Endpoint")).Dot(endpointStructName(signature.Name)).Call(Id(firstArgName(signature)), Op("&").Id("_req"))
		g.If(Err().Op("!=").Nil()).Block(
			Return(),
		)
		g.ReturnFunc(func(group *Group) {
			for _, field := range signature.Results {
				group.Id("_resp").Assert(Op("*").Id(responseStructName(signature))).Op(".").Add(structFieldName(&field))
			}
		})
	}
}

// For custom ctx in service interface (e.g. context or ctxxx).
func firstArgName(signature *types.Function) string {
	return util.ToLowerFirst(signature.Args[0].Name)
}

// Render new Endpoint body.
//
//		return func(ctx context.Context, request interface{}) (interface{}, error) {
//			req := request.(*CountRequest)
//			count, positions := svc.Count(ctx, req.Text, req.Symbol)
//			return &CountResponse{
//				Count:     count,
//				Positions: positions,
//			}, nil
//		}
//
func createEndpointBody(signature *types.Function) *Statement {
	return Return(Func().Params(
		Id(firstArgName(signature)).Qual("context", "Context"),
		Id("request").Interface(),
	).Params(
		Interface(),
		Error(),
	).BlockFunc(func(g *Group) {
		methodParams := removeContextIfFirst(signature.Args)
		if len(methodParams) > 0 {
			g.Id("_req").Op(":=").Id("request").Assert(Op("*").Id(requestStructName(signature)))
		}

		g.Add(paramNames(signature.Results).
			Op(":=").
			Id("svc").
			Dot(signature.Name).
			CallFunc(func(g *Group) {
				g.Add(Id(firstArgName(signature)))
				for _, field := range methodParams {
					g.Add(Id("_req").Dot(util.ToUpperFirst(field.Name)))
				}
			}))

		g.Return(
			Op("&").Id(responseStructName(signature)).Values(dictByVariables(removeContextIfFirst(signature.Results))),
			Nil(),
		)
	}))
}

// Render full new Endpoint function.
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count, positions := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{
//					Count:     count,
//					Positions: positions,
//				}, nil
//			}
//		}
//
func createEndpoint(signature *types.Function, info *GenerationInfo) *Statement {
	return Func().
		Id(endpointStructName(signature.Name)).Params(Id("svc").Id(info.Iface.Name)).Params(Qual(PackagePathGoKitEndpoint, "Endpoint")).
		Block(createEndpointBody(signature))
}
