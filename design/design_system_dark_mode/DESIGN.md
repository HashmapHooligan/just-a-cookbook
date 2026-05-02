---
name: Midnight Gastronomy
colors:
  surface: '#1d100d'
  surface-dim: '#1d100d'
  surface-bright: '#463531'
  surface-container-lowest: '#170b08'
  surface-container-low: '#261815'
  surface-container: '#2a1c19'
  surface-container-high: '#362623'
  surface-container-highest: '#41312d'
  on-surface: '#f7ddd7'
  on-surface-variant: '#e2bfb6'
  inverse-surface: '#f7ddd7'
  inverse-on-surface: '#3d2d29'
  outline: '#a98a82'
  outline-variant: '#5a413b'
  surface-tint: '#ffb4a2'
  primary: '#ffb4a2'
  on-primary: '#621100'
  primary-container: '#ff653f'
  on-primary-container: '#601100'
  inverse-primary: '#b12e0a'
  secondary: '#f4be53'
  on-secondary: '#412d00'
  secondary-container: '#b88921'
  on-secondary-container: '#392600'
  tertiary: '#cbbeff'
  on-tertiary: '#322663'
  tertiary-container: '#9b8ed2'
  on-tertiary-container: '#312562'
  error: '#ffb4ab'
  on-error: '#690005'
  error-container: '#93000a'
  on-error-container: '#ffdad6'
  primary-fixed: '#ffdad2'
  primary-fixed-dim: '#ffb4a2'
  on-primary-fixed: '#3c0700'
  on-primary-fixed-variant: '#8a1c00'
  secondary-fixed: '#ffdea6'
  secondary-fixed-dim: '#f4be53'
  on-secondary-fixed: '#271900'
  on-secondary-fixed-variant: '#5d4200'
  tertiary-fixed: '#e6deff'
  tertiary-fixed-dim: '#cbbeff'
  on-tertiary-fixed: '#1d0f4d'
  on-tertiary-fixed-variant: '#493e7b'
  background: '#1d100d'
  on-background: '#f7ddd7'
  surface-variant: '#41312d'
typography:
  h1:
    fontFamily: Plus Jakarta Sans
    fontSize: 40px
    fontWeight: '800'
    lineHeight: '1.2'
    letterSpacing: -0.02em
  h2:
    fontFamily: Plus Jakarta Sans
    fontSize: 32px
    fontWeight: '700'
    lineHeight: '1.3'
    letterSpacing: -0.01em
  h3:
    fontFamily: Plus Jakarta Sans
    fontSize: 24px
    fontWeight: '700'
    lineHeight: '1.4'
  body-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 18px
    fontWeight: '400'
    lineHeight: '1.6'
  body-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 16px
    fontWeight: '400'
    lineHeight: '1.6'
  label-sm:
    fontFamily: Plus Jakarta Sans
    fontSize: 13px
    fontWeight: '600'
    lineHeight: '1'
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  base: 8px
  xs: 4px
  sm: 12px
  md: 24px
  lg: 40px
  xl: 64px
  gutter: 24px
  margin: 32px
---

## Brand & Style

This design system is tailored for high-end culinary enthusiasts who appreciate a sophisticated, late-night aesthetic. The brand personality is premium, moody, and energetic, evoking the atmosphere of an exclusive dimly-lit bistro or a modern chef's kitchen after hours. 

The design style leans into **Corporate Modern with a Glassmorphism influence**. It maintains the structured, reliable feel of the Quasar framework while using layered surfaces and subtle translucency to prevent the dark interface from feeling heavy. The emotional response is one of focus and excitement, highlighting food imagery against a deep, immersive background.

## Colors

The palette is anchored by a deep navy-purple background that provides a rich canvas for high-contrast elements. 

- **Primary Action:** The pink-orange (#FF653F) is reserved for critical interactions, CTAs, and active states to ensure immediate visual recognition.
- **Accents:** The light yellow (#FFC85C) serves as a functional accent for secondary information, star ratings, and category labels.
- **Surfaces:** UI containers use slightly lighter variations of the base purple to create a sense of depth without introducing grey scales.
- **Typography:** Body text uses an off-white (#F8F9FA) for maximum legibility, while headers and highlights leverage the light yellow.

## Typography

This design system utilizes **Plus Jakarta Sans** across all levels to maintain a soft, welcoming, yet professional feel. 

Headlines are set with heavy weights and tighter letter spacing to create a bold, editorial impact. Body text is prioritized for readability with generous line heights. Labels and small captions use a semi-bold weight and increased tracking to remain legible against the dark background.

## Layout & Spacing

The system follows a **12-column fluid grid** model with fixed gutters of 24px. The layout philosophy emphasizes clear grouping through the use of "logical clusters" rather than heavy lines. 

Spacing follows an 8px rhythmic scale. Use `md` (24px) for standard padding within cards and sections, and `lg` (40px) to separate major content blocks. Vertical rhythm is essential in dark mode to prevent visual clutter; prioritize whitespace over borders whenever possible.

## Elevation & Depth

Depth is conveyed through **Tonal Layering** combined with subtle **Ambient Shadows**. Instead of traditional black shadows, elevation is represented by:

1.  **Surface Lightness:** Higher elevation levels correspond to lighter shades of the base purple.
2.  **Shadow Character:** Shadows are extremely diffused, using a low-opacity version of the background color (e.g., `rgba(15, 8, 40, 0.6)`) to create a soft glow effect rather than a harsh drop.
3.  **Inner Glow:** For primary buttons, a 1px inner stroke or "rim light" on the top edge enhances the material feel, making elements appear slightly extruded.

## Shapes

The shape language is **Rounded**, reflecting the approachable and friendly nature of the brand. 

- **Standard Elements:** Buttons, input fields, and small cards use a 0.5rem (8px) radius.
- **Large Containers:** Recipe cards and modal overlays use 1rem (16px) for a softer, more modern appearance.
- **Pill Elements:** Interactive chips and tags use a fully rounded (pill) radius to distinguish them from structural components.

## Components

### Buttons
- **Primary:** Solid #FF653F fill with white text. Apply a subtle outer glow on hover.
- **Secondary:** Outlined with #FFC85C and matching text. 
- **Ghost:** No background, #FFC85C text, used for less prominent actions.

### Input Fields
Inputs should feature a slightly lighter surface than the background (#2A1B66) with a subtle 1px border in a muted purple. Labels should use the light yellow (#FFC85C) for focus states.

### Cards
Cards use a 1px subtle border (#3D2B85) to define edges against the deep background. Imagery should be full-bleed or top-aligned with high-quality photography to leverage the dark contrast.

### Chips & Tags
Used for recipe categories or dietary labels. Use a semi-transparent fill of the accent color (e.g., 15% opacity #FF653F) with a solid colored border and text.

### Selection Controls
Checkboxes and radio buttons use the #FF653F accent for checked states. The "off" state should be a simple high-contrast outline to ensure visibility on the dark canvas.

### Additional Suggestions
- **Ingredient Toggles:** A custom list item with a strike-through animation for "completed" steps.
- **Floating Action Button (FAB):** A Quasar-inspired circular button in #FF653F for "Add to List" or "Start Cooking" features.