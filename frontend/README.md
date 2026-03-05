# Itadaki Restaurant Management System - Web Frontend

This is the official web client for the Itadaki Restaurant Management System. Built with SvelteKit and TypeScript, it provides a high-performance, real-time interface for administrators, cashiers, and customers.

## Key Modules

- **Dashboard**: Real-time analytics and performance metrics using ApexCharts and WebSocket integration.
- **Order Management**: Comprehensive workflow for tracking, updating, and managing customer orders.
- **Product & Inventory**: Admin interface for managing menu items, categories, and pricing.
- **Access Control**: Role-based access control (RBAC) interface for managing users and system permissions.
- **Settings**: User profile management and system infrastructure monitoring.

## Technical Architecture

- **Framework**: SvelteKit (Svelte 5)
- **Language**: TypeScript
- **Styling**: Tailwind CSS with custom glassmorphism design patterns.
- **State**: Reactive Svelte stores for authentication and real-time data streaming.
- **Real-time**: Gorilla WebSocket integration for live order synchronization.

## Integration

The frontend communicates with the FrostByte Go API. API configurations are managed via environment variables to support seamless deployment across environments.

---
Copyright 2026 Itadaki System
