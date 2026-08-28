import { apiFetch } from "./client";

export interface SignUpPayload {
  email: string;
  username: string;
  password: string;
}

export interface SignInPayload {
  username: string;
  password: string;
}

interface AuthResponse {
  token: string; // подправь под реальное поле, если бэк называет иначе
}

export function signUp(payload: SignUpPayload) {
  return apiFetch("/sign-up", {
    method: "POST",
    body: JSON.stringify(payload),
  }) as Promise<AuthResponse>;
}

export function signIn(payload: SignInPayload) {
  return apiFetch("/sign-in", {
    method: "POST",
    body: JSON.stringify(payload),
  }) as Promise<AuthResponse>;
}