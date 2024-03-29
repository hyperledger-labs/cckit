package mapping_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	identitytestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
	"github.com/hyperledger-labs/cckit/state/mapping"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata/schema"
	state_schema "github.com/hyperledger-labs/cckit/state/schema"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var _ = Describe(`State mapping in chaincode with default serializer`, func() {

	var (
		compositeIDCC, complexIDCC, sliceIDCC, indexesCC, configCC *testcc.MockStub
		Owner                                                      = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
	)

	Describe(`init chaincodes`, func() {
		compositeIDCC = testcc.NewMockStub(`proto`, testdata.NewCompositeIdCC(serialize.DefaultSerializer))
		compositeIDCC.From(Owner).Init()

		complexIDCC = testcc.NewMockStub(`complex_id`, testdata.NewComplexIdCC())
		complexIDCC.From(Owner).Init()

		sliceIDCC = testcc.NewMockStub(`slice_id`, testdata.NewSliceIdCC())
		sliceIDCC.From(Owner).Init()

		indexesCC = testcc.NewMockStub(`indexes`, testdata.NewIndexesCC())
		indexesCC.From(Owner).Init()

		configCC = testcc.NewMockStub(`config`, testdata.NewCCWithConfig())
		configCC.From(Owner).Init()
	})

	Describe(`Entity with composite id`, func() {
		create1 := testdata.CreateEntityWithCompositeId[0]
		create2 := testdata.CreateEntityWithCompositeId[1]
		create3 := testdata.CreateEntityWithCompositeId[2]

		It("Allow to get mapping data by namespace", func() {
			mapping, err := testdata.EntityWithCompositeIdStateMapping.GetByNamespace(testdata.EntityCompositeIdNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(mapping.Schema()).String()).To(
				Equal(reflect.TypeOf(&schema.EntityWithCompositeId{}).String()))

			key, err := mapping.PrimaryKey(&schema.EntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(key).To(Equal(
				testdata.EntityCompositeIdNamespace.Append(
					state.Key{create1.IdFirstPart, strconv.Itoa(int(create1.IdSecondPart)), testdata.Dates[0]})))
		})

		It("Allow to add data to chaincode state", func(done Done) {
			events, closer := compositeIDCC.EventSubscription()
			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, create1))

			expectcc.EventStringerEqual(<-events,
				`CreateEntityWithCompositeId`, create1, compositeIDCC.Serializer)

			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, create2))
			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, create3))

			_ = closer()
			close(done)
		})

		It("Disallow to insert entries with same primary key", func() {
			expectcc.ResponseError(compositeIDCC.Invoke(testdata.CreateFunc, create1), state.ErrKeyAlreadyExists)
		})

		It("Allow to get entry list", func() {
			// default serializer should serialize proto to / from binary representation
			entities := expectcc.PayloadIs(compositeIDCC.Query(testdata.ListFunc),
				&schema.EntityWithCompositeIdList{}, serialize.DefaultSerializer).(*schema.EntityWithCompositeIdList)
			Expect(len(entities.Items)).To(Equal(3))
			Expect(entities.Items[0].Name).To(Equal(create1.Name))
			Expect(entities.Items[0].Value).To(BeNumerically("==", create1.Value))
		})

		It("Allow to get entry raw protobuf", func() {
			dataFromCC := compositeIDCC.Query(testdata.GetFunc,
				&schema.EntityCompositeId{
					IdFirstPart:  create1.IdFirstPart,
					IdSecondPart: create1.IdSecondPart,
					IdThirdPart:  create1.IdThirdPart,
				},
			).Payload

			e := &schema.EntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,

				Name:  create1.Name,
				Value: create1.Value,
			}
			Expect(dataFromCC).To(Equal(testcc.MustProtoMarshal(e)))
		})

		It("Allow update data in chaincode state", func() {
			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.UpdateFunc, &schema.UpdateEntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
				Name:         `New name`,
				Value:        1000,
			}))

			entityFromCC := expectcc.PayloadIs(
				compositeIDCC.Query(testdata.GetFunc, &schema.EntityCompositeId{
					IdFirstPart:  create1.IdFirstPart,
					IdSecondPart: create1.IdSecondPart,
					IdThirdPart:  create1.IdThirdPart,
				}),
				&schema.EntityWithCompositeId{}, compositeIDCC.Serializer).(*schema.EntityWithCompositeId)

			// state is updated
			Expect(entityFromCC.Name).To(Equal(`New name`))
			Expect(entityFromCC.Value).To(BeNumerically("==", 1000))
		})

		It("Allow to delete entry", func() {
			toDelete := &schema.EntityCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
			}

			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.DeleteFunc, toDelete))
			ee := expectcc.PayloadIs(
				compositeIDCC.Invoke(testdata.ListFunc),
				&schema.EntityWithCompositeIdList{}, compositeIDCC.Serializer).(*schema.EntityWithCompositeIdList)

			Expect(len(ee.Items)).To(Equal(2))
			expectcc.ResponseError(compositeIDCC.Invoke(testdata.GetFunc, toDelete), state.ErrKeyNotFound)
		})

		It("Allow to insert entry once more time", func() {
			expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, create1))
		})
	})

	Describe(`Entity with complex id`, func() {
		ent1 := testdata.CreateEntityWithComplextId[0]

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(complexIDCC.From(Owner).Invoke(`entityInsert`, ent1))
			keys := expectcc.PayloadIs(complexIDCC.From(Owner).Invoke(
				`debugStateKeys`, `EntityWithComplexId`), &[]string{}, complexIDCC.Serializer).([]string)
			Expect(len(keys)).To(Equal(1))

			timeStr := time.Unix(ent1.Id.IdPart3.GetSeconds(), int64(ent1.Id.IdPart3.GetNanos())).Format(`2006-01-02`)
			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal(
				string(rune(0)) +
					`EntityWithComplexId` + string(rune(0)) +
					ent1.Id.IdPart1[0] + string(rune(0)) +
					ent1.Id.IdPart1[1] + string(rune(0)) +
					ent1.Id.IdPart2 + string(rune(0)) +
					timeStr + string(rune(0))))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(complexIDCC.Query(`entityGet`, ent1.Id)).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent1)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(complexIDCC.Query(`entityList`),
				&state_schema.List{}, complexIDCC.Serializer).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent1)))
		})
	})

	Describe(`Entity with slice id`, func() {

		ent2 := &schema.EntityWithSliceId{Id: []string{`aa`, `bb`}, SomeDate: ptypes.TimestampNow()}

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(sliceIDCC.Invoke(`entityInsert`, ent2))
			keys := expectcc.PayloadIs(sliceIDCC.From(Owner).Invoke(
				`debugStateKeys`, `EntityWithSliceId`), &[]string{}, sliceIDCC.Serializer).([]string)

			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal(
				"\x00" + `EntityWithSliceId` + string(rune(0)) + ent2.Id[0] + string(rune(0)) + ent2.Id[1] + string(rune(0))))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(sliceIDCC.Query(`entityGet`, state.StringsIdToStr(ent2.Id))).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent2)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(
				sliceIDCC.Query(`entityList`), &state_schema.List{}, sliceIDCC.Serializer).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent2)))
		})
	})

	Describe(`Entity with indexes`, func() {

		create1 := testdata.CreateEntityWithIndexes[0]
		create2 := testdata.CreateEntityWithIndexes[1]

		It("Allow to add data with single external id", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create1))
		})

		It("Disallow to add data to chaincode state with same uniq key fields", func() {
			createWithNewId := proto.Clone(create1).(*schema.CreateEntityWithIndexes)
			createWithNewId.Id = `abcdef` // id is really new

			// errored on checking uniq key
			expectcc.ResponseError(
				indexesCC.Invoke(`create`, create1),
				mapping.ErrMappingUniqKeyExists)
		})

		It("Allow finding data by uniq key", func() {
			fromCCByExtId := expectcc.PayloadIs(
				indexesCC.Query(`getByExternalId`, create1.ExternalId),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			fromCCById := expectcc.PayloadIs(
				indexesCC.Query(`get`, create1.Id),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId).To(BeEquivalentTo(fromCCById))
		})

		It("Allow to get idx state key by uniq key", func() {
			idxKey, err := testdata.EntityWithIndexesStateMapping.IdxKey(
				&schema.EntityWithIndexes{}, `ExternalId`, []string{create1.ExternalId}, indexesCC.Serializer)
			Expect(err).NotTo(HaveOccurred())

			Expect(idxKey).To(BeEquivalentTo([]string{
				mapping.KeyRefNamespace,
				strings.Join(mapping.SchemaNamespace(&schema.EntityWithIndexes{}), `-`),
				`ExternalId`,
				create1.ExternalId,
			}))
		})

		It("Disallow finding data by non existent uniq key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByExternalId`, `some-non-existent-id`),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to add data with multiple external id", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create2))
		})

		It("Allow to find data by multi key", func() {
			fromCCByExtId1 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[0]),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			fromCCByExtId2 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[1]),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			fromCCById := expectcc.PayloadIs(
				indexesCC.Query(`get`, create2.Id),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId1).To(BeEquivalentTo(fromCCById))
			Expect(fromCCByExtId2).To(BeEquivalentTo(fromCCById))
		})

		It("Allow update indexes value", func() {
			update2 := &schema.UpdateEntityWithIndexes{
				Id:                  create2.Id,
				ExternalId:          `some_new_external_id`,
				OptionalExternalIds: []string{create2.OptionalExternalIds[0], `AND SOME NEW`},
			}
			expectcc.ResponseOk(indexesCC.Invoke(`update`, update2))
		})

		It("Allow to find data by updated multi key", func() {
			fromCCByExtId1 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[0]),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			fromCCByExtId2 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, `AND SOME NEW`),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId1.Id).To(Equal(create2.Id))
			Expect(fromCCByExtId2.Id).To(Equal(create2.Id))

			Expect(fromCCByExtId2.OptionalExternalIds).To(
				BeEquivalentTo([]string{create2.OptionalExternalIds[0], `AND SOME NEW`}))
		})

		It("Disallow to find data by previous multi key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[1]),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to find data by updated uniq key", func() {
			fromCCByExtId := expectcc.PayloadIs(
				indexesCC.Query(`getByExternalId`, `some_new_external_id`),
				&schema.EntityWithIndexes{}, indexesCC.Serializer).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId.Id).To(Equal(create2.Id))
			Expect(fromCCByExtId.ExternalId).To(Equal(`some_new_external_id`))
		})

		It("Disallow to find data by previous uniq key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByExternalId`, create2.ExternalId),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to delete entry", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`delete`, create2.Id))

			ee := expectcc.PayloadIs(
				indexesCC.Invoke(`list`),
				&schema.EntityWithIndexesList{}, indexesCC.Serializer).(*schema.EntityWithIndexesList)

			Expect(len(ee.Items)).To(Equal(1))
			expectcc.ResponseError(indexesCC.Invoke(`get`, create2.Id), state.ErrKeyNotFound)
		})

		It("Allow to insert entry once more time", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create2))
		})

	})

	Describe(`Entity with static key`, func() {
		configSample := &schema.Config{
			Field1: `aaa`,
			Field2: `bbb`,
		}

		It("Disallow to get config before set", func() {
			expectcc.ResponseError(configCC.Invoke(`configGet`),
				`get state with key=Config: state entry not found`)
		})

		It("Allow to set config", func() {
			expectcc.ResponseOk(configCC.Invoke(`configSet`, configSample))
		})

		It("Allow to get config", func() {
			confFromCC := expectcc.PayloadIs(
				configCC.Invoke(`configGet`), &schema.Config{}, configCC.Serializer).(*schema.Config)
			Expect(confFromCC.Field1).To(Equal(configSample.Field1))
			Expect(confFromCC.Field2).To(Equal(configSample.Field2))
		})

	})
})

var _ = Describe(`State mapping in chaincode with JSON serializer`, func() {

	var (
		compositeIDCC *testcc.MockStub
		Owner         = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
	)

	Describe("init chaincode", func() {
		compositeIDCC = testcc.NewMockStub(`proto`, testdata.NewCompositeIdCC(serialize.PreferJSONSerializer))
		compositeIDCC.Serializer = serialize.PreferJSONSerializer // need to set for correct invoking
		compositeIDCC.From(Owner).Init()
	})

	It("Allow to add data to chaincode state", func() {
		expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, testdata.CreateEntityWithCompositeId[0]))
	})

	It("Allow to get entry list", func() {
		res := compositeIDCC.Query(testdata.ListFunc)
		Expect(string(res.Payload)[0:1]).To(Equal(`{`)) // json serialized
		entities := expectcc.JSONPayloadIs(res, &schema.EntityWithCompositeIdList{}).(*schema.EntityWithCompositeIdList)
		Expect(len(entities.Items)).To(Equal(1))
		Expect(entities.Items[0].Name).To(Equal(testdata.CreateEntityWithCompositeId[0].Name))
		Expect(entities.Items[0].Value).To(BeNumerically("==", testdata.CreateEntityWithCompositeId[0].Value))
	})
})
