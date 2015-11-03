package agent

import (
	"fmt"

	"github.com/gdamore/mangos/compat"
)

type Measurement struct {
	Name             string
	Unit             string
	Type             string
	Value            float64
	Host             string
	ResourceId       string
	Timestamp        float64
	ResourceMetadata map[string]string
}

type MetricPuller struct {
	Name     string
	ProbeId  string
	Interval float64
}

type MeteringAgent struct {
	MetricEndpoint      string
	MetricBus           *nanomsg.Socket
	PublicationEndpoint string
	PublicationBus      *nanomsg.Socket
}

func CreateAgent(metricAgentEndpoint string, metricPublicationEndpoint string) MeteringAgent {
	//
	// 1: Create a namomsg pull-based Message bus --> pull-based socket
	//
	// 2: Create a namomsg push-based Message bus --> Messages to be consumed by the published
	//
	// 3: Create an inward channel to consume the message bus incoming data
	//
	// 4: Create an outgoing channel onto which the processed metrics will be pushed through
	//
	// 5: Forward the incoming message to the pusblisher using goroutines
	//
	// TODO: do some more stuff here

	meteringAgent := MeteringAgent{metricAgentEndpoint, nil, metricPublicationEndpoint, nil}
	meteringAgent.MetricBus = meteringAgent.createMetricBus(metricAgentEndpoint)
	meteringAgent.PublicationBus = meteringAgent.createPublicationBus(metricPublicationEndpoint)

	return meteringAgent
}

func (m *MeteringAgent) createMetricBus(agentEndpoint string) *nanomsg.Socket {
	pullBus, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PULL)

	if err != nil {
		panic(err)
	}

	return pullBus
}

func (m *MeteringAgent) createPublicationBus(agentEndpoint string) *nanomsg.Socket {
	pushBus, err := nanomsg.NewSocket(nanomsg.AF_SP, nanomsg.PUSH)

	if err != nil {
		panic(err)
	}

	return pushBus
}

func (m *MeteringAgent) ReceiveMetric() <-chan Measurement {
	m.MetricBus.Connect(m.MetricEndpoint)
	agentChannel := make(chan Measurement, 8)
	go func() {
		for {
			data, err := m.MetricBus.Recv(0)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
			agentChannel <- Measurement{
				"Name",
				"Unit",
				"Type",
				12.34,
				"Host",
				"ResourceId",
				0.0,
				map[string]string{"title": "FAKE"},
			}
		}
	}()
	return agentChannel
}

func (m *MeteringAgent) SendMetric() chan<- Measurement {
	m.PublicationBus.Connect(m.PublicationEndpoint)
	agentChannel := make(chan Measurement, 8)
	go func() {
		for {
			data, err := m.PublicationBus.Recv(0)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
			// agentChannel <- data
		}
	}()
	return agentChannel

}

func (m *MeteringAgent) Run() {
	for measurement := range m.ReceiveMetric() {
		fmt.Println(measurement)
	}
}
