##Prerequisites:
1. Go 1.6.2

##Build from source instructions:

```
1. $ cd terraform 
   Terraform working dir, e.g. proudction use 'cd /opt/terraform', or 
   dev use'cd $HOME/projects/terraform' for developemnt. If it doesn't exist then create it.

2. $ mkdir -p src/github.com/hashicorp bin pkg
3. $ cd src/github.com/hashicorp
4. $ https://github.com/seizadi/terraform.git  # clone forked repo instead of original
5. $ cd terraform
6. $ export GOPATH=`pwd`
7. $ export PATH=$PATH:$GOPATH/bin
8. $ go run build.go setup # run only once for godep
```

Follow the remaining instruction from BUILDING.md
```
9. $ $godep restore    # Will pull down dependencies in your current GOPATH
.....
```

## Notes

Got Error, now I add $GOPATH/bin to path step 7 above
also see https://github.com/hashicorp/terraform/issues/1688
```
sc-l-seizadi:terraform seizadi$ make dev
==> Checking that code complies with gofmt requirements...
go generate $(go list ./... | grep -v /vendor/)
2016/05/29 09:13:40 Generated command/internal_plugin_list.go
command/hook_count_action.go:3: running "stringer": exec: "stringer": executable file not found in $PATH
config/resource_mode.go:3: running "stringer": exec: "stringer": executable file not found in $PATH
helper/schema/resource_data_get_source.go:3: running "stringer": exec: "stringer": executable file not found in $PATH
terraform/graph_config_node_type.go:3: running "stringer": exec: "stringer": executable file not found in $PATH
make: *** [generate] Error 1
```
