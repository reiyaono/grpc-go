### Protocol Buffersとは
Googleによって2008年にオープンソース化されたスキーマ言語
スキーマ言語... 要素や属性などの構造を定義するための言語 xmlとか

なぜスキーマ言語が重要か
マイクロサービス化やクライアントの多様化
スキーマ言語をしようしてどういったデータをやりとりするかを宣言的に定義しておく

### Protocol Buffersの特徴
・gRPCのデータフォーマットとして使用されている
・プログラミング言語からは独立しており、様々な言語に変換可能
・バイナリ形式にシリアライズされるので、サイズが小さく高速な通信が可能
・型安全にデータのやり取りが可能
・JSONに変換することも可能
・配列やネストなどJSONと比べると複雑な構造には不向き


### Protocol Buffersを使用した開発の進め方
① スキーマの定義
② 開発言語のオブジェクトを自動生成
③ バイナリ形式へシリアライズ


### messageとは
- 複数のフィールドを持つことができる型定義
  - それぞれのフィールドにはスカラ型もしくはコンポジット型を定義することができる
- コンパイルすると構造体やクラスに変換される

### Tag
- Protocol Bufferではフィールドはフィールド名ではなく、タグ番号によって識別される
- 重複は許されず、一意である必要がある
- タグは連番にする必要がないので、あまり使わないフィールドはあえて16版以降を割り当てて1byteで表せる1~15を確保しておくのもあり

### デフォルト値
- 定義したmessageでデータをやりとりする際に、定義したフィールドがセットされてない場合デフォルト値が設定される
- デフォルト値は型によって決められる
```
string: 空の文字列
bytes: 空のbyte
bool: false
整数型・浮動小数点数: 0
列挙型: タグ番号0の値
repeated: 空のリスト
```

### コンパイル
protoファイルをコンパイルしたファイルにはファイル名にpbがつく

### gRPCの概要
RPC(Remote Procedure Call)
HTTP/2を使用

Microservice間の通信
モバイルユーザーが利用するサービス
速度が求められる場合


### Serviceとは
- RPCの実装単位
- サービス内に定義するメソッドがエンドポイントになる
- 1サービス内に複数のメソッドを定義できる

### 4種類の通信方法
- Unary RPC
  - 1 req 1 res
  - APIなど
- Server Streaming RPC
  - 1 req N res
  - push通知など
- Client Streaming RPC
  - N req 1 res
  - file upload など
- Bidirectional Streaming RPC
  -  N req N res
  - チャットやオンライン対戦ゲームなど


### ダックタイピング
Goにおけるダックタイピングとは、オブジェクトの型ではなく、そのオブジェクトが持つメソッドの実装によって振る舞いを決定するプログラミングの概念です。Goでは、インターフェースを利用してこの概念を実現します。

ダックタイピングの由来
「もしそれがアヒルのように歩き、アヒルのように鳴くなら、それはアヒルである」という考え方に基づいています。つまり、オブジェクトが特定のメソッドを持っていれば、その型が何であれ、そのオブジェクトを特定のインターフェースとして扱うことができます。

```
package main

import "fmt"

// Dog型の定義
type Dog struct {
    Name string
}

// Dog型に関連付けられたメソッド
// (d Dog)はこのメソッドが属する構造体を表す
func (d Dog) Speak() string {
    return "Woof! My name is " + d.Name
}

func main() {
    // Dog型のインスタンスを作成
    myDog := Dog{Name: "Buddy"}

    // Speakメソッドを呼び出し
    fmt.Println(myDog.Speak())
}
```