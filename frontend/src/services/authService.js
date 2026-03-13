import api from "@/services/api";

const normalizeAuthPayload = (data) => {
	const token = data?.token || data?.accessToken || data?.data?.token;
	const user = data?.user || data?.data?.user || data?.profile || null;

	return { token, user };
};

export const authService = {
	register: async (payload) => {
		const response = await api.post("/auth/register", payload);
		return normalizeAuthPayload(response.data);
	},

	login: async (payload) => {
		const response = await api.post("/auth/login", payload);
		return normalizeAuthPayload(response.data);
	},

	getProfile: async () => {
		const response = await api.get("/auth/me");
		return response?.data?.user || response?.data?.data || response?.data;
	},
};
