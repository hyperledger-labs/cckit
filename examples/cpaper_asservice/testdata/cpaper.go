package testdata

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/hyperledger-labs/cckit/examples/cpaper_asservice"
	testcc "github.com/hyperledger-labs/cckit/testing"
)

var (
	Id1 = &cpaper_asservice.CommercialPaperId{
		Issuer:      "SomeIssuer",
		PaperNumber: "0001",
	}

	ExternalId1 = &cpaper_asservice.ExternalId{
		Id: "EXT0001",
	}

	Issue1 = &cpaper_asservice.IssueCommercialPaper{
		Issuer:       Id1.Issuer,
		PaperNumber:  Id1.PaperNumber,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
		FaceValue:    100000,
		ExternalId:   ExternalId1.Id,
	}

	Buy1 = &cpaper_asservice.BuyCommercialPaper{
		Issuer:       Id1.Issuer,
		PaperNumber:  Id1.PaperNumber,
		CurrentOwner: Id1.Issuer,
		NewOwner:     "SomeBuyer",
		Price:        95000,
		PurchaseDate: ptypes.TimestampNow(),
	}

	Redeem1 = &cpaper_asservice.RedeemCommercialPaper{
		Issuer:         Id1.Issuer,
		PaperNumber:    Id1.PaperNumber,
		RedeemingOwner: Buy1.NewOwner,
		RedeemDate:     ptypes.TimestampNow(),
	}

	CpaperInState1 = &cpaper_asservice.CommercialPaper{
		Issuer:       Id1.Issuer,
		Owner:        Id1.Issuer,
		State:        cpaper_asservice.CommercialPaper_STATE_ISSUED,
		PaperNumber:  Id1.PaperNumber,
		FaceValue:    Issue1.FaceValue,
		IssueDate:    Issue1.IssueDate,
		MaturityDate: Issue1.MaturityDate,
		ExternalId:   Issue1.ExternalId,
	}
)
