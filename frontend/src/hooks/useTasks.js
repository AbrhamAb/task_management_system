"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { taskService } from "@/services/taskService";
import { useAuthStore } from "@/store/authStore";

export const useTasks = () => {
	const queryClient = useQueryClient();
	const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

	const tasksQuery = useQuery({
		queryKey: ["tasks"],
		queryFn: taskService.getTasks,
		enabled: isAuthenticated,
	});

	const createTaskMutation = useMutation({
		mutationFn: taskService.createTask,
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tasks"] }),
	});

	const updateTaskMutation = useMutation({
		mutationFn: taskService.updateTask,
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tasks"] }),
	});

	const deleteTaskMutation = useMutation({
		mutationFn: taskService.deleteTask,
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tasks"] }),
	});

	const updateTaskStatusMutation = useMutation({
		mutationFn: taskService.updateTaskStatus,
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tasks"] }),
	});

	return {
		tasks: tasksQuery.data || [],
		isLoadingTasks: tasksQuery.isLoading,
		tasksError: tasksQuery.error,
		refetchTasks: tasksQuery.refetch,

		createTask: createTaskMutation.mutateAsync,
		isCreatingTask: createTaskMutation.isPending,

		updateTask: updateTaskMutation.mutateAsync,
		isUpdatingTask: updateTaskMutation.isPending,

		deleteTask: deleteTaskMutation.mutateAsync,
		isDeletingTask: deleteTaskMutation.isPending,

		updateTaskStatus: updateTaskStatusMutation.mutateAsync,
		isUpdatingTaskStatus: updateTaskStatusMutation.isPending,
	};
};
