package protocol

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"log"
	"time"

	"github.com/peterchanxyz/protocol/crypto"
	"github.com/peterchanxyz/protocol/gen/go/api"
	"github.com/peterchanxyz/protocol/gen/go/core"
	"github.com/peterchanxyz/protocol/gen/go/core/contract"
	"github.com/peterchanxyz/protocol/util/base58"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func main() {

	addr := "api.trongrid.io:50051"

	gc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	tronClient := api.NewWalletClient(gc)

	//transfer

	//to account address
	toAddress := ""

	amount := 1

	tronClient.TransferAsset()

}

func tronTransfer(ctx context.Context, cli api.WalletClient, ownerKey *ecdsa.PrivateKey, toAddress string,
	amount int64) *api.Return {

	transferContract := &contract.TransferContract{}
	transferContract.OwnerAddress = crypto.PubkeyToAddress(ownerKey.PublicKey).Bytes()
	transferContract.ToAddress = base58.DecodeCheck(toAddress)

	transferTransaction, err := cli.CreateTransaction(ctx, transferContract)
	if err != nil {
		log.Fatalf("transfer error: %v", err)
	}

	if transferTransaction == nil || len(transferTransaction.
		GetRawData().GetContract()) == 0 {
		log.Fatalf("transfer error: invalid transaction")
	}

	signTransaction(transferTransaction, ownerKey)

	result, err := cli.BroadcastTransaction(ctx, transferTransaction)

	if err != nil {
		log.Fatalf("transfer error: %v", err)
	}

	return result
}

func withdrawBalance(ctx context.Context, cli api.WalletClient, ownerKey *ecdsa.PrivateKey, toAddress string,
	amount int64) *api.Return {

	withdrawBalanceContract := new(contract.WithdrawBalanceContract)
	withdrawBalanceContract.OwnerAddress = crypto.PubkeyToAddress(ownerKey.PublicKey).Bytes()

	withdrawBalanceTransaction, err := cli.WithdrawBalance(ctx, withdrawBalanceContract)
	if err != nil {
		log.Fatalf("withdraw balance error: %v", err)
	}

	if withdrawBalanceTransaction == nil || len(withdrawBalanceTransaction.
		GetRawData().GetContract()) == 0 {
		log.Fatalf("withdraw balance error: invalid transaction")
	}

	signTransaction(withdrawBalanceTransaction, ownerKey)

	result, err := cli.BroadcastTransaction(ctx, withdrawBalanceTransaction)

	if err != nil {
		log.Fatalf("withdraw balance error: %v", err)
	}

	return result

}

func signTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey) {
	transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000

	rawData, err := proto.Marshal(transaction.GetRawData())

	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	contractList := transaction.GetRawData().GetContract()

	for range contractList {
		signature, err := crypto.Sign(hash, key)

		if err != nil {
			log.Fatalf("sign transaction error: %v", err)
		}

		transaction.Signature = append(transaction.Signature, signature)
	}
}
