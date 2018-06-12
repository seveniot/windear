package component

import (
	"github.com/SevenIOT/windear/agent"
	"github.com/SevenIOT/windear/core"
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/6
 *
 */

type AgentComponent struct {
	port   int
	svrCtx core.ServerContext
}

func NewAgentComponent(p int, ctx core.ServerContext) *AgentComponent {
	agent := &AgentComponent{
		port:   p,
		svrCtx: ctx,
	}

	return agent
}

func (a *AgentComponent) Start() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		log.Fatalf("failed to start agent: %v", err.Error())
	}

	s := grpc.NewServer()
	agent.RegisterRemoteAgentServer(s, a) //pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Infof("start remote agent server at port:%v", a.port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve agent: %v", err.Error())
	}
}

func (a *AgentComponent) KickRemote(clientId, hostAddr string) error {
	conn, err := grpc.Dial(hostAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := agent.NewRemoteAgentClient(conn) //pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.KickClient(ctx, &agent.Request{ClientId: clientId})
	if err != nil {
		return err
	}

	if r.Status == agent.Response_SUCCESS {
		return nil
	}

	return ex.AgentReturnFail
}

func (a *AgentComponent) PubMsgRemote(clientId string, content []byte, hostAddr string) error {
	conn, err := grpc.Dial(hostAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := agent.NewRemoteAgentClient(conn) //pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PubMsg(ctx, &agent.Request{ClientId: clientId, Content: content})
	if err != nil {
		return err
	}

	if r.Status == agent.Response_SUCCESS {
		return nil
	}

	return ex.AgentReturnFail
}

func (a *AgentComponent) KickClient(ctx context.Context, req *agent.Request) (*agent.Response, error) {
	log.Infof("receive kick client request:%v", req.ClientId)

	a.svrCtx.KickClient(req.ClientId)

	return &agent.Response{Status: agent.Response_SUCCESS}, nil
}

func (a *AgentComponent) PubMsg(ctx context.Context, req *agent.Request) (*agent.Response, error) {
	a.svrCtx.PubMsg(req.ClientId, req.Content)

	return &agent.Response{Status: agent.Response_SUCCESS}, nil
}
