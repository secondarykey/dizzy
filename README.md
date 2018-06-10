# dizzy

良く簡易的にDatastoreを使うことが多くて、簡単にメンテできるようにということで作っています。
まだstring,int,float,boolを登録できるようになっただけで、
未完成ですが、どんどん型を登録できるようにしていきます。

構造体に[ds.Meta](https://github.com/knightso/base/blob/master/gae/ds/datastore.go#L40)を埋め込んで
コメント「+dizzy」を埋め込んで実現します。

```sh
$ go get github.com/secondarykey/dizzy
$ go install github.com/secondarykey/dizzy
$ cd $GOPATH/src/github.com/secondarykey/dizzy/examples
$ dizzy gen .
$ dev_appserver.py dizzy_app.yaml
```

[管理画面](http://localhost:8080/_dizzy/) にアクセス(admin)できるはずです


予定としては

* 型を増やす
* デモサイトを作成
* もろもろ

Continue on [wiki](https://github.com/secondarykey/dizzy/wiki).
