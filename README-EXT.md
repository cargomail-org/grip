## Identity Propagation Transactions

The DNS-Bound JWT tokens issued by different issuers are chained using digital signatures where issuer and audience claims between tokens must match. This mechanism provides authenticity and integrity protection during identity propagation transactions across multiple security domains.

<div>
    <img src=./images/data_provenance.svg alt="Chain of Identity Propagation Transactions" width="500">
</div>

<p class="figure">
    Fig.&nbsp;2.&emsp;Chain of Identity Propagation Transactions
</p>

## Applications and Usage Patterns

GRIP may be used not only to reimplement existing authentication mechanisms but also to track the origin and history of data.

### Email Authentication

TBD

<div>
    <img src=./images/email_authentication.svg alt="Email Authentication" width="600">
</div>

<p class="figure">
    Fig.&nbsp;3.&emsp;An example of email authentication
</p>

### Data Provenance

TBD
(Digital Asset Transfer, Data Tampering (Deepfake) Detection). Security tokens paired with *digital resources* stored in a *correspondence ledger* detailing the origin and changes guarantee the confidence and validity of data.

### Double-Spend Prevention

TBD
(Virtual/Digital Asset Transfer, Data Tampering (Deepfake) Detection). If someone tries to send the same *digital resources* once again to other recipients, you can prevent this by using a trusted third-party *correspondence ledger* provider.