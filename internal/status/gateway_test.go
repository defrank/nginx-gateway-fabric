package status

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/nginxinc/nginx-kubernetes-gateway/internal/helpers"
	"github.com/nginxinc/nginx-kubernetes-gateway/internal/state"
)

func TestPrepareGatewayStatus(t *testing.T) {
	ipAddrType := v1beta1.IPAddressType
	podIP := v1beta1.GatewayAddress{
		Type:  &ipAddrType,
		Value: "1.2.3.4",
	}

	status := state.GatewayStatus{
		Conditions: CreateTestConditions("GatewayTest"),
		ListenerStatuses: state.ListenerStatuses{
			"listener": {
				AttachedRoutes: 3,
				Conditions:     CreateTestConditions("ListenerTest"),
			},
		},
		ObservedGeneration: 1,
	}

	transitionTime := metav1.NewTime(time.Now())

	expected := v1beta1.GatewayStatus{
		Conditions: CreateExpectedAPIConditions("GatewayTest", 1, transitionTime),
		Listeners: []v1beta1.ListenerStatus{
			{
				Name: "listener",
				SupportedKinds: []v1beta1.RouteGroupKind{
					{
						Kind: "HTTPRoute",
					},
				},
				AttachedRoutes: 3,
				Conditions:     CreateExpectedAPIConditions("ListenerTest", 1, transitionTime),
			},
		},
		Addresses: []v1beta1.GatewayAddress{podIP},
	}

	g := NewGomegaWithT(t)

	result := prepareGatewayStatus(status, "1.2.3.4", transitionTime)
	g.Expect(helpers.Diff(expected, result)).To(BeEmpty())
}
