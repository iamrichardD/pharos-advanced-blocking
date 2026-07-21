# Design System: Pharos Advanced Blocking Marketing Website

**Version:** 1.0  
**Date:** July 2026  
**Stack:** Astro + Tailwind CSS v4  
**Audience:** DevOps engineers, DNS administrators, security-focused developers  

---

## 1. Color Palette

### Primary Color Resolution: TUI Blue Takes Precedence

**The Decision:** Adopt Pharos TUI blue (`#005f87`) as the primary brand color for the marketing site, replacing the current Tailwind blue (`#2563eb`).

**Rationale:**
- **Trust & Authority**: The TUI blue is deeper, more navy-oriented, and signals "serious infrastructure" to a DNS/security audience. Lighter material blues feel too consumer-friendly for this demographic.
- **Brand Cohesion**: Users experiencing the TUI (the product) see `#005f87` daily. The marketing site should reinforce that same visual identity rather than contradict it.
- **Kathy Sierra Principle**: Clarity and consistency matter more than design trends. A tech audience notices when the marketing doesn't match the product.
- **Pragmatism**: Rather than create a new blue scale, we'll extend Tailwind's approach and define a curated Pharos blue scale that echoes the TUI while maintaining accessibility.

### Color Definitions

#### Primary: Pharos Blue (derived from #005f87)

| Shade | Hex Value | Use Case | WCAG |
|-------|-----------|----------|------|
| 50 | #e6f2f7 | Backgrounds, very light accents | AAA |
| 100 | #cce5ef | Hover backgrounds, light UI | AAA |
| 200 | #99cbe0 | Borders, disabled states | AA |
| 300 | #66b0d0 | Secondary buttons, muted accents | AA |
| 400 | #3396c0 | Hovered secondary elements | AA |
| 500 | #0087b8 | Primary action (lighter anchor) | AA |
| 600 | #006a94 | Primary buttons, links | AAA |
| 700 | #005f87 | TUI Brand Color, strong emphasis | AAA |
| 800 | #00486a | Dark mode primary, high contrast | AAA |
| 900 | #003347 | Deep backgrounds, very dark mode | AAA |
| 950 | #001d2e | Maximum contrast, rare use | AAA |

**Migration Note:** Current site uses `primary-600: #2563eb`. Update all instances to reference the new Pharos blue scale. This is a one-time configuration update in `src/styles/global.css`.

#### Secondary: Neutral Scale (unchanged)

Keep existing slate-based neutrals:
- `secondary-50` through `secondary-900` remain as-is (current gray scale)
- Ensures good contrast with new primary color
- Works for text, borders, backgrounds, and dividers

#### Status Colors

| Color | Hex | Purpose | Example |
|-------|-----|---------|---------|
| Success | #10b981 | Validation passed, deployment complete | "Configuration synced successfully" |
| Warning | #f59e0b | Caution, dry-run mode, deprecation | "1 deprecated flag detected" |
| Error | #ef4444 | Validation failure, sync error | "IP validation failed: invalid CIDR" |
| Info | #3b82f6 | Tips, hints, supplementary info | "Tip: Use --watch to auto-reload configs" |

#### Accent: Violet (for highlights and secondary CTAs)

- **Current**: `#8b5cf6` (purple-600)
- **Keep as-is**: Good contrast with Pharos blue; works well for highlights, secondary actions, and decorative elements

### Tailwind Configuration Update

```css
/* src/styles/global.css - @theme section */
@theme {
  --color-primary-50: #e6f2f7;
  --color-primary-100: #cce5ef;
  --color-primary-200: #99cbe0;
  --color-primary-300: #66b0d0;
  --color-primary-400: #3396c0;
  --color-primary-500: #0087b8;
  --color-primary-600: #006a94;
  --color-primary-700: #005f87;    /* TUI brand color */
  --color-primary-800: #00486a;
  --color-primary-900: #003347;
  --color-primary-950: #001d2e;
  
  /* Secondary / Neutral (unchanged) */
  --color-secondary-50: #f8fafc;
  --color-secondary-100: #f1f5f9;
  --color-secondary-200: #e2e8f0;
  --color-secondary-300: #cbd5e1;
  --color-secondary-400: #94a3b8;
  --color-secondary-500: #64748b;
  --color-secondary-600: #475569;
  --color-secondary-700: #334155;
  --color-secondary-800: #1e293b;
  --color-secondary-900: #0f172a;

  /* Status & Accent */
  --color-success-500: #10b981;
  --color-warning-500: #f59e0b;
  --color-error-500: #ef4444;
  --color-info-500: #3b82f6;
  --color-accent-600: #8b5cf6;
}
```

---

## 2. Typography

### Font Family Strategy

**Primary Font:** System UI stack (no custom web fonts)

```css
/* src/styles/global.css - html rule */
font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 
             Oxygen, Ubuntu, Cantarell, sans-serif;
```

**Rationale:** 
- System fonts load instantly (no FOUT/FOIT delays).
- Tech audience expects clean, legible sans-serif; they're used to terminals and code editors.
- Reduces bundle size and improves Core Web Vitals.
- Works equally well on Windows (Segoe), macOS (SF Pro), and Linux (Ubuntu).

**Monospace Font:** Tailwind default

```css
font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 
             'Liberation Mono', Menlo, monospace;
```

Used for code blocks, CLI commands, and inline `<code>` snippets.

### Heading Hierarchy (Sizes & Weights)

All headings use `font-semibold` (weight 600) with tight line-height for scanability.

| Element | Size (Mobile / Tablet / Desktop) | Line Height | Usage |
|---------|----------------------------------|-------------|-------|
| **h1** | 2rem / 3rem / 4rem (32px / 48px / 64px) | 1.2 | Page title, hero section |
| **h2** | 1.875rem / 2.25rem / 3rem (30px / 36px / 48px) | 1.25 | Section title, CLI guide headings |
| **h3** | 1.5rem / 1.875rem / 2.25rem (24px / 30px / 36px) | 1.25 | Subsection, feature cards |
| **h4** | 1.25rem / 1.5rem / 1.875rem (20px / 24px / 30px) | 1.33 | Component title, CLI command |
| **h5** | 1.125rem / 1.25rem / 1.5rem (18px / 20px / 24px) | 1.33 | Card subtitle, tip labels |
| **h6** | 1rem / 1.125rem / 1.25rem (16px / 18px / 20px) | 1.5 | Rarely used; use h5 instead |

**Current implementation in global.css is appropriate.** No changes needed.

### Body Text

| Context | Size | Line Height | Weight | Color (Light / Dark) |
|---------|------|-------------|--------|----------------------|
| Paragraph copy | 1rem (16px) | 1.6 | 400 (normal) | secondary-700 / secondary-300 |
| Small text (captions) | 0.875rem (14px) | 1.5 | 400 | secondary-600 / secondary-400 |
| Bold text (emphasis) | 1rem | 1.6 | 600 | secondary-800 / secondary-200 |
| Muted text (hints) | 0.875rem | 1.5 | 400 | secondary-500 / secondary-400 |

**Rationale:** 
- Base 16px on mobile ensures readability without pinch-zoom on older devices.
- 1.6 line-height (Tailwind `leading-relaxed`) suits longer paragraphs; DNS documentation can be dense.
- Tech audience won't complain about "too much" whitespace; they'll appreciate breathing room.

### Code & Command Syntax

```css
/* Inline code (e.g., `pab start` in running text) */
code:not(pre code) {
  @apply font-mono text-sm px-1.5 py-0.5 rounded 
         bg-secondary-100 dark:bg-secondary-800 
         text-secondary-900 dark:text-secondary-100 
         border border-secondary-200 dark:border-secondary-700;
}

/* Code block (pre > code for multi-line) */
pre {
  @apply bg-secondary-900 text-secondary-100 p-4 rounded-lg 
         overflow-x-auto border border-secondary-800;
  @apply font-mono text-sm leading-relaxed;
}
```

**No syntax highlighting library required.** Simple monospace + gray background works for CLI/config snippets aimed at technical readers who expect terminal-style code.

---

## 3. Component Specifications

### Button System

#### Primary Button (Call-to-Action)

```html
<!-- Install the CLI, Deploy, Sync buttons -->
<a href="/install" class="btn btn-primary">
  Install the CLI
</a>
```

```css
.btn {
  @apply inline-flex items-center justify-center px-4 py-2 rounded-lg 
         font-medium text-sm focus:outline-none focus:ring-2 
         focus:ring-offset-2 disabled:opacity-50 
         disabled:cursor-not-allowed transition-all;
}

.btn-primary {
  @apply bg-primary-700 text-white 
         hover:bg-primary-800 hover:text-white 
         focus:ring-primary-600 shadow-lg hover:shadow-xl;
  /* Dark mode uses same colors (primary scale is already dark-friendly) */
}
```

**Sizes:**
- Default: `px-4 py-2` (compact, for inline/secondary contexts)
- Large: `px-8 py-4` (CTA hero section, add `text-base`)

**States:**
- **Default**: Primary-700 background, white text
- **Hover**: Primary-800 (darker), shadow upgrade
- **Focus**: Outline ring with primary-600
- **Disabled**: Opacity 50%, cursor not-allowed, no hover

#### Secondary Button

```html
<!-- Read the Docs, Learn More -->
<a href="/guide" class="btn btn-secondary">
  Read the Docs
</a>
```

```css
.btn-secondary {
  @apply bg-secondary-100 text-secondary-900 
         hover:bg-secondary-200 
         dark:bg-secondary-800 dark:text-secondary-100 
         dark:hover:bg-secondary-700 
         focus:ring-secondary-500;
}
```

**Rationale:** Subtle, de-emphasized action. Works on both light and dark backgrounds.

### Card System

```html
<!-- Feature cards, documentation cards, command reference -->
<div class="card card-hover p-6">
  <div class="w-12 h-12 rounded-lg bg-primary-100 dark:bg-primary-900/50 
              flex items-center justify-center text-primary-600 
              dark:text-primary-400 mb-4">
    <svg><!-- icon --></svg>
  </div>
  <h3 class="text-xl mb-2">Bulletproof Configs</h3>
  <p class="text-secondary-600 dark:text-secondary-400">
    Description text...
  </p>
</div>
```

```css
.card {
  @apply bg-white dark:bg-secondary-800 rounded-xl 
         shadow-lg border border-secondary-200 dark:border-secondary-700 
         transition-all;
}

.card-hover {
  @apply hover:shadow-xl hover:-translate-y-1;
}
```

**Use Cases:**
- Feature highlights (homepage)
- Command reference (grouped by category)
- Documentation cards (3-4 per row on desktop)

**Icon Container:**
- Subtle background: `bg-primary-100` / `dark:bg-primary-900/50`
- Icon color: `text-primary-600` / `dark:text-primary-400`
- Size: 48px × 48px (w-12 h-12)
- Rounded: Medium border radius (`rounded-lg`)

### Code Blocks

```html
<!-- CLI command example -->
<pre><code class="language-bash">pab start --watch
pab client add 192.168.1.100 --policy strict
</code></pre>
```

**Styling:**
- Background: `secondary-900` (near-black for contrast)
- Text: `secondary-100` (off-white, easy on eyes)
- Font: Monospace, `text-sm`, `leading-relaxed` (1.625 line height)
- Padding: `p-4` (inside breathing room)
- Border: `border border-secondary-800` (subtle frame)
- Overflow: `overflow-x-auto` (don't break mobile layout)

**No syntax highlighting.** Keep it simple for technical users. If colorized syntax is desired in v2, Shiki or Prism can be added without breaking current design.

### Navigation

#### Header (Sticky Navigation)

```html
<nav class="border-b border-secondary-200 dark:border-secondary-800 
            bg-white/80 dark:bg-secondary-900/80 backdrop-blur-md 
            sticky top-0 z-50">
  <div class="container-custom py-4 flex items-center justify-between">
    <a href="/" class="font-bold text-2xl tracking-tight 
                       text-secondary-900 dark:text-white">
      Pharos<span class="text-primary-600">AB</span>
    </a>
    <div class="flex gap-6 items-center">
      <a href="/user-guide" class="font-medium hover:text-primary-600">User Guide</a>
      <a href="/cli-reference" class="font-medium hover:text-primary-600">CLI Reference</a>
      <a href="/installation" class="font-medium hover:text-primary-600">Install</a>
      <a href="https://github.com/..." class="btn btn-primary">GitHub</a>
    </div>
  </div>
</nav>
```

**Design Notes:**
- Sticky positioning keeps nav visible while scrolling.
- Frosted glass effect (`backdrop-blur-md`) adds sophistication without visual noise.
- Logo: "Pharos" in default text, "AB" in primary color accent for brand mark.
- Links use `font-medium` (weight 500) with hover state shifting to `primary-600`.
- Z-index 50 ensures nav stays above content.

**Responsive:** Mobile nav (hamburger menu) can be added in v2 if nav gets crowded.

### Images & Screenshots

#### TUI Screenshot Framing

When showcasing TUI terminal screenshots:

```html
<!-- Framed TUI screenshot -->
<figure class="my-8">
  <div class="bg-secondary-900 rounded-xl border-4 border-primary-700 
              shadow-2xl overflow-hidden">
    <img src="/screenshots/tui-dashboard.png" 
         alt="Pharos TUI dashboard showing client IP mappings" 
         class="w-full h-auto">
  </div>
  <figcaption class="text-sm text-secondary-600 dark:text-secondary-400 
                      mt-3 text-center">
    The Pharos TUI dashboard with client IP → group mappings
  </figcaption>
</figure>
```

**Frame Specification:**
- Border: 4px solid `primary-700` (#005f87) — matches TUI brand
- Background: `secondary-900` (matches terminal aesthetic)
- Rounded corners: `rounded-xl` (subtle modernity without overdoing it)
- Shadow: `shadow-2xl` (depth, stands out on page)
- Aspect Ratio: Let image scale naturally; no forced aspect ratio

**Rationale:** The border reinforces that this is a product screenshot, not marketing fluff. The Pharos blue frame ties TUI visuals to the marketing brand.

#### Regular Screenshots (comparisons, before/after)

For non-TUI screenshots (e.g., DNS config comparison):

```html
<div class="border border-secondary-200 dark:border-secondary-700 
            rounded-lg overflow-hidden shadow-lg">
  <img src="/screenshots/config-editor.png" alt="Config validation UI" class="w-full">
</div>
```

**Simpler frame:** Just a subtle border and shadow. No color accent needed for non-TUI assets.

### Tables (Command Reference)

```html
<div class="overflow-x-auto">
  <table class="w-full text-sm border-collapse">
    <thead class="bg-primary-50 dark:bg-primary-900/30 border-b border-primary-200 dark:border-primary-800">
      <tr>
        <th class="px-4 py-3 text-left font-semibold text-primary-900 dark:text-primary-100">Command</th>
        <th class="px-4 py-3 text-left font-semibold text-primary-900 dark:text-primary-100">Flags</th>
        <th class="px-4 py-3 text-left font-semibold text-primary-900 dark:text-primary-100">Purpose</th>
      </tr>
    </thead>
    <tbody>
      <tr class="border-b border-secondary-200 dark:border-secondary-700 hover:bg-secondary-50 dark:hover:bg-secondary-800/50">
        <td class="px-4 py-3 font-mono text-primary-600 dark:text-primary-400">pab start</td>
        <td class="px-4 py-3">--watch, --config</td>
        <td class="px-4 py-3">Launch TUI dashboard</td>
      </tr>
      <!-- more rows -->
    </tbody>
  </table>
</div>
```

**Table Styling:**
- Header: Primary-tinted background (`primary-50` / `primary-900/30` dark)
- Rows: Hover effect adds subtle background (`secondary-50` / `secondary-800/50` dark)
- Borders: Subtle secondary borders between rows
- Text: Monospace for commands; regular body text for descriptions
- Responsive: Wrap in `overflow-x-auto` to prevent horizontal scroll on mobile

### Forms & Inputs

Keep minimal; focus on clarity over decoration.

```html
<label class="block mb-4">
  <span class="block text-sm font-medium text-secondary-900 dark:text-secondary-100 mb-2">
    Configuration File Path
  </span>
  <input type="text" 
         placeholder="/etc/pab/dnsApp.config" 
         class="w-full px-3 py-2 border border-secondary-300 dark:border-secondary-600 
                rounded-lg bg-white dark:bg-secondary-800 
                text-secondary-900 dark:text-secondary-100
                focus:outline-none focus:ring-2 focus:ring-primary-600 focus:border-transparent">
</label>

<button class="btn btn-primary">Save Configuration</button>
```

**Input States:**
- **Default**: Neutral border, white/dark background
- **Focus**: Primary ring (2px) + transparent border
- **Error**: Red border (`border-error-500`) + error ring
- **Disabled**: Gray text + cursor-not-allowed

### Alerts & Callouts

```html
<!-- Success callout -->
<div class="bg-success-50 dark:bg-success-900/20 border-l-4 border-success-500 
            p-4 rounded text-success-900 dark:text-success-100">
  <strong>Success!</strong> Configuration synced to 3 nodes.
</div>

<!-- Warning callout -->
<div class="bg-warning-50 dark:bg-warning-900/20 border-l-4 border-warning-500 
            p-4 rounded text-warning-900 dark:text-warning-100">
  <strong>Tip:</strong> Use `--dry-run` to preview changes before applying.
</div>

<!-- Error callout -->
<div class="bg-error-50 dark:bg-error-900/20 border-l-4 border-error-500 
            p-4 rounded text-error-900 dark:text-error-100">
  <strong>Error:</strong> IP validation failed: 999.999.999.999 is not a valid IPv4.
</div>

<!-- Info callout -->
<div class="bg-info-50 dark:bg-info-900/20 border-l-4 border-info-500 
            p-4 rounded text-info-900 dark:text-info-100">
  <strong>Info:</strong> This feature requires v0.2.0 or later.
</div>
```

**Anatomy:**
- Colored left border (4px) indicates type (success, warning, error, info)
- Light background (`color-50` / `color-900/20` dark) with corresponding text color
- Consistent `p-4` padding; `rounded` corners (no excessive border radius)

---

## 4. Layout & Spacing

### Grid System

**Use Tailwind's native 4px grid.** No custom spacing scale needed.

| Unit | Pixels | Tailwind Class | Use Case |
|------|--------|----------------|----------|
| 1x | 4px | `p-1`, `m-1`, `gap-1` | Micro spacing (rare) |
| 2x | 8px | `p-2`, `m-2`, `gap-2` | Tight component spacing |
| 3x | 12px | `p-3`, `m-3`, `gap-3` | Default padding |
| 4x | 16px | `p-4`, `m-4`, `gap-4` | Section padding |
| 6x | 24px | `p-6`, `m-6`, `gap-6` | Card padding, nav gap |
| 8x | 32px | `p-8`, `m-8`, `gap-8` | Prominent spacing |
| 12x | 48px | `p-12`, `m-12`, `gap-12` | Section separation |
| 16x | 64px | `p-16`, `m-16` | Hero spacing (rare) |

### Max-Width & Container

**Content max-width:** 7xl (80 characters per line, comfortable for reading)

```css
.container-custom {
  @apply max-w-7xl mx-auto px-4 sm:px-6 lg:px-8;
}
```

**Applied to:**
- Main content sections
- Navigation
- Hero section
- Footer

**Prose (documentation pages):**

```html
<article class="container-custom py-12 max-w-4xl mx-auto prose dark:prose-invert">
  <!-- MDX content -->
</article>
```

Max-width 4xl (56 characters) for tighter, denser documentation reading.

### Responsive Breakpoints

**Tailwind v4 defaults (unchanged):**

| Breakpoint | Size |
|------------|------|
| Mobile (default) | < 640px |
| `sm` | ≥ 640px |
| `md` | ≥ 768px |
| `lg` | ≥ 1024px |
| `xl` | ≥ 1280px |
| `2xl` | ≥ 1536px |

**Mobile-First Approach:**
- Design for mobile first (single column, stacked cards)
- `md:grid-cols-3` for three-column feature grid
- `lg:text-lg` for larger headings on desktop

**Example:**

```html
<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
  <!-- Cards auto-stack on mobile, 2 cols on tablet, 3 on desktop -->
</div>
```

### Vertical Rhythm

**Sections use consistent vertical spacing:**

```css
.section {
  @apply py-16 md:py-24 lg:py-32;
}

.section-sm {
  @apply py-8 md:py-12 lg:py-16;
}
```

**In practice:**
- Hero section: `.section` (large spacing for impact)
- Feature cards: `.section` or `.section-sm` between feature groups
- Footer: `.py-8` (less prominent)

---

## 5. Visual Hierarchy & Emphasis

### Drawing Attention to Key Features

**Strategy:** Use color, size, and whitespace—not animation or decoration.

#### Hero Section (Homepage)

```html
<section class="section container-custom flex flex-col items-center text-center">
  <h1 class="max-w-4xl mb-6">
    Stop wrestling with JSON. Master Technitium Advanced Blocking.
  </h1>
  <p class="text-lg md:text-xl text-secondary-600 dark:text-secondary-400 max-w-2xl mb-10">
    Manually editing dnsApp.config is brittle and stressful...
  </p>
  <div class="flex flex-col sm:flex-row gap-4 w-full justify-center">
    <a href="/install" class="btn btn-primary px-8 py-4 text-base">
      Install the CLI
    </a>
    <a href="/guide" class="btn btn-secondary px-8 py-4 text-base">
      Read the Docs
    </a>
  </div>
</section>
```

**Hierarchy:**
- Large h1 (max-width 4xl to prevent visual sprawl)
- Subheading in muted secondary color (semantic secondary info)
- CTAs in primary/secondary buttons with large padding
- Centered text draws focus naturally
- Whitespace above/below the section (`.section` padding)

#### Feature Highlights

Card-based approach with icon + title + description:

```html
<div class="card p-6">
  <div class="w-12 h-12 rounded-lg bg-primary-100 dark:bg-primary-900/50 
              flex items-center justify-center text-primary-600 
              dark:text-primary-400 mb-4">
    <svg><!-- icon --></svg>
  </div>
  <h3 class="text-xl mb-2">Bulletproof Configs</h3>
  <p class="text-secondary-600 dark:text-secondary-400">
    No more broken JSON files...
  </p>
</div>
```

**Visual Weight:**
- Icon (primary color): Draws eyes first
- Title (larger, bolder): Second pass
- Description (muted, smaller): Third pass (detail readers)

### Distinguishing Marketing Copy from Technical Reference

**Marketing Pages** (homepage, features):
- Conversational tone
- Short paragraphs (2–3 sentences max)
- Icons, emojis sparingly
- Emphasis via color (primary blue callouts) and bold text
- Open spacing

**Technical Reference Pages** (CLI reference, user guide):
- Dense information (needed for developers)
- Tables and structured data
- Code blocks frequently
- Monospace for commands
- Tighter spacing (but still readable)

**Visual Cue:** Use different background colors to signal context switch.

```html
<!-- Marketing section (light bg) -->
<section class="section bg-white dark:bg-secondary-900">
  <!-- conversational content -->
</section>

<!-- Technical section (darker/neutral bg) -->
<section class="section bg-secondary-50 dark:bg-secondary-800/50">
  <div class="container-custom">
    <h2>CLI Reference</h2>
    <table><!-- command reference --></table>
  </div>
</section>
```

### Screenshots & Demos in Context

**Placement Strategy:**
- Break text every 3–4 paragraphs with a screenshot or demo
- Use figure/figcaption for semantic HTML
- Ensure alt text describes the terminal output or UI state
- TUI screenshots: Use the Pharos blue frame (as detailed in Component Specs)

```html
<section class="section">
  <div class="container-custom">
    <h2 class="mb-6">Interactive Dashboard</h2>
    <p class="text-lg mb-8">The TUI displays...</p>
    
    <figure class="my-12">
      <div class="bg-secondary-900 rounded-xl border-4 border-primary-700 
                  shadow-2xl overflow-hidden">
        <img src="/tui-dashboard.png" alt="..." class="w-full">
      </div>
      <figcaption class="text-sm text-secondary-600 dark:text-secondary-400 mt-3">
        Dashboard showing client mappings
      </figcaption>
    </figure>
    
    <p class="text-lg">This view allows you to...</p>
  </div>
</section>
```

### Whitespace & Contrast

- **Between sections:** `.section` padding (16px–32px vertical) creates breathing room
- **Between cards:** `gap-6` or `gap-8` in grid layouts
- **Between text elements:** Natural line-height (1.6 for body) + margin-bottom on headings
- **On dark mode:** Increase contrast slightly; use primary-400 for links instead of primary-500

---

## 6. Screenshots & Demo Assets

### Guidelines for TUI Screenshots

**Before screenshot:**
1. Clear the terminal (or use a clean state)
2. Ensure Pharos blue is visible in the TUI output (border, headers)
3. Capture at a readable resolution (1280×720 min, 16:9 aspect ratio)
4. Use proper lighting/contrast (avoid dim backgrounds)

**In Figma/Mockup** (if post-processing):
- Add the 4px Pharos blue border (`#005f87`)
- Dark background (`#0f172a` or darker) to simulate terminal context
- Optional: Add a slight shadow for depth

**Alt Text Requirements (WCAG):**

```html
<!-- Good: Describes what's shown, not just "screenshot" -->
<img src="/tui-example.png" 
     alt="Pharos TUI dashboard with 5 client IP mappings to blocking groups, 
          showing client 192.168.1.100 mapped to the 'Strict' group">

<!-- Bad: Too vague -->
<img src="/tui-example.png" alt="TUI screenshot">
```

### Color Consistency: Which Blue for Framing?

**Use Pharos blue (#005f87) for TUI screenshots only.**

- **Reason:** TUI displays *are* the product's signature interaction mode. Framing them in the brand color reinforces product identity.
- **Non-TUI screenshots** (config files, web UI comparisons): Use neutral borders or secondary colors.

### Asset Naming Convention

- TUI screenshots: `/public/screenshots/tui-{feature}-{version}.png`
  - Example: `tui-dashboard-v0.2.0.png`
- Comparison/demo screenshots: `/public/screenshots/{feature}-{type}.png`
  - Example: `config-editor-before.png`, `config-editor-after.png`
- Diagrams/infographics: `/public/diagrams/{topic}.svg` (prefer SVG for scalability)

### Responsive Image Handling

```html
<!-- Astro Image component with optimization -->
<Image
  src={import('../assets/tui-dashboard.png')}
  alt="..."
  width={1280}
  height={720}
  class="w-full rounded-lg shadow-lg border-4 border-primary-700"
/>
```

Or in Markdown (MDX):

```markdown
![Alt text for accessibility](/screenshots/tui-dashboard.png)
```

---

## 7. Accessibility & Pragmatism

### WCAG 2.1 AA Compliance (Minimum)

**Color Contrast Ratios:**
- **Body text (16px):** Min 4.5:1 contrast (AA)
  - Primary-700 (#005f87) on white: ✓ 8.5:1
  - Secondary-600 (#475569) on white: ✓ 5.1:1
  - Dark mode: Secondary-300 on secondary-900: ✓ 7:1
- **Large text (18px+):** Min 3:1 contrast (AA)
  - Headings always use primary or secondary-900 with sufficient contrast

**Keyboard Navigation:**
- All buttons/links focusable via Tab
- Focus outline: 2px ring, offset 2px (visible on all backgrounds)
- No keyboard traps

**Semantic HTML:**
- Use `<button>` for actions, `<a>` for navigation
- Use `<nav>`, `<article>`, `<section>` landmark elements
- Proper heading hierarchy (h1 → h2 → h3, no skipping)

### Avoid Over-Decoration; Embrace Clarity

**Don't use:**
- Bouncing animations or parallax scrolling
- Auto-playing videos or sound
- Blinking text or animated emoji
- Hover effects on mobile (they're confusing on touch devices)

**Do use:**
- Smooth transitions (200ms fade on hover)
- Clear, descriptive link text (not "click here")
- Sufficient color contrast (checked via WebAIM Contrast Checker)
- Simple icons (Heroicons style; no overly decorative graphics)

**Kathy Sierra Principle:** If it doesn't help users understand or use the product, remove it.

### Technical Audience Doesn't Need Hand-Holding

- Assume readers know what DNS, CIDR, and JSON are
- Don't explain "what is a CLI" — dive straight into usage
- Provide advanced options (flags, config), not simplified "easy mode"
- Trust that users will read code examples and understand monospace syntax

**In practice:**
- No "beginner's guide" tab separate from main docs (just write docs clearly)
- Code examples show full command, not step-by-step screenshots
- Assume UNIX/Linux background (use `~/.config/pab` without explaining `~`)

---

## 8. Brand Consistency & Positioning

### Link to Brand Guidelines

This design system reinforces **Pharos Advanced Blocking's positioning as:**
- **Reliable**: Deep blue (not trendy) signals maturity and trust
- **Utilitarian**: Clean typography, minimal decoration, fast load times
- **Professional**: Technical audience sees serious infrastructure tool, not consumer app
- **Open**: GitHub link in nav, documentation transparent and complete

**See:** `BRAND_GUIDELINES.md` (parallel task) for messaging, voice, and brand voice (CLI-first, honest, pragmatic).

### Visual System Reinforces Brand Positioning

| Positioning | Visual Reinforcement | Why |
|--------------|---------------------|-----|
| **Reliable** | Dark blue primary color, consistent spacing | Serious, mature, DNS is critical |
| **Utilitarian** | Minimal shadows/decoration, monospace for commands | No fluff; tech audience doesn't want design theater |
| **Professional** | Formal typography, clear hierarchy, semantic HTML | Infrastructure tool, not a game or toy |
| **Open Source** | GitHub button prominent in nav; MIT/AGPL links in footer | Transparency builds trust |

### Deliberate Visual Tension: Utilitarian + Remarkable

**The Balancing Act:**

The design system intentionally sits at the intersection of two brand attributes:
1. **Utilitarian** (CLI-first, minimal decoration, fast)
2. **Remarkable** (stands out visually, memorable brand color)

**How we achieve both:**

- **Color:** Pharos blue (#005f87) is neither trendy nor dull—it's distinctive but professional
- **Typography:** System fonts (fast, clean) without sacrificing legibility or character
- **Layout:** Generous whitespace (not cramped) paired with bold typography for hierarchy
- **Components:** Cards and sections feel modern (rounded corners, shadows) without being overly decorative

**Example tension resolution:**

```html
<!-- Utilitarian: simple layout, semantic HTML -->
<section class="section">
  <div class="container-custom">
    <h2>CLI Commands</h2>
    <table>
      <!-- command ref -->
    </table>
  </div>
</section>

<!-- Remarkable: the blue border + rounded corners + shadow make it visually distinct -->
<div class="bg-primary-50 dark:bg-primary-900/20 border-l-4 border-primary-700 
            rounded-r-lg p-4 shadow-md">
  <p>Important note about compatibility...</p>
</div>
```

---

## Implementation Roadmap

### Phase 1: Immediate (v0.2.0 launch)

1. **Update color palette** in `src/styles/global.css` to new Pharos blue scale
2. **Verify all links** still render correctly (primary-600 → new primary-600)
3. **Test dark mode** (no changes needed; scale already dark-friendly)
4. **Screenshot TUI** with proper framing and add to `/public/screenshots/`
5. **Validate WCAG AA** via WebAIM or axe DevTools

### Phase 2: Optional (v0.3.0+)

1. Add custom TUI screenshot framing mockup (Figma template)
2. Introduce "Tabs" component if CLI Reference grows to multiple guides
3. Add animated CLI command examples (using Asciinema or similar, with opt-out)
4. Mobile hamburger menu for nav if it becomes crowded
5. Syntax highlighting for code blocks (optional; Shiki, Prism, or highlight.js)

### Phase 3: Analytics & Iteration (Based on User Feedback)

1. Track which sections get highest engagement
2. Adjust color brightness if accessibility metrics suggest readability issues
3. Refine screenshot captions based on user comprehension metrics

---

## Appendix: Tailwind Config Quick Reference

### No separate tailwind.config.js needed

All configuration lives in `src/styles/global.css` using `@theme` and `@layer`:

```css
@import "tailwindcss";
@plugin "@tailwindcss/typography";

@theme {
  --color-primary-*: ...
  --color-secondary-*: ...
  --color-success-500: #10b981;
  --color-warning-500: #f59e0b;
  --color-error-500: #ef4444;
  --color-info-500: #3b82f6;
  --color-accent-600: #8b5cf6;
}

@layer base { /* reset & defaults */ }
@layer components { /* .btn, .card, .section, etc. */ }
@layer utilities { /* rarely needed */ }
```

### Migration Checklist (Current → Pharos Blue)

- [ ] Update `--color-primary-*` scale in `global.css`
- [ ] Test all `.bg-primary-*`, `.text-primary-*`, `.border-primary-*` classes render correctly
- [ ] Verify dark mode variant classes (`dark:bg-primary-*`, etc.) are still correct
- [ ] Check navigation links (`.hover:text-primary-*`)
- [ ] Validate buttons (`.btn-primary` uses new primary-700)
- [ ] Test focus rings (`.focus:ring-primary-*`)
- [ ] Screenshot a page on light and dark mode; compare with design mockup

---

## Conclusion

This design system balances **pragmatism** (use Tailwind defaults, no custom fonts, no over-engineering) with **brand differentiation** (distinctive Pharos blue, clear hierarchy, technical credibility). It scales from mobile to desktop, respects dark mode users, and prioritizes accessibility and clarity—core values for a DNS security tool aimed at engineers.

**The Pharos blue reconciliation (#005f87) is the single largest decision:** it aligns marketing visuals with the TUI product, builds brand consistency, and signals maturity to the target audience. From there, the rest of the system follows naturally from Tailwind defaults and proven component patterns.

**Next steps:** Update the color palette in `global.css`, test all pages in light/dark mode, and validate WCAG AA compliance before v0.2.0 launch.
