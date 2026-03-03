# nigopdf

PDF の結合と分割だけをシンプルに行う CLI ツール。

## nigopdf とは

[pdfcpu](https://github.com/pdfcpu/pdfcpu) は Go 製の高機能な PDF ライブラリですが、機能が非常に多く、日常的に使いたい操作にたどり着くまでのハードルが高くなりがちです。

nigopdf は **pdfcpu の中から「PDF の結合」と「PDF の分割」だけを切り出し**、コマンド一発で使えるようにしたツールです。

- 結合 → `nigopdf merge`
- 分割 → `nigopdf split`



## インストール

```bash
brew tap wadoyoka/tap
brew install nigopdf
```

ソースからビルドする場合:

```bash
go install github.com/wadoyoka/nigopdf@latest
```

## 使い方

### 結合 (merge)

ディレクトリ内の PDF ファイルを 1 つに結合します。

```bash
# カレントディレクトリの PDF を結合
nigopdf merge

# 指定ディレクトリの PDF を結合
nigopdf merge /path/to/pdfs

# 出力ファイル名を指定
nigopdf merge -o result.pdf

# 結合対象を確認（ドライラン）
nigopdf merge -n /path/to/pdfs

# サブディレクトリも含めて結合
nigopdf merge -r /path/to/pdfs
```

ファイル名の昇順で結合されます。デフォルトの出力ファイルは `merged.pdf` です。

### 分割 (split)

PDF ファイルを複数のパートに分割します。

```bash
# 3 分割
nigopdf split --parts 3 input.pdf

# 各パートが 10MB 以下になるように分割
nigopdf split --max-size 10MB input.pdf

# 出力先ディレクトリを指定
nigopdf split --parts 3 input.pdf -o ./out/
nigopdf split --max-size 1MB input.pdf -o ./out/
```

出力ファイルは `{元のファイル名}_1.pdf`, `{元のファイル名}_2.pdf`, ... の形式です。デフォルトでは入力ファイルと同じディレクトリに出力されます。

#### オプション

| オプション | 説明 |
|-----------|------|
| `--parts N` | N 個に均等分割（2 以上を指定） |
| `--max-size SIZE` | 各パートの最大サイズを指定（例: `10MB`, `500KB`） |
| `-o DIR` | 出力先ディレクトリ |

`--parts` と `--max-size` は同時に指定できません。サイズの単位は `B`, `KB`, `MB`, `GB` に対応しています。

> **補足**: `--max-size` では各ページのサイズを個別に見積もるため、フォントや画像などの共有リソースの影響で、実際のパートサイズは指定値より小さくなる傾向があります。これは安全側の挙動です。

## pdfcpu について

内部で [pdfcpu](https://github.com/pdfcpu/pdfcpu) (v0.11.1) を利用しています。pdfcpu は暗号化、最適化、フォーム操作、スタンプ、透かしなど非常に多くの機能を持つ PDF ライブラリですが、nigopdf はその中から **結合と分割に必要な API だけ** を使用しています。

pdfcpu の全機能に興味がある方は、公式リポジトリを参照してください。

## ライセンス

このプロジェクトは [MIT License](LICENSE) の下で公開されています。

### サードパーティライセンス

[THIRD_PARTY_LICENSES](THIRD_PARTY_LICENSES/) ディレクトリに全文があります。`nigopdf --credits` でも確認できます。
