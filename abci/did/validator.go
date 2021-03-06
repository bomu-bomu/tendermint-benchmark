package did

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/watcharaphat/tendermint-benchmark/abci/code"
)

const (
	ValidatorSetChangePrefix string = "val:"
)

func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ValidatorSetChangePrefix)
}

func (app *DIDApplication) Validators() (validators []types.Validator) {
	app.logger.Infof("Validators")
	// itr := app.state.db.Iterate(nil, nil)
	// for ; itr.Valid(); itr.Next() {
	// 	if isValidatorTx(itr.Key()) {
	// 		validator := new(types.Validator)
	// 		err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		validators = append(validators, *validator)
	// 	}
	// }

	// viewed := []string{}
	app.state.db.Iterate(func(key []byte, value []byte) bool {
		// viewed = append(viewed, string(key))

		validator := new(types.Validator)
		err := types.ReadMessage(bytes.NewBuffer(key), validator)
		if err != nil {
			panic(err)
		}
		validators = append(validators, *validator)

		return false
	})

	return
}

// format is "val:pubkey"tx
func (app *DIDApplication) execValidatorTx(tx []byte) types.ResponseDeliverTx {
	tx = tx[len(ValidatorSetChangePrefix):]

	// TODO change get PubKey and Power when got ValidatorTx
	// Use "@" as separator since pubKey is base64 and may contain "/"
	pubKeyAndPower := strings.Split(string(tx), "@")
	if len(pubKeyAndPower) < 1 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected 'pubkey'. Got %v", pubKeyAndPower),
		}
	}
	pubkeyS, powerS := pubKeyAndPower[0], "10"
	if len(pubKeyAndPower) > 1 {
		powerS = "0"
	}

	publicKey, err := url.PathUnescape(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) cannot unescape", pubkeyS)}
	}

	pubkey, err := base64.StdEncoding.DecodeString(publicKey)
	var pubKeyEd crypto.PubKeyEd25519
	copy(pubKeyEd[:], pubkey)

	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid hex", publicKey)}
	}
	/*_, err = crypto.PubKeyFromBytes(pubkey)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%X) is invalid go-crypto encoded", pubkey)}
	}*/

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}

	var pubKeyObj types.PubKey
	pubKeyObj.Type = "ed25519"
	pubKeyObj.Data = pubkey

	// update
	return app.updateValidator(types.Validator{pubKeyEd.Address(), pubKeyObj, power})
	// return app.updateValidator(types.Validator{pubKeyEd.Bytes(), power})
}

// add, update, or remove a validator
func (app *DIDApplication) updateValidator(v types.Validator) types.ResponseDeliverTx {
	// key := []byte("val:" + string(v.PubKey))
	key := []byte("val:" + base64.StdEncoding.EncodeToString(v.PubKey.GetData()))
	if v.Power == 0 {
		// remove validator
		if !app.state.db.Has(key) {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf("Cannot remove non-existent validator %X", key)}
		}
		app.state.db.Remove(key)
		// app.state.Size--
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Error encoding validator: %v", err)}
		}
		app.state.db.Set(key, value.Bytes())
		// app.state.Size++
	}

	// we only update the changes array if we successfully updated the tree
	app.ValUpdates = append(app.ValUpdates, v)

	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}
