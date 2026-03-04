# Product System Documentation

This document outlines the technical specifications, data model, and API integration details for the **Itadaki** product system. It is designed to assist developers in integrating product data into frontend web applications.

## 1. Data Model

The product entity is structured as follows.

### Schema

| Field Name | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `id` | `Integer` | Yes (System Assigned) | Unique identifier for the product. |
| `name` | `String` | Yes | The display name of the dish. |
| `description` | `String` | Yes | A short description of ingredients or preparation. |
| `price` | `Decimal` | Yes | Cost of the item (USD). |
| `product_image_uri` | `String` | Yes | Relative URL path to the product image (e.g., `/images/item.jpg`). |
| `categories` | `Array<String>` | Yes | List of categories the product belongs to (e.g., `["Zensai"]`). |

### TypeScript Interface (Web Frontend)

```typescript
interface Product {
 id?: number; // Optional for new creations
 name: string;
 description: string;
 price: number;
 product_image_uri: string;
 categories: string[];
}
```

---

## 2. API Endpoints

These endpoints are available via the base URL `https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1`.

### Retrieval
- **Get All Products**
 - `GET /products`
 - Returns a list of all active products.
- **Get Product by ID**
 - `GET /products/{id}`
 - Returns details for a specific product.
- **Get Products by Category**
 - `GET /categories/{category_name}/products`
 - Returns products filtered by category (e.g., `Menrui`).

### Management (Admin Only)
- **Create Product**
 - `POST /products`
 - Payload: `{ "name": "...", "price": 10.50, ... }`
- **Update Product**
 - `PUT /products/{id}`
- **Delete Product**
 - `DELETE /products/{id}`

---

## 3. Categories

The system is organized into the following standard categories:

| Category | Japanese | Description |
| :--- | :--- | :--- |
| **Zensai** | | Appetizers / Starters |
| **Sushi & Sashimi** | / | Raw fish and rice rolls |
| **Menrui** | | Noodle dishes (Ramen, Udon, Soba) |
| **Donburi** | | Rice bowls |
| **Kanmi** | | Desserts |

---

## 4. Asset Mapping

The Android application uses local drawable resources (`.png`), while the web API expects URL paths (`.jpg`). When integrating, ensure you map the assets correctly in your web project's public folder.

| Product Name | Android Resource (PNG) | Web API URI (JPG) |
| :--- | :--- | :--- |
| Pork Gyoza | `pork_gyoza.png` | `/images/pork-gyoza.jpg` |
| Edamame | `edamame.png` | `/images/edamame.jpg` |
| Shrimp Tempura | `shrimp_tempura.png` | `/images/shrimp-tempura.jpg` |
| Maguro Nigiri | `maguro_nigiri.png` | `/images/maguro-nigiri.jpg` |
| California Roll | `cali_roll.png` | `/images/cali-roll.jpg` |
| Salmon Sashimi | `salmon_sashimi.png` | `/images/salmon-sashimi.jpg` |
| Tonkotsu Ramen | `tonkotsu_ramen.png` | `/images/tonkotsu-ramen.jpg` |
| Tempura Udon | `tempura_udon.png` | `/images/tempura-udon.jpg` |
| Vegetable Yakisoba | `yakisoba.png` | `/images/yakisoba.jpg` |
| Gyu-Don | `gyudon.png` | `/images/gyudon.jpg` |
| Katsu-Don | `katsudon.png` | `/images/katsudon.jpg` |
| Unagi-Don | `unagi_don.png` | `/images/unagi-don.jpg` |
| Matcha Mochi | `matcha_mochi.png` | `/images/matcha-mochi.jpg` |
| Taiyaki | `taiyaki.png` | `/images/taiyaki.jpg` |
| Black Sesame Ice Cream | `sesame_ice_cream.png` | `/images/sesame-ice-cream.jpg` |

> **Note:** You will need to convert the PNG assets found in `app/src/main/res/drawable/` to JPG and place them in your web project's image directory to match the API response.
