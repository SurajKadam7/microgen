package template

import (
	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/parser"
	"github.com/devimteam/microgen/util"
)

type ClientTemplate struct {
}

// Renders whole client file.
//
//		// This file was automatically generated by "microgen" utility.
//		// Please, do not edit.
//		package stringsvc
//
//		import (
//			context "context"
//			transportlayer "github.com/devimteam/go-kit/transportlayer/grpc"
//		)
//
//		type client struct {
//			tc transportlayer.Client
//		}
//
//		func NewClient(tc transportlayer.Client) StringService {
//			return &client{tc}
//		}
//
//		func (c *client) Count(ctx context.Context, text string, symbol string) (count int, positions []int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := c.tc.Call(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//		}
//
func (ClientTemplate) Render(i *parser.Interface) *File {
	f := NewFile(i.PackageName)

	f.Type().Id("client").Struct(
		Id("tc").Op("*").Qual(PackagePathTransportLayerGRPC, "Client"),
	)

	f.Func().Id("NewClient").Call(Id("tc").Op("*").Qual(PackagePathTransportLayerGRPC, "Client")).Id(i.Name).Block(
		Return().Op("&").Id("client").Values(
			Id("tc"),
		),
	)
	f.Line()
	for _, signature := range i.FuncSignatures {
		f.Add(clientMethod(signature)).Line()
	}

	return f
}

func (ClientTemplate) Path() string {
	return "./client.go"
}

// Render full client method.
//
//		func (c *client) Count(ctx context.Context, text string, symbol string) (count int, positions []int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := c.tc.Call(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//		}
//
func clientMethod(signature *parser.FuncSignature) *Statement {
	return methodDefinition("client", signature).
		BlockFunc(clientMethodBody(signature))
}

// Render interface client method body.
//
//		req := CountRequest{
//			Symbol: symbol,
//			Text:   text,
//		}
//		resp, err := c.tc.Call(ctx, &req)
//		if err != nil {
//			return
//		}
//		return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
//
func clientMethodBody(signature *parser.FuncSignature) func(g *Group) {
	errName := getFirstErrorFieldName(signature.Results)
	return func(g *Group) {
		g.Id("req").Op(":=").Id(requestStructName(signature)).Values(dictByFuncFields(removeContextIfFirst(signature.Params)))
		g.List(Id("resp"), Id(errName)).Op(":=").Id(util.FirstLowerChar("client")).Dot("tc").Dot("Call").Call(Id(firstArgName(signature)), Op("&").Id("req"))
		g.If(Id(errName).Op("!=").Nil()).Block(
			Return(),
		)
		g.ReturnFunc(func(group *Group) {
			for _, field := range signature.Results {
				group.Id("resp").Assert(Op("*").Id(responseStructName(signature))).Op(".").Add(structFieldName(field))
			}
		})
	}
}

// Get from function field slice
// If field with error type not found, it return `err`
func getFirstErrorFieldName(fields []*parser.FuncField) string {
	for _, field := range fields {
		if field.Type == "error" {
			return field.Name
		}
	}
	return "err"
}
