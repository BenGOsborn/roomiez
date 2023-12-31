import { writable } from "svelte/store";

export const email = writable<string>("");

export const page = writable<number>(1);

export const location = writable<string>("");
export const radius = writable<number>(1);

export const price = writable<number>(1501);
export const bond = writable<number>(3001);

export const age = writable<string>("");
export const duration = writable<string>("");
export const gender = writable<string>("");
export const rentalType = writable<string>("");
export const tenant = writable<string>("");

export const features = writable<{ [key: string]: boolean }>({});

export const mapReady = writable<boolean>(false);
