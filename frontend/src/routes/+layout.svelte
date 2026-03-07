<script lang="ts">
    /**
     * Main Application Layout
     *
     * Manages the global shell including sidebar navigation, top header bar,
     * WebSocket connectivity, authentication guards, and RBAC-based menu
     * filtering. Bypassed entirely on the login route.
     */
    import "../app.css";
    import { LayoutDashboard, ShoppingCart, Users, Settings, LogOut, Coffee, UtensilsCrossed, Tag, Shield, Menu, X, Bell } from "lucide-svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { auth } from "$lib/stores/auth";
    import { ws } from "$lib/stores/websocket";
    import { theme } from "$lib/stores/theme";
    import { env } from "$env/dynamic/public";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";

    $: isLoginPage = $page.url.pathname === '/login';

    let sidebarOpen = true;
    let mobileMenuOpen = false;

    onMount(() => {
        auth.init();
        theme.init();
        
        const token = localStorage.getItem('token');
        if (!token && !isLoginPage) {
            goto('/login');
        }

        /** Listen for automatic token refresh events to update the auth store reactively. */
        window.addEventListener('token-refreshed', ((e: CustomEvent) => {
            auth.updateToken(e.detail);
        }) as EventListener);
    });

    /**
     * Reactively establish or reconnect the WebSocket when the auth token changes.
     * Resets the retry counter on new tokens (e.g., after login or token refresh).
     */
    $: if ($auth.token && !isLoginPage) {
        const wsUrl = env.PUBLIC_WS_URL || 'wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws';
        ws.resetRetries();
        ws.connect(wsUrl, $auth.token);
    }

    const allMenuItems = [
        { 
            name: 'Dashboard', 
            icon: LayoutDashboard, 
            href: '/', 
            allowed: ['Admin', 'Cashier'] 
        },
        { 
            name: 'Order Now', 
            icon: UtensilsCrossed, 
            href: '/orders/create', 
            allowed: ['Customer'] 
        },
        { 
            name: 'Orders', 
            icon: ShoppingCart, 
            href: '/orders', 
            allowed: ['Admin', 'Cashier', 'Customer'] 
        },
        { 
            name: 'Products', 
            icon: Coffee, 
            href: '/products', 
            allowed: ['Admin'] 
        },
        { 
            name: 'Categories', 
            icon: Tag, 
            href: '/categories', 
            allowed: ['Admin'] 
        },
        { 
            name: 'Users', 
            icon: Users, 
            href: '/users', 
            allowed: ['Admin'] 
        },
        { 
            name: 'Roles', 
            icon: Shield, 
            href: '/roles', 
            allowed: ['Admin'] 
        },
        { 
            name: 'Settings', 
            icon: Settings, 
            href: '/settings', 
            allowed: ['All'] 
        },
    ];

    $: roles = $auth.user?.roles?.map((r: any) => r.name) || [];
    
    $: visibleMenuItems = allMenuItems.filter(item => {
        if (!$auth.user) return false;
        if (item.allowed.includes('All')) return true;
        return item.allowed.some(allowedRole => roles.includes(allowedRole));
    });

    function logout() {
        auth.logout();
        goto("/login");
    }

    function getPageTitle(pathname: string): string {
        const item = allMenuItems.find(i => i.href === pathname);
        return item?.name || 'Page';
    }
</script>

{#if isLoginPage}
    <slot />
{:else}
    <div class="flex h-dvh bg-[var(--color-bg-main)] text-[var(--color-text-primary)] overflow-hidden">
        <!-- Decorative Animated Mesh Background -->
        <div class="mesh-bg">
            <div class="mesh-orb w-[250px] sm:w-[350px] lg:w-[500px] h-[250px] sm:h-[350px] lg:h-[500px] top-[-5%] right-[-5%]" style="background: var(--mesh-1);"></div>
            <div class="mesh-orb w-[200px] sm:w-[300px] lg:w-[400px] h-[200px] sm:h-[300px] lg:h-[400px] bottom-[10%] left-[5%]" style="background: var(--mesh-2); animation-delay: -10s;"></div>
            <div class="mesh-orb w-[150px] sm:w-[200px] lg:w-[300px] h-[150px] sm:h-[200px] lg:h-[300px] top-[40%] left-[40%]" style="background: var(--mesh-3); animation-delay: -15s;"></div>
        </div>

        <!-- Mobile Menu Overlay -->
        {#if mobileMenuOpen}
            <button 
                class="fixed inset-0 bg-black/40 backdrop-blur-sm z-30 lg:hidden"
                on:click={() => mobileMenuOpen = false}
                aria-label="Close navigation menu"
            ></button>
        {/if}

        <!-- Navigation Sidebar — RBAC-Filtered Menu Items -->
        <aside class="
            fixed lg:relative z-40 h-full
            w-[260px] flex flex-col
            glass-panel
            border-r border-[var(--glass-border)]
            transition-transform duration-300 ease-out
            {mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        ">
            <!-- Brand Identity -->
            <a href="/" class="px-5 py-5 flex items-center gap-3 hover:opacity-80 transition-opacity border-b border-[var(--glass-border-subtle)]">
                <div class="w-10 h-10 rounded-2xl bg-[var(--color-bg-surface)] shadow-md p-1 flex-shrink-0">
                    <img src="/images/itadaki_logo.png" alt="Logo" class="w-full h-full object-contain" />
                </div>
                <div class="text-xl font-extrabold tracking-tight text-[var(--color-primary)]">Itadaki</div>
            </a>
            
            <!-- Primary Navigation Links -->
            <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
                <div class="px-3 mb-3 text-[10px] font-bold uppercase tracking-widest text-[var(--color-text-secondary)] opacity-60">Menu</div>
                {#each visibleMenuItems as item}
                    <a href={item.href} 
                       on:click={() => mobileMenuOpen = false}
                       class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 font-semibold text-sm relative
                              {$page.url.pathname === item.href 
                                ? 'bg-[var(--color-primary)] text-[var(--color-text-inverse)] shadow-lg' 
                                : 'text-[var(--color-text-secondary)] hover:bg-[var(--glass-bg)] hover:text-[var(--color-text-primary)]'}">
                        {#if $page.url.pathname === item.href}
                            <div class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-5 bg-[var(--color-text-inverse)] rounded-r-full"></div>
                        {/if}
                        <svelte:component this={item.icon} size={18} />
                        <span>{item.name}</span>
                    </a>
                {/each}
            </nav>

            <!-- Authenticated User Profile & Sign Out -->
            <div class="p-3 border-t border-[var(--glass-border-subtle)]">
                {#if $auth.user}
                    <div class="flex items-center gap-3 px-3 py-2 mb-2">
                        <div class="w-8 h-8 rounded-xl bg-[var(--color-primary)]/10 flex items-center justify-center text-[var(--color-primary)] text-xs font-bold">
                            {$auth.user.username[0]?.toUpperCase()}
                        </div>
                        <div class="flex-1 min-w-0">
                            <div class="text-sm font-bold text-[var(--color-text-primary)] truncate">{$auth.user.username}</div>
                            <div class="text-[10px] text-[var(--color-text-secondary)] font-medium">{roles[0] || 'User'}</div>
                        </div>
                    </div>
                {/if}
                <button on:click={logout} class="flex items-center gap-3 w-full px-3 py-2.5 rounded-xl text-[var(--color-text-secondary)] hover:bg-[var(--status-error)]/10 hover:text-[var(--status-error)] transition-all duration-200 text-sm font-semibold">
                    <LogOut size={18} />
                    <span>Sign Out</span>
                </button>
            </div>
        </aside>

        <!-- Main Content Area -->
        <div class="flex-1 flex flex-col min-w-0">
            <!-- Top Header Bar — Page Title, WebSocket Status, Controls -->
            <header class="h-16 flex items-center justify-between px-6 border-b border-[var(--glass-border-subtle)] glass-panel shrink-0 z-10">
                <div class="flex items-center gap-4">
                    <!-- Mobile menu toggle -->
                    <button 
                        class="lg:hidden p-2 rounded-xl hover:bg-[var(--glass-bg)] text-[var(--color-text-secondary)] transition-colors"
                        on:click={() => mobileMenuOpen = !mobileMenuOpen}
                    >
                        {#if mobileMenuOpen}
                            <X size={20} />
                        {:else}
                            <Menu size={20} />
                        {/if}
                    </button>
                    <h2 class="text-lg font-bold text-[var(--color-text-primary)]">{getPageTitle($page.url.pathname)}</h2>
                </div>
                
                <div class="flex items-center gap-3">
                    <!-- WebSocket Status -->
                    <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] text-xs font-bold">
                        <span class="relative flex h-2 w-2">
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 {$ws.connected ? 'bg-[var(--status-success)]' : 'bg-[var(--status-error)]'}"></span>
                            <span class="relative inline-flex rounded-full h-2 w-2 {$ws.connected ? 'bg-[var(--status-success)]' : 'bg-[var(--status-error)]'}"></span>
                        </span>
                        <span class="hidden sm:inline {$ws.connected ? 'text-[var(--status-success)]' : 'text-[var(--status-error)]'}">
                            {$ws.connected ? 'Live' : 'Offline'}
                        </span>
                    </div>

                    <!-- Theme Toggle -->
                    <ThemeToggle compact={true} />

                    <!-- Notifications placeholder -->
                    <button class="p-2.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] hover:border-[var(--color-primary)] text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-all">
                        <Bell size={18} />
                    </button>
                </div>
            </header>

            <!-- Page Content -->
            <main class="flex-1 overflow-y-auto overflow-x-hidden p-4 sm:p-6 lg:p-8">
                <slot />
            </main>
        </div>
    </div>
{/if}
