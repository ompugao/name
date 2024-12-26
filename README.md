名前候補探索器
==============

条件に当てはまる名前の候補を列挙します。五格分類法に対応しています。


使い方
-----

<details>
<summary>各コマンドのヘルプ</summary>

```console
$ name -h
Usage: name [subcommand] [options]

SUBCOMMANDS
  filter    name filter related commands
  search    search for given names
  info    show information about a given name

$ name search -h
Usage: name [options] <familyName>

OPTIONS
  -dir-dict string
        Directory of MeCab dictionary (full space only) (default "")
  -max-length int
        Maximum length of a given name (default 3)
  -min-length int
        Minimum length of a given name (default 1)
  -space string
        Search spaces (available: full, common) (default "common")
  -yomi-count int
        Number of Yomi-Gana candidates (default 5)

STDIN
        Filter notated in JSON. See "name filter validate --help" for details.

EXAMPLES
        $ name search 山田 < ./filter.example.json
        評点    画数    名前    読み    天格    地格    人格    外格    総格
        15      13      一喜    イッキ  吉      大吉    大吉    大大吉  大吉
        15      13      一喜    イッキ  吉      大吉    大吉    大大吉  大吉
        ...

$ name info -h
Usage: name info [options] <familyName> <givenName> <yomi>

EXAMPLES
        $ name info 山田 太郎 タロウ
        評点    画数    名前    読み    天格    地格    人格    外格    総格
        8       13      太郎    タロウ  吉      大吉    大凶    大凶    大吉

$ name filter validate -h
Usage: name filter validate

STDIN
        Filter notated in JSON.

                filter       := true | false | and | or | not | sex | length | mora | strokes | minRank | minTotalRank |
                        yomiCount | yomi | kanjiCount | kanji
                true         := {"true": {}}
                false        := {"false": {}}
                and          := {"and": [filter...]}
                or           := {"or": [filter...]}
                not          := {"not": filter}
                sex          := {"sex": "asexual" | "male" | "female"}
                length       := {"length": count}
                mora         := {"maxMora": count}
                strokes      := {"strokes": count}
                minRank      := {"minRank": 0-4} (4=大大吉, 3=大吉, 2=吉, 1=凶, 0=大凶)
                minTotalRank := {"minTotalRank": byte}
                yomiCount    := {"yomiCount": {"rune": rune, "count": count}}
                yomi         := {"yomi": match}
                kanjiCount   := {"kanjiCount": {"rune": rune, "count": count}}
                kanji        := {"kanji": match}
                count        := {"equal": byte} | {"greaterThan": byte} | {"lessThan": byte}
                match        := {"equal": string} | {"contain": string} | {"startWith": string} | {"endWith": string}
                byte         := 0-255
                rune         := string that contains only one rune

EXAMPLES
        $ name filter validate < valid-filter.json
        $ echo $?
        0

        $ name filter validate < invalid-filter.json
        $ echo $?
        1

$ name filter test -h
Usage: name filter test <familyName> <givenName> <yomi>

STDIN
        Filter notated in JSON. See "name filter validate --help" for details.

EXAMPLES
        $ name filter test 田中 太郎 タロウ < filter.json
        $ echo $?
        0

        $ name filter test 田中 太郎 タロウ < filter.json
        $ echo $?
        1
        
$ name filter apply -h
Usage: name filter apply <familyName> --to <path>
OPTIONS
  -to name search
        path to the result file of name search

STDIN
        Filter notated in JSON. See "name filter validate --help" for details.

EXAMPLES
        $ name filter apply 山田 --to /path/to/result.tsv < ./filter.example.json
        評点    画数    名前    読み    天格    地格    人格    外格    総格
        15      13      一喜    イッキ  吉      大吉    大吉    大大吉  大吉
        15      13      一喜    イッキ  吉      大吉    大吉    大大吉  大吉
```
</details>

まずフィルタを用意します。フィルタは名前の候補について真であれば候補を残し、偽であれば候補を除去します。

| 説明         | 構文                                                                      | 例                                                                                                                     |
|:-----------|:------------------------------------------------------------------------|:----------------------------------------------------------------------------------------------------------------------|
| 真          | `{"true": {}}`                                                          | `{"true": {}}`                                                                                                        |
| 偽          | `{"false": {}}`                                                         | `{"false": {}}`                                                                                                       |
| 論理積        | `{"and": [filter...]}`                                                  | `{"and": [{"yomiCount": {"rune": "ア", "count": {"equal": 1}}}, {"yomiCount": {"rune": "イ", "count": {"equal": 1}}}]}` |
| 論理和        | `{"or": [filter...]}`                                                   | `{"or": [{"yomiCount": {"rune": "ア", "count": {"equal": 1}}}, {"yomiCount": {"rune": "イ", "count": {"equal": 1}}}]}`  |
| 否定論理       | `{"not": filter}`                                                       | `{"not": {"yomiCount": {"rune": "ア", "count": {"equal": 1}}}}`                                                        |
| 性別         | `{"sex": sex}`                                                          | `{"sex": "asexual"}`                                                                                                  |
| 長さ         | `{"length": count}`                                                     | `{"length": 3}`                                                                                                       |
| 読み仮名のモーラ数  | `{"mora": count}`                                                       | `{"mora": {"equal": 3}}`                                                                                              |
| よくある読み仮名   | `{"commonYomi": {}}`                                                    | `{"commonYomi": {}}`                                                                                                  |
| 画数         | `{"strokes": count}`                                                    | `{"strokes": {"lessThan": 25}}`                                                                                       |
| 五格それぞれの最小値 | `{"minRank": 0-4}`（`4`=大大吉, `3`=大吉, `2`=吉, `1`=凶, `0`=大凶）               | `{"minRank": 3}`                                                                                                      |
| 五格の合計値の最小値 | `{"minTotalRank": byte}`                                                | `{"minTotalRank": 11}`                                                                                                |
| 指定した読み仮名の数 | `{"yomiCount": {"rune": string, "count": count}}`                       | `{"yomiCount": {"rune": "ア", "count": {"equal": 1}}}`                                                                 |
| 読み仮名のマッチ   | `{"yomi": match}`                                                       | `{"yomiMatch": {"exactly": "タロウ"}}`                                                                                   |                                                     
| 漢字のマッチ     | `{"kanji": match}`                                                      | `{"yomiMatch": {"exactly": "タロウ"}}`                                                                                   |                                                     
| 指定した漢字の数   | `{"kanjiCount": {"rune": string, "count": count}}`                      | `{"yomiMatch": {"exactly": "タロウ"}}`                                                                                   |                                                     
| `count`    | `{"equal": byte}` or `{"lessThan": byte}` or `{"greaterThan": byte}`    | `{"equal": 1}`                                                                                                        |
| `match`    | `{"equal": string}` or `{"startWith": string}` or `{"endWith": string}` | `{"equal": "タロウ"}`                                                                                                    |
| `sex`      | `"asexual"` or `"male"` or `"female"`                                   | `{"sex": "asexual"}`                                                                                                  |

<details>
<summary>フィルタの例</summary>

```json
{
  "and": [
    {"sex": "male"},
    {"commonYomi": {}},
    {"mora": {"equal": 3}},
    {"minRank": 2},
    {"minTotalRank": 11},
    {
      "or": [
        {
          "and": [
            {"yomiCount": {"rune": "ユ", "count": {"equal": 1}}},
            {"yomiCount": {"rune": "ウ", "count": {"equal": 0}}},
            {"yomiCount": {"rune": "サ", "count": {"lessThan": 2}}},
            {"yomiCount": {"rune": "キ", "count": {"equal": 0}}}
          ]
        },
        {
          "and": [
            {"yomiCount": {"rune": "ユ", "count": {"equal": 0}}},
            {"yomiCount": {"rune": "ウ", "count": {"equal": 1}}},
            {"yomiCount": {"rune": "サ", "count": {"lessThan": 2}}},
            {"yomiCount": {"rune": "キ", "count": {"equal": 0}}}
          ]
        },
        {
          "and": [
            {"yomiCount": {"rune": "ユ", "count": {"equal": 0}}},
            {"yomiCount": {"rune": "ウ", "count": {"equal": 0}}},
            {"yomiCount": {"rune": "サ", "count": {"equal": 0}}},
            {"yomiCount": {"rune": "キ", "count": {"equal": 1}}}
          ]
        }
      ]
    }
  ]
}
```
</details>

### 頻出空間探索

よくある人名の空間から名前候補を探索します。時間はほとんどかかりません。

```console
$ name search --space common 山田 < ./filter.json | tee result.tsv
評点    画数    名前    読み    性別    天格    地格    人格    外格    総格
15      16      匠真    ショウマ        両性    吉      大吉    大吉    大吉    大大吉
14      23      奨真    ショウマ        両性    吉      大吉    吉      大吉    大大吉
...
```

### 全空間探索

常用漢字+人名用漢字の空間から名前候補を探索します。かなり時間がかかります。現実的な時間で探索を終えるために `--max-length` を指定するなら `3` 以下を推奨します。

```console
$ name search --space full 山田 --max-length 2 < ./filter.json | tee result.tsv
評点    画数    名前    読み    性別    天格    地格    人格    外格    総格
14      16      丈辞    ジョウジ        男性    吉      大吉    吉      大吉    大大吉
13      21      丈騎    タケキ  男性    吉      大吉    吉      大吉    大吉
...
```

### 名前判定

```console
$ name info 山田 太郎 タロウ
評点    画数    名前    読み    天格    地格    人格    外格    総格
8       13      太郎    タロウ  吉      大吉    大凶    大凶    大吉
```

### フィルタ検査

```console
$ name filter validate < filter.json
$ echo $?
0
```

### フィルタ試験

```console
$ name filter test 山田 太郎 タロウ < filter.json
$ echo $?
1
```


### フィルタ再適用

```console
$ name filter apply 山田 --to /path/to/result.tsv < ./filter.json
評点    画数    名前    読み    性別    天格    地格    人格    外格    総格
14      16      丈辞    ジョウジ        男性    吉      大吉    吉      大吉    大大吉
13      21      丈騎    タケキ  男性    吉      大吉    吉      大吉    大吉
...
```


インストール方法
----------------
### macOS

1. `brew install mecab mecab-ipadic` を実行
5. [neologd](https://github.com/neologd/mecab-ipadic-neologd) をインストール（推奨）
3. 以下を実行：

    ```console
    $ export CGO_LDFLAGS="`mecab-config --libs`"
    $ export CGO_CFLAGS="`mecab-config --cflags`"
    $ go install github.com/Kuniwak/name
    ```


### Debian / Ubuntu

1. `sudo apt install mecab libmecab-dev mecab-ipadic-utf8`
5. [neologd](https://github.com/neologd/mecab-ipadic-neologd) をインストール（推奨）
3. 以下を実行：

    ```console
    $ export CGO_LDFLAGS="`mecab-config --libs`"
    $ export CGO_CFLAGS="`mecab-config --cflags`"
    $ go install github.com/Kuniwak/name
    ```


### Windows
1. 管理者権限の MinGW 環境で `.\assets\bin\install-mecab-mingw` を実行
2. 環境変数に `MECABRC=C:\MeCab\etc\mecabrc` を追加
3. 環境変数 `PATH` に `C:\MeCab\bin` を追加
4. `dotnet tool install -g MecabConfig` を実行 (.NET 8 が必要)
5. [neologd](https://github.com/neologd/mecab-ipadic-neologd) をインストール（推奨）
6. 以下を実行：

    ```console
    $ $Env:CGO_LDFLAGS = mecab-config --libs"
    $ $Env:CGO_CFLAGS = mecab-config --cflags"
    $ go install github.com/Kuniwak/name
    ```


ライセンス
---------
MIT License
