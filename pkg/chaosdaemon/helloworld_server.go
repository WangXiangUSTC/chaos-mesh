package chaosdaemon

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/chaos-mesh/chaos-mesh/pkg/bpm"
	pb "github.com/chaos-mesh/chaos-mesh/pkg/chaosdaemon/pb"
)

func (s *daemonServer) ExecHelloWorldChaos(ctx context.Context, req *pb.ExecHelloWorldRequest) (*empty.Empty, error) {
	log.Info("ExecHelloWorldChaos", "request", req)

	pid, err := s.crClient.GetPidFromContainerID(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	cmd := bpm.DefaultProcessBuilder("sh", "-c", fmt.Sprintf("echo 'hello' `hostname`")).
		SetNS(pid, bpm.UtsNS).
		SetContext(ctx).
		Build()
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(out) != 0 {
		log.Info("cmd output", "output", string(out))
	}

	return &empty.Empty{}, nil
}