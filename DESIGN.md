---
version: 1.0
name: Workbench
description: A data-dense control surface for the CloudWeGo microservice console. Quiet chrome, single indigo accent (#4F46E5), one near-black ink, and a tight neutral palette where surface-color shifts (canvas / raised / sunken) replace decorative shadows. Typography on a 14px base with a 4-step weight ladder (400 / 500 / 600 / 700). Three calibrated elevation levels reserved for floating layers (dropdown / popover / dialog). No gradients, no glow, no shimmer — every pixel earns its place by carrying information.

colors:
  # ── Brand: the only interactive hue ─────────────────────────────
  primary: "#4F46E5"
  primary-hover: "#4338CA"
  primary-active: "#3730A3"
  primary-soft: "#EEF2FF"
  primary-soft-strong: "#E0E7FF"
  primary-on: "#FFFFFF"
  ring: "rgba(79, 70, 229, 0.24)"

  # ── Surface (light theme; dark theme defined in `themes.dark`) ──
  canvas: "#FFFFFF"
  raised: "#FAFAFA"
  sunken: "#F5F5F7"
  overlay: "#FFFFFF"
  inverse: "#0A0A0B"

  # ── Ink (text on light) ─────────────────────────────────────────
  ink-strong: "#0A0A0B"
  ink: "#1F2024"
  ink-muted: "#52525B"
  ink-subtle: "#9CA3AF"
  ink-disabled: "#D4D4D8"
  ink-on-primary: "#FFFFFF"
  ink-on-inverse: "#FAFAFA"

  # ── Borders / dividers ──────────────────────────────────────────
  border-subtle: "#F1F1F3"
  border: "#E4E4E7"
  border-strong: "#D4D4D8"
  border-focus: "#4F46E5"

  # ── Semantic ────────────────────────────────────────────────────
  success: "#10B981"
  success-soft: "#ECFDF5"
  success-ink: "#047857"
  warning: "#F59E0B"
  warning-soft: "#FFFBEB"
  warning-ink: "#B45309"
  danger: "#EF4444"
  danger-soft: "#FEF2F2"
  danger-ink: "#B91C1C"
  info: "#3B82F6"
  info-soft: "#EFF6FF"
  info-ink: "#1D4ED8"

  # ── Chart palette (multi-series data viz) ───────────────────────
  chart-1: "#4F46E5"
  chart-2: "#0EA5E9"
  chart-3: "#10B981"
  chart-4: "#F59E0B"
  chart-5: "#EC4899"
  chart-6: "#8B5CF6"
  chart-grid: "#E4E4E7"
  chart-axis: "#9CA3AF"

themes:
  dark:
    primary: "#818CF8"
    primary-hover: "#A5B4FC"
    primary-active: "#C7D2FE"
    primary-soft: "rgba(129, 140, 248, 0.12)"
    primary-soft-strong: "rgba(129, 140, 248, 0.20)"
    primary-on: "#0A0A0B"
    ring: "rgba(129, 140, 248, 0.32)"
    canvas: "#0A0A0B"
    raised: "#111114"
    sunken: "#16161A"
    overlay: "#1C1C21"
    inverse: "#FAFAFA"
    ink-strong: "#FAFAFA"
    ink: "#E4E4E7"
    ink-muted: "#A1A1AA"
    ink-subtle: "#71717A"
    ink-disabled: "#3F3F46"
    ink-on-primary: "#0A0A0B"
    ink-on-inverse: "#0A0A0B"
    border-subtle: "#1F1F23"
    border: "#27272A"
    border-strong: "#3F3F46"
    border-focus: "#818CF8"
    success: "#34D399"
    success-soft: "rgba(52, 211, 153, 0.12)"
    success-ink: "#6EE7B7"
    warning: "#FBBF24"
    warning-soft: "rgba(251, 191, 36, 0.12)"
    warning-ink: "#FCD34D"
    danger: "#F87171"
    danger-soft: "rgba(248, 113, 113, 0.12)"
    danger-ink: "#FCA5A5"
    info: "#60A5FA"
    info-soft: "rgba(96, 165, 250, 0.12)"
    info-ink: "#93C5FD"
    chart-1: "#818CF8"
    chart-2: "#38BDF8"
    chart-3: "#34D399"
    chart-4: "#FBBF24"
    chart-5: "#F472B6"
    chart-6: "#A78BFA"
    chart-grid: "#27272A"
    chart-axis: "#71717A"

typography:
  font-display:
    family: '"PingFang SC", "Source Han Sans CN", "Noto Sans CJK SC", system-ui, -apple-system, "Segoe UI", "Microsoft YaHei", sans-serif'
  font-text:
    family: '"PingFang SC", "Source Han Sans CN", "Noto Sans CJK SC", system-ui, -apple-system, "Segoe UI", "Microsoft YaHei", sans-serif'
  font-mono:
    family: '"JetBrains Mono", "SF Mono", "Cascadia Code", "Source Code Pro", Consolas, monospace'

  display-xl:
    size: 28px
    weight: 600
    line-height: 1.2
    letter-spacing: "-0.015em"
  display-lg:
    size: 22px
    weight: 600
    line-height: 1.25
    letter-spacing: "-0.012em"
  display-md:
    size: 18px
    weight: 600
    line-height: 1.35
    letter-spacing: "-0.008em"
  body-lg:
    size: 15px
    weight: 400
    line-height: 1.55
    letter-spacing: "0"
  body:
    size: 14px
    weight: 400
    line-height: 1.55
    letter-spacing: "0"
  body-strong:
    size: 14px
    weight: 600
    line-height: 1.55
    letter-spacing: "0"
  caption:
    size: 13px
    weight: 400
    line-height: 1.5
    letter-spacing: "0"
  caption-strong:
    size: 13px
    weight: 500
    line-height: 1.5
    letter-spacing: "0"
  micro:
    size: 12px
    weight: 500
    line-height: 1.4
    letter-spacing: "0.01em"
  overline:
    size: 11px
    weight: 600
    line-height: 1.3
    letter-spacing: "0.06em"
    transform: uppercase
  code:
    size: 13px
    weight: 400
    line-height: 1.55
    letter-spacing: "0"
  kbd:
    size: 12px
    weight: 500
    line-height: 1.0
    letter-spacing: "0"

rounded:
  none: 0px
  xs: 4px
  sm: 6px
  md: 8px
  lg: 12px
  xl: 16px
  pill: 9999px
  full: 9999px

spacing:
  px: 1px
  0.5: 2px
  1: 4px
  1.5: 6px
  2: 8px
  2.5: 10px
  3: 12px
  4: 16px
  5: 20px
  6: 24px
  8: 32px
  10: 40px
  12: 48px
  16: 64px
  20: 80px

elevation:
  none: "none"
  e1: "0 1px 2px 0 rgba(15, 23, 42, 0.04), 0 0 0 1px rgba(15, 23, 42, 0.04)"
  e2: "0 4px 12px -2px rgba(15, 23, 42, 0.08), 0 2px 4px -2px rgba(15, 23, 42, 0.04)"
  e3: "0 16px 40px -12px rgba(15, 23, 42, 0.18), 0 4px 12px -4px rgba(15, 23, 42, 0.08)"
  focus-ring: "0 0 0 3px var(--ring)"
  inner-hairline: "inset 0 0 0 1px var(--border)"

motion:
  duration-instant: 80ms
  duration-fast: 140ms
  duration-base: 200ms
  duration-slow: 280ms
  ease-out: cubic-bezier(0.16, 1, 0.3, 1)
  ease-in-out: cubic-bezier(0.4, 0, 0.2, 1)
  ease-in: cubic-bezier(0.4, 0, 1, 1)

layout:
  sidebar-width: 240px
  sidebar-collapsed: 56px
  topbar-height: 56px
  content-max-width: 1440px
  content-padding-x: 24px
  content-padding-y: 24px
  table-row-height: 44px
  table-row-compact: 36px
  form-label-gap: 6px
  form-field-gap: 16px

components:
  button-primary:
    background: "{colors.primary}"
    color: "{colors.ink-on-primary}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.sm}"
    padding: "8px 14px"
    height: 36px
  button-primary-hover:
    background: "{colors.primary-hover}"
  button-primary-active:
    background: "{colors.primary-active}"
  button-secondary:
    background: "{colors.canvas}"
    color: "{colors.ink}"
    border: "1px solid {colors.border}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.sm}"
    padding: "8px 14px"
    height: 36px
  button-secondary-hover:
    background: "{colors.sunken}"
    border: "1px solid {colors.border-strong}"
  button-ghost:
    background: transparent
    color: "{colors.ink-muted}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
    height: 36px
  button-ghost-hover:
    background: "{colors.sunken}"
    color: "{colors.ink}"
  button-danger:
    background: "{colors.danger}"
    color: "{colors.ink-on-primary}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.sm}"
    padding: "8px 14px"
    height: 36px
  button-icon:
    background: transparent
    color: "{colors.ink-muted}"
    rounded: "{rounded.sm}"
    size: 32px
  button-icon-hover:
    background: "{colors.sunken}"
    color: "{colors.ink}"
  input:
    background: "{colors.canvas}"
    color: "{colors.ink}"
    border: "1px solid {colors.border}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
    height: 36px
  input-focus:
    border: "1px solid {colors.border-focus}"
    shadow: "{elevation.focus-ring}"
  input-error:
    border: "1px solid {colors.danger}"
  card:
    background: "{colors.canvas}"
    border: "1px solid {colors.border-subtle}"
    rounded: "{rounded.md}"
    padding: 20px
  card-elevated:
    background: "{colors.canvas}"
    border: "1px solid {colors.border-subtle}"
    rounded: "{rounded.md}"
    padding: 20px
    shadow: "{elevation.e1}"
  table-row:
    background: "{colors.canvas}"
    border-bottom: "1px solid {colors.border-subtle}"
    typography: "{typography.body}"
    height: "{layout.table-row-height}"
  table-row-hover:
    background: "{colors.sunken}"
  table-row-selected:
    background: "{colors.primary-soft}"
  table-header:
    background: "{colors.canvas}"
    color: "{colors.ink-muted}"
    typography: "{typography.micro}"
    border-bottom: "1px solid {colors.border}"
    height: 40px
  sidebar:
    background: "{colors.sunken}"
    border-right: "1px solid {colors.border-subtle}"
    width: "{layout.sidebar-width}"
  sidebar-item:
    color: "{colors.ink-muted}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    padding: "7px 10px"
    height: 32px
  sidebar-item-hover:
    background: "{colors.canvas}"
    color: "{colors.ink}"
  sidebar-item-active:
    background: "{colors.canvas}"
    color: "{colors.primary}"
    border-left: "2px solid {colors.primary}"
    weight: 600
  topbar:
    background: "{colors.canvas}"
    border-bottom: "1px solid {colors.border-subtle}"
    height: "{layout.topbar-height}"
    padding: "0 24px"
  search-omnibar:
    background: "{colors.sunken}"
    border: "1px solid transparent"
    color: "{colors.ink-muted}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    height: 32px
    padding: "0 12px"
  search-omnibar-focus:
    background: "{colors.canvas}"
    border: "1px solid {colors.border-focus}"
    shadow: "{elevation.focus-ring}"
  badge:
    background: "{colors.sunken}"
    color: "{colors.ink-muted}"
    typography: "{typography.micro}"
    rounded: "{rounded.xs}"
    padding: "2px 6px"
    height: 20px
  badge-primary:
    background: "{colors.primary-soft}"
    color: "{colors.primary-active}"
  badge-success:
    background: "{colors.success-soft}"
    color: "{colors.success-ink}"
  badge-warning:
    background: "{colors.warning-soft}"
    color: "{colors.warning-ink}"
  badge-danger:
    background: "{colors.danger-soft}"
    color: "{colors.danger-ink}"
  status-dot:
    size: 6px
    rounded: "{rounded.full}"
  dropdown:
    background: "{colors.overlay}"
    border: "1px solid {colors.border}"
    rounded: "{rounded.md}"
    padding: "4px"
    shadow: "{elevation.e2}"
  dropdown-item:
    color: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.xs}"
    padding: "6px 8px"
    height: 32px
  dropdown-item-hover:
    background: "{colors.sunken}"
  popover:
    background: "{colors.overlay}"
    border: "1px solid {colors.border}"
    rounded: "{rounded.md}"
    padding: 16px
    shadow: "{elevation.e2}"
  dialog:
    background: "{colors.overlay}"
    rounded: "{rounded.lg}"
    shadow: "{elevation.e3}"
    padding: 24px
    max-width: 480px
  dialog-overlay:
    background: "rgba(15, 23, 42, 0.40)"
  toast:
    background: "{colors.inverse}"
    color: "{colors.ink-on-inverse}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    shadow: "{elevation.e2}"
    padding: "10px 14px"
    height: 40px
  kbd:
    background: "{colors.canvas}"
    color: "{colors.ink-muted}"
    border: "1px solid {colors.border}"
    typography: "{typography.kbd}"
    rounded: "{rounded.xs}"
    padding: "2px 6px"
    height: 20px
  pagination-button:
    background: transparent
    color: "{colors.ink-muted}"
    border: "1px solid {colors.border}"
    typography: "{typography.body}"
    rounded: "{rounded.sm}"
    size: 32px
  pagination-button-active:
    background: "{colors.primary}"
    color: "{colors.ink-on-primary}"
    border: "1px solid {colors.primary}"
  empty-state:
    background: transparent
    color: "{colors.ink-muted}"
    typography: "{typography.body}"
    padding: "48px 24px"
  skeleton:
    background: "{colors.sunken}"
    rounded: "{rounded.xs}"
    shimmer: "linear-gradient(90deg, transparent, rgba(15, 23, 42, 0.04), transparent)"
---

## 总览

CloudWeGo 微服务管理后台是一个**数据密集的工作台**——管理员每天打开它做的事是查表、改配置、读审计、调权限。设计的唯一目标是让数据自己说话，UI 只是托住数据的骨架。任何"营销味"装饰（径向渐变、流光、光晕、卡通形象、彩色统计图标）一律视为污染。

视觉语言围绕三件事建立：

1. **唯一交互色**：深靛蓝 `{colors.primary}`（#4F46E5）。整个站点中"可点"的元素只用这一个颜色——主按钮填充、次按钮文字、链接、Focus 环、激活态边线、选中行底色（淡靛蓝 `{colors.primary-soft}`）。语义色（success/warning/danger/info）只用于**状态徽标与提示**，不参与交互意图。
2. **表面色制造层级**：白底 `{colors.canvas}`、抬升 `{colors.raised}`、下沉 `{colors.sunken}` 三档表面色之间相差 1–4%，肉眼几乎难分辨，但足以替代大多数阴影。Sidebar 用 `sunken`、内容区用 `canvas`、表格 hover 行回到 `sunken`——层级靠"色温"而非"投影"。
3. **三档阴影只用于浮层**：默认所有元素 elevation = none。`{elevation.e1}` 仅用于卡片轻微抬升、`{elevation.e2}` 用于 dropdown / popover、`{elevation.e3}` 用于 dialog / drawer。**永远不会**给按钮、输入框、表格行、徽标加阴影。

字体回到中文优先：`PingFang SC` 在 macOS 自动命中，`Microsoft YaHei` 在 Windows 自动命中，思源黑体作为开源回退。基础字号是 14px——不是 Apple 营销页的 17px，因为后台一屏要装下表格、表单、详情双栏，14px 是数据密度与可读性的平衡点。权重梯队 400 / 500 / 600 / 700，**没有 300**（300 适合营销大字号，后台用不到）。

**关键特征：**
- 一种交互色（`{colors.primary}`），一种墨色（`{colors.ink}`），三档表面（canvas/raised/sunken），四档语义（success/warning/danger/info）。
- Sidebar 240px 展开 / 56px 折叠；Topbar 56px 高度；内容区 24px 内边距，最大宽度 1440px。
- 表格行高 44px（紧凑模式 36px），表头 40px 大写小字 + tracking +0.01em。
- 圆角语法：`{rounded.xs}` 4px（badge/tag/code）→ `{rounded.sm}` 6px（按钮/输入/dropdown 项）→ `{rounded.md}` 8px（卡片/dropdown/popover）→ `{rounded.lg}` 12px（dialog/drawer）→ `{rounded.pill}`（status dot 等极少处）。
- Focus 环统一为 3px 半透明靛蓝 (`{elevation.focus-ring}`)，不是边框加粗——边框加粗会撑开布局。
- 一切动效 ≤ 280ms，缓动函数只用 `{motion.ease-out}`（默认）和 `{motion.ease-in-out}`（双向状态）。**禁止旋转、翻转、shimmer 装饰、parallax、spotlight、跟随鼠标光晕**。
- 主题：亮色为默认，暗色为对偶。`themes.dark` 中的所有 token 一一映射，组件实现层只用语义 token（如 `var(--ink)`），不直接引用 hex。

## 颜色

### 品牌交互色 — 深靛蓝（唯一）

| Token | Hex | 用途 |
|---|---|---|
| `{colors.primary}` | #4F46E5 | 主按钮、链接默认色、激活 Sidebar 项左侧 2px 边线、Tab 激活态、Switch 开启态、Checkbox 勾选态 |
| `{colors.primary-hover}` | #4338CA | 主按钮 hover |
| `{colors.primary-active}` | #3730A3 | 主按钮按下；徽标 primary 文字色（在淡靛蓝底上保证 4.5:1 对比） |
| `{colors.primary-soft}` | #EEF2FF | Sidebar 激活项背景、表格选中行、徽标 primary 背景、Tab 激活态背景（可选） |
| `{colors.primary-soft-strong}` | #E0E7FF | hover 强化、avatar 占位 |
| `{colors.ring}` | rgba(79,70,229,0.24) | 所有 Focus 环（按钮、输入、链接、Tab）|

**铁律：交互色只有这一个**。不要为"看起来更生动"在某个图标上加紫色或青色——视觉一致性比"丰富"重要 100 倍。

### 表面色 — 三档梯队

| Token | Hex | 用途 |
|---|---|---|
| `{colors.canvas}` | #FFFFFF | 主内容区底；卡片底；Topbar 底；Dialog/Popover 内层 |
| `{colors.raised}` | #FAFAFA | 极少使用——仅当卡片需要从 canvas 上轻微抬起且不希望加阴影时；当前项目可不用 |
| `{colors.sunken}` | #F5F5F7 | Sidebar 底；表格 hover 行；输入框 disabled 底；search omnibar 默认底；徽标 neutral 底 |
| `{colors.overlay}` | #FFFFFF | dropdown / popover / dialog 内容层（与 canvas 等值，但语义独立，便于暗色主题切换） |
| `{colors.inverse}` | #0A0A0B | toast 通知背景；视觉强调反转区域 |

> **核心理念**：白与近白之间的 4% 色差就是层级语言。亮色主题里 sidebar 不需要加边框（只需色差 + 1px `border-subtle`），dropdown 浮在 canvas 之上靠 `e2` 阴影，按钮 hover 不"抬起"，只是底色变成 `sunken`。

### 文字

| Token | Hex | 用途 |
|---|---|---|
| `{colors.ink-strong}` | #0A0A0B | 页面主标题（display-xl/lg） |
| `{colors.ink}` | #1F2024 | 默认正文、表格单元、表单值 |
| `{colors.ink-muted}` | #52525B | 次级说明、表头、描述文字 |
| `{colors.ink-subtle}` | #9CA3AF | 占位符、辅助提示、面包屑分隔符 |
| `{colors.ink-disabled}` | #D4D4D8 | disabled 文字 |
| `{colors.ink-on-primary}` | #FFFFFF | 主按钮文字、Tag-primary 文字 |

不用纯黑（#000）。`#0A0A0B` 是带 1% 蓝调的近黑，与 `#FFFFFF` 配对时观感比纯黑更"光学正确"。

### 边线

| Token | Hex | 用途 |
|---|---|---|
| `{colors.border-subtle}` | #F1F1F3 | 卡片外边线、表格行分隔、面包屑分隔 |
| `{colors.border}` | #E4E4E7 | 输入框/按钮/dropdown 外边线、表头底线 |
| `{colors.border-strong}` | #D4D4D8 | 输入 hover 边线、强调分隔（极少） |
| `{colors.border-focus}` | #4F46E5 | 输入 focus 时的实色边线（搭配 ring 阴影） |

### 语义色

| 语义 | 实色 / 软底 / 文字色 | 用途 |
|---|---|---|
| Success | `{colors.success}` / `{colors.success-soft}` / `{colors.success-ink}` | 操作成功 toast；用户/服务状态"激活"徽标；进度条完成 |
| Warning | `{colors.warning}` / `{colors.warning-soft}` / `{colors.warning-ink}` | 配额预警；待审批徽标；非破坏性提示横幅 |
| Danger | `{colors.danger}` / `{colors.danger-soft}` / `{colors.danger-ink}` | 删除按钮；锁定/禁用徽标；表单校验错误；破坏性二次确认 |
| Info | `{colors.info}` / `{colors.info-soft}` / `{colors.info-ink}` | 中性通知；新功能横幅；提示性 callout |

**语义色不能用作"主色替代"**——不要因为某个表格"想突出"就把 success 当主色用。它们只描述**状态**，不描述**操作意图**。

### 暗色主题

`themes.dark` 中的每个 token 与亮色主题一一对应。实现层（CSS 变量）只暴露语义名（`var(--color-canvas)`、`var(--color-ink)` 等），切换主题时通过修改 `:root[data-theme="dark"]` 覆盖变量值。组件代码**不允许直接写 hex**。

## 字体

### 字体族

```
font-display / font-text:
  "PingFang SC", "Source Han Sans CN", "Noto Sans CJK SC",
  system-ui, -apple-system, "Segoe UI", "Microsoft YaHei", sans-serif

font-mono:
  "JetBrains Mono", "SF Mono", "Cascadia Code",
  "Source Code Pro", Consolas, monospace
```

display 与 text 共用同一字体栈——CJK 字体在不同字号下渲染稳定，没必要做 SF Pro Display / Text 那样的双栈切换。Display 与 text 的差异完全靠**字号 + 字重 + 字距**控制。

### 字号梯队

| Token | Size | Weight | Line-height | Letter-spacing | 用途 |
|---|---|---|---|---|---|
| `{typography.display-xl}` | 28px | 600 | 1.2 | -0.015em | 页面主标题（每页≤1 处，如 Dashboard 页） |
| `{typography.display-lg}` | 22px | 600 | 1.25 | -0.012em | 卡片大标题、Drawer 标题 |
| `{typography.display-md}` | 18px | 600 | 1.35 | -0.008em | 区块小标题、Dialog 标题 |
| `{typography.body-lg}` | 15px | 400 | 1.55 | 0 | 详情页阅读型描述、空态说明 |
| `{typography.body}` | 14px | 400 | 1.55 | 0 | **默认正文**——表格、表单值、按钮文字、菜单项 |
| `{typography.body-strong}` | 14px | 600 | 1.55 | 0 | 表单 label、按钮文字（实质上 body 与 body-strong 共用按钮，按钮以 strong 为准）、强调正文 |
| `{typography.caption}` | 13px | 400 | 1.5 | 0 | 次级描述、表格辅助列、表单 helper text |
| `{typography.caption-strong}` | 13px | 500 | 1.5 | 0 | 次级强调（如表单错误提示） |
| `{typography.micro}` | 12px | 500 | 1.4 | 0.01em | 徽标、Tag、状态标签、表头小写文字 |
| `{typography.overline}` | 11px | 600 | 1.3 | 0.06em + uppercase | 区块小标签（"OVERVIEW"、"SYSTEM"），罕用 |
| `{typography.code}` | 13px | 400 | 1.55 | 0 | 行内代码、API 端点、ID |
| `{typography.kbd}` | 12px | 500 | 1.0 | 0 | 快捷键提示（如 ⌘K） |

### 原则

- **基础正文 14px**，不是 16/17px。后台密度优先；用 14px 才能让一屏装下完整表格 + 筛选器 + 分页。
- **微负字距只用于 display 级别**（≥18px）。14px 及以下保持 0，避免中英文混排时的卡顿。
- **唯一加正字距的是 micro 与 overline**（+0.01em / +0.06em）。这两类 token 用于"系统语气"标签（表头、UPPERCASE overline），加宽字距让它们读起来像"标签"而不是"句子"。
- **权重梯队 400 / 500 / 600 / 700，无 300**。500 仅用于 micro 与 caption-strong（标签语境）；正文体量 strong 一律走 600。
- **Line-height 1.55 是正文铁律**。读表格、读详情，行高低于 1.5 会挤压；高于 1.6 又显得松散。
- **数字用 tabular-nums**：所有表格数字、统计值、价格、计数应用 `font-variant-numeric: tabular-nums` 让数字等宽对齐。CSS 实现层会通过工具类 `.tabular` 提供。

## 布局

### 间距系统

基础单位 4px。`{spacing}` 提供 `0.5 (2px) → 20 (80px)` 共 13 档。规则：

- 组件内紧密元素（icon ↔ 文字、表单 label ↔ 输入）用 `{spacing.1.5}`–`{spacing.2}`（6–8px）。
- 组件内常规间距（输入框 padding、按钮内边距）用 `{spacing.2}`–`{spacing.3}`（8–12px）。
- 卡片内边距用 `{spacing.5}`（20px）；卡片与卡片之间垂直间距 `{spacing.4}`（16px）。
- 区块（section）之间垂直间距 `{spacing.6}`（24px）。
- 页面主内容到 Topbar 的留白 `{spacing.6}`（24px）。
- **不要使用 32px 以上的间距来"营造高级感"**——后台密度优先。

### 网格

- **内容最大宽度**：`{layout.content-max-width}` = 1440px。超过此宽度的屏幕，左右留白吸收溢出。
- **内容左右内边距**：`{layout.content-padding-x}` = 24px。
- **表格列网格**：固定列宽优先；最大灵活列（描述、备注）用 `min-content` + `max-content` 约束；操作列固定 160px 右对齐。
- **表单网格**：单栏宽度 480px（单字段编辑）；详情页双栏 240px label + 1fr 内容。

### Sidebar

- 展开 240px / 折叠 56px。**不要再用 64px**——56px 才是工程后台的"刚好够"宽度（GitHub、Linear、Vercel 都是 56）。
- 顶部 logo 区高度 56px（与 Topbar 对齐）。
- 菜单项高度 32px（紧凑），padding 7×10px，激活态左侧 2px `{colors.primary}` 边线 + 底色 `{colors.canvas}`（在 `sunken` sidebar 上抬起）。
- 二级菜单缩进 20px，不再加缩进图标。

### Topbar

- 固定高度 56px。背景 `{colors.canvas}`，下边线 `{colors.border-subtle}`。
- 左侧：菜单折叠按钮（icon-only ghost）+ 面包屑。
- 中间：全局搜索 omnibar（默认占位 "搜索菜单、用户、组织…  ⌘K"）。
- 右侧：主题切换、语言切换、用户菜单。**不要在 Topbar 上挂统计徽标**（"3 条新通知"那种），通知放进用户菜单内层。

## 阴影

| Token | 值 | 唯一用途 |
|---|---|---|
| `{elevation.none}` | none | **默认值**——按钮、输入、表格行、徽标、Tag、卡片（除非"抬升卡片"语义） |
| `{elevation.e1}` | 0 1px 2px / 0 0 0 1px | 极少——抬升型卡片（如 Dashboard 主卡片） |
| `{elevation.e2}` | 0 4px 12px -2px / 0 2px 4px -2px | dropdown、popover、Toast |
| `{elevation.e3}` | 0 16px 40px -12px / 0 4px 12px -4px | dialog、drawer |
| `{elevation.focus-ring}` | 0 0 0 3px var(--ring) | 所有 focus 状态 |
| `{elevation.inner-hairline}` | inset 0 0 0 1px var(--border) | 替代 border 时使用（不撑开布局） |

**阴影哲学**：默认元素 elevation = none。需要强调时**优先换表面色**（canvas → sunken），其次才考虑加阴影。

## 形状

### 圆角梯队

| Token | Value | 唯一用途 |
|---|---|---|
| `{rounded.none}` | 0 | 全屏遮罩、Drawer 直边 |
| `{rounded.xs}` | 4px | Badge、Tag、Code 行内代码、Kbd 键帽 |
| `{rounded.sm}` | 6px | 按钮、输入框、Select、Switch、Dropdown 内项 |
| `{rounded.md}` | 8px | 卡片、Dropdown 容器、Popover、Tooltip |
| `{rounded.lg}` | 12px | Dialog、Drawer 模态、大型空态卡片 |
| `{rounded.xl}` | 16px | 极少——某些落地页/示意图 |
| `{rounded.pill}` | 9999px | Status dot、计数徽标（带数字，如"3"）|
| `{rounded.full}` | 9999px | Avatar、icon-only 浮起按钮 |

**与 Apple 的差异**：Apple 把 `pill` 用作主按钮形状，因为它要"召唤"点击。Workbench 后台用 `sm`（6px）按钮——按钮在工作台是高频操作，不需要"召唤感"，只要"清晰可点"。pill 反而留给徽标和 status dot。

### 图像几何

- **Avatar**：圆形，三档 24/32/40px。无字符 avatar 用 `{colors.primary-soft-strong}` 底 + `{colors.primary-active}` 首字母（不要用渐变）。
- **Logo**：固定纵横比，最大高度 32px。Sidebar 顶部 logo 区域 56×56 内居中。
- **数据图表**：禁用 3D、阴影、渐变填充。色板见下文 **Charts** 章节。

## Charts（数据可视化）

后台需要展示折线、柱状、饼图、漏斗、热力等多系列图表时，**统一使用预定义 6 色梯队**，避免每图各搞一套色。

### 色板

| Token | 亮色 | 暗色 | 语义 |
|---|---|---|---|
| `{colors.chart-1}` | #4F46E5 | #818CF8 | 主系列（与 primary 同源，确保"主指标"识别度） |
| `{colors.chart-2}` | #0EA5E9 | #38BDF8 | 第二系列（青蓝，与 primary 拉开 90° 色相） |
| `{colors.chart-3}` | #10B981 | #34D399 | 第三系列（success 绿，可用于"健康/达成"语义） |
| `{colors.chart-4}` | #F59E0B | #FBBF24 | 第四系列（warning 黄，可用于"预警/待处理"语义） |
| `{colors.chart-5}` | #EC4899 | #F472B6 | 第五系列（粉，纯视觉区分用） |
| `{colors.chart-6}` | #8B5CF6 | #A78BFA | 第六系列（紫，纯视觉区分用） |
| `{colors.chart-grid}` | #E4E4E7 | #27272A | 网格线（与 border 同色） |
| `{colors.chart-axis}` | #9CA3AF | #71717A | 轴线、刻度文字（与 ink-subtle 同源） |

### 使用规则

1. **单系列图**：用 `chart-1`（与品牌主色一致，强化"主数据"识别）。
2. **2 系列**：`chart-1` + `chart-2`（指标 vs 对比基线，例如本周 vs 上周）。
3. **3-6 系列**：按 `chart-1 → chart-6` 顺序递增。**不要跳号**，不要把 `chart-3`（绿）当装饰色用——它在多系列里出现时读者会下意识把绿当作"达成 / 正向"语义。
4. **>6 系列**：先尝试合并次要系列（"其他 N 项"），如果必须铺开，循环回 `chart-1` 但加 30% 透明度做区分。永远不引入第 7 个 hex。
5. **饼图/堆叠图**：优先按数值大小降序映射 `chart-1..6`，让最重要的切片自然成为主色。
6. **空态**：图表无数据时用 `{colors.ink-subtle}` 画一条横向虚线 + 居中"暂无数据"文案，**不要**用 chart 色板表示"0"。

### 图表样式守则

- **网格线**：1px solid `{colors.chart-grid}`；只画水平网格，不画垂直网格（除非是热力图）。
- **轴**：轴线 1px solid `{colors.chart-grid}`；刻度文字 12px / 400 / `{colors.chart-axis}`。
- **数据点 / 柱顶 / 切片**：禁用阴影、禁用渐变填充、禁用 3D 透视。
- **悬浮 tooltip**：复用 `{component.popover}` 形状（`{rounded.md}` + `{elevation.e2}`）；内部数字用 `tabular`。
- **图例**：放在图右侧或底部，用 6×6 圆点 + `{typography.caption}` 文字；禁止把图例伪装成可点按钮（除非真的能切换显隐）。
- **图表容器**：放在 `{component.card}` 内（`{rounded.md}` + `{colors.border-subtle}`），不要再叠一层阴影。

### 推荐库与适配

项目暂无图表库引入。需要时优先级：

1. **Recharts** / **VChart**（声明式，对 Tailwind 主题对接友好）—— `<Line stroke="var(--color-chart-1)" />` 即可。
2. **ECharts**（功能最全，但需要传 theme 对象）—— 提供一个 `wbChartTheme` 工具函数把 `chart-1..6` 注入 ECharts theme。
3. **D3**（最自由，最重）—— 直接用 CSS 变量。

**禁止**：Chart.js（默认色板和动效与 DESIGN.md 冲突，定制成本高）。

## 组件

### 按钮

后台只需要 4 种按钮 + 1 种 icon-only：

- **`button-primary`** — 高斯密度 36px / 6px 圆角 / `{colors.primary}` 实色 / 白字 14/600。每个区块**最多 1 个**主按钮。
- **`button-secondary`** — 36px / 6px / `{colors.canvas}` 底 + `{colors.border}` 边线 + ink 文字。次操作。
- **`button-ghost`** — 36px / 6px / 透明 + ink-muted 文字。第三梯队，filter 区 / Toolbar 内。
- **`button-danger`** — 36px / 6px / `{colors.danger}` 实色 / 白字。仅破坏性二次确认。
- **`button-icon`** — 32×32 / 6px / 透明 → hover sunken。Topbar、表格行操作、卡片右上角"更多"。

按下态**不要 transform: scale**——后台高频操作不需要"按压感"，做颜色加深 (`button-primary-active`) 即可。

### 输入与表单

- 输入框统一 36px 高，6px 圆角，14px body 字号。Focus 时 1px 实色靛蓝边线 + 3px ring 阴影。
- 表单 label 14/600 在输入上方，间距 6px。Helper text 13px caption 在输入下方，间距 6px。
- 错误态：边线变 `{colors.danger}`，下方红字 `{typography.caption-strong}` + 红色 `{colors.danger}`。
- **不再自创 search-input pill**：搜索输入和普通输入用同一形状（6px 圆角），靠左侧 leading icon 区分用途。

### 表格

后台的核心组件。规范：

- 表头 40px 高，`{typography.micro}`（12/500/+0.01em），背景 `{colors.canvas}`，下方 1px `{colors.border}`。
- 行 44px（默认）/ 36px（紧凑切换）。行间分隔 `{colors.border-subtle}`。
- Hover 行底色 `{colors.sunken}`；选中行底色 `{colors.primary-soft}`，**不要加边框**——色块本身就是层级。
- 排序图标：默认 `↕` (`{colors.ink-subtle}`)，激活 `↑/↓` (`{colors.primary}`)。
- 操作列固定右对齐，160–200px 宽，里面只放 1–2 个 ghost button 或 1 个 icon button + dropdown。
- **数字列用 `tabular-nums` + 右对齐**。
- 第一列若是名称/标题，加 24×24 avatar 或图标 + 名称两栏布局。

### Sidebar 菜单

参考 `{component.sidebar-item}` / `{component.sidebar-item-active}`。要点：

- 一级项 32px 高，`{typography.body}`（14/400）。
- 激活态：左侧 2px `{colors.primary}` 边线 + 背景 `{colors.canvas}` + 文字 `{colors.primary}` + 字重 600。**不要再用金色或大面积 primary 底色**。
- 二级项缩进 20px（不加图标），高度同 32px。
- 折叠态（56px sidebar）：只显示图标居中，hover tooltip 显示文字。

### Dropdown / Popover / Dialog

- Dropdown 与 Popover 共用形状（`{rounded.md}` + `{elevation.e2}`），区别在内容：dropdown 是命令列表，popover 是富内容（表单、说明）。
- Dropdown 内项 32px 高，`{rounded.xs}`，hover `{colors.sunken}`。
- Dialog 圆角 12px，最大宽度 480px（窄表单）/ 640px（标准）/ 800px（详情/向导）。Overlay 用 `rgba(15,23,42,0.40)`——不要 backdrop-blur，那是营销页装饰。

### 徽标 / Tag

- Badge 高度 20px，`{rounded.xs}` 4px 圆角，2×6 padding，`{typography.micro}`。
- 5 种语义：neutral / primary / success / warning / danger，分别对应 `sunken+ink-muted` / `primary-soft+primary-active` / 类推。
- Status dot：6px 圆点 + 文字。dot 用语义实色（success/warning/danger/info），文字用 ink。

### 空态

`{component.empty-state}` 规范：48×24 padding，居中。包含：

1. 一个 32×32 单色 lucide 图标（`{colors.ink-subtle}`）。
2. 一行 `{typography.body-lg}` ink 主文案（"暂无数据" / "尚未配置 OIDC"）。
3. 一行 `{typography.caption}` ink-muted 辅助文案。
4. 可选：一个 `button-secondary` 引导操作。

**禁止**：插图、SVG 卡通、emoji、动画、渐变背景。

### 骨架屏 / 加载

- 骨架块用 `{colors.sunken}` 底色 + `shimmer` 渐变扫光（线性 90deg，1.6s 循环），4px 圆角。
- 全局 loading 指示器：Topbar 下方 2px 进度条，靛蓝 + ease-out。
- **不要使用旋转 spinner 占满屏幕**——骨架屏是后台的标准态。

### 通知 Toast

`{component.toast}` 反色块（黑底白字 + e2 阴影），右下角弹出，4 秒后自动消失。一次最多 3 条，旧的上推。语义图标（成功/失败/警告/信息）放在文字左侧。

## 模式（Patterns）

后台高频出现的复合形态：

### 列表页（List Page）

```
[页面标题][主操作按钮]
[筛选条 chip / 搜索 omnibar]
[卡片：表格 + 分页]
[空态 / 骨架（when applicable）]
```

- 标题与主操作同行，标题左对齐 display-md，主操作右对齐 button-primary。
- 筛选条**横向单行**（不要立成卡片），左对齐多个 select / chip / 搜索框，右侧"重置 / 搜索"两个 ghost/secondary 按钮。
- 表格本体放在一个 `card`（无阴影，仅 1px border-subtle），内部 0 padding，让表格贴齐卡片边缘。
- 分页固定在卡片底部，左侧总数、中间页码、右侧"每页 N 条"select。

### 详情页（Detail Page）

```
[面包屑]
[标题区：avatar + 主信息 + 状态徽标 + 操作按钮组]
[标签页：概览 / 角色 / 审计 / ...]
[双栏布局：左 240px label / 右 1fr 值]
```

- 标题区**单层**（不要弄"卡片中卡片"），底部 1px border 分隔标签页。
- 标签页激活态：底部 2px `{colors.primary}` 下划线 + 文字加深 ink-strong。
- 双栏布局每行高 36px（与输入一致），label `{typography.caption-strong}` muted，值 `{typography.body}` ink。

### 表单页（Form Page）

- 单栏宽度 480px 居中（创建/编辑场景）。
- 区块之间用 `{typography.display-md}` 小标题 + 24px 间距分隔，**不要画卡片包起来**。
- 提交按钮区固定在表单底部，主按钮右对齐 + 取消 ghost 左侧。

### 全局搜索（Cmd+K）

- 触发：⌘K / Ctrl+K，键盘焦点直接进入搜索 input。
- 形态：居中 dialog，640×480px，`{rounded.lg}`。
- 内容：搜索 input + 分组结果列表（页面 / 用户 / 组织 / 角色 / 审计记录）。
- 每项 40px 高，左 icon + 主文字 + 路径 caption + 右侧 kbd 提示（`↵` 进入）。

## Do's and Don'ts

### Do
- 用 `{colors.primary}` 做唯一交互色——按钮、链接、Focus、激活、选中都共用它。
- 用表面色（canvas / sunken）切换替代阴影。Sidebar 用 sunken、内容用 canvas、表格 hover 回到 sunken。
- 字号守住 14px 正文。display-xl/lg 用于真正的页面/卡片标题。
- Focus 用 `{elevation.focus-ring}`（3px 半透明 ring），不要用粗边框。
- 数字列用 `tabular-nums` + 右对齐。
- 空态用 `{component.empty-state}` 范式（icon + 主辅文案，可选按钮），无插图无动画。
- 暗色主题切换通过 `:root[data-theme]` 覆盖语义 token。组件代码里不写 hex。

### Don't
- 不引入第二个交互色。哪怕"图标想丰富一些"——丰富的代价是失去识别度。
- 不给按钮、卡片、表格行、徽标加阴影。阴影只属于浮层。
- 不在背景挂 radial-gradient / 点状网格 / 模糊光斑 / shimmer 动画 / spotlight 跟随光晕。**全部删除**。
- 不硬编码字体（`font-['Inter']`）。所有字体走 `{typography.font-text}` 字体族变量。
- 不混用圆角语法。按钮全 6px，卡片全 8px，dialog 全 12px——任何"特例"都要先质疑。
- 不在 Sidebar 高亮态用大面积底色 + 字色变化。靛蓝左边线 + 白底 + ink 文字够用。
- 不用 transform: scale 做按钮按压。颜色加深更克制。
- 不在 Dialog overlay 用 backdrop-blur。模糊滤镜在 4K 屏 + 大区域时性能成本高，且后台不需要那种"惊艳感"。
- 不允许"统计图标用 4 种渐变色"那种装饰。统计卡片用单色 + 数字 + 描述即可。

## 响应式

| 名称 | 宽度 | 关键变化 |
|---|---|---|
| Mobile | ≤640px | 单栏、Sidebar 抽屉化（drawer-from-left）、Topbar 隐藏面包屑、表格水平滚动 |
| Tablet | 641–1024px | Sidebar 默认折叠（56px），双栏详情页变单栏 |
| Desktop | 1025–1440px | Sidebar 默认展开（240px），完整布局 |
| Wide | ≥1441px | 内容锁定 1440px，左右吸收余量 |

后台主战场是 ≥1024px，移动端只做"勉强可用"——管理员极少在手机上做权限配置。

### 触摸目标

- 最小 36×36（按钮、输入）。Mobile 表格行操作改为右滑显示，避免单手够不到右侧操作列。

## 迭代指南

1. **新组件先找 token**。需要某个色、某个间距、某个阴影时，先翻 `colors / spacing / elevation`——95% 的需求已有 token。
2. **新增 token 必须在 `DESIGN.md` 这里登记**，并且解释为什么旧 token 不够。私自加变量是设计失控的开端。
3. **变体（active / hover / disabled）**作为同一组件下的 sibling 条目（`button-primary` / `button-primary-hover` / `button-primary-active`）。
4. **永远只有一个交互色**。如果某个新场景"需要"绿色按钮——它需要的是 success 徽标，不是绿色按钮。
5. **写 CSS 时只用语义 token**：`var(--color-primary)`、`var(--color-canvas)`。不允许 `#4F46E5` 出现在组件代码里。
6. **暗色主题校验**：每加一个组件，立刻在 `[data-theme="dark"]` 下打开看一遍。亮色"看起来对"的色差，在暗色下经常完全消失或翻转——必须实测，不能想当然。

## 已知缺口

- **图表色板**已在 v1.0 补齐 `chart-1..6` + `chart-grid` / `chart-axis`（详见上文 Charts 章节）。
- **打印样式**（PDF 导出审计日志、用户报告）未规范。
- **i18n 字体回退**：日文 / 韩文场景未覆盖；当前优先简体中文 + 英文双语。需要时在 font 栈追加 `"Hiragino Sans"`、`"Apple SD Gothic Neo"`。
- **Drawer 组件**（侧滑抽屉）规范暂未列出，复用 dialog 形状 + 右贴边变体即可。
- **可访问性对比度**：本规范的所有色对均通过 WCAG AA（4.5:1 正文、3:1 大字号）；AAA 级别（7:1）未覆盖，未来如果有合规要求，需要把 ink-subtle 与 ink-muted 收紧。
