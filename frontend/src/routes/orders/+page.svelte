<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Order } from "$lib/types";
    import { ws } from "$lib/stores/websocket";
    import { Search, Filter, Eye, X, Plus, ChevronDown, Check, AlertCircle } from "lucide-svelte";
    import { goto } from "$app/navigation";

    let orders: Order[] = [];
    let loading = true;
    let searchQuery = "";
    let statusFilter = "ALL";
    let filterType = "ALL";
    let filterValue = "";

    // RBAC
    let canUpdateStatus = false;
    let isCustomer = false;

    // Modal State
    let selectedOrder: Order | null = null;
    let showDetailsModal = false;
    let updatingId: number | null = null;

    // --- Function Definitions First ---

    async function fetchOrders() {
        loading = true;
        try {
            let endpoint = "/orders";
            if (filterType === "USER" && filterValue) endpoint = `/orders/user/${filterValue}`;
            else if (filterType === "ROLE" && filterValue) endpoint = `/orders/role/${filterValue}`;

            const res = await api.get(endpoint);
            orders = res.data;
        } catch (e: any) {
            console.error("Failed to fetch orders:", e);
            orders = [];
        } finally {
            loading = false;
        }
    }

    async function handleStatusChange(event: Event, order: Order) {
        const newStatus = (event.target as HTMLSelectElement).value as Order['status'];
        if (!newStatus || newStatus === order.status) return;

        updatingId = order.order_id;
        try {
            await api.put(`/orders/${order.order_id}/status`, { status: newStatus });
            
            // Optimistic Update
            const idx = orders.findIndex(o => o.order_id === order.order_id);
            if (idx !== -1) {
                orders[idx] = { ...orders[idx], status: newStatus };
            }
            if (selectedOrder && selectedOrder.order_id === order.order_id) {
                selectedOrder.status = newStatus;
            }
        } catch (e: any) {
            alert("Failed to update: " + (e.response?.data?.error || "Unknown error"));
            fetchOrders();
        } finally {
            updatingId = null;
        }
    }

    async function openOrderDetails(orderId: number) {
        try {
            const res = await api.get(`/orders/${orderId}`);
            selectedOrder = res.data;
            showDetailsModal = true;
        } catch (e) {
            console.error("Failed details load", e);
        }
    }

    function closeDetailsModal() {
        showDetailsModal = false;
        selectedOrder = null;
    }

    function getStatusColor(status: string) {
        switch (status.toUpperCase()) {
            case 'COMPLETED': return 'text-status-success bg-status-success/10 border-status-success/20';
            case 'READY': return 'text-status-info bg-status-info/10 border-status-info/20';
            case 'COOKING': return 'text-orange-500 bg-orange-500/10 border-orange-500/20';
            case 'ACCEPTED': return 'text-blue-500 bg-blue-500/10 border-blue-500/20';
            case 'PENDING': return 'text-status-warning bg-status-warning/10 border-status-warning/20';
            case 'CANCELLED': return 'text-status-error bg-status-error/10 border-status-error/20';
            default: return 'text-text-secondary';
        }
    }

    // --- Valid Transitions ---
    const nextStatus: Record<string, string[]> = {
        'PENDING': ['ACCEPTED', 'CANCELLED'],
        'ACCEPTED': ['COOKING', 'CANCELLED'],
        'COOKING': ['READY'],
        'READY': ['COMPLETED'],
        'COMPLETED': [],
        'CANCELLED': []
    };

    // --- Lifecycle & Reactivity ---

    onMount(() => {
        const userStr = localStorage.getItem("user");
        if (userStr) {
            const user = JSON.parse(userStr);
            const roles = user.roles?.map((r: any) => r.name.toLowerCase()) || [];
            canUpdateStatus = roles.includes('admin') || roles.includes('cashier');
            isCustomer = roles.includes('customer');
        }
        fetchOrders();
    });

    // Real-time updates
    $: if ($ws.data && ($ws.data.event === 'NEW_ORDER' || $ws.data.event === 'ORDER_UPDATED')) {
        if ($ws.data.event === 'ORDER_UPDATED' && $ws.data.data?.order_id) {
            const idx = orders.findIndex(o => o.order_id === $ws.data.data.order_id);
            if (idx !== -1) {
                orders[idx] = { ...orders[idx], status: $ws.data.data.status };
            }
        } else {
            fetchOrders(); 
        }
    }

    $: filteredOrders = orders.filter(order => {
        const matchesSearch = order.order_id.toString().includes(searchQuery) || 
                              order.user?.username.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus = statusFilter === "ALL" || order.status.toUpperCase() === statusFilter.toUpperCase();
        return matchesSearch && matchesStatus;
    });
</script>

<div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-accent font-serif">Orders</h1>
            <p class="text-text-secondary mt-1">Manage and track customer orders.</p>
        </div>
        <div class="flex gap-3">
             {#if isCustomer}
                <button on:click={() => goto("/orders/create")} class="flex items-center gap-2 px-5 py-2.5 bg-status-success hover:bg-green-600 text-white rounded-xl font-bold transition-all shadow-lg shadow-green-500/20 active:scale-95">
                    <Plus size={18} />
                    New Order
                </button>
             {/if}
            <button on:click={fetchOrders} class="px-5 py-2.5 bg-surface hover:bg-main text-text-primary rounded-xl font-bold transition-colors border border-glass-border shadow-sm">
                Refresh
            </button>
        </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-col xl:flex-row gap-4">
        <div class="glass-card p-2 flex-1 flex items-center px-4">
            <Search class="text-text-secondary mr-3" size={18} />
            <input bind:value={searchQuery} placeholder="Local Search..." class="bg-transparent border-none text-text-primary focus:ring-0 w-full outline-none" />
        </div>
        <div class="glass-card p-2 flex items-center gap-2 px-4">
            <Filter size={16} class="text-text-secondary" />
            <select bind:value={statusFilter} class="bg-transparent border-none text-text-primary text-sm font-bold focus:ring-0 outline-none">
                <option value="ALL">All Status</option>
                <option value="PENDING">Pending</option>
                <option value="ACCEPTED">Accepted</option>
                <option value="COOKING">Cooking</option>
                <option value="READY">Ready</option>
                <option value="COMPLETED">Completed</option>
                <option value="CANCELLED">Cancelled</option>
            </select>
        </div>
    </div>

    <!-- Table -->
    <div class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
            <table class="w-full text-left">
                <thead class="bg-primary/5 border-b border-glass-border">
                    <tr>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider">Order ID</th>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider">Customer</th>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider">Date</th>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider">Total</th>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider">Status</th>
                        <th class="px-6 py-4 font-bold text-text-secondary text-sm uppercase tracking-wider text-right">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-glass-border">
                    {#if loading}
                        <tr><td colspan="6" class="px-6 py-12 text-center text-text-secondary">Loading...</td></tr>
                    {:else}
                        {#each filteredOrders as order}
                            <tr 
                                class="hover:bg-primary/5 transition-colors group" 
                                on:click={() => openOrderDetails(order.order_id)}
                                tabindex="0"
                                on:keydown={(e) => e.key === 'Enter' && openOrderDetails(order.order_id)}
                            >
                                <td class="px-6 py-4 font-bold text-text-primary">#{order.order_id}</td>
                                <td class="px-6 py-4 text-text-primary font-medium">{order.user?.username}</td>
                                <td class="px-6 py-4 text-text-secondary text-sm">{order.created_at}</td>
                                <td class="px-6 py-4 font-mono font-bold text-text-primary">${order.total_amount.toFixed(2)}</td>
                                <td class="px-6 py-4">
                                    <!-- Status Dropdown for quick edit -->
                                    {#if canUpdateStatus && nextStatus[order.status.toUpperCase()]?.length > 0}
                                        <div class="relative inline-block">
                                            <select 
                                                class="appearance-none pl-3 pr-8 py-1.5 rounded-lg text-xs font-bold border cursor-pointer outline-none focus:ring-2 focus:ring-primary/50 disabled:opacity-50 {getStatusColor(order.status)}"
                                                value={order.status.toUpperCase()}
                                                on:change={(e) => handleStatusChange(e, order)}
                                                disabled={updatingId === order.order_id}
                                                on:click|stopPropagation
                                            >
                                                <option value={order.status.toUpperCase()} disabled>{order.status}</option>
                                                {#each nextStatus[order.status.toUpperCase()] || [] as next}
                                                    <option value={next} class="text-text-primary bg-surface">{next}</option>
                                                {/each}
                                            </select>
                                            <ChevronDown size={12} class="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none opacity-50" />
                                        </div>
                                    {:else}
                                        <span class="px-3 py-1.5 rounded-lg text-xs font-bold border {getStatusColor(order.status)}">
                                            {order.status}
                                        </span>
                                    {/if}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <button 
                                        on:click|stopPropagation={() => openOrderDetails(order.order_id)} 
                                        class="p-2 text-text-secondary hover:text-primary transition-colors bg-surface hover:bg-main rounded-lg border border-transparent hover:border-glass-border"
                                    >
                                        <Eye size={18} />
                                    </button>
                                </td>
                            </tr>
                        {/each}
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>

<!-- Details Modal -->
{#if showDetailsModal && selectedOrder}
    <div 
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md" 
        role="dialog"
        aria-modal="true"
        tabindex="-1"
    >
        <!-- Backdrop Button -->
        <button 
            type="button"
            class="absolute inset-0 w-full h-full cursor-default focus:outline-none"
            on:click={closeDetailsModal}
            on:keydown={(e) => e.key === 'Escape' && closeDetailsModal()}
            aria-label="Close modal"
        ></button>

        <div 
            class="glass-card w-full max-w-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh] relative z-10" 
            role="document"
        >
            <div class="p-6 border-b border-glass-border flex justify-between items-start bg-surface/30 backdrop-blur-md">
                <div>
                    <h2 class="text-2xl font-bold text-text-primary">Order #{selectedOrder.order_id}</h2>
                    <p class="text-text-secondary text-sm mt-1">Customer: {selectedOrder.user?.username}</p>
                </div>
                <button on:click={closeDetailsModal} class="text-text-secondary hover:text-primary"><X size={24} /></button>
            </div>
            
            <div class="p-6 overflow-y-auto flex-1">
                <div class="space-y-3">
                    {#if selectedOrder.products}
                        {#each selectedOrder.products as item}
                            <div class="flex justify-between p-4 bg-surface/40 border border-glass-border rounded-xl">
                                <span class="font-bold text-text-primary">{item.product?.name} <span class="text-text-secondary font-normal">x{item.quantity}</span></span>
                                <span class="font-mono text-primary font-bold">${item.line_total.toFixed(2)}</span>
                            </div>
                        {/each}
                    {/if}
                </div>
                <div class="mt-6 pt-6 border-t border-glass-border text-right">
                    <span class="text-2xl font-black text-primary">Total: ${selectedOrder.total_amount.toFixed(2)}</span>
                </div>
            </div>

            <!-- Modal Actions (Alternative to Table Actions) -->
            {#if canUpdateStatus && nextStatus[selectedOrder.status.toUpperCase()]?.length > 0}
                <div class="p-6 bg-surface/30 border-t border-glass-border flex gap-3">
                    {#each nextStatus[selectedOrder.status.toUpperCase()] as status}
                        <button 
                            on:click={(e) => handleStatusChange({ target: { value: status } } as any, selectedOrder!)}
                            class="flex-1 py-3 rounded-xl font-bold text-white shadow-lg transition-all active:scale-95
                                   {status === 'CANCELLED' ? 'bg-status-error hover:bg-red-600' : 'bg-primary hover:bg-primary-dark'}"
                        >
                            Mark as {status}
                        </button>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
{/if}