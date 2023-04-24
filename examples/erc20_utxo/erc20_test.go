package erc20_utxo_test

import (
	"encoding/base64"
	"math/big"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/examples/erc20_utxo"
	"github.com/hyperledger-labs/cckit/examples/erc20_utxo/service/allowance"
	"github.com/hyperledger-labs/cckit/examples/erc20_utxo/service/config"
	"github.com/hyperledger-labs/cckit/extensions/account"
	"github.com/hyperledger-labs/cckit/extensions/token"
	"github.com/hyperledger-labs/cckit/extensions/token/decimal"
	"github.com/hyperledger-labs/cckit/identity"
	"github.com/hyperledger-labs/cckit/identity/testdata"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
)

func TestERC20(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ERC20 Test suite")
}

var (
	ownerIdentity = testdata.Certificates[0].MustIdentity(testdata.DefaultMSP)
	user1Identity = testdata.Certificates[1].MustIdentity(testdata.DefaultMSP)
	user2Identity = testdata.Certificates[2].MustIdentity(testdata.DefaultMSP)

	ownerAddress = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(ownerIdentity.Cert.PublicKey))
	user1Address = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(user1Identity.Cert.PublicKey))
	user2Address = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(user2Identity.Cert.PublicKey))

	cc *testcc.MockStub
)

var _ = Describe(`ERC`, func() {

	BeforeSuite(func() {
		chaincode, err := erc20_utxo.NewChaincode()
		Expect(err).NotTo(HaveOccurred())
		cc = testcc.NewMockStub(`erc20`, chaincode)

		expectcc.ResponseOk(cc.From(ownerIdentity).Init())
	})

	It(`Allow to call init once more time `, func() {
		expectcc.ResponseOk(cc.From(ownerIdentity).Init())
	})

	Context(`token info`, func() {

		It(`Allow to get token name`, func() {
			name := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(config.ConfigServiceChaincode_GetName, nil),
				&config.NameResponse{}).(*config.NameResponse)

			Expect(name.Name).To(Equal(erc20_utxo.Token.Name))
		})
	})

	Context(`initial balance`, func() {

		It(`Allow to know invoker address `, func() {
			address := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(account.AccountServiceChaincode_GetInvokerAddress, nil),
				&account.AddressId{}).(*account.AddressId)

			Expect(address.Address).To(Equal(user1Address))

			address = expectcc.PayloadIs(
				cc.From(user2Identity).
					Query(account.AccountServiceChaincode_GetInvokerAddress, nil),
				&account.AddressId{}).(*account.AddressId)

			Expect(address.Address).To(Equal(user2Address))
		})

		It(`Allow to get owner balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity). // call by any user
							Query(token.BalanceServiceChaincode_GetBalance,
						&token.BalanceId{Address: ownerAddress, Symbol: erc20_utxo.Token.Symbol}),
				&token.Balance{}).(*token.Balance)

			Expect(b.Address).To(Equal(ownerAddress))
			Expect(b.Amount).To(Equal(erc20_utxo.Token.TotalSupply))
		})

		It(`Allow to get zero balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(token.BalanceServiceChaincode_GetBalance,
						&token.BalanceId{Address: user1Address, Symbol: erc20_utxo.Token.Symbol}),
				&token.Balance{}).(*token.Balance)

			Expect(b.Amount).To(Equal(decimal.New(big.NewInt(0))))
		})

	})

	Context(`transfer`, func() {
		var transferAmount = decimal.New(big.NewInt(100))

		It(`Disallow to transfer balance by user with zero balance`, func() {
			expectcc.ResponseError(
				cc.From(user1Identity).
					Invoke(token.BalanceServiceChaincode_Transfer,
						&token.TransferRequest{
							Recipient: user2Address,
							Symbol:    erc20_utxo.Token.Symbol,
							Amount:    transferAmount,
						}), token.ErrAmountInsuficcient)

		})

		It(`Allow to transfer balance by owner`, func() {
			r := expectcc.PayloadIs(
				cc.From(ownerIdentity).
					Invoke(token.BalanceServiceChaincode_Transfer,
						&token.TransferRequest{
							Recipient: user1Address,
							Symbol:    erc20_utxo.Token.Symbol,
							Amount:    transferAmount,
						}),
				&token.TransferResponse{}).(*token.TransferResponse)

			Expect(r.Sender).To(Equal(ownerAddress))
			Expect(r.Amount).To(Equal(transferAmount))
		})

		It(`Allow to get new non zero balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(token.BalanceServiceChaincode_GetBalance,
						&token.BalanceId{
							Address: user1Address,
							Symbol:  erc20_utxo.Token.Symbol,
						}),
				&token.Balance{}).(*token.Balance)

			Expect(b.Amount).To(Equal(transferAmount))
		})

	})

	Context(`Allowance`, func() {

		var allowAmount = decimal.New(big.NewInt(50))

		It(`Allow to approve amount by owner for spender even if balance is zero`, func() {
			a := expectcc.PayloadIs(
				cc.From(user2Identity).
					Invoke(allowance.AllowanceServiceChaincode_Approve,
						&allowance.ApproveRequest{
							Owner:   user2Address,
							Spender: user1Address,
							Symbol:  erc20_utxo.Token.Symbol,
							Amount:  allowAmount,
						}),
				&allowance.Allowance{}).(*allowance.Allowance)

			Expect(a.Owner).To(Equal(user2Address))
			Expect(a.Spender).To(Equal(user1Address))
			Expect(a.Amount).To(Equal(allowAmount))
		})
		It(`Disallow to approve amount by non owner`, func() {
			expectcc.ResponseError(
				cc.From(user2Identity).
					Invoke(allowance.AllowanceServiceChaincode_Approve,
						&allowance.ApproveRequest{
							Owner:   ownerAddress,
							Spender: user1Address,
							Symbol:  erc20_utxo.Token.Symbol,
							Amount:  allowAmount,
						}), allowance.ErrOwnerOnly)
		})

		It(`Allow to approve amount by owner for spender if amount is sufficient`, func() {
			a := expectcc.PayloadIs(
				cc.From(ownerIdentity).
					Invoke(allowance.AllowanceServiceChaincode_Approve,
						&allowance.ApproveRequest{
							Owner:   ownerAddress,
							Spender: user2Address,
							Symbol:  erc20_utxo.Token.Symbol,
							Amount:  allowAmount,
						}),
				&allowance.Allowance{}).(*allowance.Allowance)

			Expect(a.Owner).To(Equal(ownerAddress))
			Expect(a.Spender).To(Equal(user2Address))
			Expect(a.Amount).To(Equal(allowAmount))
		})

		It(`Allow to transfer from`, func() {
			spenderIdentity := user2Identity
			spenderAddress := user2Address

			t := expectcc.PayloadIs(
				cc.From(spenderIdentity).
					Invoke(allowance.AllowanceServiceChaincode_TransferFrom,
						&allowance.TransferFromRequest{
							Owner:     ownerAddress,
							Recipient: spenderAddress,
							Symbol:    erc20_utxo.Token.Symbol,
							Amount:    allowAmount,
						}),
				&allowance.TransferFromResponse{}).(*allowance.TransferFromResponse)

			Expect(t.Owner).To(Equal(ownerAddress))
			Expect(t.Recipient).To(Equal(spenderAddress))
			Expect(t.Amount).To(Equal(allowAmount))
		})
	})
})
