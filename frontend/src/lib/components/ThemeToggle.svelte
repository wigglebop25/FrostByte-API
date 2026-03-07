<script lang="ts">
    /**
     * Theme Toggle Component
     *
     * Provides light/dark/system theme switching with two display modes:
     * compact (single cycle button) and full (segmented control).
     */
    import { theme, type Theme } from '$lib/stores/theme';
    import { Sun, Moon, Monitor } from 'lucide-svelte';

    export let compact = false;

    const modes: { value: Theme; label: string }[] = [
        { value: 'light', label: 'Light' },
        { value: 'dark', label: 'Dark' },
        { value: 'system', label: 'System' }
    ];

    function cycle() {
        const order: Theme[] = ['light', 'dark', 'system'];
        const idx = order.indexOf($theme);
        const next = order[(idx + 1) % order.length];
        /** Applies a brief CSS transition class during theme switch */
        document.documentElement.classList.add('theme-transitioning');
        theme.set(next);
        setTimeout(() => document.documentElement.classList.remove('theme-transitioning'), 500);
    }

    function setMode(mode: Theme) {
        document.documentElement.classList.add('theme-transitioning');
        theme.set(mode);
        setTimeout(() => document.documentElement.classList.remove('theme-transitioning'), 500);
    }
</script>

{#if compact}
    <!-- Compact Mode — Single Cycle Button -->
    <button
        on:click={cycle}
        class="relative p-2.5 rounded-xl transition-all duration-200
               bg-[var(--glass-bg)] border border-[var(--glass-border)]
               hover:border-[var(--color-primary)] hover:shadow-md
               backdrop-blur-sm group"
        title="Toggle theme: {$theme}"
    >
        {#if $theme === 'light'}
            <Sun size={18} class="text-amber-500 group-hover:rotate-45 transition-transform duration-300" />
        {:else if $theme === 'dark'}
            <Moon size={18} class="text-indigo-400 group-hover:-rotate-12 transition-transform duration-300" />
        {:else}
            <Monitor size={18} class="text-text-secondary group-hover:scale-110 transition-transform duration-300" />
        {/if}
    </button>
{:else}
    <!-- Full Mode — Segmented Control -->
    <div class="inline-flex items-center gap-1 p-1 rounded-2xl bg-[var(--glass-bg)] border border-[var(--glass-border)] backdrop-blur-sm">
        {#each modes as mode}
            <button
                on:click={() => setMode(mode.value)}
                class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-semibold transition-all duration-200
                       {$theme === mode.value
                         ? 'bg-[var(--color-primary)] text-[var(--color-text-inverse)] shadow-lg'
                         : 'text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] hover:bg-[var(--color-bg-surface)]'}"
            >
                {#if mode.value === 'light'}
                    <Sun size={15} />
                {:else if mode.value === 'dark'}
                    <Moon size={15} />
                {:else}
                    <Monitor size={15} />
                {/if}
                <span>{mode.label}</span>
            </button>
        {/each}
    </div>
{/if}
