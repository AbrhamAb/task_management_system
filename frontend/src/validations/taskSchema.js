import * as yup from "yup";

export const TASK_STATUS = ["Pending", "In Progress", "Completed"];

export const taskSchema = yup.object({
	title: yup
		.string()
		.trim()
		.required("Task title is required")
		.max(100, "Task title cannot exceed 100 characters"),
	description: yup
		.string()
		.max(500, "Task description cannot exceed 500 characters")
		.nullable(),
	status: yup
		.string()
		.oneOf(TASK_STATUS, "Please choose a valid status")
		.required("Status is required"),
});
