# 奢华摩羯座配色规范

## 设计理念

**配色主题**：奢华摩羯座 (Luxury Capricorn Color Scheme)

**核心价值**：稳重、优雅、追求卓越 | Steady, Elegant, Excellence

**设计哲学**：以深空为底，金线为骨，星光点缀，打造兼具沉稳与奢华的视觉体验

---

## 官方配色定义

### CSS 变量定义

```css
:root {
    /* Palette Definition - Luxury Capricorn */
    --bg-base: #141416;              /* 深岩灰 - 背景基础 */
    --c-accent: #D4AF37;             /* 香槟金 - 核心高亮/图标 */
    --c-text-main: #F2F0E4;         /* 羊皮纸白 - 主标题 */
    --c-text-sub: #8B9bb4;          /* 矿石灰 - 副文本 */
    --c-btn-bg: #2C2E33;            /* 青铜褐 - 按钮底色 */
    --c-highlight: #FFF8E7;          /* 亮光色 - 极亮部 */
}
```

### Tailwind CSS 配置

**tailwind.config.js** 或 **tailwind.config.ts**：

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  theme: {
    extend: {
      colors: {
        // 背景系统
        bg: {
          base: '#141416',
          radial: {
            from: '#2a2d35',
            to: '#000000',
          },
        },
        // 香槟金 - 核心高亮
        gold: {
          DEFAULT: '#D4AF37',
          light: '#E8C56B',
          dark: '#C4A033',
        },
        // 羊皮纸白 - 主文本
        vellum: {
          DEFAULT: '#F2F0E4',
          light: '#F8FAFC',
        },
        // 矿石灰 - 副文本
        mineral: {
          DEFAULT: '#8B9bb4',
          light: '#A8C4D9',
          dark: '#6E8A9F',
        },
        // 青铜褐 - 按钮/卡片
        bronze: {
          DEFAULT: '#2C2E33',
          light: '#3D4047',
          dark: '#1E1F22',
        },
        // 亮光色 - 极亮部
        highlight: {
          DEFAULT: '#FFF8E7',
        },
      },
      fontFamily: {
        // 字体系统
        cinzel: ['Cinzel', 'serif'],
        'noto-serif-sc': ['"Noto Serif SC"', 'serif'],
        lato: ['Lato', 'sans-serif'],
      },
      backgroundImage: {
        // 深空径向渐变
        'deep-space': 'radial-gradient(circle at 50% 10%, #2a2d35 0%, #000000 100%)',
        // 金色渐变
        'gold-gradient': 'linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37)',
        // 青铜渐变
        'bronze-gradient': 'linear-gradient(90deg, #2C2E33, #D4AF37)',
      },
      animation: {
        // 星星闪烁
        'twinkle': 'twinkle 3s ease-in-out infinite',
        // 文字光泽
        'shine': 'shine 5s linear infinite',
      },
      keyframes: {
        twinkle: {
          '0%, 100%': { opacity: '1', filter: 'brightness(1)' },
          '50%': { opacity: '0.6', filter: 'brightness(0.7)' },
        },
        shine: {
          'to': { backgroundPosition: '200% center' },
        },
      },
    },
  },
}
```

### Tailwind CSS 颜色类名映射

| 用途 | Tailwind 类名 | 原始颜色值 |
|------|--------------|-----------|
| 背景基础 | `bg-bg-base` | #141416 |
| 径向渐变背景 | `bg-deep-space` | radial-gradient |
| 香槟金文本 | `text-gold` | #D4AF37 |
| 香槟金边框 | `border-gold` | #D4AF37 |
| 羊皮纸白文本 | `text-vellum` | #F2F0E4 |
| 矿石灰文本 | `text-mineral` | #8B9bb4 |
| 青铜褐背景 | `bg-bronze` | #2C2E33 |
| 青铜褐边框 | `border-bronze` | #2C2E33 |
| 金色渐变文字 | `bg-gold-gradient` | linear-gradient |
| 青铜渐变进度条 | `bg-bronze-gradient` | linear-gradient |

---

## 配色详解

### 1. 深空径向渐变 / Deep Space Radial Gradient

**颜色值**：`radial-gradient(circle at 50% 10%, #2a2d35 0%, #000000 100%)`

**设计意图**：模拟深空宇宙，营造沉浸式、宁静的视觉氛围

**用途**：
- 全局背景色
- 深色模式基础
- 卡片容器背景

**示例**：
```css
body {
    background-color: #000;
    background-image: radial-gradient(circle at 50% 10%, #2a2d35 0%, #000000 100%);
}
```

**Tailwind**：
```html
<body class="bg-black bg-deep-space">
```

---

### 2. 香槟金 / Champagne Gold

**颜色值**：`#D4AF37`

**设计意图**：象征尊贵、优雅，如同星座中的星光连线

**用途**：
- 核心高亮色
- 图标颜色
- 边框和分割线
- 星星和装饰元素
- 强调文本

**示例**：
```css
.accent-text {
    color: #D4AF37;
}

.icon {
    color: #D4AF37;
}

.border-gold {
    border-color: #D4AF37;
}
```

**Tailwind**：
```html
<span class="text-gold">金色文本</span>
<div class="border border-gold">金色边框</div>
<svg class="fill-current text-gold">...</svg>
```

---

### 3. 羊皮纸白 / Vellum White

**颜色值**：`#F2F0E4`

**设计意图**：柔和的米白色，模拟羊皮纸质感，温暖且易读

**用途**：
- 主要文本
- 标题
- 内容文字
- 渐变文字（与金色配合）

**示例**：
```css
h1, h2, h3 {
    color: #F2F0E4;
}

p, span, div {
    color: #F2F0E4;
}

.gradient-text {
    background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}
```

**Tailwind**：
```html
<h1 class="text-vellum">标题</h1>
<p class="text-vellum/80">正文内容</p>
<span class="bg-gold-gradient text-transparent bg-clip-text">渐变文字</span>
```

---

### 4. 青铜褐 / Bronze Brown

**颜色值**：`#2C2E33`

**设计意图**：稳重厚实的青铜色调，体现摩羯座脚踏实地的特质

**用途**：
- 主按钮背景
- 卡片背景
- 次要 UI 元素
- 进度条填充

**示例**：
```css
.btn-primary {
    background: #2C2E33;
    border: 1px solid #D4AF37;
}

.card {
    background: rgba(44, 46, 51, 0.4);
}

.progress-fill {
    background: linear-gradient(90deg, #2C2E33, #D4AF37);
}
```

**Tailwind**：
```html
<button class="bg-bronze border border-gold">按钮</button>
<div class="bg-bronze/40 backdrop-blur">卡片</div>
<div class="bg-bronze-gradient h-2 rounded-full">进度条</div>
```

---

### 5. 矿石灰 / Mineral Gray

**颜色值**：`#8B9bb4`

**设计意图**：低饱和度的灰蓝色，如同夜空中的暗云，不抢眼但不可或缺

**用途**：
- 副文本
- 辅助信息
- 图标标签
- 描述性文字
- 次要装饰

**示例**：
```css
.subtitle, .description {
    color: #8B9bb4;
}

.icon-label {
    color: #8B9bb4;
}
```

**Tailwind**：
```html
<p class="text-mineral">副文本</p>
<span class="text-mineral/60">辅助信息</span>
```

---

### 6. 亮光色 / Highlight White

**颜色值**：`#FFF8E7`

**设计意图**：接近白色的暖色调，用于极亮部高光

**用途**：
- 极亮部高光
- 悬浮状态提示
- 特别强调

**示例**：
```css
.highlight {
    color: #FFF8E7;
}
```

**Tailwind**：
```html
<span class="text-highlight">高亮文本</span>
```

---

## 字体系统

### 官方字体搭配

```css
/* 英文标题 */
font-family: 'Cinzel', serif;

/* 中文标题 */
font-family: 'Noto Serif SC', serif;

/* 正文 */
font-family: 'Lato', 'Noto Serif SC', sans-serif;
```

### 字体选择理由

- **Cinzel**: 古典风格的衬线字体，完美契合摩羯座的历史感和优雅气质
- **Noto Serif SC**: 优雅的中文衬线字体，与 Cinzel 形成和谐的中英文搭配
- **Lato**: 现代无衬线字体，提供良好的可读性和简洁感

### 字体引入

```html
<link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&family=Noto+Serif+SC:wght@300;700&family=Lato:wght@300;400&display=swap" rel="stylesheet">
```

### Tailwind 配置

```javascript
fontFamily: {
  cinzel: ['Cinzel', 'serif'],
  'noto-serif-sc': ['"Noto Serif SC"', 'serif'],
  lato: ['Lato', 'sans-serif'],
}
```

### 使用示例

```html
<!-- 英文标题 -->
<h1 class="font-cinzel text-4xl text-gold">CAPRICORN</h1>

<!-- 中文标题 -->
<h2 class="font-noto-serif-sc text-2xl text-vellum">摩羯座</h2>

<!-- 正文 -->
<p class="font-lato text-mineral">正文内容</p>
```

---

## UI 组件配色规范

### 按钮组件

#### CSS 方式

```css
/* 主按钮 - 青铜褐底色 + 金色边框 */
.btn-primary {
    background: #2C2E33;
    color: #F2F0E4;
    border: 1px solid #D4AF37;
    padding: 14px 40px;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 2px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: 'Cinzel', serif;
}

.btn-primary:hover {
    background: #D4AF37;
    color: #000;
    box-shadow: 0 0 20px rgba(212, 175, 55, 0.4);
}

/* 次要按钮 - 透明 + 金色边框 */
.btn-secondary {
    background: transparent;
    color: #D4AF37;
    border: 2px solid #D4AF37;
    padding: 12px 30px;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 2px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.btn-secondary:hover {
    background: rgba(212, 175, 55, 0.1);
}

/* 强调按钮 - 金色渐变 */
.btn-accent {
    background: linear-gradient(135deg, #D4AF37 0%, #C4A963 100%);
    color: #000;
    padding: 14px 40px;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 2px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: 'Cinzel', serif;
    box-shadow: 0 4px 15px rgba(212, 175, 55, 0.4);
}
```

#### Tailwind CSS 方式

```html
<!-- 主按钮 -->
<button class="bg-bronze text-vellum border border-gold px-10 py-3.5 text-xs uppercase tracking-wider hover:bg-gold hover:text-black transition-all duration-300 shadow-lg font-cinzel">
  主按钮
</button>

<!-- 次要按钮 -->
<button class="bg-transparent text-gold border-2 border-gold px-8 py-3 text-xs uppercase tracking-wider hover:bg-gold/10 transition-all duration-300">
  次要按钮
</button>

<!-- 强调按钮 -->
<button class="bg-gradient-to-r from-gold to-gold-dark text-black px-10 py-3.5 text-xs uppercase tracking-wider hover:shadow-gold/40 transition-all duration-300 shadow-lg font-cinzel">
  强调按钮
</button>
```

---

### 卡片组件

#### CSS 方式

```css
.card {
    background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 20px;
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(10px);
    padding: 40px 30px;
    transition: transform 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.card:hover {
    transform: translateY(-5px);
    border-color: rgba(212, 175, 55, 0.3);
    box-shadow: 0 30px 60px rgba(212, 175, 55, 0.1);
}
```

#### Tailwind CSS 方式

```html
<!-- 卡片 -->
<div class="bg-gradient-to-br from-[rgba(30,32,36,0.9)] to-[rgba(20,20,22,0.95)] border border-white/5 rounded-[20px] shadow-[0_20px_50px_rgba(0,0,0,0.5)] backdrop-blur-md p-[30px_40px] hover:-translate-y-1 hover:border-gold/30 hover:shadow-[0_30px_60px_rgba(212,175,55,0.1)] transition-all duration-[400ms]">
  <!-- 卡片内容 -->
</div>
```

---

### 标签组件

#### CSS 方式

```css
/* 金色标签 */
.tag-gold {
    background: rgba(44, 46, 51, 0.6);
    color: #D4AF37;
    border: 1px solid #D4AF37;
    padding: 8px 16px;
    border-radius: 9999px;
    font-size: 12px;
}

/* 次要标签 */
.tag-secondary {
    background: rgba(212, 175, 55, 0.15);
    color: #F2F0E4;
    border: 1px solid #D4AF37;
    padding: 8px 16px;
    border-radius: 9999px;
    font-size: 12px;
}

/* 暗色标签 */
.tag-dim {
    background: rgba(212, 175, 55, 0.1);
    color: #8B9bb4;
    border: 1px solid rgba(212, 175, 55, 0.3);
    padding: 8px 16px;
    border-radius: 9999px;
    font-size: 12px;
}
```

#### Tailwind CSS 方式

```html
<!-- 金色标签 -->
<span class="bg-bronze/60 text-gold border border-gold px-4 py-2 rounded-full text-xs">
  ♑ 摩羯座
</span>

<!-- 次要标签 -->
<span class="bg-gold/15 text-vellum border border-gold px-4 py-2 rounded-full text-xs">
  ✨ 稳重
</span>

<!-- 暗色标签 -->
<span class="bg-gold/10 text-mineral border border-gold/30 px-4 py-2 rounded-full text-xs">
  🎯 追求卓越
</span>
```

---

### 输入框组件

#### CSS 方式

```css
.input-wrapper {
    position: relative;
    margin-bottom: 20px;
}

.input-label {
    position: absolute;
    left: 16px;
    top: -8px;
    background: rgba(26, 26, 46, 0.95);
    padding: 0 4px;
    font-size: 12px;
    color: #8B9bb4;
}

.input-field {
    width: 100%;
    padding: 14px 16px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: #F2F0E4;
    font-size: 14px;
    transition: all 0.3s ease;
}

.input-field:focus {
    outline: none;
    border-color: #D4AF37;
    box-shadow: 0 0 0 3px rgba(212, 175, 55, 0.1);
}

.input-field::placeholder {
    color: rgba(242, 240, 228, 0.4);
}
```

#### Tailwind CSS 方式

```html
<!-- 输入框 -->
<div class="relative mb-5">
  <label class="absolute left-4 top-[-8px] bg-[rgba(26,26,46,0.95)] px-1 text-xs text-mineral">
    用户名 / Username
  </label>
  <input
    type="text"
    placeholder="请输入用户名"
    class="w-full px-4 py-3.5 border border-white/10 rounded-xl bg-white/5 text-vellum text-sm transition-all duration-300 focus:outline-none focus:border-gold focus:shadow-[0_0_0_3px_rgba(212,175,55,0.1)] placeholder:text-vellum/40"
  />
</div>
```

---

### 进度条组件

#### CSS 方式

```css
.progress-bar {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 9999px;
    height: 8px;
    overflow: hidden;
}

.progress-fill {
    height: 100%;
    border-radius: 9999px;
    transition: all 0.5s ease;
}

/* 金色进度条 */
.progress-fill-gold {
    background: linear-gradient(90deg, #2C2E33, #D4AF37);
}

/* 次金进度条 */
.progress-fill-secondary {
    background: linear-gradient(90deg, #D4AF37, #C4A963);
}

/* 暗色进度条 */
.progress-fill-dim {
    background: linear-gradient(90deg, #2a2d35, #3a3d45);
}
```

#### Tailwind CSS 方式

```html
<!-- 进度条容器 -->
<div class="h-2 bg-white/10 rounded-full overflow-hidden">
  <!-- 金色进度条 -->
  <div class="h-full bg-bronze-gradient rounded-full transition-all duration-500" style="width: 75%"></div>
</div>

<!-- 次金进度条 -->
<div class="h-2 bg-white/10 rounded-full overflow-hidden">
  <div class="h-full bg-gold-gradient rounded-full transition-all duration-500" style="width: 60%"></div>
</div>

<!-- 暗色进度条 -->
<div class="h-2 bg-white/10 rounded-full overflow-hidden">
  <div class="h-full bg-gradient-to-r from-gray-800 to-gray-700 rounded-full transition-all duration-500" style="width: 90%"></div>
</div>
```

---

### 提示框组件

#### CSS 方式

```css
.alert {
    border-radius: 12px;
    padding: 16px;
    display: flex;
    align-items: flex-start;
    gap: 12px;
}

/* 信息提示 */
.alert-info {
    background: rgba(44, 46, 51, 0.4);
    border: 1px solid #8B9bb4;
    color: #F2F0E4;
}

.alert-info svg {
    color: #D4AF37;
    flex-shrink: 0;
}

/* 警告提示 */
.alert-warning {
    background: rgba(212, 175, 55, 0.15);
    border: 1px solid #D4AF37;
    color: #F2F0E4;
}

.alert-warning svg {
    color: #D4AF37;
    flex-shrink: 0;
}
```

#### Tailwind CSS 方式

```html
<!-- 信息提示 -->
<div class="bg-bronze/40 border border-mineral rounded-xl p-4 flex items-start gap-3">
  <svg class="text-gold flex-shrink-0" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
    <circle cx="12" cy="12" r="10"/>
    <path d="M12 16v-4M12 8h.01"/>
  </svg>
  <div>
    <h5 class="text-vellum text-sm">信息提示</h5>
    <p class="text-mineral text-xs">摩羯座的人通常具有很强的责任感和进取心</p>
  </div>
</div>

<!-- 警告提示 -->
<div class="bg-gold/15 border border-gold rounded-xl p-4 flex items-start gap-3">
  <svg class="text-gold flex-shrink-0" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
    <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
    <path d="M12 9v4M12 17h.01"/>
  </svg>
  <div>
    <h5 class="text-vellum text-sm">注意</h5>
    <p class="text-mineral text-xs">Capricorn: Dec 22 - Jan 19 | Earth Element</p>
  </div>
</div>
```

---

## 动画效果规范

### 文字光泽动画

**效果**：金色光泽从左到右缓慢流动

**CSS**：
```css
@keyframes shine {
    to { background-position: 200% center; }
}

.shine-text {
    background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
    background-size: 200% auto;
    animation: shine 5s linear infinite;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}
```

**Tailwind**：
```html
<h1 class="bg-gold-gradient bg-clip-text text-transparent bg-[length:200%_auto] animate-shine">
  摩羯座
</h1>
```

---

### 星星闪烁动画

**效果**：星星的透明度和亮度周期性变化

**CSS**：
```css
@keyframes twinkle {
    0%, 100% { opacity: 1; filter: brightness(1); }
    50% { opacity: 0.6; filter: brightness(0.7); }
}

.star {
    animation: twinkle 3s ease-in-out infinite;
}
```

**Tailwind**：
```html
<span class="animate-twinkle">★</span>
```

---

### 悬浮效果

**效果**：元素悬浮时向上浮动并增强阴影

**CSS**：
```css
.hover-lift {
    transition: all 0.3s ease;
}

.hover-lift:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
}
```

**Tailwind**：
```html
<div class="transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_15px_40px_rgba(0,0,0,0.4)]">
  内容
</div>
```

---

### 浮动动画

**效果**：元素在垂直方向上缓慢浮动

**CSS**：
```css
@keyframes float {
    0% { transform: translateY(0px); }
    50% { transform: translateY(-5px); }
    100% { transform: translateY(0px); }
}

.float-animation {
    animation: float 4s ease-in-out infinite;
}
```

**使用场景**：星星图标、装饰元素

---

### 连线绘制动画

**效果**：星座连线从起点到终点逐渐绘制

**CSS**：
```css
.constellation-lines {
    stroke-dasharray: 1000;
    stroke-dashoffset: 1000;
    animation: drawLine 3s ease-out forwards 0.5s;
}

@keyframes drawLine {
    to { stroke-dashoffset: 0; }
}
```

**SVG**：
```html
<path class="constellation-lines" d="M 60,80 L 140,50 L 220,90" />
```

---

## 图标颜色规范

### 图标颜色分层

| 图标类型 | 颜色值 | 用途 |
|---------|-------|------|
| 主要图标 | `#D4AF37` (香槟金) | 核心功能、重点强调 |
| 次要图标 | `#2C2E33` (青铜褐) | 辅助功能、背景装饰 |
| 辅助图标 | `#8B9bb4` (矿石灰) | 提示信息、次要装饰 |

### 使用示例

```html
<!-- 主要图标 -->
<svg class="text-gold" viewBox="0 0 24 24">
  <!-- 图标路径 -->
</svg>

<!-- 次要图标 -->
<svg class="text-bronze" viewBox="0 0 24 24">
  <!-- 图标路径 -->
</svg>

<!-- 辅助图标 -->
<svg class="text-mineral" viewBox="0 0 24 24">
  <!-- 图标路径 -->
</svg>
```

---

## 配色使用原则

### 1. 主次分明
- **金色**用于强调和引导注意力
- **青铜色**用于基础 UI 元素
- **灰色**用于辅助信息

### 2. 对比度控制
- 确保文本与背景对比度符合 WCAG AA 标准（4.5:1）
- 香槟金 (#D4AF37) 与深岩灰 (#141416) 对比度：12.6:1 ✓
- 羊皮纸白 (#F2F0E4) 与深岩灰 (#141416) 对比度：15.2:1 ✓

### 3. 一致性
- 所有 UI 组件必须遵循官方配色
- 不要引入未经批准的额外颜色

### 4. 渐变使用
- 仅使用金色渐变 (`#D4AF37` → `#F2F0E4` → `#D4AF37`)
- 不使用绿色、紫色等其他颜色渐变
- 青铜渐变仅用于进度条等特定场景

### 5. 动画节制
- 使用金色闪烁和光泽动画
- 避免过度动画，保持页面流畅（60fps）
- 动画时长控制在 3-5 秒之间

### 6. 留白与层次
- 利用留白突出重要内容
- 通过阴影和透明度创造层次感
- 避免过度拥挤的布局

---

## 实施检查清单

开发 UI 组件时，请检查：

- [ ] 背景使用深空径向渐变
- [ ] 标题使用羊皮纸白 (#F2F0E4)
- [ ] 高亮元素使用香槟金 (#D4AF37)
- [ ] 按钮背景使用青铜褐 (#2C2E33)
- [ ] 副文本使用矿石灰 (#8B9bb4)
- [ ] 字体使用 Cinzel + Noto Serif SC + Lato
- [ ] 进度条使用金色系渐变
- [ ] 标签使用青铜色系
- [ ] 使用 Tailwind CSS 配色时，确保使用官方定义的颜色变量
- [ ] 确保所有自定义颜色都已添加到 Tailwind 配置
- [ ] 文本与背景对比度符合 WCAG AA 标准
- [ ] 动画流畅且不影响性能
- [ ] 在不同设备上验证配色效果

---

## 常见配色场景

### 页面标题

**CSS**：
```css
.page-title {
    font-family: 'Cinzel', serif;
    font-size: 2.5rem;
    letter-spacing: 4px;
    color: #D4AF37;
    background: linear-gradient(to right, #D4AF37, #F2F0E4, #D4AF37);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    animation: shine 5s linear infinite;
}
```

**Tailwind**：
```html
<h1 class="font-cinzel text-4xl tracking-widest text-gold bg-gold-gradient bg-clip-text text-transparent bg-[length:200%_auto] animate-shine">
  摩羯座
</h1>
```

---

### 导航菜单

**CSS**：
```css
.nav-link {
    color: #8B9bb4;
    transition: color 0.3s ease;
}

.nav-link:hover,
.nav-link.active {
    color: #D4AF37;
}
```

**Tailwind**：
```html
<nav>
  <a class="text-mineral transition-colors duration-300 hover:text-gold">首页</a>
  <a class="text-mineral transition-colors duration-300 hover:text-gold active:text-gold">关于</a>
</nav>
```

---

### 卡片标题

**CSS**：
```css
.card-title {
    font-family: 'Cinzel', serif;
    color: #D4AF37;
    font-size: 1.1rem;
    letter-spacing: 2px;
}
```

**Tailwind**：
```html
<div class="card-title font-cinzel text-gold tracking-wider">
  标题文本
</div>
```

---

### 分割线

**CSS**：
```css
.divider {
    border-top: 1px solid rgba(212, 175, 55, 0.2);
}
```

**Tailwind**：
```html
<div class="border-t border-gold/20"></div>
```

---

### 星星图标

**CSS**：
```css
.star-icon {
    color: #D4AF37;
    animation: twinkle 3s ease-in-out infinite;
}
```

**Tailwind**：
```html
<span class="text-gold animate-twinkle">★</span>
```

---

### 日期标签

**CSS**：
```css
.date-label {
    font-size: 12px;
    color: #D4AF37;
    text-transform: uppercase;
    letter-spacing: 3px;
    opacity: 0.8;
    position: relative;
    display: inline-block;
}

.date-label::before,
.date-label::after {
    content: '';
    position: absolute;
    top: 50%;
    width: 30px;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--c-accent));
}

.date-label::before {
    right: 100%;
    margin-right: 15px;
}

.date-label::after {
    left: 100%;
    margin-left: 15px;
    transform: scaleX(-1);
}
```

**Tailwind**：
```html
<div class="relative inline-block text-xs text-gold uppercase tracking-[3px] opacity-80">
  <span class="before:absolute before:right-full before:top-1/2 before:h-px before:w-[30px] before:mr-[15px] before:bg-gradient-to-r before:from-transparent before:to-gold after:absolute after:left-full after:top-1/2 after:h-px after:w-[30px] after:ml-[15px] after:bg-gradient-to-r after:from-transparent after:to-gold after:scale-x-[-1]">
    Dec 22 — Jan 19
  </span>
</div>
```

---

## 配色测试

在提交 UI 设计前，请进行以下测试：

### 1. 对比度测试
- 使用 Chrome DevTools 检查文本对比度
- 确保所有文本达到 WCAG AA 标准（4.5:1）
- 优先检查小号文本和低对比度区域

### 2. 响应式测试
- 在桌面端（1920x1080）、平板（768x1024）、手机（375x667）测试
- 确保配色在不同尺寸下保持一致性

### 3. 暗色模式测试
- 确保在深色背景下的可读性
- 测试不同亮度环境下的显示效果

### 4. 动画性能
- 检查动画是否流畅（60fps）
- 使用 Chrome DevTools Performance 分析
- 确保动画不阻塞主线程

### 5. 浏览器兼容性
- 在 Chrome、Firefox、Safari、Edge 中验证
- 测试不同浏览器版本的渲染效果
- 特别检查渐变和动画的兼容性

### 6. Tailwind 配置验证
- 确保所有自定义颜色都已在 `tailwind.config.js` 中定义
- 检查动画 keyframes 是否正确配置
- 验证字体配置是否生效

### 7. 可访问性测试
- 使用屏幕阅读器测试
- 确保颜色不是唯一的区分方式
- 测试高对比度模式

---

## 设计资源

### 参考文件

- **UI 示例文件**：`docs/ui-preview/capricorn-ui.html` - 完整的配色方案展示
- **Tailwind 配置**：`tailwind.config.js` 或 `tailwind.config.ts` - 项目 Tailwind 配置文件

### 在线工具

- **颜色对比度检查器**：https://webaim.org/resources/contrastchecker/
- **Tailwind CSS 文档**：https://tailwindcss.com/docs
- **Google Fonts**：https://fonts.google.com/

### 设计灵感

- **星座主题**：参考真实星座的视觉表现
- **奢华品牌**：学习高端品牌的配色运用
- **深色界面**：研究成功的深色模式设计案例

---

## 版本历史

- **v1.0.0** (2026-02-01)
  - 初始版本
  - 定义核心配色方案
  - 提供完整的组件规范和示例
  - 支持 CSS 和 Tailwind CSS 两种实现方式

---

## 贡献指南

如需修改配色规范，请遵循以下流程：

1. 提交 PR 到文档仓库
2. 说明修改原因和影响范围
3. 提供视觉示例
4. 等待团队审核和批准

---

**文档维护**：设计团队
**最后更新**：2026-02-01
**适用版本**：v1.0.0
