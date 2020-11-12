package helloworldchaos

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/chaos-mesh/chaos-mesh/controllers/common"
	pb "github.com/chaos-mesh/chaos-mesh/pkg/chaosdaemon/pb"
	"github.com/chaos-mesh/chaos-mesh/pkg/router"
	ctx "github.com/chaos-mesh/chaos-mesh/pkg/router/context"
	end "github.com/chaos-mesh/chaos-mesh/pkg/router/endpoint"
	"github.com/chaos-mesh/chaos-mesh/pkg/utils"
)

// endpoint is dns-chaos reconciler
type endpoint struct {
	ctx.Context
}

// Apply applies helloworld chaos
func (r *endpoint) Apply(ctx context.Context, req ctrl.Request, chaos v1alpha1.InnerObject) error {
	r.Log.Info("Apply helloworld chaos")
	helloworldchaos, ok := chaos.(*v1alpha1.HelloWorldChaos)
	if !ok {
		return errors.New("chaos is not helloworldchaos")
	}

	pods, err := utils.SelectPods(ctx, r.Client, r.Reader, helloworldchaos.Spec.GetSelector())
	if err != nil {
		return err
	}

	for _, pod := range pods {
		daemonClient, err := utils.NewChaosDaemonClient(ctx, r.Client,
			&pod, common.ControllerCfg.ChaosDaemonPort)
		if err != nil {
			r.Log.Error(err, "get chaos daemon client")
			return err
		}
		defer daemonClient.Close()
		if len(pod.Status.ContainerStatuses) == 0 {
			return fmt.Errorf("%s %s can't get the state of container", pod.Namespace, pod.Name)
		}

		containerID := pod.Status.ContainerStatuses[0].ContainerID

		_, err = daemonClient.ExecHelloWorldChaos(ctx, &pb.ExecHelloWorldRequest{
			ContainerId: containerID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Recover means the reconciler recovers the chaos action
func (r *endpoint) Recover(ctx context.Context, req ctrl.Request, chaos v1alpha1.InnerObject) error {
	return nil
}

// Object would return the instance of chaos
func (r *endpoint) Object() v1alpha1.InnerObject {
	return &v1alpha1.HelloWorldChaos{}
}

func init() {
	router.Register("helloworldchaos", &v1alpha1.HelloWorldChaos{}, func(obj runtime.Object) bool {
		return true
	}, func(ctx ctx.Context) end.Endpoint {
		return &endpoint{
			Context: ctx,
		}
	})
}
