<!-- @import "style.less" -->

# Global Reference Identity Protocol

## Abstract

Global Reference Identity Protocol (GRIP) is a token-based security protocol that authenticates service requests between untrusted hosts across the Internet.

## Introduction

With the growing popularity of communication systems, there is a need for an interoperable standard that specifies how to issue and convey information about the user from one service to another across security domain boundaries.

## Goals and Objectives

Enhance email authentication with a token-based protocol. The resulting concept should support the data provenance mechanism.

## GRIP Acronyms

GRIP uses special jargon. For the sake of brevity of this document, the following list of acronyms will be used:
<pre>
DKIM    DomainKeys Identified Mail
DNS     Domain Name System
CA      Certificate Authority
CN      Common Name
TLS     Transport Layer Security
mTLS    Mutual Transport Layer Security

HTTP(S) Hypertext Transfer Protocol (Secure)

JWT     JSON Web Token
JWK     JSON Web Key
</pre>

## Impersonation and Delegation

The proposed mechanism allows clients with the appropriate security token to impersonate users by delegated signing authority. The client obtains a security token that allows it to act as a specific user. The security token may carry the logical name of the target service for which it is constrained.

## Assertions

Assertions are statements from a token producer to a token consumer that contain information about the principal. In the Identity Propagation scenario, the target server uses the information in the assertion to identify the client and user to make authorization decisions about their access to the service running on that server.

## DNS-Bound Tokens

In some service-to-service communication scenarios, three identities are employed: user, client, and server identities. Fundamentally, mutual TLS (mTLS) and TLS certificates resolve client and server identities, while tokens resolve client and user identities. A DNS-bound token is a self-issued assertion in a JWT format signed by an mTLS private key that the first service uses to authenticate to the second service. The mTLS public key hash is published in the first service domain using the DNS TXT record, where the CN attribute of the mTLS public key certificate is used as a global client identifier with respect to the service it represents.

## Service Discovery

The client typically connects to the server using service-specific protocols like SMTP or HTTP(S). These protocols require a connection to a specific port in addition to connecting to a particular server. A DNS SRV record defines a symbolic name, the transport protocol, port, and hostname to connect to when accessing the service. Therefore, DNS SRV records are the recommended way to discover service-specific servers.

## Identity Propagation

In most security concepts and mechanisms, the user's security context propagation stops at the user's security domain boundaries. In end-to-end identity propagation, the user's security context is extended outside the user's security perimeter.

Using self-signed certificates ensures you can quickly start with the most straightforward identity propagation mechanism. The sequence diagram illustrated in Figure 1 shows the identity propagation flow without end-user involvement, where the client requests access to the server on behalf of the impersonated user using a self-issued security token.

The sequence diagram is self-explanatory.


<div class="diagram">
    <img src=./images/self-issued_identity_propagation_flow.svg alt="Identity Propagation Flow">
</div>

<p class="figure">
    Fig.&nbsp;1.&emsp;Identity Propagation Flow
</p>

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

## Implementation

[Cargomail](https://github.com/cargomail-org/cargomail) — an implementation of the [revised Internet Mail architecture](https://github.com/cargomail-org/cargomail/raw/main/whitepaper/Cargomail.pdf), stands as proof of concept for the GRIP mechanism. Cargomail uses GRIP for email authentication instead of the DKIM authentication method. GRIP employs identity propagation and assertion apparatus to convey identity information about the end user across different administrative authorities of the email services. Furthermore, instead of using MX records to discover the communication services, Cargomail relies on DNS SRV records.

## Conclusion

Using mTLS and DNS technology is an effective option to secure service-to-service communication between unrelated security domains in an untrusted environment. Generally speaking, GRIP allows actions to be taken on behalf of users across the Internet. Being application-protocol agnostic, GRIP can be applied to any communication service protected by TLS, including email.