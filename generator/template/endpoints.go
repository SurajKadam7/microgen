package template

import (
	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/parser"
	"github.com/devimteam/microgen/util"
)

type EndpointsTemplate struct {
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
//		context "context"
//		endpoint "github.com/go-kit/kit/endpoint"
//		)
//
//		type Endpoints struct {
//			CountEndpoint endpoint.Endpoint
//		}
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count
//		}
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{Count: count}, nil
//			}
//		}
//
func (EndpointsTemplate) Render(i *parser.Interface) *File {
	f := NewFile(i.PackageName)

	f.Type().Id("Endpoints").StructFunc(func(g *Group) {
		for _, signature := range i.FuncSignatures {
			g.Id(endpointStructName(signature.Name)).Qual(PackagePathGoKitEndpoint, "Endpoint")
		}
	})

	for _, signature := range i.FuncSignatures {
		f.Add(serviceEndpointMethod(signature)).Line() // .Line() means \n
	}
	f.Line() // Blank line
	for _, signature := range i.FuncSignatures {
		f.Add(createEndpoint(signature, i)).Line() // .Line() means \n
	}

	return f
}

func (EndpointsTemplate) Path() string {
	return "./endpoints.go"
}

// Render full endpoints method.
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count
//		}
//
func serviceEndpointMethod(signature *parser.FuncSignature) *Statement {
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
//		return resp.(*CountResponse).Count
//
func serviceEndpointMethodBody(signature *parser.FuncSignature) func(g *Group) {
	req := "req"
	resp := "resp"
	return func(g *Group) {
		//	req := CountRequest{
		//		Symbol: symbol,
		//		Text:   text,
		//	}
		g.Id(req).Op(":=").Id(requestStructName(signature)).Values(dictByFuncFields(signature.Params))
		//  resp, err := e.CountEndpoint(ctx, &req)
		g.List(Id(resp), Err()).Op(":=").Id(util.FirstLowerChar("Endpoint")).Dot(endpointStructName(signature.Name)).Call(Id(firstArgName(signature)), Op("&").Id(req))
		//  if err != nil {
		//	    return
		//  }
		g.If(Err().Op("!=").Nil()).Block(
			Return(),
		)
		//  return resp.(*CountResponse).Count, ...
		g.ReturnFunc(func(group *Group) {
			for _, field := range signature.Results {
				group.Id(resp).Assert(Op("*").Id(responseStructName(signature))).Op(".").Add(structFieldName(field))
			}
		})
	}
}

// For custom ctx in service interface (e.g. context or ctxxx).
func firstArgName(signature *parser.FuncSignature) string {
	return util.FirstLowerChar(signature.Params[0].Name)
}

// Render new Endpoint body.
//
//		return func(ctx context.Context, request interface{}) (interface{}, error) {
//			req := request.(*CountRequest)
//			count := svc.Count(ctx, req.Text, req.Symbol)
//			return &CountResponse{Count: count}, nil
//		}
//
func createEndpointBody(signature *parser.FuncSignature) *Statement {
	return Return(Func().Params(
		Id(firstArgName(signature)).Qual("context", "Context"),
		Id("request").Interface(),
	).Params(
		Interface(),
		Error(),
	).BlockFunc(func(g *Group) {
		g.Id("req").Op(":=").Id("request").Assert(Op("*").Id(requestStructName(signature)))
		g.Add(serviceMethodCallWithReceivers("svc", "req", signature))
		g.Return(
			Op("&").Id(responseStructName(signature)).Values(dictByFuncFields(signature.Results)),
			Nil(),
		)
	}))
}

// Render full new Endpoint function.
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{Count: count}, nil
//			}
//		}
//
func createEndpoint(signature *parser.FuncSignature, svcInterface *parser.Interface) *Statement {
	return Func().
		Id(endpointStructName(signature.Name)).Params(Id("svc").Id(svcInterface.Name)).Params(Qual(PackagePathGoKitEndpoint, "Endpoint")).
		Block(createEndpointBody(signature))
}
