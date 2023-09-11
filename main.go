package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"time"

	"buf.build/gen/go/twisp/api/grpc/go/twisp/core/v1/corev1grpc"
	corev1 "buf.build/gen/go/twisp/api/protocolbuffers/go/twisp/core/v1"
	typev1 "buf.build/gen/go/twisp/api/protocolbuffers/go/twisp/type/v1"
	"github.com/gofrs/uuid"
	"github.com/twisp/grpc/auth/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const TwispAccountID = "TwispAccountID"

func hi(uuid uuid.UUID) uint64 {
	return binary.BigEndian.Uint64(uuid[0:8])
}

func lo(uuid uuid.UUID) uint64 {
	return binary.BigEndian.Uint64(uuid[8:16])
}

func UUID(uuid uuid.UUID) *typev1.UUID {
	return &typev1.UUID{
		Hi: hi(uuid),
		Lo: lo(uuid),
	}
}

func main() {
	generator := token.NewTokenGeneratorIAM("cloud", "us-east-1")
	refresher, err := token.NewTokenRefresherTTL(generator, .95, .9, time.Now)
	if err != nil {
		panic(err)
	}
	defer refresher.Stop()

	conn, err := grpc.Dial("api.us-east-1.cloud.twisp.com:50051", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	// conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	journalClient := corev1grpc.NewJournalServiceClient(conn)

	token, _, err := refresher.Token()
	if err != nil {
		panic(err)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-twisp-account-id", TwispAccountID, "Authorization", "Bearer "+string(token))

	resp, err := journalClient.ListJournals(ctx, &corev1.ListJournalsRequest{
		Where: &corev1.ListJournalsRequest_Filters{
			Code: &corev1.FilterValue{
				Value: &corev1.FilterValue_Eq{
					Eq: "ACTIVE",
				},
			},
		},
		Index: corev1.ListJournalsRequest_INDEX_STATUS,
		Paging: &corev1.Paginate{
			First: 10,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", resp)

	// strP := func(str string) *string { return &str }

	// _, err = schemaClient.CreateIndex(ctx, &corev1.CreateIndexRequest{
	// 	Name:   "tx",
	// 	On:     corev1.Table_TABLE_TRANSACTION,
	// 	Unique: false,
	// 	Partition: []*corev1.PartitionKey{
	// 		{
	// 			Alias: "entityID",
	// 			Value: "string(document.metadata.entityID)",
	// 		},
	// 		{
	// 			Alias: "carrierEntityID",
	// 			Value: "string(document.metadata.carrierEntityID)",
	// 		},
	// 	},
	// 	Sort: []*corev1.IndexKey{
	// 		{
	// 			Alias: "creation",
	// 			Value: "document.created",
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// journalID := uuid.Must(uuid.NewV4())

	// journal, err := journalClient.CreateJournal(ctx, &corev1.CreateJournalRequest{
	// 	JournalId:   UUID(journalID),
	// 	Name:        "Journal",
	// 	Description: strP("a journal"),
	// 	Status:      corev1.JournalStatus_JOURNAL_STATUS_ACTIVE_UNSPECIFIED,
	// 	Code:        strP("TEST_JOURNAL"),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set1, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set1",
	// 	Description:       strP("set1 parent"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set1a, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set1.a",
	// 	Description:       strP("set1.a"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set1a1, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set1.a.1",
	// 	Description:       strP("set1.a.1"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set1a1a, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set1.a.1.a",
	// 	Description:       strP("set1.a.1.a"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set2, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set2",
	// 	Description:       strP("set2"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set2a, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set2.a",
	// 	Description:       strP("set2.a"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set3, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set3",
	// 	Description:       strP("set3"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// set3a, err := accountSetClient.CreateAccountSet(ctx, &corev1.CreateAccountSetRequest{
	// 	AccountSetId:      UUID(uuid.Must(uuid.NewV4())),
	// 	JournalId:         journal.Journal.JournalId,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Name:              "set3.a",
	// 	Description:       strP("set3.a"),
	// 	Config: &corev1.AccountSetConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set1.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_SET,
	// 	MemberId:     set1a.AccountSet.AccountSetId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set1a.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_SET,
	// 	MemberId:     set1a1.AccountSet.AccountSetId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set1a1.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_SET,
	// 	MemberId:     set1a1a.AccountSet.AccountSetId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set2.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_SET,
	// 	MemberId:     set2a.AccountSet.AccountSetId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set3.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_SET,
	// 	MemberId:     set3a.AccountSet.AccountSetId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// debtAcct, err := accountClient.CreateAccount(ctx, &corev1.CreateAccountRequest{
	// 	AccountId:         UUID(uuid.Must(uuid.FromString("32997079-82d9-5f13-b5ea-312fedd63cfc"))),
	// 	Name:              "debtAcct",
	// 	Code:              "debtAcct",
	// 	Description:       strP("debt"),
	// 	Status:            corev1.AccountStatus_ACCOUNT_STATUS_ACTIVE_UNSPECIFIED,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set1a1a.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_UNSPECIFIED,
	// 	MemberId:     debtAcct.Account.AccountId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set2a.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_UNSPECIFIED,
	// 	MemberId:     debtAcct.Account.AccountId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// carrierAcct, err := accountClient.CreateAccount(ctx, &corev1.CreateAccountRequest{
	// 	AccountId:         UUID(uuid.Must(uuid.FromString("39fd1d0f-6fc6-5313-b44a-eb941a85a3fb"))),
	// 	Name:              "carrierDebtAcct",
	// 	Code:              "carrierDebtAcct",
	// 	Description:       strP("debt"),
	// 	Status:            corev1.AccountStatus_ACCOUNT_STATUS_ACTIVE_UNSPECIFIED,
	// 	NormalBalanceType: corev1.DebitOrCredit_DEBIT_OR_CREDIT_CREDIT,
	// 	Config: &corev1.AccountConfig{
	// 		EnableConcurrentPosting: true,
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = accountSetClient.AddAccountSetMember(ctx, &corev1.AddAccountSetMemberRequest{
	// 	AccountSetId: set3a.AccountSet.AccountSetId,
	// 	MemberType:   corev1.AccountSetMemberType_ACCOUNT_SET_MEMBER_TYPE_ACCOUNT_UNSPECIFIED,
	// 	MemberId:     carrierAcct.Account.AccountId,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = trancodeClient.CreateTranCode(ctx, &corev1.CreateTranCodeRequest{
	// 	TranCodeId:  UUID(uuid.Must(uuid.NewV4())),
	// 	Code:        "ADVANCE_COMMISSION",
	// 	Description: strP("Advance Commission"),
	// 	Status:      corev1.TranCodeStatus_TRAN_CODE_STATUS_ACTIVE_UNSPECIFIED,
	// 	Params: []*corev1.TranCodeParam{
	// 		{
	// 			Name: "carrierDebtAcct",
	// 			Type: typev1.Type_TYPE_UUID,
	// 		},
	// 		{
	// 			Name: "debtAcct",
	// 			Type: typev1.Type_TYPE_UUID,
	// 		},
	// 		{
	// 			Name: "amount",
	// 			Type: typev1.Type_TYPE_DECIMAL,
	// 		},
	// 		{
	// 			Name: "correlation",
	// 			Type: typev1.Type_TYPE_STRING,
	// 		},
	// 		{
	// 			Name: "effective",
	// 			Type: typev1.Type_TYPE_DATE,
	// 		},
	// 		{
	// 			Name: "metadata",
	// 			Type: typev1.Type_TYPE_JSON,
	// 		},
	// 	},
	// 	Transaction: &corev1.TranCodeTransaction{
	// 		Effective:     "params.effective",
	// 		JournalId:     fmt.Sprintf("uuid('%s')", journalID.String()),
	// 		CorrelationId: "params.correlation",
	// 		Metadata:      "params.metadata",
	// 	},
	// 	Entries: []*corev1.TranCodeEntry{
	// 		{
	// 			EntryType: "'COMMISSION_ADVANCE_CR'",
	// 			AccountId: "params.debtAcct",
	// 			Layer:     "PENDING",
	// 			Direction: "CREDIT",
	// 			Units:     "params.amount",
	// 			Currency:  "'USD'",
	// 		},
	// 		{
	// 			EntryType: "'COMMISSION_ADVANCE_DR'",
	// 			AccountId: "params.carrierDebtAcct",
	// 			Layer:     "PENDING",
	// 			Direction: "DEBIT",
	// 			Units:     "params.amount",
	// 			Currency:  "'USD'",
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < 150; i++ {
	// 	st, _ := json.Marshal(map[string]any{
	// 		"payoutID":        "PAYOUT#1#1157#2023-05-03",
	// 		"month":           i,
	// 		"policyNumber":    "iwyuxz",
	// 		"entityID":        1157,
	// 		"carrierEntityID": "1",
	// 		"layer":           "PENDING",
	// 	})

	// 	params := map[string]string{
	// 		"debtAcct":        "32997079-82d9-5f13-b5ea-312fedd63cfc",
	// 		"carrierDebtAcct": "39fd1d0f-6fc6-5313-b44a-eb941a85a3fb",
	// 		"amount":          "10.00",
	// 		"effective":       time.Now().Format(time.DateOnly),
	// 		"correlation":     "",
	// 		"metadata":        string(st),
	// 	}

	// 	in := &corev1.PostTransactionRequest{
	// 		TransactionId: uuid.New(),
	// 		TranCode:      "ADVANCE_COMMISSION",
	// 		Params:        params,
	// 	}

	// 	fmt.Printf("step: doing step %d\n", i)
	// 	_, err := transactionClient.PostTransaction(ctx, in)
	// 	if err != nil {
	// 		st, ok := status.FromError(err)
	// 		if !ok {
	// 			panic(err)
	// 		}

	// 		if !strings.Contains(st.Message(), "unique constraint violation") {
	// 			panic(fmt.Errorf("got error on month %d: %w", i, err))
	// 		} else {
	// 			fmt.Printf("unique constraint violation %d\n", i)
	// 		}
	// 	}

	//}
}
