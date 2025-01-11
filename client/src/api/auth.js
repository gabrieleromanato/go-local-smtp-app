import api from "./base";

export async function login(data) {
    try {
        const response = await api.post("/login", data);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}

export async function register(data) {
    try {
        const response = await api.post("/register", data);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}

export function logout() {
    localStorage.removeItem("app-token");
    localStorage.removeItem("app-user-id");
}

export function setToken(token) {
    localStorage.setItem("app-token", token);
}

export function setUserId(userId) {
    localStorage.setItem("app-user-id", userId);
}

export function getUserId() {
    return localStorage.getItem("app-user-id");
}

export function getToken() {
    return localStorage.getItem("app-token");
}

export function isAuthenticated() {
    return getToken() !== null;
}

export async function checkToken() {
    try {
        const response = await api.get("/check-token");
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}