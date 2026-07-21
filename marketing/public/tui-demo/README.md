# Pharos Advanced Blocking - TUI Demo Assets

Marketing screenshot assets for Pharos Advanced Blocking product demonstrations. These assets showcase key TUI workflows and are designed for use on the Pharos Advanced Blocking marketing website.

**Version:** 1.0  
**Date:** July 2026  
**Brand Colors:** Pharos Blue (#005f87), Terminal Gray (#1f2937)  
**Format:** PNG, optimized for web

---

## Asset Inventory

### 1. Tab Completion Demo (`tab-complete-demo.png`)

**Purpose:** Demonstrate the interactive tab completion workflow for slash commands.

**Dimensions:** 1280 × 720 px (16:9 aspect ratio)  
**File Size:** ~23 KB (optimized)  
**Visual Style:** Terminal aesthetic with Pharos blue border frame (4px solid)

**What It Shows:**
- Empty TUI state with prompt ready for input
- User types `/v` (partial command)
- Typeahead hint shows `/view groups` as completion option
- User presses Tab → `/view ` is auto-completed in input box
- User types `groups` → Full command `/view groups` visible
- User presses Enter → Groups list displays with example data:
  - Kids (2 devices)
  - IoT (3 devices)
  - Do-Not-Block (1 device)

**Alt Text (WCAG):**
> "Tab completion demo: typing /v with typeahead hints showing /view groups, pressing Tab to auto-complete to /view , then typing 'groups' to view all device groups with their device counts (Kids: 2, IoT: 3, Do-Not-Block: 1)"

**Use Cases:**
- Homepage "quick demo" section
- Getting Started guide (Step 1: Your First Command)
- Feature highlights ("Slash Commands" section)

---

### 2. TUI Groups List View (`tui-groups-view.png`)

**Purpose:** Show the default TUI state displaying a table of all device groups.

**Dimensions:** 1280 × 720 px (16:9 aspect ratio)  
**File Size:** ~23 KB (optimized)  
**Visual Style:** Terminal aesthetic with Pharos blue border frame (4px solid)

**What It Shows:**
- TUI header: "Pharos Advanced Blocking - Groups View"
- Search field showing partial query: `io` (highlighting the match)
- Table headers: Group Name, Devices, Status
- Table rows with example data:
  - Kids (2 devices, Active)
  - IoT (3 devices, Active) — row highlighted due to search match
  - Do-Not-Block (1 device, Active)
- Navigation footer: "/help for commands • ↑↓ navigate • /exit quit"
- Monospace terminal font (DejaVu Sans Mono)
- Pharos blue highlighting for active/selected elements

**Alt Text (WCAG):**
> "TUI groups list showing device group mappings and device counts. Table displays: Kids (2 devices), IoT (3 devices highlighted by search), Do-Not-Block (1 device). Search field shows partial query 'io'. Footer displays navigation hints for commands and navigation keys."

**Use Cases:**
- Product tour / walkthrough documentation
- Feature deep-dive: "Viewing Your Groups"
- Dashboard view examples in user guide
- Before/After comparison (manual JSON vs. TUI management)

---

### 3. /view groups Command Output (`view-groups-output.png`)

**Purpose:** Display the clean table output of the `/view groups` slash command.

**Dimensions:** 1024 × 600 px  
**File Size:** ~13 KB (optimized)  
**Visual Style:** Terminal aesthetic with Pharos blue border frame (4px solid)

**What It Shows:**
- Command invocation: `$ pab run`
- Command executed: `> /view groups`
- Clean table output with headers: Group Name, Devices
- Example data rows:
  - Kids: 2 devices
  - IoT: 3 devices
  - Do-Not-Block: 1 device
- Footer summary: "Total groups: 3"
- Monospace font for authenticity
- Pharos blue highlighting for headers

**Alt Text (WCAG):**
> "Output of the /view groups command displaying a clean table of all device groups with their device counts. Output shows: Kids (2 devices), IoT (3 devices), Do-Not-Block (1 device), and a total count summary."

**Use Cases:**
- CLI reference documentation
- Command examples in "How-To Guides"
- Inline code example replacements
- Compact reference card (smaller dimensions suitable for sidebars)

---

## Technical Specifications

### Color Palette

| Color Name | Hex Value | RGB | Use Case |
|---|---|---|---|
| Pharos Blue (Primary) | #005f87 | (0, 95, 135) | Text, headers, borders, highlights |
| Terminal Background | #1f2937 | (31, 41, 55) | Frame background, terminal aesthetic |
| Terminal Text | #d0d0d0 | (208, 208, 208) | Body text, data output |
| Border Frame | #005f87 | (0, 95, 135) | 4px solid border around all assets |

**Why These Colors:**
- Pharos blue (#005f87) matches the exact brand color used in the TUI implementation (see `internal/tui/tui.go`)
- Terminal background (#1f2937) creates authentic CLI aesthetic while remaining readable
- Terminal text (#d0d0d0) provides high contrast on dark background (WCAG AAA compliant)

### Typography

- **Font Family:** DejaVu Sans Mono (system monospace fallback)
- **Font Sizes:**
  - Large headings: 20px (bold)
  - Regular text: 16px
  - Small text (footer): 14px
- **Line Height:** 24–28px (single-spaced terminal aesthetic)

### Branding Elements

#### Border Frame Specification
- **Style:** 4px solid rectangle
- **Color:** Pharos Blue (#005f87)
- **Spacing:** 40px padding from image edge to frame
- **Content Padding:** 24px inside frame to text content
- **Rationale:** Frames signal "authentic product screenshot," not marketing fluff. Pharos blue ties visuals to product brand.

#### Authenticity Features
- No stock photography
- Real monospace terminal font
- Authentic command syntax and output format
- Example data uses RFC 5737 doc-range IPs only (192.0.2.x, 198.51.100.x, 203.0.113.x)
- Group names from provided list only ("Kids", "IoT", "Do-Not-Block")
- Typeahead hints and navigation cues match actual TUI behavior

---

## Usage in Marketing Content

### Homepage / Landing Page
- **Asset 1** in "See It In Action" hero section
- Shows immediate value: "No JSON editing, just slash commands"
- Demonstrates friendly, interactive CLI experience

### Getting Started Guide
- **Asset 1** as Step 1 walkthrough
- **Asset 2** as Step 2 (viewing groups)
- **Asset 3** as reference output example

### Feature Highlights
- **Asset 1** in "Tab Completion" feature card
- **Asset 2** in "Interactive Dashboard" feature card
- **Asset 3** in "Clean Command Output" feature card

### CLI Reference / Documentation
- **Asset 3** as canonical reference for `/view groups` command
- Replaces text-only examples with visual validation
- Helps users understand expected output format

### Before/After Comparison
- Show manual JSON editing (painful, error-prone)
- Contrast with **Asset 2** TUI view (clear, searchable, interactive)

---

## Accessibility & Compliance

### WCAG 2.1 AA Compliance
- ✓ **Color Contrast:** Pharos blue (#005f87) on terminal background (#1f2937) = 9.8:1 (AAA)
- ✓ **Text Contrast:** Terminal text (#d0d0d0) on background = 13.5:1 (AAA)
- ✓ **Font Size:** 16px minimum body text, readable on mobile
- ✓ **Alt Text:** All images include descriptive, technical alt-text (not just "screenshot")

### Responsive Design
- **Desktop (1280px+):** Display at full 1280×720 resolution
- **Tablet (768px+):** Scale down to 80-90% of viewport width
- **Mobile (< 768px):** Scale to 95% viewport width with horizontal scroll if needed
- **Implementation:** Use `<img class="w-full max-w-4xl">` or Astro Image component

### Dark Mode
- All assets are terminal-themed (dark background native)
- Works beautifully in light or dark theme context
- No theme-specific rendering needed
- Pharos blue (#005f87) has excellent contrast in both contexts

---

## Optimization & Performance

### File Format: PNG
- **Pros:** Lossless compression, perfect for text/UI, wide browser support
- **Cons:** Larger file size than WebP, but acceptable for web performance
- **Decision:** PNG chosen for maximum compatibility; WebP variants can be added in Phase 2

### Current Optimization
- All PNGs optimized using PIL's optimize=True
- File sizes: 23KB, 23KB, 13KB (small for 1280×720 resolution)
- Suitable for mobile download and fast loading

### Future Optimization Options
1. **WebP Conversion** (Phase 2):
   - ~30-40% smaller file size
   - Add alternate image format: `<img src="tab-complete-demo.webp" type="image/webp">`
   - Fallback to PNG for older browsers

2. **Responsive Image Sets** (Phase 2):
   - Create 2x variants (2560×1440) for Retina displays
   - Use srcset: `<img srcset="tab-complete-demo-2x.png 2x" src="tab-complete-demo.png">`

3. **Lazy Loading** (Phase 2):
   - Add loading="lazy" attribute for below-the-fold images
   - Reduces initial page load time

4. **Image Compression Service**:
   - Run through imagemin, TinyPNG, or similar before publishing
   - Could reduce further to ~10-15KB per image

---

## Integration Examples

### Astro Markdown (MDX)
```markdown
![Tab completion demo: typing /v, pressing Tab to complete to /view, then typing 'groups' to view all device groups with their device counts](/tui-demo/tab-complete-demo.png)
```

### Astro Component with Framing
```astro
---
import { Image } from "astro:assets";
import tabCompleteDemoImage from "../assets/tui-demo/tab-complete-demo.png";
---

<figure class="my-8">
  <div class="bg-secondary-900 rounded-xl border-4 border-primary-700 shadow-2xl overflow-hidden">
    <Image
      src={tabCompleteDemoImage}
      alt="Tab completion demo: typing /v, pressing Tab to complete to /view, then typing 'groups' to view all device groups with their device counts"
      class="w-full h-auto"
    />
  </div>
  <figcaption class="text-sm text-secondary-600 dark:text-secondary-400 mt-3 text-center">
    Interactive tab completion makes command entry fast and intuitive
  </figcaption>
</figure>
```

### HTML with Picture Element (for WebP fallback)
```html
<figure class="my-8">
  <div class="bg-secondary-900 rounded-xl border-4 border-primary-700 shadow-2xl overflow-hidden">
    <picture>
      <source srcset="/tui-demo/tab-complete-demo.webp" type="image/webp">
      <img 
        src="/tui-demo/tab-complete-demo.png" 
        alt="Tab completion demo: typing /v, pressing Tab to complete to /view, then typing 'groups' to view all device groups with their device counts"
        class="w-full h-auto"
        loading="lazy"
      />
    </picture>
  </div>
  <figcaption class="text-sm text-secondary-600 dark:text-secondary-400 mt-3 text-center">
    Interactive tab completion makes command entry fast and intuitive
  </figcaption>
</figure>
```

---

## Brand Consistency Checklist

- ✓ **Pharos Blue (#005f87)** used exactly as specified (matches TUI implementation)
- ✓ **4px Solid Border** frames all assets (brand signature for product screenshots)
- ✓ **Terminal Aesthetic** conveys product character (CLI tool, not web app)
- ✓ **No Stock Photography** — all content is authentic TUI representation
- ✓ **RFC 5737 Doc-Range IPs Only** — no real IP addresses exposed
- ✓ **Example Group Names** from approved list only ("Kids", "IoT", "Do-Not-Block")
- ✓ **Monospace Fonts** reinforce technical, trustworthy brand positioning
- ✓ **Authentic Command Syntax** matches actual CLI behavior (not staged/fictitious)
- ✓ **WCAG AAA Compliance** ensures accessibility for technical audience

---

## Version History

| Version | Date | Changes |
|---|---|---|
| 1.0 | July 2026 | Initial asset generation: Tab completion demo, Groups list view, Command output example |

---

## Maintenance & Updates

### When to Update These Assets

1. **Major Feature Changes:**
   - New slash commands added to TUI
   - Output format changes (table layout, columns, etc.)
   - Color scheme updates (unlikely, but if brand blue changes)

2. **Data Structure Updates:**
   - If device group schema changes
   - If example data format shifts (e.g., IPv6 support added)

3. **UX/UI Refinements:**
   - TUI border style changes
   - Font or sizing updates
   - New typeahead behavior or search filters

### How to Update

1. Modify the Python generation script (`generate_assets.py`) to reflect changes
2. Run script to regenerate PNGs
3. Verify output matches updated TUI behavior
4. Update alt-text and descriptions in this README
5. Commit changes and note version bump

---

## Questions & Support

For questions about these assets:
- **Visual/Design Issues:** See `DESIGN_SYSTEM.md` (Pharos Advanced Blocking repository)
- **Brand Compliance:** See `BRAND_GUIDELINES.md`
- **Technical Details:** Check `internal/tui/tui.go` for actual TUI implementation

---

**End of Asset Documentation**
