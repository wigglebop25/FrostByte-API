/**
 * Product Image Resolver — Azure Blob Storage Integration
 *
 * Centralizes product image URL resolution across the application.
 * Supports multiple URI formats and gracefully falls back to Azure Blob
 * Storage when only a filename or product name is available.
 *
 * @module utils/image
 */

import type { Product } from "$lib/types";

/** Azure Blob Storage base URL for the product images container. */
export const AZURE_BLOB_BASE = "https://frostbytedata.blob.core.windows.net/products/";

/** Local fallback image displayed when remote assets fail to load. */
const FALLBACK_IMAGE = "/images/itadaki_logo.png";

/**
 * Resolves the display URL for a product image.
 *
 * Resolution priority:
 *  1. Absolute HTTP(S) URL — used as-is (e.g., full Azure Blob URL from DB).
 *  2. Relative path with separator — filename extracted and mapped to Blob Storage.
 *  3. Bare filename or slug — normalized and appended to Blob Storage base.
 *  4. No URI provided — product name converted to slug and resolved via Blob Storage.
 *
 * @param p - The product whose image URL should be resolved.
 * @returns A fully-qualified image URL suitable for an `<img>` element.
 */
export function getProductImageUrl(p: Product): string {
    if (p.product_image_uri && p.product_image_uri.startsWith('http')) {
        return p.product_image_uri;
    }

    if (p.product_image_uri && p.product_image_uri.includes('/')) {
        const filename = p.product_image_uri.split('/').pop()?.split('.')[0].replace(/-/g, '_') + '.png';
        return `${AZURE_BLOB_BASE}${filename}`;
    }

    if (p.product_image_uri) {
        const slug = p.product_image_uri.toLowerCase().replace(/[-\s]/g, '_').split('.')[0] + '.png';
        return `${AZURE_BLOB_BASE}${slug}`;
    }

    const nameSlug = p.name.toLowerCase().replace(/[-\s]/g, '_') + '.png';
    return `${AZURE_BLOB_BASE}${nameSlug}`;
}

/**
 * Global image error handler for product thumbnails.
 * Replaces the broken source with the local fallback logo to prevent
 * infinite error loops when the fallback itself is already set.
 *
 * @param e - The native browser error event from the `<img>` element.
 */
export function handleImageError(e: Event) {
    const target = e.currentTarget as HTMLImageElement;
    if (target.src !== window.location.origin + FALLBACK_IMAGE) {
        target.src = FALLBACK_IMAGE;
    }
}
