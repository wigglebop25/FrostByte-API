<script lang="ts">
    import "../app.css";
    import { LayoutDashboard, ShoppingCart, Users, Settings, LogOut, Coffee, UtensilsCrossed, Tag, Shield } from "lucide-svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { auth } from "$lib/stores/auth";

    $: isLoginPage = $page.url.pathname === '/login';

    // Initialize auth on first load
    onMount(() => {
        auth.init();
    });

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

    // Reactive Menu: Updates automatically when $auth.user changes
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
</script>

{#if isLoginPage}
    <slot />
{:else}
    <div class="flex h-screen bg-main text-text-primary">
        <!-- Sidebar -->
        <aside class="w-64 glass-card rounded-none border-y-0 border-l-0 border-r flex flex-col z-20">
            <div class="p-6 flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-surface shadow-md p-1 flex-shrink-0">
                    <img src="/images/itadaki_logo.png" alt="Logo" class="w-full h-full object-contain" />
                </div>
                <div class="text-2xl font-bold tracking-tight text-accent italic">Itadaki</div>
            </div>
            
            <nav class="flex-1 px-4 space-y-2 mt-6">
                {#each visibleMenuItems as item}
                    <a href={item.href} 
                       class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-300 font-semibold
                              {$page.url.pathname === item.href 
                                ? 'bg-primary text-text-inverse shadow-xl shadow-primary/30 scale-[1.02]' 
                                : 'text-text-secondary hover:bg-surface/50 hover:text-primary'}">
                        <svelte:component this={item.icon} size={20} />
                        <span>{item.name}</span>
                    </a>
                {/each}
            </nav>

            <div class="p-4 border-t border-glass-border">
                {#if $auth.user}
                    <div class="px-4 py-2 mb-2 text-xs font-bold text-text-secondary uppercase tracking-wider">
                        Signed in as <span class="text-primary">{$auth.user.username}</span>
                    </div>
                {/if}
                <button on:click={logout} class="flex items-center gap-3 w-full px-4 py-3 rounded-xl text-text-secondary hover:bg-status-error/10 hover:text-status-error transition-all duration-300">
                    <LogOut size={20} />
                    <span>Logout</span>
                </button>
            </div>
        </aside>

        <!-- Main Content -->
        <main class="flex-1 overflow-auto p-8 relative scroll-smooth">
            <!-- Background Orbs for visual depth -->
            <div class="fixed top-0 right-0 w-[500px] h-[500px] bg-primary/5 rounded-full blur-[120px] pointer-events-none -z-10 translate-x-1/3 -translate-y-1/3"></div>
            <div class="fixed bottom-0 left-64 w-[400px] h-[400px] bg-status-info/5 rounded-full blur-[100px] pointer-events-none -z-10 -translate-x-1/2 translate-y-1/2"></div>
            
            <slot />
        </main>
    </div>
{/if}
