# Itadaki Web Color Palette

This document defines the color system for the Itadaki web application, ensuring consistency with the Android client in both Light and Dark modes.

## Brand Identity

### Light Mode (Default)
| Color Name | Hex Code | Usage |
| :--- | :--- | :--- |
| **Itadaki Red** | `#E63946` | Primary Actions, Buttons |
| **Itadaki Crimson** | `#9B1C1C` | Headers, Strong Accents |
| **Itadaki Cream** | `#FFF8E1` | Main Background |
| **Itadaki Charcoal** | `#1D3557` | Primary Text |
| **Surface White** | `#FFFFFF` | Cards, Modals |

### Dark Mode
| Color Name | Hex Code | Usage |
| :--- | :--- | :--- |
| **Itadaki Red (Dark)** | `#FF8A80` | Primary Actions (High Contrast) |
| **Itadaki Crimson (Dark)** | `#E57373` | Headers (Softer Red) |
| **Itadaki Dark Gray** | `#1E1E1E` | Main Background (replaces Cream) |
| **Itadaki Off-White** | `#E0E0E0` | Primary Text (replaces Charcoal) |
| **Surface Dark** | `#121212` | Cards, Modals |

---

## Glassmorphism System

### Light Glass
- **Overlay**: `rgba(255, 255, 255, 0.2)`
- **Border**: `rgba(255, 255, 255, 0.27)`
- **Shadow**: `0 4px 30px rgba(0, 0, 0, 0.1)`

### Dark Glass
- **Overlay**: `rgba(255, 255, 255, 0.05)` (More subtle)
- **Border**: `rgba(255, 255, 255, 0.1)`
- **Shadow**: `0 4px 30px rgba(0, 0, 0, 0.5)` (Stronger shadow)

---

## CSS Variables (Svelte Implementation)

Copy these into your `app.css`. This system automatically switches based on the user's system preference.

```css
:root {
 /* --- Functional Colors (Shared) --- */
 --status-success: #10B981;
 --status-warning: #F59E0B;
 --status-info: #3B82F6;
 --status-error: #EF4444;

 /* --- Light Mode (Default) --- */
 --color-primary: #E63946;
 --color-primary-dark: #D32F2F;
 --color-accent: #9B1C1C;
 
 --color-bg-main: #FFF8E1;
 --color-bg-surface: #FFFFFF;
 
 --color-text-primary: #1D3557;
 --color-text-secondary: #455A64;
 --color-text-inverse: #F1FAEE;

 /* Glassmorphism - Light */
 --glass-bg: rgba(255, 255, 255, 0.2);
 --glass-border: 1px solid rgba(255, 255, 255, 0.27);
 --glass-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
 --glass-blur: blur(20px);
}

/* --- Dark Mode Overrides --- */
@media (prefers-color-scheme: dark) {
 :root {
 /* Brand Colors (Adjusted for contrast) */
 --color-primary: #FF8A80;
 --color-primary-dark: #E57373;
 --color-accent: #FFCDD2;

 /* Backgrounds */
 --color-bg-main: #121212; /* Deep dark background */
 --color-bg-surface: #1E1E1E; /* Slightly lighter surface */

 /* Text */
 --color-text-primary: #E6E1E5;
 --color-text-secondary: #CAC4D0;
 --color-text-inverse: #1C1B1F;

 /* Glassmorphism - Dark */
 --glass-bg: rgba(30, 30, 30, 0.6); /* Darker, more opaque for readability */
 --glass-border: 1px solid rgba(255, 255, 255, 0.1);
 --glass-shadow: 0 4px 30px rgba(0, 0, 0, 0.5);
 }
}
```

### Usage Example

```css
.card {
 background-color: var(--color-bg-surface);
 color: var(--color-text-primary);
 border: var(--glass-border);
}

.button {
 background-color: var(--color-primary);
 color: var(--color-text-inverse);
}
```