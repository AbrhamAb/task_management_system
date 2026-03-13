"use client";

import { useEffect, useMemo, useState } from "react";
import { TASK_STATUS, taskSchema } from "@/validations/taskSchema";

const defaultValues = {
	title: "",
	description: "",
	status: "Pending",
};

export default function TaskForm({
	initialValues,
	onSubmit,
	onCancel,
	submitLabel,
	isSubmitting,
}) {
	const mode = useMemo(
		() => (initialValues?.id ? "edit" : "create"),
		[initialValues]
	);

	const [formValues, setFormValues] = useState(defaultValues);
	const [errors, setErrors] = useState({});
	const [submitError, setSubmitError] = useState("");

	useEffect(() => {
		setFormValues({
			title: initialValues?.title || "",
			description: initialValues?.description || "",
			status: initialValues?.status || "Pending",
		});
		setErrors({});
		setSubmitError("");
	}, [initialValues]);

	const handleChange = (event) => {
		const { name, value } = event.target;
		setFormValues((prev) => ({ ...prev, [name]: value }));
		if (submitError) {
			setSubmitError("");
		}
	};

	const handleSubmit = async (event) => {
		event.preventDefault();

		try {
			const validated = await taskSchema.validate(formValues, {
				abortEarly: false,
			});
			setErrors({});
			setSubmitError("");
			await onSubmit(validated);

			if (mode === "create") {
				setFormValues(defaultValues);
			}
		} catch (error) {
			if (!error?.inner) {
				setSubmitError(error?.message || "Unable to save task. Please try again.");
				return;
			}

			const nextErrors = {};
			error.inner.forEach((item) => {
				if (item.path && !nextErrors[item.path]) {
					nextErrors[item.path] = item.message;
				}
			});
			setErrors(nextErrors);
		}
	};

	return (
		<form onSubmit={handleSubmit} className="space-y-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
			<div>
				<label htmlFor="title" className="mb-1 block text-sm font-medium text-slate-700">
					Title
				</label>
				<input
					id="title"
					name="title"
					type="text"
					value={formValues.title}
					onChange={handleChange}
					placeholder="Design dashboard layout"
					className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 focus:border-emerald-500 focus:outline-none"
				/>
				{errors.title && <p className="mt-1 text-xs text-rose-600">{errors.title}</p>}
			</div>

			<div>
				<label htmlFor="description" className="mb-1 block text-sm font-medium text-slate-700">
					Description
				</label>
				<textarea
					id="description"
					name="description"
					rows={4}
					value={formValues.description}
					onChange={handleChange}
					placeholder="Add details for this task"
					className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 focus:border-emerald-500 focus:outline-none"
				/>
				{errors.description && (
					<p className="mt-1 text-xs text-rose-600">{errors.description}</p>
				)}
			</div>

			<div>
				<label htmlFor="status" className="mb-1 block text-sm font-medium text-slate-700">
					Status
				</label>
				<select
					id="status"
					name="status"
					value={formValues.status}
					onChange={handleChange}
					className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-emerald-500 focus:outline-none"
				>
					{TASK_STATUS.map((status) => (
						<option key={status} value={status}>
							{status}
						</option>
					))}
				</select>
				{errors.status && <p className="mt-1 text-xs text-rose-600">{errors.status}</p>}
			</div>

			<div className="flex gap-2">
				<button
					type="submit"
					disabled={isSubmitting}
					className="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-70"
				>
					{isSubmitting ? "Saving..." : submitLabel || "Save Task"}
				</button>
				{onCancel && (
					<button
						type="button"
						onClick={onCancel}
						className="rounded-lg border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-100"
					>
						Cancel
					</button>
				)}
			</div>
			{submitError && (
				<p className="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700">
					{submitError}
				</p>
			)}
		</form>
	);
}
