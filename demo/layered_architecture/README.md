# アーキテクチャパターン比較学習ガイド

レイヤードアーキテクチャ、ヘキサゴナルアーキテクチャ、オニオンアーキテクチャ、クリーンアーキテクチャの違いを理解するための学習ドキュメントです。

---

## 4つのアーキテクチャの進化

```
レイヤードアーキテクチャ (1990年代〜)
        │
        │  問題: ビジネスロジックがDBに依存
        ▼
ヘキサゴナルアーキテクチャ (2005年, Alistair Cockburn)
        │
        │  Port/Adapter という概念の導入
        ▼
オニオンアーキテクチャ (2008年, Jeffrey Palermo)
        │
        │  同心円モデル、ドメイン中心
        ▼
クリーンアーキテクチャ (2012年, Robert C. Martin)
        │
        └  上記を統合・一般化
```

---

## ドキュメント構成

| # | ドキュメント | 内容 |
|---|-------------|------|
| 01 | [歴史と概要](docs/01_history_and_overview.md) | 4つのアーキテクチャの進化の流れ |
| 02 | [レイヤードアーキテクチャ](docs/02_layered_architecture.md) | 伝統的な水平分割アーキテクチャ |
| 03 | [ヘキサゴナルアーキテクチャ](docs/03_hexagonal_architecture.md) | Ports & Adapters パターン |
| 04 | [オニオンアーキテクチャ](docs/04_onion_architecture.md) | ドメイン中心の同心円モデル |
| 05 | [クリーンアーキテクチャ](docs/05_clean_architecture.md) | 統合された4層モデル |
| 06 | [比較まとめ](docs/06_comparison.md) | 4つのアーキテクチャの違いを整理 |

---

## クイック比較

| 観点 | レイヤード | ヘキサゴナル | オニオン | クリーン |
|-----|-----------|-------------|---------|---------|
| **年代** | 1990年代 | 2005年 | 2008年 | 2012年 |
| **依存の方向** | 上→下 | 外→内 | 外→内 | 外→内 |
| **中心概念** | 層 | Port/Adapter | ドメインモデル | Entity/UseCase |
| **DB層の位置** | 最下層 | 外側(Adapter) | 最外層 | 最外層 |
| **複雑さ** | 低 | 中 | 中 | 高 |

---

## 共通の目的

4つのアーキテクチャすべてに共通する目的：

- **関心の分離** (Separation of Concerns)
- **テスタビリティの向上**
- **保守性・変更容易性の確保**
- **ビジネスロジックの保護**

---

## 学習の進め方

1. **[01_history_and_overview.md](docs/01_history_and_overview.md)** で全体像を把握
2. **02〜05** で各アーキテクチャを個別に理解
3. **[06_comparison.md](docs/06_comparison.md)** で違いを整理
4. 実装例で実際のコードを確認（オプション）

---

## ディレクトリ構成

```
layered_architecture/
├── README.md                    # このファイル
├── docs/                        # 学習ドキュメント
│   ├── 01_history_and_overview.md
│   ├── 02_layered_architecture.md
│   ├── 03_hexagonal_architecture.md
│   ├── 04_onion_architecture.md
│   ├── 05_clean_architecture.md
│   └── 06_comparison.md
└── examples/                    # 実装例（オプション）
    ├── layered/
    ├── hexagonal/
    ├── onion/
    └── clean/
```

---

## 参考資料

- [The Clean Architecture - Robert C. Martin (2012)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture - Alistair Cockburn (2005)](https://alistair.cockburn.us/hexagonal-architecture/)
- [The Onion Architecture - Jeffrey Palermo (2008)](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/)
- [Patterns of Enterprise Application Architecture - Martin Fowler](https://martinfowler.com/eaaCatalog/)
