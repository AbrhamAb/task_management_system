"use client";

import { create } from "zustand";

const AUTH_KEY = "tms_auth";
const TOKEN_KEY = "tms_token";

const readAuthFromStorage = () => {
	if (typeof window === "undefined") {
		return { token: null, user: null };
	}

	try {
		const raw = localStorage.getItem(AUTH_KEY);
		if (!raw) {
			return { token: null, user: null };
		}

		const parsed = JSON.parse(raw);
		return {
			token: parsed?.token || null,
			user: parsed?.user || null,
		};
	} catch {
		return { token: null, user: null };
	}
};

const persistAuth = ({ token, user }) => {
	if (typeof window === "undefined") {
		return;
	}

	if (!token) {
		localStorage.removeItem(AUTH_KEY);
		localStorage.removeItem(TOKEN_KEY);
		return;
	}

	localStorage.setItem(AUTH_KEY, JSON.stringify({ token, user }));
	localStorage.setItem(TOKEN_KEY, token);
};

export const useAuthStore = create((set) => ({
	token: null,
	user: null,
	isAuthenticated: false,
	hydrated: false,

	hydrate: () => {
		const { token, user } = readAuthFromStorage();
		set({
			token,
			user,
			isAuthenticated: Boolean(token),
			hydrated: true,
		});
	},

	setAuth: ({ token, user }) => {
		persistAuth({ token, user });
		set({ token, user, isAuthenticated: Boolean(token) });
	},

	logout: () => {
		persistAuth({ token: null, user: null });
		set({ token: null, user: null, isAuthenticated: false });
	},
}));
