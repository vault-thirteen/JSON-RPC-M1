# JSON-RPC M1

A remote procedure call framework written in _Go_ programming language. 

Yes, this is a framework, not a library, while it does not allow to automate the process of creation of an RPC infrastructure. 
Unfortunately, _Go_ language does not fully support dynamic typing. It would be possible to automate the process of encoding and decoding of requests and responses, but it would need to use the "Reflection" feature of the language, which works about a hundred times slower than a normal code.

This framework helps in implementing an RPC server and RPC client using the _JSON-RPC M1_ protocol. More information about the _JSON-RPC M1_ protocol can be found in its repository: https://github.com/vault-thirteen/JSON-RPC-M1-Specification

## Usage Example

A very simple usage example is available in the `example\simple` folder of this repository. A description and more information can be found there.

## Features

The framework offers several features useful in practice.

* Settings of this framework are configurable. For example, you can set your own _HTTP_ client using _TLS_, etc.
* The RPC server is able to catch and log exceptions (called "panic" in _Go_ language).
* The framework can count the requests.
* The framework can measure time taken to perform function calls on the server side.
* The framework allows user's function to see an ID of a request.
* The framework allows to set additional meta information in request and response.
* The framework uses a simple and robust protocol, which is focused on data safety and reliability.
* The framework is very simple and does not require external tools. 

As opposed to many other RPC protocols, this framework has some limits, which are the result of its simplicity.

* One-side messages are forbidden. 
  * Every request must be "acknowledged" with a response. 
  * If you need an RPC for game servers, use the _UDP_ protocol and do not cry when someone de-synchs.
* Batch function calls are forbidden for safety reasons.
  * If you need to call for several functions, make several function calls.
* The client makes one request at a time.
  * If you need to send spam to the server, use something else.
* Error codes are not compatible with _Google_'s _JSON RPC_ and _XML RPC_ protocols.
  * We are not _Google_.
* This framework is not going to be as fast as _GRPC_ with _Protocol Buffers_.
  * As with all the protocols using _JSON_ format, textual format is always slower than a binary one.
