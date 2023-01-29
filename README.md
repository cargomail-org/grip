<!-- @import "style.less" -->

# Identity Propagation and Assertions

## Introduction

With the growing popularity of protocols based on the OAuth 2.0 specification, there is a need for an interoperable standard that specifies how to convey information about the user from an identity provider (IdP) to a resource server (RS) across security domain boundaries.

## Motivation

Use OAuth 2.0 mechanism in MTA-to-MTA communication in the email system.

## Identity Propagation

In most security concepts and mechanisms, the user's security context propagation stops at the IdP/Client security domain boundaries. In end-to-end identity propagation, the user's security context is extended to the RS across security domain boundaries, as illustrated in Figure&nbsp;1.

![Model](./images/identity_propagation_model.svg)

<p class="figure">
Fig.&nbsp;1.&emsp;End-to-End Identity Propagation Model
</p>

The user authenticates at the IdP using an authorization code flow. After successful authentication, the Relying Party (RP)/Client obtains an access token, which exchanges at the JSON-based Security Token Service (STS) for assertion in a JSON Web Token (JWT) format that carries information about the client and user.

## Impersonation and Delegation

As a specific form of identity propagation, the [OAuth 2.0 Token Exchange RFC](https://datatracker.ietf.org/doc/html/rfc8693) describes impersonation and delegation, where the client obtains a security token that allows it to act as a user in the case of impersonation or, in the case of delegation, allows it to act on behalf of the user. The security token may carry the logical name of the target service for which it is constrained in the aud (audience) claim.

## Assertions

Assertions are statements from an IdP to an RS that contain information about a client and user. The RS uses the information in the assertion to identify the client and user and make authorization decisions about their access to resources controlled by the RS.

## Sequence Diagram

The sequence diagram illustrated in Figure&nbsp;2 shows an identity propagation flow for the user authenticated at the IdP requesting access to resources stored on the RS using a client with a public identifier.

The sequence diagram is self-explanatory; the OIDC authentication flow is omitted for clarity.

<div class="diagram">
    <img src=./images/identity_propagation_flow.svg alt="Sequence Diagram">
</div>

<p class="figure">
Fig.&nbsp;2.&emsp;Identity Propagation Flow
</p>

## Client to Resource Server Authentication

I addition to using [mTLS Certificate-Bound Access Tokens](https://www.rfc-editor.org/rfc/rfc8705#name-mutual-tls-client-certifica), it is recommended to use one of the following means of proving ownership of the client identifier:

1. [DNS TXT](https://github.com/cargomail-org/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_dns_txt.svg)
2. [WebFinger](https://github.com/cargomail-org/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_webfinger.svg)
3. [DANE—(DANCE WG)](https://github.com/cargomail-org/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_dane.svg)

## Resource Server Discovery

The resource server is usually accessed using a service-specific protocol such as email, instant messaging, etc. These protocols need to connect to a specific port in addition to connecting with a specific server.

DNS SRV records defines a symbolic name, the transport protocol, and the port and hostname to connect to for accessing the service. Therefore, DNS SRV records are the recommended way to enable the discovery of service-specific resource servers.

## OpenID Connect Discovery

In order to avoid running an HTTP service that responds to the Webfinger requests as specified in the [OpenID Connect Dynamic Client Registration](https://openid.net/specs/openid-connect-registration-1_0.html) document, it is worth considering whether an [OpenID Connect DNS-based Discovery](https://datatracker.ietf.org/doc/html/draft-sanz-openid-dns-discovery-01) mechanism is more appropriate.

## Usability Considerations

The primary benefit of the identity propagation and assertions in the form of the constrained delegation concept is that it addresses the trust problem between the MTAs. Using OAuth 2.0 and OpenID Connect technologies is an effective option to secure MTA-to-MTA communication. From an OAuth 2.0 perspective, the outbound MTA is an OAuth 2.0 client, and the inbound MTA is an OAuth 2.0 resource server.

## Conclusion

[Cargomail](https://github.com/cargomail-org/cargomail), a revised email system, stands as a proof of concept of the Identity Propagation and Assertions security mechanism.

## Acknowledgment

[NIST Special Publication 800-63C](https://pages.nist.gov/800-63-3/sp800-63c.html), Digital Identity Guidelines: Federation and Assertions, has proven to be an abstract framework for the Identity Propagation and Assertions security mechanism.
