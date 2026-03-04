import axios from "axios";
import { env } from "$env/dynamic/public";

const api = axios.create({
baseURL: env.PUBLIC_API_URL || "http://localhost:8080",
headers: {
"Content-Type": "application/json"
}
});

// Automatically add JWT token if it exists in localStorage
if (typeof window !== "undefined") {
api.interceptors.request.use((config) => {
const token = localStorage.getItem("token");
if (token) {
config.headers.Authorization = `Bearer ${token}`;
}
return config;
});
}

export default api;
