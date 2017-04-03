go-keepassrpc
=============

go-keepassrpc is an implementation of the protocol provided by the KeePassRPC
plugin shipped with [KeeFox](http://keefox.org), a popular Firefox integration
for [KeePass Password Safe](http://keepass.info).

kp
--

`kp` is a command-line client which uses the `keepassrpc` package to talk to
a running KeePass instance.

keepassrpc
----------

`keepassrpc` is the client library which implements the SRP protocol, the
challenge/response protocol, and eventually the JSON RPC protocol. A good
place to start would be to look at the `Client` struct and the `NewClient()`
function.

Additional Links
----------------

From the original author:
* https://github.com/kee-org/KeeFox/wiki/en-%7C-Technical-%7C-KeePassRPC
* https://github.com/kee-org/KeeFox/wiki/en-%7C-Technical-%7C-KeePassRPC-detail

Client implementation (Firefox browser plugin):
* https://github.com/kee-org/browser-addon/tree/master/background
  * `kprpcClient.js` covers the KeePassRPC protocol implementation itself
  * `SRP.js` is their take on SRP

Server implementation (plugin running within KeePass):
* https://github.com/kee-org/keepassrpc
