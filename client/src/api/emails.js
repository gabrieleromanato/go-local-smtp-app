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

export async function sendEmail(email) {
    const data = new FormData();
    for (const key in email) {
        if (key === "attachments") {
            for (let i = 0; i < email.attachments.length; i++) {
                data.append("attachments", email.attachments[i]);
            }
            continue;
        }
        data.append(key, email[key]);
    }
    try {
        const response = await api.post("/api/emails", data);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}

export async function searchEmails(query = '', page = 1) {
    try {
        const response = await api.get(`/api/search?query=${query}&page=${page}`);
        return response.data;
    } catch (error) {
        throw error.response.data;
    }
}