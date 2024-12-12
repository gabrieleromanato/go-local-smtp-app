import api from "./base";

export async function login(data) {
    try {
        const response = await api.post("/login", data);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}

export function logout() {
    localStorage.removeItem("app-token");
}

export function setToken(token) {
    localStorage.setItem("app-token", token);
}

export function getToken() {
    return localStorage.getItem("app-token");
}

export function isAuthenticated() {
    return !!getToken();
}