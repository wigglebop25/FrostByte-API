<script lang="ts">
    /**
     * Settings & Preferences Page
     *
     * Provides user profile overview, theme/display customization,
     * and system infrastructure details. Accessible to all
     * authenticated roles.
     */
    import { onMount } from "svelte";
    import { ws } from "$lib/stores/websocket";
    import { theme } from "$lib/stores/theme";
    import api from "$lib/utils/api";
    import { User, Shield, Globe, Cpu, LogOut, Info, Palette } from "lucide-svelte";
    import { goto } from "$app/navigation";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";

    let user: any = null;
    let roles: string[] = [];
    let apiVersion = "Checking...";
    let activeTab = 'profile';

    const tabs = [
        { id: 'profile', label: 'Account Profile', icon: User },
        { id: 'theme', label: 'Theme & Display', icon: Palette },
        { id: 'system', label: 'About System', icon: Info }
    ];

    onMount(async () => {
        const userStr = localStorage.getItem("user");
        if (userStr) {
            user = JSON.parse(userStr);
            roles = user.roles?.map((r: any) => r.name) || [];
        }

        try {
            const res = await api.get('/health');
            apiVersion = "v" + res.data.version;
        } catch (e) {
            console.error("Failed to fetch version", e);
            apiVersion = "Unknown";
        }
    });

    function logout() {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        goto("/login");
    }
</script>

<div class="space-y-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div>
        <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">Settings</h1>
        <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Manage your profile and system preferences.</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 md:grid-cols-4 gap-6">
        <!-- Settings Navigation Tabs -->
        <div class="space-y-1">
            {#each tabs as tab}
                <button 
                    on:click={() => activeTab = tab.id}
                    class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl font-semibold text-sm transition-all duration-200
                           {activeTab === tab.id 
                             ? 'bg-[var(--color-primary)] text-[var(--color-text-inverse)] shadow-lg' 
                             : 'text-[var(--color-text-secondary)] hover:bg-[var(--glass-bg)] hover:text-[var(--color-text-primary)]'}"
                >
                    <svelte:component this={tab.icon} size={18} />
                    <span>{tab.label}</span>
                </button>
            {/each}
        </div>

        <!-- Tab Content Panels -->
        <div class="sm:col-span-2 md:col-span-3 space-y-5">
            <!-- Account Profile Panel -->
            {#if activeTab === 'profile'}
                <div class="glass-card p-6">
                    <div class="flex items-center gap-5 mb-6 pb-6 border-b border-[var(--glass-border-subtle)]">
                        <div class="w-16 h-16 rounded-2xl bg-[var(--color-primary)]/10 flex items-center justify-center text-[var(--color-primary)] text-2xl font-bold">
                            {user?.username?.[0]?.toUpperCase() || 'U'}
                        </div>
                        <div>
                            <h2 class="text-xl font-bold text-[var(--color-text-primary)]">{user?.username || 'Guest User'}</h2>
                            <div class="flex gap-1.5 mt-1.5">
                                {#each roles as role}
                                    <span class="text-[10px] font-bold uppercase tracking-widest text-[var(--color-primary)] px-2 py-0.5 bg-[var(--color-primary)]/10 rounded-lg border border-[var(--color-primary)]/20">
                                        {role}
                                    </span>
                                {/each}
                            </div>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <div class="flex justify-between items-center py-2.5">
                            <div class="flex items-center gap-3">
                                <Shield size={16} class="text-[var(--color-text-secondary)]" />
                                <span class="font-semibold text-sm text-[var(--color-text-primary)]">Role Permissions</span>
                            </div>
                            <span class="text-sm text-[var(--color-text-secondary)]">
                                {user?.roles?.[0]?.permissions || 'Standard Access'}
                            </span>
                        </div>
                        <div class="flex justify-between items-center py-2.5">
                            <div class="flex items-center gap-3">
                                <Globe size={16} class="text-[var(--color-text-secondary)]" />
                                <span class="font-semibold text-sm text-[var(--color-text-primary)]">Language</span>
                            </div>
                            <span class="text-sm text-[var(--color-text-secondary)]">English / Japanese</span>
                        </div>
                    </div>
                </div>

                <!-- Session Termination (Danger Zone) -->
                <div class="glass-card p-6 border-[var(--status-error)]/20">
                    <div class="flex items-center justify-between">
                        <div>
                            <h3 class="font-bold text-[var(--status-error)]">Sign Out</h3>
                            <p class="text-sm text-[var(--color-text-secondary)] mt-0.5">End your current session.</p>
                        </div>
                        <button on:click={logout} class="px-5 py-2.5 bg-[var(--status-error)] hover:opacity-90 text-white font-bold rounded-xl transition-all shadow-lg flex items-center gap-2 active:scale-[0.97] text-sm">
                            <LogOut size={16} />
                            Sign Out
                        </button>
                    </div>
                </div>
            {/if}

            <!-- Theme & Display Preferences Panel -->
            {#if activeTab === 'theme'}
                <div class="glass-card p-6">
                    <h3 class="text-lg font-bold text-[var(--color-text-primary)] mb-1">Appearance</h3>
                    <p class="text-sm text-[var(--color-text-secondary)] mb-5">Choose your preferred color scheme.</p>
                    
                    <ThemeToggle compact={false} />

                    <div class="mt-6 pt-5 border-t border-[var(--glass-border-subtle)]">
                        <div class="text-xs font-semibold text-[var(--color-text-secondary)] mb-3 uppercase tracking-wider">Preview</div>
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                            <div class="glass-card p-4 text-center">
                                <div class="w-8 h-8 rounded-xl bg-[var(--color-primary)] mx-auto mb-2"></div>
                                <span class="text-xs font-semibold text-[var(--color-text-secondary)]">Primary</span>
                            </div>
                            <div class="glass-card p-4 text-center">
                                <div class="w-8 h-8 rounded-xl bg-[var(--color-bg-surface)] border border-[var(--glass-border)] mx-auto mb-2"></div>
                                <span class="text-xs font-semibold text-[var(--color-text-secondary)]">Surface</span>
                            </div>
                            <div class="glass-card p-4 text-center">
                                <div class="w-8 h-8 rounded-xl bg-[var(--color-bg-main)] border border-[var(--glass-border)] mx-auto mb-2"></div>
                                <span class="text-xs font-semibold text-[var(--color-text-secondary)]">Background</span>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}

            <!-- System Infrastructure Details Panel -->
            {#if activeTab === 'system'}
                <div class="glass-card p-6">
                    <h3 class="text-lg font-bold text-[var(--color-text-primary)] mb-5 flex items-center gap-2">
                        <Cpu size={18} class="text-[var(--color-primary)]" />
                        System Infrastructure
                    </h3>
                    
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div class="p-4 rounded-2xl bg-[var(--glass-bg)] border border-[var(--glass-border-subtle)]">
                            <div class="text-[10px] font-semibold text-[var(--color-text-secondary)] uppercase mb-1 tracking-wider">API Version</div>
                            <div class="font-mono text-sm font-bold text-[var(--color-text-primary)]">{apiVersion}</div>
                        </div>
                        <div class="p-4 rounded-2xl bg-[var(--glass-bg)] border border-[var(--glass-border-subtle)]">
                            <div class="text-[10px] font-semibold text-[var(--color-text-secondary)] uppercase mb-1 tracking-wider">WebSocket</div>
                            <div class="font-mono text-sm font-bold {$ws.connected ? 'text-[var(--status-success)]' : 'text-[var(--status-error)]'}">
                                {$ws.connected ? 'Connected' : 'Offline'}
                            </div>
                        </div>
                        <div class="p-4 rounded-2xl bg-[var(--glass-bg)] border border-[var(--glass-border-subtle)]">
                            <div class="text-[10px] font-semibold text-[var(--color-text-secondary)] uppercase mb-1 tracking-wider">Platform</div>
                            <div class="font-mono text-sm font-bold text-[var(--color-text-primary)]">Azure VM</div>
                        </div>
                        <div class="p-4 rounded-2xl bg-[var(--glass-bg)] border border-[var(--glass-border-subtle)]">
                            <div class="text-[10px] font-semibold text-[var(--color-text-secondary)] uppercase mb-1 tracking-wider">Region</div>
                            <div class="font-mono text-sm font-bold text-[var(--color-text-primary)]">SE Asia</div>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</div>