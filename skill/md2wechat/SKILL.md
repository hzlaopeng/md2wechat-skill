---
name: md2wechat
description: Convert Markdown articles to WeChat Official Account formatted HTML using AI. Default mode: AI (no API key needed). Default theme: ocean-calm (professional blue). Supports AI image generation with DashScope (z-image-turbo), auto-saves to ./output/images/. Also supports assisted writing with customizable creator styles (default: Dan Koe). Use for writing articles, converting markdown to WeChat format, or generating cover images.
---

# MD to WeChat

Converts Markdown articles to WeChat Official Account formatted HTML with inline CSS using AI-powered themed layouts.

**✅ Pre-configured**: AI mode enabled, ocean-calm theme, DashScope image generation ready.

## Quick Start

```bash
# Preview HTML (uses ocean-calm theme by default)
bash skill/md2wechat/scripts/run.sh convert article.md --preview

# Use different theme (autumn-warm / spring-fresh / ocean-calm)
bash skill/md2wechat/scripts/run.sh convert article.md --theme autumn-warm --preview

# Generate cover image (auto-saves to ./output/images/)
bash skill/md2wechat/scripts/run.sh generate_image "Minimalist book cover, ocean blue and gold"

# Generate 16:9 cover for WeChat (recommended)
bash skill/md2wechat/scripts/run.sh generate_image --size 2560x1440 "prompt"
```

### Natural Language Image Generation

You can also ask me to generate images using natural language:

#### Generate Image for Article (Insert into Markdown)

```
"Help me generate a product concept image at the beginning of article.md"
"Add an image showing the product features after the second paragraph"
"Create a diagram for the comparison section in article.md"
```

I will:
1. Read the article to understand the context
2. Insert the AI image generation syntax at the appropriate location
3. Call the conversion command to generate and upload the image

#### Generate Standalone Image (Not for Article)

```
"Generate an image of a cute cat sitting on a windowsill"
"Create a product concept image: modern smart home device, white design"
"Make a diagram showing the user flow"
```

I will:
1. Call the image generation command directly
2. Return the generated image URL and WeChat material ID

### Natural Language Writing Assistance

You can also ask me to help write articles using creator styles:

#### Write Article from Idea

```
"Write an article about self-discipline using Dan Koe style"
"Help me write a post about productivity with a sharp, grounded tone"
"Create a story-style article about my travel experience"
```

I will:
1. Understand your idea or topic
2. Use the appropriate writing style (default: Dan Koe)
3. Generate article structure and content
4. Extract memorable quotes
5. Optionally generate matching cover image

#### Refine Existing Content

```
"Rewrite this article with a more engaging style"
"Polish my article.md with Dan Koe's writing style"
"Make this content more profound and sharp"
```

I will:
1. Read your existing content
2. Apply the selected writing style
3. Maintain original meaning while improving expression

#### Generate Cover Only

```
"Generate a cover image for my article about self-discipline"
"Create a Victorian woodcut style cover for my philosophy piece"
```

#### List Available Styles

```
"Show me all available writing styles"
"What writing styles can I use?"
```

**Available Writing Styles:**
- **Dan Koe** (default): Profound, sharp, grounded - great for personal growth and opinion pieces

Users can add custom styles in `writers/` directory. See `writers/README.md` for details.

## Workflow Checklist

```
Progress:
- [ ] Step 1: Analyze Markdown structure
- [ ] Step 2: Select theme (default: ocean-calm)
- [ ] Step 3: Generate HTML with AI (returns prompt for Claude)
- [ ] Step 4: Process images (AI生成→本地保存→上传微信)
- [ ] Step 5: Replace image URLs
- [ ] Step 6: Preview or upload to draft
```

**Default Settings** (no need to specify):
- Mode: `ai` (AI-powered, no API key needed)
- Theme: `ocean-calm` (professional blue, great for tech/business)

---

## Step 1: Analyze Markdown

Read the markdown file and extract:

| Element | How to Extract |
|---------|----------------|
| **Title** | First `# heading` or filename |
| **Author** | Look for `Author:` or `作者:` in frontmatter |
| **Digest** | First paragraph or generate from content (max 120 chars) |
| **Images** | Collect all `![alt](src)` references |
| **Structure** | Headings, lists, code blocks, quotes, tables |

**Image Reference Types**:

| Type | Syntax | Processing |
|------|--------|------------|
| Local | `![alt](./path/image.png)` | Upload to WeChat |
| Online | `![alt](https://example.com/image.png)` | Download then upload |
| AI Generate | `![alt](__generate:prompt__)` | Generate via AI then upload |

---

## Step 2: Select Theme

### AI Themes

| Theme | Description | Best For |
|-------|-------------|----------|
| **ocean-calm** (default) | Calm blue tones, professional, rational | Tech articles, business analysis |
| **autumn-warm** | Warm orange tones, emotional, literary | Stories, lifestyle, essays |
| **spring-fresh** | Fresh green tones, natural, vibrant | Travel, nature, outdoor |
| **custom** | Use custom prompt | Brand customization |

**Ask the user**: "Which theme would you like? (ocean-calm/autumn-warm/spring-fresh)" - Only ask if the user doesn't specify in their request.

**Default**: Use `ocean-calm` theme if user doesn't specify.

Read detailed style prompts from [references/themes.md](references/themes.md)

---

## Step 3: Generate HTML

Read the selected style prompt from `references/themes.md` and generate HTML with **inline CSS**.

**Important Rules**:

1. All CSS must be **inline** (in `style` attributes)
2. No external stylesheets or scripts
3. Use WeChat-safe HTML tags only
4. Image placeholder format: `<!-- IMG:0 -->`, `<!-- IMG:1 -->`, etc.

**Safe HTML Tags**:
- `<p>`, `<br>`, `<strong>`, `<em>`, `<u>`, `<a>`
- `<h1>` to `<h6>`
- `<ul>`, `<ol>`, `<li>`
- `<blockquote>`, `<pre>`, `<code>`
- `<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, `<td>`
- `<section>`, `<span>` (with inline styles)

**Avoid**:
- `<script>`, `<iframe>`, `<form>`
- External CSS/JS references
- Complex positioning (fixed, absolute)

**Critical for WeChat**:
- Create a main `<div>` container immediately after `<body>` to hold all global styles
- Specify `color` explicitly for each `<p>` tag (WeChat resets to black otherwise)
- Use two `<span>` tags for heading symbols: one with color+text-shadow, one with solid color

---

## Step 4: Process Images

### Image Generation Methods

There are **three ways** to generate AI images:

#### Method 1: Natural Language - For Article (Recommended)

Simply describe what you want in plain language:

```
User: "Generate a product concept image at the beginning of article.md"

User: "Add a comparison chart after the third paragraph"

User: "Create an image showing the workflow diagram in article.md"
```

**How I process natural language requests:**

1. **Understand the intent** - Identify where to insert the image
2. **Read the article** - Analyze context to create an appropriate prompt
3. **Insert the syntax** - Add `![alt](__generate:prompt__)` at the correct location
4. **Show the prompt** - Display the generated prompt for transparency
5. **Generate and upload** - Call the conversion command to complete

**Note**: Proceed directly with generation. Only ask for confirmation if the prompt is complex or ambiguous.

**Example conversation:**

```
User: "Add a product image at the start of my article"
Claude: "I'll add a product concept image at the beginning of article.md.
Based on your article about 'Smart Home Hub', I'll use this prompt:
'A modern smart home hub device, sleek white design with LED indicator
lights, minimalist product photography on a clean white background'
I'll proceed with generating the image."
```

#### Method 2: Natural Language - Standalone Image

Generate an image without any article:

```
User: "Generate an image of a cute cat sitting on a windowsill"
User: "Create a product concept: modern smart home device"
User: "Make a diagram showing user signup flow"
```

**I will:**
1. Create an appropriate prompt based on your description
2. Call: `bash skill/md2wechat/scripts/run.sh generate_image "prompt"`
3. Return the WeChat URL and media ID

**Use when:** You just need an image, not for any article.

#### Method 3: Manual Syntax

Write the image generation syntax directly in Markdown:

```markdown
![Product Concept](__generate:A futuristic smart home hub device, sleek design__)
```

**Syntax format:** `![alt text](__generate:prompt__)`

---

### Processing Images by Type

For each image reference in order:

#### Local Image

```bash
bash skill/md2wechat/scripts/run.sh upload_image "/path/to/image.png"
```

Response:
```json
{"success": true, "wechat_url": "https://mmbiz.qpic.cn/...", "media_id": "xxx"}
```

#### Online Image

```bash
bash skill/md2wechat/scripts/run.sh download_and_upload "https://example.com/image.png"
```

#### AI Generated Image (via CLI)

```bash
# Generate with default size (1024x1024 square)
bash skill/md2wechat/scripts/run.sh generate_image "A cute cat sitting on a windowsill"

# Generate with 16:9 ratio for WeChat cover (recommended)
bash skill/md2wechat/scripts/run.sh generate_image --size 2560x1440 "prompt"
```

**WeChat Cover Images**: For article covers, use 16:9 horizontal ratio (2560x1440 recommended) as it displays better in WeChat's feed and article list. Square images (1024x1024) are cropped in preview.

**Note**: AI image generation uses DashScope API (already configured).

**Image Processing Pipeline**:
1. If AI generation: Call DashScope API → get URL
2. **Save local copy** to `./output/images/{timestamp}-{keyword}.png`
3. If online: Download image to temp
4. If local: Read file
5. Compress if width > 1920px (configurable)
6. Upload to WeChat material API (if configured)
7. Return `wechat_url` and `media_id`
8. Store result for HTML replacement

**Local Save Feature**:
- Generated images are automatically saved to `./output/images/` directory
- Filename format: `{timestamp}-{keyword}.png` (e.g., `20260118-113848-Minimalistbookc.png`)
- Keyword is extracted from the prompt (max 15 alphanumeric characters)
- Local copies are preserved even if WeChat upload fails

---

## Step 5: Replace Image URLs

Replace placeholders in HTML:

```html
<!-- Before -->
<!-- IMG:0 -->
<!-- IMG:1 -->

<!-- After -->
<img src="https://mmbiz.qpic.cn/..." />
<img src="https://mmbiz.qpic.cn/..." />
```

Use the WeChat URLs returned from image processing.

---

## Step 6: Preview or Upload

Ask user:

1. **Preview only** - Show HTML for review
2. **Upload to draft** - Create WeChat draft article

### Preview Mode

Display HTML in markdown code block for user to copy.

### Upload Mode

Create draft and run:

```bash
bash skill/md2wechat/scripts/run.sh convert article.md --draft --cover cover.jpg
```

**Required for draft**:
- `WECHAT_APPID` environment variable
- `WECHAT_SECRET` environment variable
- Cover image (use `--cover` or first image in content)

Response:
```json
{"success": true, "media_id": "draft_media_id", "draft_url": "https://mp.weixin.qq.com/..."}
```

---

## Configuration

### Required for WeChat API

| Variable | Description | Required |
|----------|-------------|----------|
| `WECHAT_APPID` | WeChat Official Account AppID | Yes, for draft upload |
| `WECHAT_SECRET` | WeChat API Secret | Yes, for draft upload |

### AI Features (Already Configured)

| Variable | Description | Status |
|----------|-------------|--------|
| `image_provider` | Image service provider (dashscope/tuzi/openai) | ✅ Set to dashscope |
| `image_key` | DashScope API key for z-image-turbo model | ✅ Configured |
| `image_base_url` | DashScope API base URL | ✅ Set to https://dashscope.aliyuncs.com/api/v1 |
| `image_model` | Model name (z-image-turbo) | ✅ Configured |
| `convert_mode` | Conversion mode (always `ai`) | ✅ Set to AI |
| `default_theme` | Default theme (ocean-calm) | ✅ Set |

**DashScope Image Generation:**
- Uses `/services/aigc/multimodal-generation/generation` endpoint
- Supports `z-image-turbo` model (fast, high quality)
- Image size format: `"1024*1024"` (use `*` not `x`)
- Requires API key in config file (already configured)

### How to Get AppID and Secret

1. Visit [WeChat Developer Platform](https://developers.weixin.qq.com/platform)
2. Login and select your Official Account
3. Go to **Settings & Development** → **Basic Configuration**
4. Find in **Developer ID** section:
   - **Developer ID (AppID)**: Copy directly
   - **Developer Password (AppSecret)**: Click "Reset" to get
5. Add these values to config file: `~/.config/md2wechat/config.yaml`

> **Warning**: AppSecret is very important, keep it secure!

### Config File Location

```
~/.config/md2wechat/config.yaml    # Global config
```

---

## Error Handling

| Error | Action |
|-------|--------|
| Missing config | Ask user to set environment variable or run `md2wechat config init` |
| Image upload fails | Log error, continue with placeholder |
| WeChat API fails | Show error message, return HTML for manual upload |
| Markdown parse error | Ask user to check file format |
| IP not in whitelist | Guide user to add IP to WeChat whitelist (see Troubleshooting) |

---

## Complete Examples

### Example 1: Simple Article (No Images)

**Input**: `simple.md`
```markdown
# My First Article

This is a simple article with no images.
```

**Process**:
1. Generate HTML with ocean-calm theme (default)
2. Skip image processing
3. Ask: preview or upload?
4. If upload → create draft

### Example 2: Article with Local Images

**Input**: `with-images.md`
```markdown
# Travel Diary

Day 1 in Paris:

![Eiffel Tower](./photos/eiffel.jpg)
```

**Process**:
1. Analyze: 1 local image
2. Generate HTML with `<!-- IMG:0 -->` placeholder
3. Run: `upload_image "./photos/eiffel.jpg"`
4. Replace placeholder with WeChat URL
5. Preview or upload

### Example 3: AI Theme with Custom Style

**Input**: `story.md`
```markdown
# The Old Library

A story about memories...
```

**Process**:
1. User selects autumn-warm theme
2. Read theme prompt from references/themes.md
3. Generate themed HTML with inline CSS
4. Preview or upload

### Example 4: AI Image Generation via Natural Language

**User Request:**
```
"Help me add a product concept image at the beginning of article.md"
```

**Process**:
1. Read article.md to understand the product
2. Create an appropriate image prompt based on context
3. Confirm with user: "I'll use this prompt: '...'"
4. Insert `![Product Concept](__generate:...)` at line 2
5. Run conversion command to generate and upload

**Result:** Image generated and uploaded to WeChat

---

### Example 5: Article with Pre-written Image Syntax

**Input**: `mixed.md`
```markdown
# Tech Review

![Product Photo](./product.jpg)

![Comparison Chart](https://example.com/chart.png)

![Concept Art](__generate:Futuristic gadget design__)
```

**Process**:
1. Process 3 images in order
2. Each returns WeChat URL
3. Replace all placeholders
4. Final HTML with all WeChat-hosted images

---

## References

- [Style Themes](references/themes.md) - Detailed style prompts for AI themes
- [HTML Guide](references/html-guide.md) - WeChat HTML constraints and best practices
- [Image Syntax](references/image-syntax.md) - Image reference syntax and generation
- [Writing Guide](references/writing-guide.md) - Writer style assistant documentation
- [WeChat API](references/wechat-api.md) - API reference

---

## Troubleshooting

### Configuration Issues

**Q: Do I need to configure anything?**
A: No! The skill comes pre-configured:
- ✅ AI mode enabled (no md2wechat.cn API key needed)
- ✅ ocean-calm theme set as default
- ✅ DashScope image generation configured
- ✅ Images auto-save to `./output/images/`

Only configure WeChat AppID/Secret if you want to upload drafts.

**Q: "AppID not configured" error**
A: Set `WECHAT_APPID` and `WECHAT_SECRET` in config file `~/.config/md2wechat/config.yaml`

**Q: Config file not working**
A: Check config file location: `~/.config/md2wechat/config.yaml`

### Image Issues

**Q: Image upload fails with "invalid filetype"**
A: WeChat supports JPG, PNG, GIF. Ensure image is in correct format:
```bash
# Convert using ImageMagick
convert input.tiff output.jpg
```

**Q: Images not showing in draft**
A: Images must use WeChat-hosted URLs (`mmbiz.qpic.cn`), not external URLs.

**Q: AI image generation fails**
A: Check that DashScope API key is configured correctly (should already be set up).

### WeChat API Issues

**Q: "IP not in whitelist" error**
A: Add your server IP to WeChat whitelist:

1. Get your public IP:
```bash
curl ifconfig.me
# or
curl ip.sb
```

2. Add IP to WeChat:
   - Visit [WeChat Developer Platform](https://developers.weixin.qq.com/platform)
   - Go to **Settings & Development** → **Basic Configuration**
   - Find **IP Whitelist** section
   - Click "Set" and add your IP
   - Wait a few minutes for changes to take effect

**Q: "access_token expired" error**
A: Program auto-refreshes tokens. If persists:
```bash
# Check config
bash skill/md2wechat/scripts/run.sh config show
```

**Q: "create draft failed" error**
A: Possible causes:
1. Insufficient permissions - ensure account is verified
2. Sensitive content - check article content
3. Draft limit reached - check existing drafts

### HTML/Style Issues

**Q: Styles not working in WeChat editor**
A: Check:
1. CSS uses inline `style` attributes (not `<style>` tags)
2. CSS properties are in allowed list (see HTML Guide)
3. No syntax errors (unclosed tags, etc.)

**Q: Background color lost in WeChat**
A: WeChat strips `<body>` styles. Use main container:
```html
<div style="background-color: #faf9f5; padding: 40px 10px;">
  <!-- All content here -->
</div>
```

**Q: Text color not as expected**
A: WeChat resets `<p>` color to black. Always specify:
```html
<p style="color: #4a413d;">Your text here</p>
```

### Command Issues

**Q: "command not found: md2wechat"**
A: The `run.sh` script will auto-download the binary on first run.

**Q: AI mode is slow**
A: AI mode requires API call and takes 10-30 seconds. This is normal for generating high-quality themed layouts.

---

## CLI Commands Reference

All commands go through the `run.sh` wrapper, which handles auto-installation:

```bash
# Show help
bash skill/md2wechat/scripts/run.sh --help

# Convert and preview (default: ocean-calm theme)
bash skill/md2wechat/scripts/run.sh convert article.md --preview

# Convert with different theme
bash skill/md2wechat/scripts/run.sh convert article.md --theme autumn-warm --preview

# Convert and upload to draft
bash skill/md2wechat/scripts/run.sh convert article.md --draft --cover cover.jpg

# Upload single image
bash skill/md2wechat/scripts/run.sh upload_image photo.jpg

# Download and upload online image
bash skill/md2wechat/scripts/run.sh download_and_upload https://example.com/image.jpg

# Generate AI image (DashScope - already configured)
bash skill/md2wechat/scripts/run.sh generate_image "A cute cat sitting on a windowsill"

# Generate with 16:9 ratio for WeChat cover (recommended)
bash skill/md2wechat/scripts/run.sh generate_image --size 2560x1440 "prompt"

# Show config
bash skill/md2wechat/scripts/run.sh config show

# List available writing styles
bash skill/md2wechat/scripts/run.sh write --list

# Write with creator style (interactive)
bash skill/md2wechat/scripts/run.sh write

# Write with specific style
bash skill/md2wechat/scripts/run.sh write --style dan-koe

# Generate cover prompt only
bash skill/md2wechat/scripts/run.sh write --style dan-koe --cover-only
```
