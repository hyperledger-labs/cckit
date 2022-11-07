package expect

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	g "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state/mapping"
	"github.com/hyperledger-labs/cckit/testing/gomega"
)

// EventIs expects ChaincodeEvent name is equal to expectName and event payload can be marshaled to expectPayload
func EventIs(event *peer.ChaincodeEvent, expectName string, expectPayload interface{}, converter serialize.FromBytesConverter) interface{} {
	g.Expect(event.EventName).To(g.Equal(expectName), `event name not match`)

	return EventPayloadIs(event, expectPayload, converter)
}

// EventStringerEqual expects ChaincodeEvent name is equal to expectName and
// event payload String() equal expectPayload String()
func EventStringerEqual(event *peer.ChaincodeEvent, expectName string, expectPayload interface{}, fromBytesConverters ...serialize.FromBytesConverter) {
	payload := EventIs(event, expectName, expectPayload, defaultFromBytesConverter(fromBytesConverters...))

	g.Expect(payload).To(gomega.StringerEqual(expectPayload))
}

// EventPayloadIs expects peer.ChaincodeEvent payload can be marshaled to
// target interface{} and returns converted value
func EventPayloadIs(event *peer.ChaincodeEvent, target interface{}, fromBytesConverters ...serialize.FromBytesConverter) interface{} {
	g.Expect(event).NotTo(g.BeNil())
	data, err := defaultFromBytesConverter(fromBytesConverters...).FromBytesTo(event.Payload, target)
	description := ``
	if err != nil {
		description = err.Error()
	}
	g.Expect(err).To(g.BeNil(), description)
	return data
}

// EventEqual expects that *peer.ChaincodeEvent stringer equal to mapping.Event
func EventEqual(event *peer.ChaincodeEvent, expected *mapping.Event, fromBytesConverters ...serialize.FromBytesConverter) {
	EventStringerEqual(event, expected.Name, expected.Payload, fromBytesConverters...)
}

// EventPayloadEqual checks event payload equality
func EventPayloadEqual(event *peer.ChaincodeEvent, expectedPayload interface{}, converter serialize.FromBytesConverter) {
	EventEqual(event, mapping.EventFromPayload(expectedPayload), converter)
}
