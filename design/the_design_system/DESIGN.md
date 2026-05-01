---
name: The Design System
colors:
  surface: '#fff8f3'
  surface-dim: '#ffd486'
  surface-bright: '#fff8f3'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#fff2e0'
  surface-container: '#ffebce'
  surface-container-high: '#ffe5ba'
  surface-container-highest: '#ffdea6'
  on-surface: '#271900'
  on-surface-variant: '#48454f'
  inverse-surface: '#412d00'
  inverse-on-surface: '#ffefd7'
  outline: '#797580'
  outline-variant: '#c9c4d1'
  surface-tint: '#615594'
  primary: '#020013'
  on-primary: '#ffffff'
  primary-container: '#1e104e'
  on-primary-container: '#887cbe'
  inverse-primary: '#cbbeff'
  secondary: '#b12e0a'
  on-secondary: '#ffffff'
  secondary-container: '#fd643e'
  on-secondary-container: '#5e1000'
  tertiary: '#03000a'
  on-tertiary: '#ffffff'
  tertiary-container: '#28113c'
  on-tertiary-container: '#957aab'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#e6deff'
  primary-fixed-dim: '#cbbeff'
  on-primary-fixed: '#1d0f4d'
  on-primary-fixed-variant: '#493e7b'
  secondary-fixed: '#ffdad2'
  secondary-fixed-dim: '#ffb4a2'
  on-secondary-fixed: '#3c0700'
  on-secondary-fixed-variant: '#8a1c00'
  tertiary-fixed: '#f1dbff'
  tertiary-fixed-dim: '#d9bbf1'
  on-tertiary-fixed: '#27103b'
  on-tertiary-fixed-variant: '#543d6a'
  background: '#fff8f3'
  on-background: '#271900'
  surface-variant: '#ffdea6'
typography:
  headline-xl:
    fontFamily: Plus Jakarta Sans
    fontSize: 40px
    fontWeight: '800'
    lineHeight: 48px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 32px
    fontWeight: '700'
    lineHeight: 40px
    letterSpacing: -0.01em
  headline-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 24px
    fontWeight: '700'
    lineHeight: 32px
  body-lg:
    fontFamily: Be Vietnam Pro
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Be Vietnam Pro
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Be Vietnam Pro
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-lg:
    fontFamily: Plus Jakarta Sans
    fontSize: 14px
    fontWeight: '600'
    lineHeight: 20px
  label-md:
    fontFamily: Plus Jakarta Sans
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
  instruction-step:
    fontFamily: Be Vietnam Pro
    fontSize: 20px
    fontWeight: '500'
    lineHeight: 32px
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
  lg: 48px
  xl: 64px
  gutter: 24px
  margin: 32px
---

## Brand & Style

The design system is crafted for a premium, high-energy home cooking experience. It blends the efficiency of a professional kitchen tool with the warmth of a domestic culinary environment. The brand personality is optimistic, vibrant, and highly organized. 

The visual style follows a **Modern Corporate** aesthetic with strong **Material Design** influences, drawing inspiration from the Quasar Framework's structural logic. It prioritizes clarity and responsiveness, ensuring that complex recipe data remains accessible even in a busy kitchen environment. The interface uses high-contrast color blocks and generous whitespace to guide the user’s eye through instructions and ingredient lists without friction.

## Colors

The color palette utilizes a high-contrast interplay between deep nocturnal tones and warm, culinary-inspired highlights. 

- **Primary (#1E104E):** A deep, authoritative purple used for navigation, headers, and primary branding elements to provide a stable anchor for the UI.
- **Secondary (#FF653F):** A vibrant, high-visibility pink used exclusively for calls to action, interactive states, and critical updates (e.g., "Start Cooking" or active timers).
- **Tertiary (#452E5A):** A softer purple used for secondary actions, metadata labels, and borders to provide depth without the weight of the primary color.
- **Neutral/Background (#FFC85C):** A light cream/yellow that serves as the primary canvas. This reduces the harshness of pure white while maintaining high readability for recipe text.

Functional states (Success, Warning, Error) should follow standard Material conventions but utilize the primary purple for iconography to maintain brand cohesion.

## Typography

The typography system is engineered for utility. **Plus Jakarta Sans** provides a welcoming, modern feel for headlines and UI labels, while **Be Vietnam Pro** is used for body content and recipe steps due to its exceptional legibility and friendly character.

A specialized `instruction-step` level is included specifically for recipe directions, featuring increased line height and a slightly heavier weight to ensure the user can read the screen from a distance while cooking. Hierarchy is established through aggressive weight differences rather than just size, ensuring that important nutritional or timing data stands out.

## Layout & Spacing

The design system employs a **12-column fluid grid** for desktop and a single-column fluid layout for mobile. A strict 8px base unit (the "quanta") governs all padding and margins to ensure visual harmony.

Content is organized into logical "zones":
- **Preparation Zone:** High-density layouts using `sm` and `md` spacing for ingredient lists.
- **Action Zone:** Generous `lg` spacing around instruction steps to minimize cognitive load.
- **Navigation:** Fixed top or side bars using the primary purple to provide a permanent frame for the application.

Gutters are fixed at 24px to provide clear separation between recipe cards and utility modules.

## Elevation & Depth

This design system uses a **Tonal Layering** approach to depth, supplemented by soft, tinted shadows. Surfaces "lift" off the cream background using pure white containers.

- **Level 0 (Floor):** The Cream background (#FFC85C).
- **Level 1 (Cards/Lists):** White surfaces with a subtle 1px border in #452E5A at 10% opacity.
- **Level 2 (Active States/Modals):** White surfaces with an ambient shadow (Blur: 12px, Y: 4px) tinted with the Primary Purple at 8% opacity.
- **Level 3 (Pop-overs/Tooltips):** Deep Purple (#1E104E) surfaces with high-contrast cream text, creating a reverse-depth effect for immediate attention.

Avoid heavy black shadows; all elevation should feel airy and integrated into the warm palette.

## Shapes

The shape language is consistently **Rounded**, reflecting a friendly and approachable tool. 

- **Standard Components:** Buttons, input fields, and small cards use a 0.5rem (8px) radius.
- **Container Elements:** Recipe cards and main content areas use a "rounded-lg" 1rem (16px) radius to create a soft, modern container.
- **Interactive Pill:** Search bars and status chips use full pill-shaping (rounded-xl) to distinguish them from structural content.

This roundedness should be applied to image containers as well, ensuring that food photography feels integrated into the UI rather than boxed in.

## Components

The component library is built on Quasar-inspired primitives, focusing on touch-friendly targets and high-contrast feedback.

- **Buttons:** Use the Secondary Pink (#FF653F) for "Primary" actions with white text. "Secondary" buttons use the Primary Purple (#1E104E) as an outline. All buttons feature a 2px vertical offset on hover to simulate tactile feedback.
- **Chips:** Used for dietary tags (e.g., Vegan, GF). These should use the Tertiary Blue (#452E5A) with white text, using the "pill" shape.
- **Cards:** White backgrounds with `rounded-lg` corners. Headers within cards use the Primary Purple to create a clear anchor.
- **Input Fields:** Filled style with a 2px bottom border in Tertiary Blue. When focused, the border transitions to Secondary Pink.
- **Recipe Step List:** A custom list component where the step number is displayed in a large, bold Primary Purple circle, ensuring the user never loses their place.
- **Checkboxes:** Larger than standard (24px x 24px) to accommodate kitchen use, using Secondary Pink for the "checked" state.
- **Timers:** A specialized floating component using a high-contrast Deep Purple background and bold Cream typography for the countdown.