import type { Config } from 'tailwindcss';
import forms from '@tailwindcss/forms';

export default {
	darkMode: 'class',
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {
			colors: {
				primary: {
					DEFAULT: 'var(--color-primary)',
					dark: 'var(--color-primary-dark)'
				},
				accent: 'var(--color-accent)',
				surface: 'var(--color-bg-surface)',
				main: 'var(--color-bg-main)',
				text: {
					primary: 'var(--color-text-primary)',
					secondary: 'var(--color-text-secondary)',
					inverse: 'var(--color-text-inverse)'
				},
				status: {
					success: 'var(--status-success)',
					warning: 'var(--status-warning)',
					info: 'var(--status-info)',
					error: 'var(--status-error)'
				},
				glass: {
					border: 'var(--glass-border)',
					bg: 'var(--glass-bg)'
				}
			},
			borderRadius: {
				'2xl': '20px',
				'3xl': '24px'
			},
			backdropBlur: {
				'glass': '40px'
			}
		}
	},

	plugins: [forms]
} satisfies Config;