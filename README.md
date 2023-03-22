<!-- @import "style.less" -->

# Global Reference Identity Protocol

## Abstract

Global Reference Identity Protocol (GRIP) is a token-based security protocol that authenticates service requests between untrusted hosts across the Internet.

## Introduction

With the growing popularity of communication systems, there is a need for an interoperable standard that specifies how to convey information about the user from one service to another across security domain boundaries.

## GRIP Acronyms

GRIP uses a special jargon. For the sake of brevity of this document, the following list of acronyms will be used:

* SMTP Simple Mail Transfer Protocol
* FTPS File Transfer Protocol Secure
* EDI Electronic Data Interchange
* AS2 a protocol for transmission EDI
* gRPC Google Remote Procedure Call
* WSS WebSocket Secure
</br>
* TLS Transport Layer Security
* mTLS Mutual Transport Layer Security
</br>
* IdP Identity Provider
* OIDC OpenID Connect
* RP Relying Party
* AS Authorization Server
* RS Resource Server
* STS Security Token Service
* JWT JSON Web Token
* JWK JSON Web Key
</br>
* DNS Domain Name System
</br>
* CA Certificate Authority
* CN Common Name
</br>
* WG Working Group
## Motivation

To enhance SMTP, FTPS, and AS2 protocols with a cross-domain authentication mechanism. Also, consider using this mechanism for other communication technologies, e.g., RESTful, gRPC, WSS, and WebTransport. The authentication mechanism should be application-protocol agnostic.

## Identity Propagation

In most security concepts and mechanisms, the user's security context propagation stops at the user's security domain boundaries. In end-to-end identity propagation, the user's security context is extended outside the user's security perimeter.

## Impersonation and Delegation

The OAuth 2.0 intrinsic delegation mechanism allows clients with the appropriate security token to impersonate the user or being delegated by that user. As a specific form of identity propagation, the [OAuth 2.0 Token Exchange](https://datatracker.ietf.org/doc/html/rfc8693) specification describes impersonation and delegation, where the Client obtains a security token that allows it to act as a user in the case of impersonation or, in the case of delegation, allows it to act on behalf of the user. The output security token may carry the logical name of the target service for which it is constrained.

## Assertions

Assertions are statements from a token producer to a token consumer that contain information about the principal. In the Identity Propagation scenario, the resource server uses the information in the assertion to identify the Client and user to make authorization decisions about their access to resources controlled by that resource server.

## Identities and Certificate-Bound Tokens

In most client-service-to-server-service communication scenarios, three identities are employed: user-identity, client-identity, and server-identity. mTLS certificates resolve client-identity and server-identity, while tokens resolve user-identity. mTLS during protected resource access also serves as a proof-of-possession of the token mechanism, as stated in [section 4](https://www.rfc-editor.org/rfc/rfc8705#section-4) of the [RFC8705](https://www.rfc-editor.org/rfc/rfc8705) OAuth 2.0 mTLS Client Authentication and Certificate-Bound Access Tokens specification.

## Self-Issued Identity Propagation

The sequence diagram illustrated in Figure&nbsp;1 shows the self-issued identity propagation flow without AS and end-user involvement, where the Client requests access to resources stored on the RS on behalf of the impersonated user using a self-issued token.

The sequence diagram is self-explanatory.

<div class="diagram">
    <img src=./images/self-issued_identity_propagation_flow.svg alt="Sequence Diagram">
</div>

<p class="figure">
Fig.&nbsp;1.&emsp;Self-Issued Identity Propagation flow
</p>

## 2-Legged Identity Propagation

The sequence diagram illustrated in Figure&nbsp;2 shows the 2-legged identity propagation flow without end-user involvement, where the Client requests access to resources stored on the RS on behalf of the impersonated user using a token generated on the AS.

The sequence diagram is self-explanatory.

<div class="diagram">
    <img src=./images/2-legged_identity_propagation_flow.svg alt="Sequence Diagram">
</div>

<p class="figure">
Fig.&nbsp;2.&emsp;2-Legged Identity Propagation flow
</p>

## 3-Legged Identity Propagation

The sequence diagram illustrated in Figure&nbsp;3 shows the 3-legged identity propagation flow for the user authenticated at the IdP, where the Client requests access to resources stored on the RS on behalf of the authenticated user using a token generated on the AS.

The sequence diagram is self-explanatory; the OIDC authentication flow is omitted for clarity.

<div class="diagram">
    <img src=./images/3-legged_identity_propagation_flow.svg alt="Sequence Diagram">
</div>

<p class="figure">
Fig.&nbsp;3.&emsp;3-Legged Identity Propagation flow
</p>

## Client to Resource Server Authentication

In addition to using the [mTLS Certificate-Bound Access Tokens](https://www.rfc-editor.org/rfc/rfc8705#section-4) mechanism, it is recommended to use one of the following means of proving ownership of the client identifier:

1. DNS TXT
2. WebFinger
3. DANCE WG

## Resource Server Discovery

The resource server is usually accessed using a service-specific protocol such as email or instant messaging. These protocols need to connect to a specific port in addition to connecting with a specific server.

DNS SRV record defines a symbolic name, the transport protocol, and the port and hostname to connect to for accessing the service. Therefore, DNS SRV records are the recommended way to enable the discovery of service-specific resource servers.

## Usability Considerations

The primary benefit of Identity Propagation and Assertions in the form of the constrained delegation concept is that it addresses the zero-trust between unrelated security domains. Using an OAuth 2.0 technology is an effective option to secure service-to-service communication. From an OAuth 2.0 perspective, the outbound service is an OAuth 2.0 client, and the inbound service is an OAuth 2.0 resource server.

## Implementation

[Cargomail](https://github.com/cargomail-org/cargomail) stands as proof of the concept of the GRIP mechanism.

## Conclusion

Given that GRIP is application-protocol agnostic, it can be applied to any TLS-protected communication protocol, including SMTP and FTPS. Generally speaking, GRIP allows identity-to-identity communication in a secure manner across the Internet.