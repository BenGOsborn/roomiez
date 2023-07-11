import { writable } from "svelte/store";

export const location = writable<string>("");
export const radius = writable<number>(1);

export const age = writable<string>("");
export const duration = writable<string>("");
export const gender = writable<string>("");
export const rentalType = writable<string>("");
export const tenant = writable<string>("");

export const features = writable<{ [key: string]: boolean }>({});
