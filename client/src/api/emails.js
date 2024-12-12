import api from "./base";

export async function getEmails(page = 1) {
    try {
        const response = await api.get(`/api/emails?page=${page}`);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}

export async function deleteEmail(id) {
    try {
        const response = await api.delete(`/api/emails/${id}`);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}