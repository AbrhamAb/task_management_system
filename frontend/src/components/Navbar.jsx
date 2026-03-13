"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";

export default function Navbar() {
	const router = useRouter();
	const user = useAuthStore((state) => state.user);
	const logout = useAuthStore((state) => state.logout);

	const handleLogout = () => {
		logout();
		router.replace("/login");
	};

	return (
		<header className="sticky top-0 z-30 border-b border-emerald-100 bg-white/95 backdrop-blur">
			<div className="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-4 sm:px-6">
				<div>
					<p className="text-sm font-semibold uppercase tracking-[0.2em] text-emerald-700">
						Task Management
					</p>
					<h1 className="text-xl font-bold text-slate-900">Dashboard</h1>
				</div>
				<div className="flex items-center gap-3">
					<p className="hidden text-sm text-slate-600 sm:block">
						Signed in as <span className="font-semibold">{user?.name || "User"}</span>
					</p>
					<button
						type="button"
						onClick={handleLogout}
						className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-700"
					>
						Logout
					</button>
				</div>
			</div>
		</header>
	);
}
