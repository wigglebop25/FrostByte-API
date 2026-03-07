<script lang="ts">
    /**
     * Order Management Page
     *
     * Provides real-time order tracking with search, status filtering,
     * and workflow transitions (Pending → Accepted → Cooking → Ready →
     * Completed). Updates live via WebSocket. Admin/Cashier can update
     * status; Customers have read-only access.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Order } from "$lib/types";
    import { ws } from "$lib/stores/websocket";
    import { Eye, X, Plus, ChevronDown, Check, AlertCircle, RefreshCw } from "lucide-svelte";
    import { goto } from "$app/navigation";

    let orders: Order[] = [];
    let loading = true;
    let searchQuery = "";
    let statusFilter = "ALL";
    let filterType = "ALL";
    let filterValue = "";

    let canUpdateStatus = false;
    let isCustomer = false;

    let selectedOrder: Order | null = null;
    let showDetailsModal = false;
    let updatingId: number | null = null;

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
            case 'COMPLETED': return 'text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 border-emerald-500/20';
            case 'READY': return 'text-blue-600 dark:text-blue-400 bg-blue-500/10 border-blue-500/20';
            case 'COOKING': return 'text-orange-600 dark:text-orange-400 bg-orange-500/10 border-orange-500/20';
            case 'ACCEPTED': return 'text-indigo-600 dark:text-indigo-400 bg-indigo-500/10 border-indigo-500/20';
            case 'PENDING': return 'text-amber-600 dark:text-amber-400 bg-amber-500/10 border-amber-500/20';
            case 'CANCELLED': return 'text-red-600 dark:text-red-400 bg-red-500/10 border-red-500/20';
            default: return 'text-[var(--color-text-secondary)]';
        }
    }

    const nextStatus: Record<string, string[]> = {
        'PENDING': ['ACCEPTED', 'CANCELLED'],
        'ACCEPTED': ['COOKING', 'CANCELLED'],
        'COOKING': ['READY'],
        'READY': ['COMPLETED'],
        'COMPLETED': [],
        'CANCELLED': []
    };

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

<div class="space-y-5">
    <!-- Page Header & Actions -->
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div>
            <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">Orders</h1>
            <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Manage and track customer orders in real-time.</p>
        </div>
        <div class="flex gap-2">
             {#if isCustomer}
                <button on:click={() => goto("/orders/create")} class="glass-button flex items-center gap-2 text-sm bg-[var(--status-success)] shadow-emerald-500/20">
                    <Plus size={16} />
                    New Order
                </button>
             {/if}
            <button on:click={fetchOrders} class="flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm text-[var(--color-text-secondary)] bg-[var(--glass-bg)] border border-[var(--glass-border)] hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] transition-all backdrop-blur-sm">
                <RefreshCw size={15} />
                Refresh
            </button>
        </div>
    </div>

    <!-- Search & Status Filters -->
    <div class="flex flex-col sm:flex-row gap-3">
        <div class="flex-1">
            <input bind:value={searchQuery} placeholder="Search by ID or customer..." class="glass-input w-full px-4 text-sm" />
        </div>
        <div>
            <select bind:value={statusFilter} class="glass-input px-4 pr-8 text-sm font-semibold appearance-none cursor-pointer min-w-[150px]">
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

    <!-- Orders Data Table -->
    <div class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
            <table class="w-full text-left">
                <thead>
                    <tr class="border-b border-[var(--glass-border-subtle)]">
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider">Order</th>
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider hidden sm:table-cell">Customer</th>
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider hidden md:table-cell">Date</th>
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider">Total</th>
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider">Status</th>
                        <th class="px-3 sm:px-5 py-3 sm:py-3.5 font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider text-right">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-[var(--glass-border-subtle)]">
                    {#if loading}
                        <tr><td colspan="6" class="px-5 py-12 text-center text-[var(--color-text-secondary)] text-sm">Loading orders...</td></tr>
                    {:else if filteredOrders.length === 0}
                        <tr><td colspan="6" class="px-5 py-12 text-center text-[var(--color-text-secondary)] text-sm">No orders found.</td></tr>
                    {:else}
                        {#each filteredOrders as order}
                            <tr 
                                class="hover:bg-[var(--color-primary)]/[0.03] transition-colors cursor-pointer group" 
                                on:click={() => openOrderDetails(order.order_id)}
                                tabindex="0"
                                on:keydown={(e) => e.key === 'Enter' && openOrderDetails(order.order_id)}
                            >
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5 font-bold text-[var(--color-text-primary)] text-sm">#{order.order_id}</td>
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5 text-[var(--color-text-primary)] text-sm font-medium hidden sm:table-cell truncate max-w-[120px]">{order.user?.username}</td>
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5 text-[var(--color-text-secondary)] text-xs hidden md:table-cell">{order.created_at}</td>
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5 font-mono font-bold text-[var(--color-text-primary)] text-sm">${order.total_amount.toFixed(2)}</td>
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5">
                                    {#if canUpdateStatus && nextStatus[order.status.toUpperCase()]?.length > 0}
                                        <div class="relative inline-block">
                                            <select 
                                                class="appearance-none pl-2 sm:pl-3 pr-6 sm:pr-7 py-1.5 rounded-lg text-xs font-bold border cursor-pointer outline-none transition-all disabled:opacity-50 {getStatusColor(order.status)}"
                                                value={order.status.toUpperCase()}
                                                on:change={(e) => handleStatusChange(e, order)}
                                                disabled={updatingId === order.order_id}
                                                on:click|stopPropagation
                                            >
                                                <option value={order.status.toUpperCase()} disabled>{order.status}</option>
                                                {#each nextStatus[order.status.toUpperCase()] || [] as next}
                                                    <option value={next} class="text-[var(--color-text-primary)] bg-[var(--color-bg-surface)]">{next}</option>
                                                {/each}
                                            </select>
                                            <ChevronDown size={11} class="absolute right-1.5 sm:right-2 top-1/2 -translate-y-1/2 pointer-events-none opacity-50" />
                                        </div>
                                    {:else}
                                        <span class="px-2 sm:px-3 py-1.5 rounded-lg text-xs font-bold border {getStatusColor(order.status)}">
                                            {order.status}
                                        </span>
                                    {/if}
                                </td>
                                <td class="px-3 sm:px-5 py-3 sm:py-3.5 text-right">
                                    <button 
                                        on:click|stopPropagation={() => openOrderDetails(order.order_id)} 
                                        class="p-2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors rounded-lg hover:bg-[var(--glass-bg)]"
                                    >
                                        <Eye size={16} />
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

<!-- Order Details Modal — Line Items & Status Actions -->
{#if showDetailsModal && selectedOrder}
    <div 
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md" 
        role="dialog"
        aria-modal="true"
        tabindex="-1"
    >
        <button 
            type="button"
            class="absolute inset-0 w-full h-full cursor-default focus:outline-none"
            on:click={closeDetailsModal}
            on:keydown={(e) => e.key === 'Escape' && closeDetailsModal()}
            aria-label="Close modal"
        ></button>

        <div class="glass-card-solid w-full max-w-2xl shadow-2xl overflow-hidden flex flex-col max-h-[85vh] relative z-10">
            <div class="p-5 border-b border-[var(--glass-border-subtle)] flex justify-between items-start">
                <div>
                    <h2 class="text-xl font-bold text-[var(--color-text-primary)]">Order #{selectedOrder.order_id}</h2>
                    <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Customer: {selectedOrder.user?.username}</p>
                </div>
                <button on:click={closeDetailsModal} class="p-2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-xl hover:bg-[var(--glass-bg)] transition-colors">
                    <X size={20} />
                </button>
            </div>
            
            <div class="p-5 overflow-y-auto flex-1">
                <div class="space-y-2.5">
                    {#if selectedOrder.products}
                        {#each selectedOrder.products as item}
                            <div class="flex justify-between items-center p-3.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border-subtle)]">
                                <span class="font-semibold text-[var(--color-text-primary)] text-sm">{item.product?.name} <span class="text-[var(--color-text-secondary)] font-normal">×{item.quantity}</span></span>
                                <span class="font-mono text-[var(--color-primary)] font-bold text-sm">${item.line_total.toFixed(2)}</span>
                            </div>
                        {/each}
                    {/if}
                </div>
                <div class="mt-5 pt-5 border-t border-[var(--glass-border-subtle)] text-right">
                    <span class="text-2xl font-extrabold text-[var(--color-primary)]">${selectedOrder.total_amount.toFixed(2)}</span>
                </div>
            </div>

            {#if canUpdateStatus && nextStatus[selectedOrder.status.toUpperCase()]?.length > 0}
                <div class="p-5 border-t border-[var(--glass-border-subtle)] flex gap-2.5">
                    {#each nextStatus[selectedOrder.status.toUpperCase()] as status}
                        <button 
                            on:click={(e) => handleStatusChange({ target: { value: status } } as any, selectedOrder!)}
                            class="flex-1 py-2.5 rounded-xl font-bold text-white text-sm shadow-lg transition-all active:scale-[0.97]
                                   {status === 'CANCELLED' ? 'bg-[var(--status-error)] hover:opacity-90' : 'bg-[var(--color-primary)] hover:opacity-90'}"
                        >
                            Mark as {status}
                        </button>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
{/if}