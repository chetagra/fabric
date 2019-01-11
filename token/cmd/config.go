/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package token

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric/bccsp/sw"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/token"
	"github.com/hyperledger/fabric/token/client"
	"github.com/pkg/errors"
)

// LoadConfig converts tha passed string to a ClientConfig.
// The string can be a json string representing the ClientConfig or a path to
// a file containing a ClientConfig in json format
func LoadConfig(s string) (*client.ClientConfig, error) {
	// Try to decode as a Json string
	config := &client.ClientConfig{}
	jsonDecoder := json.NewDecoder(strings.NewReader(s))
	err := jsonDecoder.Decode(config)
	if err == nil {
		return config, nil
	}

	// Try s as a path
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}

	// Try to decode as a Json string
	config = &client.ClientConfig{}
	jsonDecoder = json.NewDecoder(file)
	err = jsonDecoder.Decode(config)
	if err == nil {
		return config, nil
	}

	return nil, errors.New("cannot load configuration, not a json, nor a file containing a json configuration")
}

// GetSigningIdentity retrieves a signing identity from the passed arguments
func GetSigningIdentity(mspConfigPath, mspID, mspType string) (msp.SigningIdentity, error) {
	mspInstance, err := LoadLocalMSPAt(mspConfigPath, mspID, mspType)
	if err != nil {
		return nil, err
	}

	signingIdentity, err := mspInstance.GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}

	return signingIdentity, nil
}

// LoadTokenOwner converts the passed string to a TokenOwner.
// The string can be the path of the msp configuration, or
// the path to a file that contains a serialised identity.
func LoadTokenOwner(s string) (*token.TokenOwner, error) {
	res, err := LoadLocalMspRecipient(s, "MSP_ID")
	if err == nil {
		return res, nil
	}

	raw, err := LoadSerialisedRecipient(s)
	if err == nil {
		return &token.TokenOwner{Raw: raw}, nil
	}

	return &token.TokenOwner{Raw: []byte(s)}, nil
}

// LoadLocalMspRecipient constructs a TokenOwner from the signing identity
// at the passed msp location.
func LoadLocalMspRecipient(mspPath, mspID string) (*token.TokenOwner, error) {
	localMSP, err := LoadLocalMSPAt(mspPath, mspID, "bccsp")
	if err != nil {
		return nil, err
	}

	signer, err := localMSP.GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}

	raw, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	return &token.TokenOwner{Raw: raw}, nil
}

// LoadLocalMSPAt loads an MSP whose configuration is stored at 'dir', and whose
// id and type are the passed as arguments.
func LoadLocalMSPAt(dir, id, mspType string) (msp.MSP, error) {
	if mspType != "bccsp" {
		return nil, errors.Errorf("invalid msp type, expected 'bccsp', got %s", mspType)
	}
	conf, err := msp.GetLocalMspConfig(dir, nil, id)
	if err != nil {
		return nil, err
	}
	ks, err := sw.NewFileBasedKeyStore(nil, filepath.Join(dir, "keystore"), true)
	if err != nil {
		return nil, err
	}
	thisMSP, err := msp.NewBccspMspWithKeyStore(msp.MSPv1_0, ks)
	if err != nil {
		return nil, err
	}
	err = thisMSP.Setup(conf)
	if err != nil {
		return nil, err
	}
	return thisMSP, nil
}

// LoadSerialisedRecipient loads a serialised identity from file
func LoadSerialisedRecipient(serialisedRecipientPath string) ([]byte, error) {
	fileCont, err := ioutil.ReadFile(serialisedRecipientPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", serialisedRecipientPath)
	}

	return fileCont, nil
}

// LoadTokenIDs converts the passed string to an array of TokenIDs.
// The string can be a json representing the TokenIDs, or a path
// to a fail containing the json representation.
func LoadTokenIDs(s string) ([]*token.TokenId, error) {
	// s can be a
	// - json string representing the token IDs
	// - a path containing a json string representing the token IDs
	res, err := LoadTokenIDsFromJson(s)
	if err == nil {
		return res, nil
	}

	return LoadTokenIDsFromFile(s)
}

// LoadTokenIDsFromJson interprets the passed string as a json representation of TokenIDs
func LoadTokenIDsFromJson(s string) ([]*token.TokenId, error) {
	var tokenIDs []*token.TokenId
	err := json.Unmarshal([]byte(s), &tokenIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling json")
	}

	return tokenIDs, nil
}

// LoadTokenIDsFromFile loads TokenIDs from the passed file supposed to contain a json
// representation of the TokenIDs.
func LoadTokenIDsFromFile(s string) ([]*token.TokenId, error) {
	fileCont, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", s)
	}

	return LoadTokenIDsFromJson(string(fileCont))
}

// LoadShares converts the passed string to an array of Shares.
// The string is either a json representing the shares, or a path to a file
// containing the json representation.
func LoadShares(s string) ([]*token.RecipientTransferShare, error) {
	// s can be a
	// - json string representing the shares
	// - a path containing a json string representing the shares
	res, err := LoadSharesFromJson(s)
	if err == nil {
		return SubstituteShareRecipient(res)
	}
	res, err = LoadSharesFromFile(s)
	if err == nil {
		return SubstituteShareRecipient(res)
	}
	return nil, errors.New("failed loading shares")
}

// LoadSharesFromJson converts the passed json string to shares
func LoadSharesFromJson(s string) ([]*token.RecipientTransferShare, error) {
	var shares []*token.RecipientTransferShare
	err := json.Unmarshal([]byte(s), &shares)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling json")
	}

	return shares, nil
}

// LoadSharesFromFile loads from file shares in json representation.
func LoadSharesFromFile(s string) ([]*token.RecipientTransferShare, error) {
	fileCont, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", s)
	}

	return LoadSharesFromJson(string(fileCont))
}

// SubstituteShareRecipient scans the recipients to see if they need additional
// post-processing. For example, a recipient can contain the path of a file containing
// the serialised identity to be loaded.
func SubstituteShareRecipient(shares []*token.RecipientTransferShare) ([]*token.RecipientTransferShare, error) {
	for _, share := range shares {
		if share == nil {
			continue
		}
		if share.Recipient == nil {
			return nil, errors.New("invalid recipient share")
		}
		recipient, err := LoadTokenOwner(string(share.Recipient.Raw))
		if err != nil {
			return nil, err
		}
		share.Recipient = recipient
	}

	return shares, nil
}

// JsonLoader implements the Loader interface
type JsonLoader struct {
}

func (*JsonLoader) TokenOwner(s string) (*token.TokenOwner, error) {
	return LoadTokenOwner(s)
}

func (*JsonLoader) TokenIDs(s string) ([]*token.TokenId, error) {
	return LoadTokenIDs(s)
}

func (*JsonLoader) Shares(s string) ([]*token.RecipientTransferShare, error) {
	return LoadShares(s)
}
