# Key learnings
* While building this project, I also read the book [The Lean Startup](http://theleanstartup.com/)
* Blockchain is a database.
* Every blockchain starts with a Genesis file which is used to distribute the first tokens to early members of the chain. It is never updated afterwards.
* A whitepaper servers to outline the specifications of how the particular blockchain will look and behave.
* It uses a event based architecture where the dB is the final aggregated stare after replaying all the transactions in a specific queue.
* Implemented CLI using [Cobra](https://github.com/spf13/cobra)
* Blockchain dB is immutable. The system is transparent, auditable and well defined.
* It is hashed by a secure crypto hash function. A specific hash represets a particular dB state.
* A batch of transactions make a block. Each block is encoded and hashed. Block has header(parent block metadata) and payload (new dB transactions)
* Byzantine Fault Tolerance, Proof of Work.
* As blockchain is public, we dont save any sensitive data in it. A transaction is represented in terms of inputs and outputs. 
* UTXO Model in Transactions. (there are no accounts or wallets at the protocol layer)
* Unspent Transaction Outputs or UTXOs serve as globally-accessible evidence that you have Bitcoin in your digital wallet.

# Future Scope

* Implement Proof of Concept for consensus
* Create a distributed P2P network to understand the blockchain better
* Implement Clean Architecture in the code