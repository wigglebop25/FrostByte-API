<script lang="ts">
    import { onMount } from "svelte";
    import { ws } from "$lib/stores/websocket";
    import { User, Shield, Globe, Cpu, LogOut, Info, Palette } from "lucide-svelte";
    import { goto } from "$app/navigation";

    let user: any = null;
    let roles: string[] = [];
    let apiVersion = "v1.0.2"; // Example version

    onMount(() => {
        const userStr = localStorage.getItem("user");
        if (userStr) {
            user = JSON.parse(userStr);
            roles = user.roles?.map((r: any) => r.name) || [];
        }
    });

    function logout() {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        goto("/login");
    }
</script>

<div class="space-y-8 max-w-4xl mx-auto">
    <!-- Header -->
    <div>
        <h1 class="text-4xl font-bold tracking-tight text-accent font-serif">Settings</h1>
        <p class="text-text-secondary mt-1 text-lg">Manage your profile and system preferences.</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- Sidebar Navigation (Internal) -->
        <div class="space-y-2">
            <button class="w-full flex items-center gap-3 px-4 py-3 bg-primary text-text-inverse rounded-xl font-bold shadow-lg shadow-primary/20 transition-all">
                <User size={20} />
                <span>Account Profile</span>
            </button>
            <button class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface text-text-secondary rounded-xl font-bold transition-all">
                <Palette size={20} />
                <span>Theme & Display</span>
            </button>
            <button class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface text-text-secondary rounded-xl font-bold transition-all">
                <Info size={20} />
                <span>About System</span>
            </button>
        </div>

        <!-- Main Settings Panels -->
        <div class="md:col-span-2 space-y-6">
            <!-- Profile Section -->
            <div class="glass-card p-8 bg-surface/40">
                <div class="flex items-center gap-6 mb-8 pb-8 border-b border-glass-border">
                    <div class="w-20 h-20 rounded-3xl bg-primary/10 flex items-center justify-center text-primary text-3xl font-bold shadow-inner">
                        {user?.username?.[0]?.toUpperCase() || 'U'}
                    </div>
                    <div>
                        <h2 class="text-2xl font-bold text-text-primary">{user?.username || 'Guest User'}</h2>
                        <div class="flex gap-2 mt-1">
                            {#each roles as role}
                                <span class="text-[10px] font-bold uppercase tracking-widest text-primary px-2 py-0.5 bg-primary/10 rounded-full border border-primary/20">
                                    {role}
                                </span>
                            {/each}
                        </div>
                    </div>
                </div>

                <div class="space-y-4">
                    <div class="flex justify-between items-center py-2">
                        <div class="flex items-center gap-3">
                            <Shield size={18} class="text-text-secondary" />
                            <span class="font-bold text-text-primary">Role Permissions</span>
                        </div>
                        <span class="text-sm font-medium text-text-secondary">
                            {user?.roles?.[0]?.permissions || 'Standard Access'}
                        </span>
                    </div>
                    <div class="flex justify-between items-center py-2">
                        <div class="flex items-center gap-3">
                            <Globe size={18} class="text-text-secondary" />
                            <span class="font-bold text-text-primary">Language</span>
                        </div>
                        <span class="text-sm font-medium text-text-secondary">English / Japanese</span>
                    </div>
                </div>
            </div>

            <!-- System Status -->
            <div class="glass-card p-8 bg-surface/40">
                <h3 class="text-lg font-bold text-text-primary mb-6 flex items-center gap-2">
                    <Cpu size={20} class="text-accent" />
                    System Infrastructure
                </h3>
                
                <div class="grid grid-cols-2 gap-4">
                    <div class="p-4 bg-main/30 rounded-2xl border border-glass-border">
                        <div class="text-[10px] font-bold text-text-secondary uppercase mb-1">API Version</div>
                        <div class="font-mono text-sm font-bold text-text-primary">{apiVersion}</div>
                    </div>
                    <div class="p-4 bg-main/30 rounded-2xl border border-glass-border">
                        <div class="text-[10px] font-bold text-text-secondary uppercase mb-1">WebSocket Status</div>
                        <div class="font-mono text-sm font-bold {$ws.connected ? 'text-status-success' : 'text-status-error'}">
                            {$ws.connected ? 'STABLE' : 'OFFLINE'}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Danger Zone -->
            <div class="glass-card p-8 border-status-error/20 bg-status-error/5">
                <div class="flex items-center justify-between">
                    <div>
                        <h3 class="text-lg font-bold text-status-error">Logout</h3>
                        <p class="text-sm text-text-secondary">Sign out from the current session.</p>
                    </div>
                    <button on:click={logout} class="px-6 py-3 bg-status-error hover:bg-red-600 text-white font-bold rounded-xl transition-all shadow-lg shadow-red-500/20 flex items-center gap-2 active:scale-95">
                        <LogOut size={18} />
                        <span>Sign Out</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>