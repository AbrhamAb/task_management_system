"use client";

import { useMemo, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Navbar from "@/components/Navbar";
import TaskCard from "@/components/TaskCard";
import TaskForm from "@/components/TaskForm";
import { useTasks } from "@/hooks/useTasks";
import { useAuthStore } from "@/store/authStore";

export default function DashboardPage() {
  const router = useRouter();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const hydrated = useAuthStore((state) => state.hydrated);
  const [isPreviewMode, setIsPreviewMode] = useState(false);

  const {
    tasks,
    isLoadingTasks,
    tasksError,
    createTask,
    isCreatingTask,
    updateTask,
    isUpdatingTask,
    deleteTask,
    isDeletingTask,
    updateTaskStatus,
    isUpdatingTaskStatus,
  } = useTasks();

  const [editingTask, setEditingTask] = useState(null);
  const [previewTasks, setPreviewTasks] = useState([
    {
      id: "demo-1",
      title: "Build login page",
      description: "Create form fields and Yup validation for login flow.",
      status: "Completed",
      createdAt: new Date().toISOString(),
    },
    {
      id: "demo-2",
      title: "Design dashboard",
      description: "Implement task list layout with status chips and action buttons.",
      status: "In Progress",
      createdAt: new Date().toISOString(),
    },
    {
      id: "demo-3",
      title: "Integrate backend API",
      description: "Connect auth and task endpoints after backend is ready.",
      status: "Pending",
      createdAt: new Date().toISOString(),
    },
  ]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    const params = new URLSearchParams(window.location.search);
    setIsPreviewMode(params.get("preview") === "1");
  }, []);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    if (!isAuthenticated && !isPreviewMode) {
      router.replace("/login");
    }
  }, [hydrated, isAuthenticated, isPreviewMode, router]);

  const activeTasks = isPreviewMode ? previewTasks : tasks;
  const isLoading = isPreviewMode ? false : isLoadingTasks;
  const activeError = isPreviewMode ? null : tasksError;
  const isSubmitting = isPreviewMode ? false : isCreatingTask || isUpdatingTask;
  const isDeleting = isPreviewMode ? false : isDeletingTask;
  const isUpdatingStatus = isPreviewMode ? false : isUpdatingTaskStatus;

  const sortedTasks = useMemo(() => {
    return [...activeTasks].sort((a, b) => {
      const aDate = new Date(a.updatedAt || a.createdAt || 0).getTime();
      const bDate = new Date(b.updatedAt || b.createdAt || 0).getTime();
      return bDate - aDate;
    });
  }, [activeTasks]);

  const handleCreate = async (values) => {
    if (isPreviewMode) {
      const now = new Date().toISOString();
      setPreviewTasks((prev) => [
        {
          id: `demo-${Date.now()}`,
          ...values,
          createdAt: now,
          updatedAt: now,
        },
        ...prev,
      ]);
      return;
    }

    await createTask(values);
  };

  const handleUpdate = async (values) => {
    if (!editingTask?.id) {
      return;
    }

    if (isPreviewMode) {
      const now = new Date().toISOString();
      setPreviewTasks((prev) =>
        prev.map((task) =>
          task.id === editingTask.id ? { ...task, ...values, updatedAt: now } : task
        )
      );
      setEditingTask(null);
      return;
    }

    await updateTask({ id: editingTask.id, payload: values });
    setEditingTask(null);
  };

  const handleDelete = async (id) => {
    if (isPreviewMode) {
      setPreviewTasks((prev) => prev.filter((task) => task.id !== id));
      if (editingTask?.id === id) {
        setEditingTask(null);
      }
      return;
    }

    await deleteTask(id);
  };

  const handleStatusChange = async (id, status) => {
    if (isPreviewMode) {
      const now = new Date().toISOString();
      setPreviewTasks((prev) =>
        prev.map((task) =>
          task.id === id ? { ...task, status, updatedAt: now } : task
        )
      );
      return;
    }

    await updateTaskStatus({ id, status });
  };

  if (!hydrated || (!isAuthenticated && !isPreviewMode)) {
    return (
      <main className="grid min-h-screen place-items-center">
        <p className="text-sm font-medium text-slate-600">Checking your session...</p>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-slate-50">
      <Navbar />

      <div className="mx-auto grid w-full max-w-6xl gap-6 px-4 py-6 sm:px-6 lg:grid-cols-[360px_1fr]">
        <section>
          <h2 className="mb-3 text-lg font-bold text-slate-900">
            {editingTask ? "Edit Task" : "Create Task"}
          </h2>
          {isPreviewMode && (
            <p className="mb-3 rounded-lg border border-sky-200 bg-sky-50 px-3 py-2 text-xs font-medium text-sky-700">
              Preview mode: data is local and temporary.
            </p>
          )}
          <TaskForm
            initialValues={editingTask}
            onSubmit={editingTask ? handleUpdate : handleCreate}
            onCancel={editingTask ? () => setEditingTask(null) : undefined}
            submitLabel={editingTask ? "Update Task" : "Add Task"}
            isSubmitting={isSubmitting}
          />
        </section>

        <section>
          <div className="mb-4 flex items-end justify-between gap-3">
            <div>
              <h2 className="text-lg font-bold text-slate-900">Your Tasks</h2>
              <p className="text-sm text-slate-600">
                {sortedTasks.length} task{sortedTasks.length === 1 ? "" : "s"}
              </p>
            </div>
          </div>

          {isLoading && (
            <div className="rounded-xl border border-slate-200 bg-white p-5 text-sm text-slate-600">
              Loading tasks...
            </div>
          )}

          {activeError && (
            <div className="rounded-xl border border-rose-200 bg-rose-50 p-5 text-sm text-rose-700">
              {activeError.message}
            </div>
          )}

          {!isLoading && !activeError && sortedTasks.length === 0 && (
            <div className="rounded-xl border border-dashed border-slate-300 bg-white p-8 text-center text-sm text-slate-600">
              No tasks yet. Create your first task from the form.
            </div>
          )}

          <div className="space-y-3">
            {sortedTasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={setEditingTask}
                onDelete={handleDelete}
                onStatusChange={handleStatusChange}
                isDeleting={isDeleting}
                isUpdatingStatus={isUpdatingStatus}
              />
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
