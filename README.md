# packer-builder-selectel
This plugin is packer builder for Selectel Cloud. It was forked from packer-builder-openstack.

## How to build
1. Clone repo.
```
git clone git@github.com:express42/packer-builder-selectel.git
cd packer-builder-selectel
```
1. Prepare your Go workspace and get dependencies.
```
export GOPATH=$HOME/go
go get
```
1. Workaround for https://github.com/hashicorp/packer/pull/3176
```
rm -rf $GOPATH/src/github.com/hashicorp/packer/vendor
go get
```
1. Run build
```
go build
```
1. Or run build with output into packer plugins directory.
```
go build -o  ~/.packer.d/plugins/packer-builder-selectel
```

## How to install
Place `packer-builder-selectel` executable into `~/.packer.d/plugins/` directory or run build with `-o` flag.
```
go build -o  ~/.packer.d/plugins/packer-builder-selectel
```

## How to use
### Prepare your environment
1. Log into your VPC control panel `https://<%NNNNN%>.selvpc.ru/auth/`, where <%NNNNN%> is numeric ID of your VPC, e.g. 12345, 54321, etc.
1. Select `Access` tab `https://<%NNNNN%>.selvpc.ru/access/`.
1. Click `Download RC-file for user <%username%>`, where <%username%> is your username.
1. Run command and enter your password for selectel's vpc.
```
source ~/Downloads/rc.sh
```
1. Now you can use Selectel API and run `packer` or `openstack` commands.

### Packer builder configuration
1. See documentation for [packer-builder-openstack](https://www.packer.io/docs/builders/openstack.html). But you should use `"type": "selectel"`.
1. See [examples](./examples).

## How to debug
### General
1. Run `packer` with env var `PACKER_LOG=1`.
```
PACKER_LOG=1 packer build examples/ubuntu.json
```
1. You will see debug info inside your console.

### API calls
1. Run [mitmproxy](https://mitmproxy.org/).
```
docker run -t -i -p 8080:8080 mitmproxy/mitmproxy
```
1. Copy mitmproxy CA cert from container to your local computer.
```
docker cp $(docker ps --filter='expose=8080/tcp' -q):/home/mitmproxy/.mitmproxy/mitmproxy-ca-cert.cer ./
```
1. Install `mitmproxy-ca-cert.cer` as trusted CA cerificate into your OS.
1. Run `packer` with env vars `HTTP_PROXY` and `HTTPS_PROXY`.
```
HTTP_PROXY=http://127.0.0.1:8080/ HTTPS_PROXY=$http_proxy packer build examples/ubuntu.json
```
1. You will see requests to api in your container's shell.

## How to contribute
1. Look for issue you want to fix in [bug tracker](https://github.com/express42/packer-builder-selectel/issues).
1. Fork repo.
1. Fix code.
1. Make Pull Request.
1. Mention fixed issues.
