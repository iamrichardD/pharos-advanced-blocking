# Pharos Advanced Blocking (pab) — Brand Guidelines

**Version:** 1.0  
**Last Updated:** July 2026  
**Audience:** Marketing, Content, Product teams

---

## 1. Brand Positioning

**Core Position:**  
Pharos Advanced Blocking is the pragmatic CLI for self-hosting admins who want bulletproof control over Technitium DNS blocking rules—no manual JSON editing, no brittle config crashes, no reinventing the wheel. It transforms DNS policy management from error-prone drudgery into a clean, validated workflow.

**Why it matters:**  
Technitium DNS is powerful but its JSON configuration (`dnsApp.config`) is a brittle choke point. Typos crash the service. Manual edits accumulate technical debt. Admins waste time on syntax validation instead of strategy. pab solves this by inserting a safety layer (schema validation, TUI workflows, dry-run preview) between the admin's intention and the DNS engine, letting them move fast without breaking things.

---

## 2. Tone & Voice

### How pab Sounds
- **Confident but pragmatic**: We know our audience is competent. No hand-holding, no fluff. Respect their time and expertise.
- **Clear over clever**: Technical jargon is fine (CIDR blocks, API tokens, JSON schema)—our audience speaks it. But always define what we're *doing*, not just what we're *using*.
- **Supportive, not patronizing**: Acknowledge the pain of manual JSON editing. Frame pab as the expert friend who's solved the same problem, not the genius telling you that you're doing it wrong.

### Emotional Tone
- **Relief**: "No more DNS crashes from typos." Admins spend mental energy worrying about configuration errors. pab removes that worry.
- **Control**: "Your devices, your rules." Empower admins to define exactly which hosts get which blocklists, without negotiating with UI wizards.
- **Confidence**: "Built by admins, for admins." Acknowledge that we share the same problems and constraints.

### Language Guidelines
- **Do embrace technical terms**: IPv4/IPv6, CIDR notation, API endpoints, JSON schema, validation. This is our audience's native language.
- **Do explain *why***: "Schema validation catches malformed client IPs before they reach Technitium" is better than just "validates schemas."
- **Avoid**: Marketing fluff ("cutting-edge," "revolutionary," "game-changing"), false empathy ("we know how hard DNS is"), or over-simplification ("just use our tool!").
- **Lean into**: Honesty about scope ("works with Technitium Advanced Blocking App"), pragmatism ("choose your own Git workflow"), and transparency ("open-source, inspect as you wish").

---

## 3. Key Messages

1. **Schema validation prevents crashes**: Malformed IPs, invalid CIDR blocks, and syntax errors are caught before they reach Technitium. No more failed deployments at 2 AM.

2. **Per-device blocking policies, without JSON hell**: Map specific client IPs or ranges to blocking groups via a clean TUI or CLI—no manual JSON array editing, no merge conflicts.

3. **Built by and for self-hosting admins**: We don't pretend to solve everyone's DNS problems. pab is laser-focused on admins who run Technitium in homelab or small business settings and want to stop writing JSON by hand.

4. **AI-agent friendly and automation-ready**: The `--json` flag and structured CLI make pab easy to integrate with scripts, ChatGPT-powered assistants, or existing automation workflows. The tool gets out of your way.

5. **Open-source you can audit**: Full transparency. No vendor lock-in. No surprise changes. Inspect the source, fork it, contribute back, or run your own version—you're in control.

---

## 4. Audience Targeting

### Who Is pab For
- **Self-hosting Technitium admins**: Home lab enthusiasts, small business IT teams, tiny networks that need DNS-level ad-blocking without paying for expensive managed DNS services.
- **Tech-savvy operators**: Comfortable with CLIs, APIs, and config files. Want to understand *why* things work, not just that they work.
- **Reliability-focused teams**: Admins who've experienced DNS outages from bad configuration. Willing to invest setup time to prevent future incidents.
- **GitOps practitioners**: Teams already using version control for infrastructure. Want to manage DNS policies via Git workflows, not web UI clicks.

### Who pab Is NOT For
- **Non-technical end-users**: pab requires CLI fluency. It's not a desktop app. If you're looking for a "set and forget" DNS ad-blocker, Technitium's web UI is probably fine.
- **Enterprise DNS engineers**: If you're managing Technitium at scale across 50+ nodes with SLA contracts, you likely need deeper tooling (multi-region failover, audit logging, RBAC). pab assumes small-to-medium deployments.
- **Casual configuration tweakers**: If you change your blocklist rules once a year, manual JSON editing is probably fine. pab shines when you're iterating frequently or managing many devices.

### What This Audience Cares About (In Order)
1. **Reliability**: DNS outages are painful. They want confidence that their configuration won't crash the service.
2. **Control**: They self-host specifically to avoid vendor lock-in. They want to understand and modify every rule.
3. **Ease of use**: They're okay with learning a CLI, but they don't want to learn fifteen different config formats.
4. **Automation**: Many run scripts or use AI assistants. They want tools that play nicely with automation stacks.
5. **Transparency**: They want to inspect the source code and know exactly what the tool is doing with their configuration and credentials.

### How They Currently Solve This Problem
- **Manual JSON editing**: Hand-editing `dnsApp.config`, validating with trial-and-error (deploy, if it breaks, rollback).
- **Not at all**: Some accept Technitium's default settings because the configuration burden feels too high.
- **Shell scripts + jq**: Technically fluent teams write custom scripts to manipulate JSON. pab saves them that engineering overhead.

---

## 5. Visual & Design Philosophy

### Design Principles
- **Utilitarian, not austere**: The visual identity should feel clean and modern, but never sterile. Terminal tools are powerful—design should reflect that power without being intimidating.
- **Accessibility first**: TUI must work in light and dark themes, with clear contrast. CLI output should be readable at any terminal width. Screenshots should show real workflows, not stylized renderings.
- **Density without clutter**: Admins appreciate information density (showing many rows of data without scrolling), but the information hierarchy must be clear. Use color and whitespace strategically.

### Color Philosophy
- **Primary: Pharos Blue (#005f87)**: Conveys trust, stability, and depth. Use for primary actions (buttons), focus states in the TUI, and brand accents in marketing.
- **Secondary: Clean grays and neutrals**: Let the data speak. Avoid color-coding that distracts from the core information.
- **Accent: Minimal use of red/yellow for warnings and validation errors**: Reserve color for semantic meaning (error, warning, success), not decoration.
- **Dark mode native**: The TUI runs in dark terminals. Ensure marketing assets work beautifully in both light and dark theme contexts.

### Typography
- **Modern sans-serif**: Use clean, readable fonts (e.g., Inter, SF Mono for code). Avoid serif fonts in marketing.
- **Clear hierarchy**: Headers should be noticeably larger. Code blocks should be monospaced and distinct. Labels should be subtle but scannable.
- **Readability over decoration**: No fancy fonts, no overlays that reduce contrast. This is a tool for reading and acting on information quickly.

### Visual Assets & Screenshots
- **Show the TUI in action**: Demonstrate the actual terminal interface. Real data (anonymized device IPs, blocklist group names) grounds the tool in reality.
- **Before/after workflows**: Show a manual JSON file (with syntax errors circled) next to a pab TUI workflow. The contrast sells the value proposition.
- **No stock photography**: Use authentic screenshots, diagrams, and terminal recordings. Admins distrust glossy marketing. Authenticity builds credibility.
- **Command-line examples**: Always show real CLI invocations with realistic flags and outputs. Make copy-paste easy.

---

## 6. Content Pillars

### Pillar 1: "Your DNS, Your Rules" (Ownership & Control)
**Story Arc**: From "I'm locked into Technitium's defaults" to "I have fine-grained control over every device."

Content themes:
- How to map specific client IPs to custom blocking groups.
- Per-device policies (e.g., "kids' tablets get strict filtering, work laptop unrestricted").
- Decoupled architecture: pab edits config, *you* decide when to Git commit and push.
- Open-source inspection: "You own the source. You own the rules."

**Content formats**: Quickstart guides, tutorial videos showing TUI workflows, success stories from homelab admins.

---

### Pillar 2: "Zero JSON Errors" (Reliability & Confidence)
**Story Arc**: From "I'm terrified my config will break DNS" to "My configuration is validated before it ever touches the server."

Content themes:
- How schema validation works and what errors it catches (invalid IPs, malformed CIDR blocks, missing groups).
- Dry-run deployments: "See what will change before it changes."
- Rollback safety: "One command undoes the last deployment."
- Real incidents: "Here's how pab would have prevented a 2 AM DNS outage."

**Content formats**: Technical deep-dives, error prevention guides, case studies, troubleshooting walkthroughs.

---

### Pillar 3: "From Chaos to Clarity" (Quick Wins & Power User Progression)
**Story Arc**: Progressive skill-building from "I just want to add one device" to "I'm managing 50+ devices across multiple groups."

Content themes:
- Five-minute setup: getting pab installed and synced with your Technitium instance.
- Baby steps: adding your first device mapping via the TUI.
- Leveling up: CIDR blocks, group templates, scripting with `--json` output.
- Automation: integrating pab into CI/CD pipelines or AI-assisted workflows.
- Performance: managing large configurations (1000+ device mappings) efficiently.

**Content formats**: Getting-started guide, video tutorials (5-10 min each), recipe collections, power-user tips.

---

### Pillar 4: "Built by Admins, For Admins" (Open-Source Trust & Community)
**Story Arc**: From "Can I trust this tool?" to "I helped build this tool."

Content themes:
- Why pab exists: the founder's own frustration with JSON editing.
- Open-source first: every line of code is visible, auditable, and owned by the community.
- How to contribute: reporting issues, submitting PRs, building plugins.
- No vendor surprises: development roadmap is public, decisions are transparent.
- Single-purpose tool philosophy: pab does one thing well (manage Technitium blocking configs), not 100 things poorly.

**Content formats**: Founding story, contributing guide, plugin development guide, changelog deep-dives, community spotlights.

---

## 7. Messaging Guardrails

### What We Will Say
- "Bulletproof configurations"
- "Schema validation catches errors before deployment"
- "Per-device blocking policies without manual JSON editing"
- "Built for Technitium Advanced Blocking App"
- "AI-agent friendly and automation-ready"
- "Open-source, fully auditable"

### What We Will NOT Say
- "The only DNS manager you'll ever need" (too broad; we're not a general DNS tool)
- "Suitable for enterprises" (we're focused on small-to-medium self-hosted deployments)
- "Works with any DNS server" (only Technitium)
- "No technical knowledge required" (we target tech-savvy admins)
- "One-click setup" (setup involves API tokens, config files, understanding CIDR notation)
- "As good as paid solutions" (we don't compare; we're an open-source alternative, not a competitor)

---

## 8. Implementation Checklist

Use these guidelines to audit marketing content, documentation, and community communications:

- [ ] **Tone Check**: Does this sound confident but pragmatic? Free of marketing fluff?
- [ ] **Accuracy Check**: Does it describe pab's actual capabilities, not aspirational features?
- [ ] **Audience Check**: Would a self-hosting Technitium admin find this useful and relevant?
- [ ] **Scope Check**: Have we been honest about who pab is and is not for?
- [ ] **Visual Check**: Do screenshots show real workflows? Are colors and typography consistent with the brand?
- [ ] **Content Pillar Check**: Does this content fall into one of the four pillars, or is it a distraction?
- [ ] **Reciprocal Check**: Does this build trust with open-source users and potential contributors?

---

## 9. Maintaining Brand Consistency

- **Documentation > Marketing**: If the docs contradict marketing claims, the docs are right. Fix the marketing.
- **Feature Parity**: Only market features that are released and stable in the current version. Road-map features can be mentioned in the context of "coming soon," not "pab does X."
- **Community Input**: Major messaging changes should be discussed with contributors and users via GitHub discussions or issues.
- **Version Sensitivity**: If pab changes significantly (e.g., drops Technitium v10 support), update marketing materials and bump the guidelines version.

---

**End of Brand Guidelines**

*For questions or suggestions, open an issue on GitHub or start a discussion in the Pharos community.*
