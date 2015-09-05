# 「※希望はありません」の誕生を観察してみよう

がっこうぐらし！のOPのフレーズ「ここには夢がちゃんとある」の時に起こる弾幕、「※希望はありません」がどういうふうに生まれてきたのか観察してみました。

![](https://raw.githubusercontent.com/ledyba/gakko-analyzer/master/screenshot.png)

ニコニコ動画の弾幕はあくまで歌詞の内容を元にした「空耳」だったりすることが多く、この「希望はありません」の弾幕は面白い存在だと思います。

# 方法

 - 配信開始から一週間分の過去ログを集める
 - 「夢」「希望」（わかば＊ガールば「鳴」）のどちらかが含まれるコメントを抽出
 - 正規化されたレーベンシュタイン距離を使って適当にクラスタリング
   - 似ているコメント同士でも1000コメント以上離れたコメントは別グループに分け、グループ同士をエッジで接続しました。
   - 書かれているコメントを真似して書くコメントと、独立に思いついて偶然にたコメントを区別するため。
 - 時系列にまとめて図に書く
   - デカいグループ（結果として「※希望はありません」になる）を真ん中に書く
   - そのグループに近い順に周辺に配置する

# 結果
## 1話

 - [こちらです](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/gakko_gurashi!_1.svg)

![](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/gakko_gurashi!_1.svg)

観察：

 - かなり初期に「※希望はありません」が出現
 - しかし最初から多かったわけではない。２万コメント後に大ヒット。
   - うまいタイミングでバズることが重要（？）
   - 一旦少なくなるタイミングもある
     - コメントは1000件のキューであることを考えると、待ち行列でよく見られるゆらぎの効果があるのか？
 - 初期は「夢がちゃんとある（〜とは言ってない）」のような従来的な歌詞に基づいたコメントが多い
   - その後も継続的に出現
 - おなじ「まんがタイムきらら」作品が元ネタの「夢もキボーもありゃしない」コメントも独立に散発的に出現している

## 2話

 - [こちらです](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/gakko_gurashi!_2.svg)

![](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/gakko_gurashi!_2.svg)

- 一話の流行を受けて「※希望はありません」は最初から大人気

# おまけ
## 「わかば＊ガール」1話

　わかば＊ガールの弾幕「TNP鳴らして」も調べました。

 - [こちらです](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/wakaba_girs_1.svg)

![](https://cdn.rawgit.com/ledyba/gakko-analyzer/master/wakaba_girs_1.svg)

 - 初期は「GONG鳴らして」、ないし素直に「ピンポン鳴らして」
 - 猛烈に書かれている瞬間があるが、よく見ると連投回避のためにちょっと文字変えてるのが固まってる
   - 実はクソ下品な弾幕書いてるのは一部の人間だけ説
   - 弾幕へ面白おかしく反応するコメントが「※」に比べて目立つのも気になる

# 使い方
```
go get github.com/ledyba/gakko-analyzer

/bin/gakko-download \
		-user "*****" \
		-pass "*****" \
		-video "sm****" \
		-when "コメント取得時刻のunixtime" > log.json

gakko-analyze -file log.json -words "夢,希望"
```
