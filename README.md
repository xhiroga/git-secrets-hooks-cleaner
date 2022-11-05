# git-secrets-hooks-cleaner

After uninstalling [`git-secrets`](https://github.com/awslabs/git-secrets), [developer may have to clean hooks generated by git-templates.](https://piruty2.hatenablog.jp/entry/%3Fp%3D521)

## Usage

```shell
./git-secrets-hooks-cleaner list
```


## Note

```shell
go mod init git-secrets-hooks-cleaner
cobra-cli init
go run main.go list

go build
./git-secrets-hooks-cleaner

make install
which git-secrets-hooks-cleaner
git-secrets-hooks-cleaner list

make uninstall
which git-secrets-hooks-cleaner
```

## Reference

