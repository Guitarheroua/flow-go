using Go = import "/go.capnp";
using Transaction = import "/messages/transactionbody/transaction_body_message.capnp";
@0x9e050064d5c9cb20;
$Go.package("capnp");
$Go.import("github.com/onflow/flow-go/model/encoding/capnp/messages/clusterblockproposal");

using Collection = List(Transaction.TransactionBodyMessage);
using Signature = Data;

struct QuorumCertificateMessage {
    view @0 :UInt64;
    blockID @1 :Data;
    signerIndices @2 :Data;
    sigData @3 :Data;
}

struct TimeoutCertificateMessage {
    view @0 :UInt64;
    newestQCViews @1 :List(UInt64);
    newestQC @2 :QuorumCertificateMessage;
    signerIndices @3 :Data;
    sigData @4 :Signature;
}

struct HeaderMessage {
    chainID @0 :Text;
    parentID @1 :Data;
    height @2 :UInt64;
    payloadHash @3 :Data;
    timestamp @4 :Int64;
    view @5 :UInt64;
    parentView @6 :UInt64;
    parentVoterIndices @7 :Data;
    parentVoterSigData @8 :Data;
    proposerID @9 :Data;
    proposerSigData @10 :Data;
    lastViewTC @11 :TimeoutCertificateMessage;
}

struct UntrustedClusterBlockPayloadMessage {
    collection @0 :Collection;
    referenceBlockID @1 :Data;
}

struct ClusterBlockProposalMessage {
    header @0 :HeaderMessage;
    payload @1 :UntrustedClusterBlockPayloadMessage;
}
