package avro

import (
	"encoding/binary"
	"fmt"
	"github.com/linkedin/goavro/v2"
	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/messages"
	"reflect"
	"time"
)

const TransactionBodySchema = `{
    "type": "record",
    "name": "TransactionBody",
    "fields": [
        {
            "name": "ReferenceBlockID",
            "type": {
                "namespace": "flow",
                "name": "Identifier",
                "type": "fixed",
                "size": 32
            }
        },
        {
            "name": "Script",
            "type": "bytes"
        },
        {
            "name": "Arguments",
            "type": {
                    "type": "array",
                    "items": "bytes"
                }
        },
        {
            "name": "GasLimit",
            "type": {
                "name": "uint64",
                "type": "fixed",
                "size": 8
            }
        },
        {
            "name": "ProposalKey",
            "type": {
                "type": "record",
                "name": "ProposalKey",
                "fields": [
                    {
                        "name": "Address",
                        "type": "fixed",
                        "size": 8
                    },
                    {
                        "name": "KeyIndex",
                        "type": {
                            "name": "uint64",
                            "type": "fixed",
                            "size": 8
                        }
                    },
                    {
                        "name": "SequenceNumber",
                        "type": {
                            "name": "uint64",
                            "type": "fixed",
                            "size": 8
                        }
                    }
                ]
            }
        },
        {
            "name": "Payer",
            "type": {
                "name": "Address",
                "type": "fixed",
                "size": 8
            }
        },
        {
            "name": "Authorizers",
            "type": {
                "type": "array",
                "items": {
                    "name": "Address",
                    "type": "fixed",
                    "size": 8
                }
            }
        },
        {
            "name": "PayloadSignatures",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "TransactionSignature",
                    "fields": [
                        {
                            "name": "Address",
                            "type": "fixed",
                            "size": 8
                        },
                        {
                            "name": "SignerIndex",
                            "type": "int"
                        },
                        {
                            "name": "KeyIndex",
                            "type": {
                                "name": "uint64",
                                "type": "fixed",
                                "size": 8
                            }
                        },
                        {
                            "name": "Signature",
                            "type": "bytes"
                        }
                    ]
                }
            }
        }
    ]
}`

const ClusterBlockProposal = `{
    "type": "record",
    "name": "ClusterBlockProposal",
    "fields": [
        {
            "type": "record",
            "name": "Block",
            "fields": [
                {
                    "name": "Header",
                    "type": {
                        "type": "record",
                        "name": "Header",
                        "fields": [
                            {
                                "name": "ChainID",
                                "type": "string"
                            },
                            {
                                "name": "ParentID",
                                "type": {
                                    "namespace": "flow",
                                    "name": "Identifier",
                                    "type": "fixed",
                                    "size": 32
                                }
                            },
                            {
                                "name": "Height",
                                "type": {
                                    "name": "uint64",
                                    "type": "fixed",
                                    "size": 8
                                }
                            },
                            {
                                "name": "PayloadHash",
                                "type": "bytes"
                            },
                            {
                                "name": "Timestamp",
                                "type": "string"
                            },
                            {
                                "name": "View",
                                "type": {
                                    "name": "uint64",
                                    "type": "fixed",
                                    "size": 8
                                }
                            },
                            {
                                "name": "ParentView",
                                "type": {
                                    "name": "uint64",
                                    "type": "fixed",
                                    "size": 8
                                }
                            },
                            {
                                "name": "ParentVoterIndices",
                                "type": "bytes"
                            },
                            {
                                "name": "ParentVoterSigData",
                                "type": "bytes"
                            },
                            {
                                "name": "ProposerID",
                                "type": {
                                    "namespace": "flow",
                                    "name": "Identifier",
                                    "type": "fixed",
                                    "size": 32
                                }
                            },
                            {
                                "name": "ProposerSigData",
                                "type": "bytes"
                            },
                            {
                                "name": "LastViewTC",
                                "type": [ "null", {
                                    "type": "record",
                                    "name": "TimeoutCertificate",
                                    "fields": [
                                        {
                                            "name": "View",
                                            "type": {
                                                "name": "uint64",
                                                "type": "fixed",
                                                "size": 8
                                            }
                                        },
                                        {
                                            "name": "NewestQCViews",
                                            "type": {
                                                "type": "array",
                                                "items": {
                                                    "name": "uint64",
                                                    "type": "fixed",
                                                    "size": 8
                                                }
                                            }
                                        },
                                        {
                                            "name": "NewestQC",
                                            "type": ["null", {
                                                "type": "record",
                                                "name": "QuorumCertificate",
                                                "fields": [
                                                    {
                                                        "name": "View",
                                                        "type": {
                                                            "name": "uint64",
                                                            "type": "fixed",
                                                            "size": 8
                                                        }
                                                    },
                                                    {
                                                        "name": "BlockID",
                                                        "type": {
                                                            "namespace": "flow",
                                                            "name": "Identifier",
                                                            "type": "fixed",
                                                            "size": 32
                                                        }
                                                    },
                                                    {
                                                        "name": "SignerIndices",
                                                        "type": "bytes"
                                                    },
                                                    {
                                                        "name": "SigData",
                                                        "type": "bytes"
                                                    }
                                                ]
                                            }],
                                            "default": null
                                        },
                                        {
                                            "name": "SignerIndices",
                                            "type": "bytes"
                                        },
                                        {
                                            "name": "SigData",
                                            "type": "bytes"
                                        }
                                    ]
                                }],
                                "default": null
                            }
                        ]
                    }
                },
                {
                    "name": "Payload",
                    "type": {
                        "type": "record",
                        "name": "UntrustedClusterBlockPayload",
                        "fields": [
                            {
                                "name": "Collection",
                                "type": {
                                    "type": "array",
                                    "items": {
                                        "type": "record",
                                        "name": "TransactionBody",
                                        "fields": [
                                            {
                                                "name": "ReferenceBlockID",
                                                "type": {
                                                    "namespace": "flow",
                                                    "name": "Identifier",
                                                    "type": "fixed",
                                                    "size": 32
                                                }
                                            },
                                            {
                                                "name": "Script",
                                                "type": "bytes"
                                            },
                                            {
                                                "name": "Arguments",
                                                "type": {
                                                    "type": "array",
                                                    "items": "bytes"
                                                }
                                            },
                                            {
                                                "name": "GasLimit",
                                                "type": {
                                                    "name": "uint64",
                                                    "type": "fixed",
                                                    "size": 8
                                                }
                                            },
                                            {
                                                "name": "ProposalKey",
                                                "type": {
                                                    "type": "record",
                                                    "name": "ProposalKey",
                                                    "fields": [
                                                        {
                                                            "name": "Address",
                                                            "type": "fixed",
                                                            "size": 8
                                                        },
                                                        {
                                                            "name": "KeyIndex",
                                                            "type": {
                                                                "name": "uint64",
                                                                "type": "fixed",
                                                                "size": 8
                                                            }
                                                        },
                                                        {
                                                            "name": "SequenceNumber",
                                                            "type": {
                                                                "name": "uint64",
                                                                "type": "fixed",
                                                                "size": 8
                                                            }
                                                        }
                                                    ]
                                                }
                                            },
                                            {
                                                "name": "Payer",
                                                "type": {
                                                    "name": "Address",
                                                    "type": "fixed",
                                                    "size": 8
                                                }
                                            },
                                            {
                                                "name": "Authorizers",
                                                "type": {
                                                    "type": "array",
                                                    "items": {
                                                        "name": "Address",
                                                        "type": "fixed",
                                                        "size": 8
                                                    }
                                                }
                                            },
                                            {
                                                "name": "PayloadSignatures",
                                                "type": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "record",
                                                        "name": "TransactionSignature",
                                                        "fields": [
                                                            {
                                                                "name": "Address",
                                                                "type": "fixed",
                                                                "size": 8
                                                            },
                                                            {
                                                                "name": "SignerIndex",
                                                                "type": "int"
                                                            },
                                                            {
                                                                "name": "KeyIndex",
                                                                "type": {
                                                                    "name": "uint64",
                                                                    "type": "fixed",
                                                                    "size": 8
                                                                }
                                                            },
                                                            {
                                                                "name": "Signature",
                                                                "type": "bytes"
                                                            }
                                                        ]
                                                    }
                                                }
                                            }
                                        ]
                                    }
                                }
                            },
                            {
                                "name": "ReferenceBlockID",
                                "type": {
                                    "namespace": "flow",
                                    "name": "Identifier",
                                    "type": "fixed",
                                    "size": 32
                                }
                            }
                        ]
                    }
                }
            ]
        }
    ]
}`

func ConvertStructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldValue := inputValue.Field(i).Interface()

		switch fieldValue.(type) {
		case uint64:
			result[field.Name] = Uint64ToBytes(fieldValue.(uint64))
		case []uint64:
			views := fieldValue.([]uint64)
			byteSlices := make([][]byte, len(views))
			for i, view := range views {
				byteSlices[i] = Uint64ToBytes(view)
			}
			result[field.Name] = byteSlices
		case flow.Address:
			result[field.Name] = fieldValue.(flow.Address).Bytes()
		case []flow.Address:
			addresses := fieldValue.([]flow.Address)
			byteSlices := make([][]byte, len(addresses))
			for i, addr := range addresses {
				byteSlices[i] = addr.Bytes()
			}
			result[field.Name] = byteSlices
		case crypto.Signature:
			result[field.Name] = fieldValue.(crypto.Signature).Bytes()
		case flow.ProposalKey:
			result[field.Name] = ConvertStructToMap(fieldValue)
		case messages.UntrustedClusterBlock:
			result[field.Name] = ConvertStructToMap(fieldValue)
		case flow.Header:
			result[field.Name] = ConvertStructToMap(fieldValue)
		case messages.UntrustedClusterBlockPayload:
			result[field.Name] = ConvertStructToMap(fieldValue)
		case flow.TransactionBody:
			result[field.Name] = ConvertStructToMap(fieldValue)
		case []flow.TransactionBody:
			transactionBodies := fieldValue.([]flow.TransactionBody)
			tbSlices := make([]map[string]interface{}, len(transactionBodies))
			for i, tb := range transactionBodies {
				tbSlices[i] = ConvertStructToMap(tb)
			}
			result[field.Name] = tbSlices
		case *flow.TimeoutCertificate:
			if timeoutCert, ok := fieldValue.(*flow.TimeoutCertificate); ok {
				if timeoutCert == nil {
					result[field.Name] = goavro.Union("null", nil)
				} else {
					result[field.Name] = goavro.Union("TimeoutCertificate", ConvertStructToMap(*timeoutCert))
				}
			}
		case *flow.QuorumCertificate:
			if quorumCert, ok := fieldValue.(*flow.QuorumCertificate); ok {
				if quorumCert == nil {
					result[field.Name] = goavro.Union("null", nil)
				} else {
					result[field.Name] = goavro.Union("QuorumCertificate", ConvertStructToMap(*quorumCert))
				}
			}
		case flow.Identifier:
			if id, ok := fieldValue.(flow.Identifier); ok {
				result[field.Name] = id[:]
			}
		case flow.ChainID:
			if chainID, ok := fieldValue.(flow.ChainID); ok {
				result[field.Name] = string(chainID)
			}
		case time.Time:
			result[field.Name] = fieldValue.(time.Time).String()
		default:
			result[field.Name] = fieldValue
		}
	}

	return result
}

func Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	return buf
}

func ConvertMapToStruct(data map[string]interface{}, output interface{}) error {
	outputValue := reflect.ValueOf(output)

	if outputValue.Kind() != reflect.Ptr || outputValue.IsNil() {
		return fmt.Errorf("output must be a non-nil pointer")
	}

	// Dereference the pointer to get the underlying struct value
	outputValue = outputValue.Elem()
	outputType := outputValue.Type()

	for i := 0; i < outputType.NumField(); i++ {
		field := outputType.Field(i)
		fieldValue, found := data[field.Name]
		if !found {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Uint64:
			outputValue.Field(i).SetUint(BytesToUint64(fieldValue.([]byte)))
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Uint64 {
				byteSlices := fieldValue.([][]byte)
				views := make([]uint64, len(byteSlices))
				for i, bytes := range byteSlices {
					views[i] = BytesToUint64(bytes)
				}
				outputValue.Field(i).Set(reflect.ValueOf(views))
			} else if field.Type.Elem() == reflect.TypeOf(flow.Address{}) {
				byteSlices := fieldValue.([][]byte)
				addresses := make([]flow.Address, len(byteSlices))
				for i, bytes := range byteSlices {
					addresses[i] = flow.BytesToAddress(bytes)
				}
				outputValue.Field(i).Set(reflect.ValueOf(addresses))
			} else if field.Type.Elem() == reflect.TypeOf(messages.UntrustedClusterBlock{}) {
				mapSlices := fieldValue.([]map[string]interface{})
				clusterBlocks := make([]messages.UntrustedClusterBlock, len(mapSlices))
				for i, data := range mapSlices {
					clusterBlock := messages.UntrustedClusterBlock{}
					ConvertMapToStruct(data, &clusterBlock)
					clusterBlocks[i] = clusterBlock
				}
				outputValue.Field(i).Set(reflect.ValueOf(clusterBlocks))
			} else {
				return fmt.Errorf("unsupported slice type for field %s", field.Name)
			}
		case reflect.Ptr:
			ptrType := field.Type.Elem()
			ptrValue := reflect.New(ptrType).Elem()
			if err := ConvertMapToStruct(fieldValue.(map[string]interface{}), ptrValue.Addr().Interface()); err != nil {
				return err
			}
			outputValue.Field(i).Set(ptrValue.Addr())
		case reflect.Struct:
			if err := ConvertMapToStruct(fieldValue.(map[string]interface{}), outputValue.Field(i).Addr().Interface()); err != nil {
				return err
			}
		case reflect.String:
			outputValue.Field(i).SetString(fieldValue.(string))
		default:
			return fmt.Errorf("unsupported field type for field %s", field.Name)
		}
	}

	return nil
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}
