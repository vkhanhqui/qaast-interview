"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function AdminPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");

  useEffect(() => {
    const storedEmail = localStorage.getItem("email");
    const storedToken = localStorage.getItem("token");

    if (!storedEmail || !storedToken) {
      // Not signed in → go back to login
      router.replace("/");
      return;
    }
    setEmail(storedEmail);
  }, [router]);

  function handleLogout() {
    localStorage.removeItem("email");
    localStorage.removeItem("token");
    router.replace("/");
  }

  return (
    <main className="flex h-screen items-center justify-center bg-gray-100">
      <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow text-center">
        <div className="mb-6 flex items-center justify-center gap-4">
          <h1 className="text-xl font-semibold">Welcome {email}</h1>
          <button
            onClick={handleLogout}
            className="rounded bg-red-500 px-3 py-1 text-white hover:bg-red-600"
          >
            Logout
          </button>
        </div>

        <nav className="flex flex-col gap-4">
          <Link
            href="/admin/users"
            className="rounded bg-blue-500 p-2 text-white hover:bg-blue-600"
          >
            View Users
          </Link>
          <Link
            href="/admin/userlogs"
            className="rounded bg-green-500 p-2 text-white hover:bg-green-600"
          >
            View User Logs
          </Link>
        </nav>
      </div>
    </main>
  );
}
