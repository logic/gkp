gkp
===

gkp is a collection of KeePass-related functionality in Go.

keepassrpc
----------

`keepassrpc` is an implementation of the protocol provided by the KeePassRPC
plugin shipped with [KeeFox](http://keefox.org), a popular Firefox integration
for [KeePass Password Safe](http://keepass.info).

The library implements the three protocols spoken by KeePassRPC: the SRP
protocol which handles initial registration, the key auth challenge/response
protocol which is used for authentication once SRP has been negotiated, and
the encrypted JSON-RPC protocol for post-authentication communication.

We use `jsonenums` to generate marshal/unmarshal helpers for a couple of the
enum values passed to us from the KeePassRPC service. To build anything based
on `keepassrpc`, you'll need to install `jsonenums` first:

    `go get github.com/campoy/jsonenums`

Once you have `jsonenums` in your `PATH`, run `go generate ./...` to create the
needed files.

keepassrpc/cli
--------------

`keepassrpc/cli` provides a number of utilities that make building CLI tools
around `keepassrpc` easier. See `kp` and `git-credential-keepassrpc` for
examples.

kp
--

`kp` is a command-line client which uses the `keepassrpc` package to talk to a
running KeePass instance.

To use this, you'll need [KeePass Password Safe](http://keepass.info/)
installed and running with an open database.

To get a copy of the KeePass plugin, you can either:
* install the [KeeFox Firefox plugin](https://addons.mozilla.org/en-US/firefox/addon/keefox/)
  and follow the installation instructions to install the KeePassRPC plugin
* download the `.xpi` file from addons.mozilla.org, unzip it, and copy
  `deps/KeePassRPC.plgx` to your KeePass plugins directory.

See the KeeFox
[getting started](https://github.com/kee-org/KeeFox/wiki/en-%7C-Getting-started)
instructions for more information.

Once the plugin is installed, the easiest next step is to install `kp`:
* `go get github.com/campoy/jsonenums`
* `go get -tags=gnome_keyring github.com/logic/gkp/kp`

(If you don't need GNOME keyring support, you can skip the `-tags` argument
on the second line.)

Once you've installed `kp`, make sure KeePass is running, and run `kp` with no
arguments to start the first-use authentication step with KeePass. Follow the
on-screen instructions, and you'll be all set.

`kp` will store a session key in your keystore (on OSX, it uses keychain; on
Linux, SecretService or GNOME Keyring), and a configuration file with your
instance username in (probably) `$HOME/.config/gkp/settings.json`.

git-credential-keepassrpc
-------------------------

`git-credential-keepassrpc` provides a `git-credential`-compatible helper for
looking up credentials. Just build it, drop it into your PATH somewhere, and
run:

    `git config --global credential.helper keepassrpc`

Build with `-tags gnome_keyring` for support for storing the auth secret in
GNOME keyring. If you use OSX, or a SecretService-compatible secrets backend,
you don't need to do anything special. So, for example:

    `go get -tags=gnome_keyring github.com/logic/gkp/git-credential-keepassrpc`

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
