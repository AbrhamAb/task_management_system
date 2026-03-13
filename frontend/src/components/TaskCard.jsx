const statusClassMap = {
	Pending: "bg-amber-100 text-amber-800",
	"In Progress": "bg-sky-100 text-sky-800",
	Completed: "bg-emerald-100 text-emerald-800",
};

const statuses = ["Pending", "In Progress", "Completed"];

export default function TaskCard({
	task,
	onEdit,
	onDelete,
	onStatusChange,
	isDeleting,
	isUpdatingStatus,
}) {
	return (
		<article className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow-md">
			<div className="mb-3 flex items-start justify-between gap-3">
				<h3 className="text-lg font-semibold text-slate-900">{task.title}</h3>
				<span
					className={`rounded-full px-3 py-1 text-xs font-semibold ${
						statusClassMap[task.status] || "bg-slate-100 text-slate-700"
					}`}
				>
					{task.status}
				</span>
			</div>

			<p className="mb-4 min-h-12 text-sm leading-6 text-slate-600">
				{task.description || "No description provided"}
			</p>

			<div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
				<select
					value={task.status}
					onChange={(event) => onStatusChange(task.id, event.target.value)}
					disabled={isUpdatingStatus}
					className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-700 focus:border-emerald-500 focus:outline-none sm:w-44"
				>
					{statuses.map((status) => (
						<option key={status} value={status}>
							{status}
						</option>
					))}
				</select>

				<div className="flex items-center gap-2">
					<button
						type="button"
						onClick={() => onEdit(task)}
						className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100"
					>
						Edit
					</button>
					<button
						type="button"
						onClick={() => onDelete(task.id)}
						disabled={isDeleting}
						className="rounded-lg bg-rose-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-70"
					>
						Delete
					</button>
				</div>
			</div>
		</article>
	);
}
