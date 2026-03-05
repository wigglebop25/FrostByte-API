export interface Category {
    category_id: number;
    name: string;
    description: string;
}

export interface Product {
    product_id: number;
    name: string;
    product_image_uri: string;
    description: string;
    price: number;
    categories?: Category[];
}

export interface CreateProductPayload {
    name: string;
    description: string;
    price: number;
    product_image_uri: string;
    categories: string[];
}

export interface OrderProduct {
    order_id: number;
    product_id: number;
    product?: Product;
    quantity: number;
    unit_price: number;
    line_total: number;
}

export interface User {
    user_id: number;
    username: string;
    roles?: Role[];
}

export interface Role {
    role_id: number;
    name: string;
}

export interface Order {
    order_id: number;
    user_id: number;
    user?: User;
    total_amount: number;
    status: 'PENDING' | 'READY' | 'COMPLETED' | 'CANCELLED';
    products?: OrderProduct[];
    created_at: string;
}