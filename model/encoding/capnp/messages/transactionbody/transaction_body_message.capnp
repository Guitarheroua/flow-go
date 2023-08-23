using Go = import "/go.capnp";
@0xb5a009291978d69a;
$Go.package("capnp");
$Go.import("github.com/onflow/flow-go/model/encoding/capnp/messages/transactionbody");

using Address = Data;
using Identifier = Data;

struct ProposalKeyMessage {
    address @0 :Address;
    keyIndex @1 :UInt64;
    sequenceNumber @2 :UInt64;
}

struct TransactionSignatureMessage {
    address @0 :Address;
    signerIndex @1 :Int32;
    keyIndex @2 :UInt64;
    signature @3 :Data;
}

struct TransactionBodyMessage {
    referenceBlockID @0 :Identifier;
    script @1 :Data;
    arguments @2 :List(Data);
	gasLimit @3 :UInt64;
    proposalKey @4 :ProposalKeyMessage;
    payer @5 :Address;
    authorizers @6 :List(Address);
	payloadSignatures @7 :List(TransactionSignatureMessage);
	envelopeSignatures @8 :List(TransactionSignatureMessage);
}
