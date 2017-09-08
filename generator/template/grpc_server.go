package template

import (
	"path/filepath"

	"github.com/devimteam/microgen/util"
	"github.com/vetcher/godecl/types"
	. "github.com/vetcher/jennifer/jen"
)

type GRPCServerTemplate struct {
	packageName        string
	ServicePackageName string
	PackagePath        string
}

func serverStructName(iface *types.Interface) string {
	return iface.Name + "Server"
}

func privateServerStructName(iface *types.Interface) string {
	return util.ToLowerFirst(iface.Name) + "Server"
}

func pathToConverter(servicePath string) string {
	return filepath.Join(servicePath, "transport/converter/protobuf")
}

// Render whole grpc server file.
//
//		// This file was automatically generated by "microgen" utility.
//		// Please, do not edit.
//		package transportgrpc
//
//		import (
//			transportlayer "github.com/devimteam/go-kit/transportlayer/grpc"
//			stringsvc "gitlab.devim.team/protobuf/stringsvc"
//			context "golang.org/x/net/context"
//		)
//
//		type server struct {
//			ts transportlayer.Server
//		}
//
//		func NewServer(endpoints []transportlayer.Endpoint) stringsvc.StringServiceServer {
//			return &server{transportlayer.NewServer(endpoints)}
//		}
//
//		func (s *server) Count(ctx context.Context, req *stringsvc.CountRequest) (*stringsvc.CountResponse, error) {
//			_, resp, err := s.ts.Serve(ctx, req)
//			if err != nil {
//				return nil, err
//			}
//			return resp.(*stringsvc.CountResponse), nil
//		}
//
func (t *GRPCServerTemplate) Render(i *types.Interface) *Statement {
	t.packageName = "transportgrpc"
	f := Statement{}

	f.Type().Id(privateServerStructName(i)).StructFunc(func(g *Group) {
		for _, method := range i.Methods {
			g.Id(util.ToLowerFirst(method.Name)).Qual(PackagePathGoKitTransportGRPC, "Handler")
		}
	}).Line()

	f.Func().Id("NewGRPCServer").
		Params(
			Id("endpoints").Op("*").Qual(t.PackagePath, "Endpoints"),
			Id("opts").Op("...").Qual(PackagePathGoKitTransportGRPC, "ServerOption"),
		).Params(
		Qual(protobufPath(t.ServicePackageName), serverStructName(i)),
	).
		Block(
			Return().Op("&").Id(privateServerStructName(i)).Values(DictFunc(func(g Dict) {
				for _, m := range i.Methods {
					g[(&Statement{}).Id(util.ToLowerFirst(m.Name))] = Qual(PackagePathGoKitTransportGRPC, "NewServer").
						Call(
							Line().Id("endpoints").Dot(endpointStructName(m.Name)),
							Line().Qual(pathToConverter(t.PackagePath), decodeRequestName(m)),
							Line().Qual(pathToConverter(t.PackagePath), encodeResponseName(m)),
							Line().Id("opts").Op("...").Line(),
						)
				}
			}),
			),
		)
	f.Line()

	for _, signature := range i.Methods {
		f.Line()
		f.Add(t.grpcServerFunc(signature, i)).Line()
	}

	return &f
}

func (GRPCServerTemplate) Path() string {
	return "./transport/grpc/server.go"
}

func (t *GRPCServerTemplate) PackageName() string {
	return t.packageName
}

// Render service interface method for grpc server.
//
//		func (s *server) Count(ctx context.Context, req *stringsvc.CountRequest) (*stringsvc.CountResponse, error) {
//			_, resp, err := s.ts.Serve(ctx, req)
//			if err != nil {
//				return nil, err
//			}
//			return resp.(*stringsvc.CountResponse), nil
//		}
//
func (t *GRPCServerTemplate) grpcServerFunc(signature *types.Function, i *types.Interface) *Statement {
	return Func().
		Params(Id(util.FirstLowerChar(privateServerStructName(i))).Op("*").Id(privateServerStructName(i))).
		Id(signature.Name).
		Call(Id("ctx").Qual(PackagePathNetContext, "Context"), Id("req").Op("*").Qual(protobufPath(t.ServicePackageName), requestStructName(signature))).
		Params(Op("*").Qual(protobufPath(t.ServicePackageName), responseStructName(signature)), Error()).
		BlockFunc(t.grpcServerFuncBody(signature, i))
}

// Render service method body for grpc server.
//
//		_, resp, err := s.ts.Serve(ctx, req)
//		if err != nil {
//			return nil, err
//		}
//		return resp.(*stringsvc.CountResponse), nil
//
func (t *GRPCServerTemplate) grpcServerFuncBody(signature *types.Function, i *types.Interface) func(g *Group) {
	return func(g *Group) {
		g.List(Id("_"), Id("resp"), Err()).
			Op(":=").
			Id(util.FirstLowerChar(privateServerStructName(i))).Dot(util.ToLowerFirst(signature.Name)).Dot("ServeGRPC").Call(Id("ctx"), Id("req"))

		g.If(Err().Op("!=").Nil()).Block(
			Return().List(Nil(), Err()),
		)

		g.Return().List(Id("resp").Assert(Op("*").Qual(protobufPath(t.ServicePackageName), responseStructName(signature))), Nil())
	}
}
