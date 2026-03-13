import api from "@/services/api";

const unwrapTasks = (data) => {
	if (Array.isArray(data)) {
		return data;
	}

	if (Array.isArray(data?.tasks)) {
		return data.tasks;
	}

	if (Array.isArray(data?.data)) {
		return data.data;
	}

	return [];
};

export const taskService = {
	getTasks: async () => {
		const response = await api.get("/tasks");
		return unwrapTasks(response.data);
	},

	createTask: async (payload) => {
		const response = await api.post("/tasks", payload);
		return response?.data?.task || response?.data?.data || response?.data;
	},

	updateTask: async ({ id, payload }) => {
		const response = await api.put(`/tasks/${id}`, payload);
		return response?.data?.task || response?.data?.data || response?.data;
	},

	deleteTask: async (id) => {
		await api.delete(`/tasks/${id}`);
		return id;
	},

	updateTaskStatus: async ({ id, status }) => {
		const response = await api.patch(`/tasks/${id}/status`, { status });
		return response?.data?.task || response?.data?.data || response?.data;
	},
};
