<!-- @import "style.less" -->

# Identity Propagation and Assertions

## Introduction

With the growing popularity of protocols based on the OAuth 2.0 specification, there is a need for an interoperable standard that specifies how to convey information about the user from an identity provider (IdP) to a resource server (RS) across security domain boundaries.

## Motivation

Use OAuth 2.0 mechanism in MTA-to-MTA communication in the email system.

## Identity Propagation

In most architectures, the user's security context propagation stops at the IdP/Client security domain boundaries. In an end-to-end identity propagation, the user's security context is extended to the RS, as illustrated in Figure&nbsp;1.

![Model](./images/identity_propagation_model.svg)

<p class="figure">
Fig.&nbsp;1.&emsp;End-to-End Identity Propagation Model
</p>

The user authenticates at the IdP using an authorization code flow. After successful authentication, the RP/Client obtains an access token, which exchanges at the STS for assertion in a JWT format that carries information about the client and user.

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

If the client service domain (typically a third-party client) is not identical to the IdP host domain, one of the following means of proving ownership of the client identifier must be used:

1. [WebFinger](https://github.com/umalabs/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_webfinger.svg)
2. [DANE—(DANCE WG)](https://github.com/umalabs/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_dane.svg)
3. [DNS TXT](https://github.com/umalabs/identity-propagation-and-assertions/blob/main/images/identity_propagation_flow_dns_txt.svg)

## Usability Considerations

The primary benefit of the identity propagation and assertions concept is that it addresses the trust problem between the MTAs. Using OAuth 2.0 and OpenID Connect technologies is an effective option to secure MTA-to-MTA communication. From the OAuth 2.0 point of view to the email system, the outbound MTA is an OAuth 2.0 client, and the inbound MTA is an OAuth 2.0 resource server.

## Conclusion

[Federizer](https://github.com/umalabs/federizer), a message transfer agent (MTA), stands as a proof of concept of the Identity Propagation and Assertions architecture.

## Acknowledgment

[NIST Special Publication 800-63C](https://pages.nist.gov/800-63-3/sp800-63c.html), Digital Identity Guidelines: Federation and Assertions, has proven to be an abstract framework for the Identity Propagation and Assertions architecture.