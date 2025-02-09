import type { Name } from "$lib/api/types";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function formatName(name: Name) {
  let s = name.default;
  if (name.other) {
    s += ` - ${name.other}`;
  }

  return s;
}

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isRoleAdmin(role: string) {
  switch (role) {
    case "super_user":
    case "admin":
      return true;
    default:
      return false;
  }
}

export function getPagedQueryOptions(searchParams: URLSearchParams) {
  const query: Record<string, string> = {};
  const filter = searchParams.get("filter");
  if (filter) {
    query["filter"] = filter;
  }

  const sort = searchParams.get("sort");
  if (sort) {
    query["sort"] = sort;
  }

  const page = searchParams.get("page");
  if (page) {
    query["page"] = page;
  }

  return query;
}
