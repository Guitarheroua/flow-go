package avro

const transactionBodySchema = `{
		"type":	"record",
		"name": "TransactionBody",
		"fields": [
			{ "name": "ReferenceBlockID", "type": "bytes"},
			{ "name": "Script", "type": "bytes"},
			{ "name": "Arguments", "type": ["null", {"type":"array", "items":"bytes"}], "default": null },
            { "name": "GasLimit", "type": {"name": "uint64", "type": "fixed", "size": 8 }},
            { "name": "ProposalKey", "type": {
                      "type": "record",
                      "name": "ProposalKey",
                      "fields": [
                      { "name": "Address", "type": "fixed", "size": 8 },
                      { "name": "KeyIndex", "type": {"name": "uint64", "type": "fixed", "size": 8 }},
                      { "name": "SequenceNumber", "type": {"name": "uint64", "type": "fixed", "size": 8 }}
                      ]}
            },
            { "name": "Payer", "type": {"name": "Address", "type": "fixed", "size": 8 }},
            { "name": "Authorizers", "type": ["null", {"type":"array", "items": {"name": "Address", "type": "fixed", "size": 8 }}], "default": null},
            { "name": "PayloadSignatures", "type": {"type":"array", "items": {
                      "type": "record",
                      "name": "TransactionSignature",
                      "fields": [
                      { "name": "Address", "type": "fixed", "size": 8 },
                      { "name": "SignerIndex", "type": "int"},
                      { "name": "KeyIndex", "type": {"name": "uint64", "type": "fixed", "size": 8 }},
                      { "name": "Signature", "type": "bytes"}
                      ]}}
            }
		]
}`
